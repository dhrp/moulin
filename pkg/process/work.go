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

var loadTaskTimeOut = client.ClientConfig.LoadTaskTimeOut
var heartBeatInterval = client.ClientConfig.HeartBeatInterval
var serverUnavailableTimeOut = client.ClientConfig.ServerUnavailableTimeOut

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

	// Create a master context for the worker
	workCtx, workCancel := context.WithCancel(context.Background())
	defer workCancel()

	// forever loop, until exit == true
	for {

		// Check if the one-hour timeout context is done
		select {
		case <-workCtx.Done():
			fmt.Println("Context stopped, not loading new tasks.")
			return 0, nil
		default:
			// Continue loading work
		}

		fmt.Printf("Loading a task from the queue.\n")
		var loadCtx context.Context
		var loadCancel context.CancelFunc

		if setTimeOut {
			loadCtx, loadCancel = context.WithTimeout(workCtx, loadTaskTimeOut)
		} else {
			loadCtx, loadCancel = context.WithCancel(workCtx)
		}

		task, err := grpcDriver.LoadTask(loadCtx, queueID)
		loadCancel() // cancel the load context

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

		go RepeatHeartBeat(workCtx, workCancel, grpcDriver, queueID, task.TaskID)

		// let the exec function do the hard work
		result, err := Exec(task)
		if err != nil {
			fmt.Printf("task failed with result %v!!", result)
		}

		if result == 0 {
			fmt.Printf("Task completed with exit code %d.\n", result)
			fmt.Printf("Marking task as completed: %s\n", task)
			status := grpcDriver.Complete(workCtx, queueID, task.TaskID)
			fmt.Println(status.Detail)
		} else {
			fmt.Printf("Task failed with exit code %d.\n", result)
			fmt.Printf("Marking task as failed: %s\n", task)
			status := grpcDriver.Fail(workCtx, queueID, task.TaskID)
			fmt.Println(status.Detail)
		}

		if exit == true {
			return 0, nil
		}
	}
}

// RepeatHeartBeat calls the heartbeat function repeatedly, until the task completes
func RepeatHeartBeat(
	workCtx context.Context,
	workCancel context.CancelFunc,
	grpcDriver *client.GRPCDriver,
	queueID string, taskID string,
) {
	var expires = time.Now().Add(serverUnavailableTimeOut)

	for {
		select {
		case <-time.After(time.Until(expires)):
			// we close the work context so other dependent processes, notably
			// the complete() function also stop retrying
			log.Println("The server timeout has exceeded, cancelling work context.")
			log.Println("Ongoing work will be finished, but will never be marked as completed.")
			workCancel()
		case <-workCtx.Done():
			log.Println("Stopping heartbeat for taskID", taskID)
			return
		case <-time.After(heartBeatInterval):
			_, err := grpcDriver.HeartBeat(queueID, taskID)
			if err != nil {
				log.Println("could not complete heartbeat:", err)
			} else {
				log.Println("heartbeat beat")
				expires = time.Now().Add(serverUnavailableTimeOut)
			}
		}
	}
}
