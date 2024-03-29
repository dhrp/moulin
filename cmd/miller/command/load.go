package command

import (
	"context"
	"fmt"
	"log"

	"github.com/dhrp/moulin/client"
	"github.com/mitchellh/cli"
)

// LoadCommand is for loading
type LoadCommand struct {
	UI cli.Ui
}

// Run (LoadCommand) executes the actual action
func (c *LoadCommand) Run(args []string) int {
	c.UI.Output("Loading item from queue")

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	if len(args) > 1 {
		fmt.Println("received too many arguments for queue")
		return -1
	} else if len(args) < 1 {
		fmt.Println("received too few arguments for queue")
		return -1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task, err := grpcDriver.LoadTask(ctx, args[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("received task %s from queue\n", task)

	return 0
}

// Help (LoadCommand) shows help
func (c *LoadCommand) Help() string {
	return "Load and return an item"
}

// Synopsis is the short description
func (c *LoadCommand) Synopsis() string {
	return "Load an item"
}
