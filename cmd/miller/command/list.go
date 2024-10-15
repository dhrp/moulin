package command

import (
	"fmt"
	"log"

	"github.com/dhrp/moulin/client"
	"github.com/mitchellh/cli"
)

// List shows all the lists that exist
type List struct {
	UI cli.Ui
}

// Run (LoadCommand) executes the actual action
func (c *List) Run(args []string) int {
	c.UI.Output("Listing the queues")

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	if len(args) > 0 {
		fmt.Println("received too many arguments")
		return -1
	}

	queueMap, err := grpcDriver.ListQueues()
	if err != nil {
		log.Printf("Error: Could not list the queues: %v\n", err)
		return -1
	}

	for key, status := range queueMap {
		fmt.Printf("%s\n", key)

		fmt.Printf("  incoming:  %d\n", status.IncomingCount)
		fmt.Printf("  running:   %d\n", status.RunningCount)
		fmt.Printf("  expired:   %d\n", status.ExpiredCount)
		fmt.Printf("  completed: %d\n", status.CompletedCount)
		fmt.Printf("  failed:    %d\n", status.FailedCount)
	}
	return 0
}

// Help (LoadCommand) shows help
func (c *List) Help() string {
	return "Get the progress of a queue, this shows the lengths of the various parts."
}

// Synopsis is the short description
func (c *List) Synopsis() string {
	return "Get the progress of a queue"
}
