package process

import (
	"fmt"
	"log"

	"github.com/dhrp/moulin/pkg/client"
)

// Work manages getting, heartbeating, and completing or failing
// items, in a loop
func Work(grpcDriver *client.GRPCDriver, queueID, workType string) (result int, err error) {

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
		fmt.Printf("  received taskID %s from queue\n", task.TaskID)
		// fmt.Println("  nap")
		// time.Sleep(2000 * time.Millisecond)
		// fmt.Println("  /nap")

		// let the exec function do the hard work
		result, err := Exec(task)
		if err != nil {
			fmt.Printf("  task failed with result %v!!", result)
			// ToDo: mark as failed
		}

		if result == 0 {
			fmt.Printf("  Task completed with code %d. Marking as complete.\n", result)
		} else {
			fmt.Printf("  Task failed with code %d ?!? (still marking as complete for now)", result)
		}

		status := grpcDriver.Complete(queueID, task.TaskID)
		fmt.Println(status)

		if exit == true {
			return 0, nil
		}
	}

}
