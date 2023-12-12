package process

import (
	"fmt"
	"log"
	"time"

	"github.com/dhrp/moulin/client"
)

// Work manages getting, heartbeating, and completing or failing
// items, in a loop
func Work(grpcDriver *client.GRPCDriver, queueID string, workType string) (result int, err error) {

	var exit bool
	if workType == "once" {
		exit = true
	} else if workType == "until-finished" {
		exit = false
	} else if workType == "forever" {
		exit = false
	}

	// forever loop, until exit == true
	for {
		task, err := grpcDriver.LoadTask(queueID)
		if err != nil {
			log.Panic("failed loading task")
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
		time.Sleep(30 * time.Second)

		if isWorking.active == false {
			fmt.Println("   done waiting")
			break
		}

		grpcDriver.HeartBeat(queueID, taskID, 300)
		fmt.Println("   heartbeat beat")
	}
}
