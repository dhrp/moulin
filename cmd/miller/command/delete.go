package command

import (
	"fmt"
	"log"

	"github.com/dhrp/moulin/client"
	"github.com/mitchellh/cli"
	"google.golang.org/grpc/status"
)

// DeleteQueue is for creating a single task
type DeleteQueue struct {
	UI cli.Ui
}

// Run (LoadCommand) executes the actual action
func (c *DeleteQueue) Run(args []string) int {
	c.UI.Output("Delete a queue and all related tasks")

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	if len(args) > 1 {
		fmt.Println("received too many arguments")
		return -1
	} else if len(args) < 1 {
		fmt.Println("received too few arguments")
		return -1
	}

	res, err := grpcDriver.DeleteQueue(args[0])

	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Error: Could not delete queue: %v\n", st.Message())
		return -1
	}

	fmt.Println(res.Detail)
	return 0
}

// Synopsis is the short description
func (c *DeleteQueue) Synopsis() string {
	return "Delete a queue and all related tasks"
}

// Help (LoadCommand) shows help
func (c *DeleteQueue) Help() string {
	return `Delete a queue and all related tasks

Usage: moulin-cli delete-queue <queueID>

## Basic example
  moulin-cli delete-queue my-queue

Note that this action is immediate and irreversible. All tasks in the queue will be lost.
`
}
