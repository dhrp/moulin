package main

import (
	"fmt"

	"github.com/nerdalize/moulin/client"
)

func main() {

	grpcDriver := client.NewGRPCDriver()

	defer grpcDriver.Connection.Close()

	status := grpcDriver.GetHealth()
	fmt.Println(status)

	task := grpcDriver.LoadTask(queueID)
	fmt.Println("loaded: " + task.TaskID)
	fmt.Println("body: " + task.Body)

}
