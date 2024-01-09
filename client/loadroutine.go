package client

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
)

// LoadRoutine is not a test on its own.
func LoadRoutine(grpcDriver *GRPCDriver, number int, channel chan bool) {
	fmt.Printf("[goroutine %d] loadroutine started\n", number)

	connectionState := grpcDriver.Connection.GetState()
	fmt.Printf("[goroutine %d] connection state: %v \n", number, connectionState)
	_, err := grpcDriver.GetHealth()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[goroutine %d] health reveived\n", number)

	connectionState = grpcDriver.Connection.GetState()
	fmt.Printf("[goroutine %d] connection state: %v \n", number, connectionState)

	returnedTask, err := grpcDriver.LoadTask(context.Background(), "clientTest2")
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Printf("[goroutine %d] loadroutine done, loaded: %s\n", number, returnedTask.Body)
	// suite.Equal(len("0vNrL62AGAdIzRZ9pReEnKeMu4x"), len(returnedTask.TaskID), "TaskID doesn't look valid")
	// suite.Equal("Task #1", returnedTask.Body, "Body doesn't match")

	if returnedTask.Body != "Task #1" {
		channel <- false
	} else {
		channel <- true
	}
}
