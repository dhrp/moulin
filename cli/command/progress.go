package command

import (
	"fmt"

	"github.com/dhrp/moulin/client"
	pb "github.com/dhrp/moulin/protobuf"
	"github.com/mitchellh/cli"
)

// Progress getting details of the progress of a given queue
type Progress struct {
	UI cli.Ui
}

// Run (LoadCommand) executes the actual action
func (c *Progress) Run(args []string) int {
	c.UI.Output("Listing the progress of a queue")

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	if len(args) > 1 {
		fmt.Println("received too many arguments")
		return -1
	} else if len(args) < 1 {
		fmt.Println("received too few arguments")
		return -1
	}

	task := new(pb.Task)
	task.QueueID = args[0]
	status, err := grpcDriver.Progress(args[0])
	if err != nil {
		fmt.Println("an error occured")
		return -1
	}

	fmt.Printf("incoming:  %d\n", status.IncomingCount)
	fmt.Printf("running:   %d\n", status.RunningCount)
	fmt.Printf("expired:   %d\n", status.ExpiredCount)
	fmt.Printf("completed: %d\n", status.CompletedCount)

	return 0
}

// Help (LoadCommand) shows help
func (c *Progress) Help() string {
	return "Get the progress of a queue, this shows the lenghts of the various parts."
}

// Synopsis is the short description
func (c *Progress) Synopsis() string {
	return "Get the progress of a queue"
}
