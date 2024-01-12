package process

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dhrp/moulin/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var loadTaskTimeOut = 30 * time.Second
var heartBeatInterval = 30 * time.Second
var heartBeatTimeOutSeconds int32 = 300 // 5 minutes

// Work manages getting, heartbeating, and completing or failing
// items, in a loop
func Work(grpcDriver *client.GRPCDriver, queueID string, workType string) (result int, err error) {

	var exit bool // whether to stop after completion of a task.
	var setTimeOut = false

	if workType == "once" {
		exit = true
	} else if workType == "until-finished" {
		exit = false
		setTimeOut = true
	} else if workType == "forever" {
		exit = false
	}

	// forever loop, until exit == true
	for {
		fmt.Printf("Loading a task from the queue.\n")
		ctx, cancel := context.WithCancel(context.Background())

		if setTimeOut {
			ctx, cancel = context.WithTimeout(ctx, loadTaskTimeOut)
		}
		defer cancel()

		task, err := grpcDriver.LoadTask(ctx, queueID)
		if status.Code(err) == codes.DeadlineExceeded {
			fmt.Printf("LoadTask Timeout (%s) has exceeded. The queue is empty, stopping.\n", loadTaskTimeOut)
			return 0, err
		}
		if err != nil {
			log.Printf("%s", status.Code(err))
			log.Printf("failed loading task")
			return 1, err
		}
		fmt.Printf("received task %s from queue\n", task)

		var isWorking = struct{ active bool }{active: true}
		go RepeatHeartBeat(grpcDriver, queueID, task.TaskID, &isWorking)

		// let the exec function do the hard work
		result, err := Exec(task)
		if err != nil {
			fmt.Printf("task failed with result %v!!", result)
		}

		isWorking.active = false

		if result == 0 {
			fmt.Printf("Task completed with exit code %d.\n", result)
			fmt.Printf("Marking task as completed: %s\n", task)
			status := grpcDriver.Complete(queueID, task.TaskID)
			fmt.Println(status.Detail)
		} else {
			fmt.Printf("Task failed with exit code %d.\n", result)
			fmt.Printf("Marking task as failed: %s\n", task)
			status := grpcDriver.Fail(queueID, task.TaskID)
			fmt.Println(status.Detail)
		}

		if exit == true {
			return 0, nil
		}
	}
}

// RepeatHeartBeat calls the heartbeat function repeatedly, until the task completes
func RepeatHeartBeat(grpcDriver *client.GRPCDriver, queueID string, taskID string, isWorking *struct{ active bool }) {
	for {
		time.Sleep(heartBeatInterval)

		if isWorking.active == false {
			fmt.Println("   done waiting")
			break
		}

		grpcDriver.HeartBeat(queueID, taskID, heartBeatTimeOutSeconds)
		fmt.Println("   heartbeat beat")
	}
}
