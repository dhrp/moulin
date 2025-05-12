package command

import (
	"flag"
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

	// Define a flag set for parsing arguments
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	progressFlag := fs.Bool("progress", false, "Show progress details for each queue")
	sortFlag := fs.String("sort", "", "Sort queues by specified criteria (e.g., 'alpha', 'created')")

	// Parse the flags
	if err := fs.Parse(args); err != nil {
		fmt.Println("Error parsing flags:", err)
		return -1
	}

	// Fetch the list of queues
	queueList, err := grpcDriver.ListQueues(*sortFlag)
	if err != nil {
		log.Printf("Error: Could not list the queues: %v\n", err)
		return -1
	}

	// Print the queues
	for _, queue := range queueList.Queues {
		fmt.Printf("%s\n", queue.QueueID)

		// Only print progress details if --progress is specified
		if *progressFlag {
			fmt.Printf("  incoming:  %d\n", queue.Progress.IncomingCount)
			fmt.Printf("  running:   %d\n", queue.Progress.RunningCount)
			fmt.Printf("  expired:   %d\n", queue.Progress.ExpiredCount)
			fmt.Printf("  completed: %d\n", queue.Progress.CompletedCount)
			fmt.Printf("  failed:    %d\n", queue.Progress.FailedCount)
		}
	}
	return 0
}

// Help (LoadCommand) shows help
func (c *List) Help() string {
	return `List all queues in the system, optionally include their progress.

Usage: moulin-cli list [--progress] [--sort <sort>]

Options:
  --progress       Show progress details for each queue
  --sort <sort>    Sort queues by specified criteria (e.g., 'alpha', 'created')
                   Prefix with '-' to reverse the order

Examples:
  # List alphabetically
  moulin-cli list --sort alpha

  # List by time created (most recent first)
  moulin-cli list --sort -created

  # Include progress details
  moulin-cli list --progress --sort alpha
`
}

// Synopsis is the short description
func (c *List) Synopsis() string {
	return "List all queues in the system"
}
