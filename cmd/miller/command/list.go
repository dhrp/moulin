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

	sortBy := ""

	if len(args) > 1 {
		fmt.Println("received too many arguments")
		return -1
	}

	if len(args) > 0 {
		sortBy = args[0]
	}

	queueList, err := grpcDriver.ListQueues(sortBy)
	if err != nil {
		log.Printf("Error: Could not list the queues: %v\n", err)
		return -1
	}

	for _, queue := range queueList.Queues {
		fmt.Printf("%s\n", queue.QueueID)

		fmt.Printf("  incoming:  %d\n", queue.Progress.IncomingCount)
		fmt.Printf("  running:   %d\n", queue.Progress.RunningCount)
		fmt.Printf("  expired:   %d\n", queue.Progress.ExpiredCount)
		fmt.Printf("  completed: %d\n", queue.Progress.CompletedCount)
		fmt.Printf("  failed:    %d\n", queue.Progress.FailedCount)
	}
	return 0
}

// Help (LoadCommand) shows help
func (c *List) Help() string {
	return `List all queues in the system, includes their progress.

Usage: moulin-cli list <sort>

Example:
  # List alphabetically
  moulin-cli list alpha
  moulin-cli list -alpha (descending order)
  
  # List by time created
  moulin-cli list created
  moulin-cli list -created (descending order)
`
}

// Synopsis is the short description
func (c *List) Synopsis() string {
	return "Get the progress of a queue"
}
