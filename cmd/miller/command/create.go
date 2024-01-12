package command

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dhrp/moulin/client"
	pb "github.com/dhrp/moulin/pkg/protobuf"
	"github.com/mitchellh/cli"
)

// CreateTask is for creating a single task
type CreateTask struct {
	UI cli.Ui
}

// Run (LoadCommand) executes the actual action
func (c *CreateTask) Run(args []string) int {
	c.UI.Output("Creating a single task")

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	if len(args) > 2 {
		fmt.Println("received too many arguments")
		return -1
	} else if len(args) < 2 {
		fmt.Println("received too few arguments")
		return -1
	}

	task := new(pb.Task)
	task.QueueID = args[0]

	if args[1] == "-" {
		fmt.Println("reading from stdin")
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			task.Body = s.Text()
			status, err := grpcDriver.PushTask(task)
			if err != nil {
				log.Fatalf("failed to push task: %v\n", err)
			}
			fmt.Printf("submit task %v\n", status)
		}
	} else {
		fmt.Println("reading argument")
		task.Body = args[1]
		status, err := grpcDriver.PushTask(task)
		if err != nil {
			log.Fatalf("failed to push task: %v\n", err)
		}
		fmt.Printf("submit task %v\n", status)
	}
	return 0
}

// Help (LoadCommand) shows help
func (c *CreateTask) Help() string {
	return `Create a task with a command to execute

Usage: moulin-cli create <queueID> <body>

## Basic usage

Create a task on any queue with a body. The body is the command that will be
executed by the worker. Quote it if it has multiple words. You are free to choose
any name for the queue, and using something like a '.' can help you organize your
queues. 

Example:
  moulin-cli create myQueue "echo hello"
  moulin-cli create myQueue.v2 "echo hello"


## Formatting your command

Keep in mind that the full body should appear and work as a single command on the 
shell of your worker. In practice this means that if you want to execute multiple 
commands you should wrap your commands in one command like. 
  
Example:  
	moulin-cli create myQueue "bash -c \"echo hello && sleep 5 && echo world\""


## Batch creating tasks

Often times you want to create a large number of tasks. You can do this by piping
newline separated plain text into this cli. Each line will be created as a separate
task. 

Examples:
  echo "echo hello" | moulin-cli create myQueue -
  cat commands.txt | moulin-cli create myQueue -

`
}

// Synopsis is the short description
func (c *CreateTask) Synopsis() string {
	return "Create a task with body"
}
