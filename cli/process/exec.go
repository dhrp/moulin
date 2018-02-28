package process

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/mattn/go-shellwords"

	pb "github.com/dhrp/moulin/protobuf"
)

const defaultFailedCode = 1

var exitCode int

// Exec executes a program
func Exec(task *pb.Task) (result int, err error) {

	var stdin []byte
	// var args []string

	// build up a context that is passed to the arg. If this context ends,
	// the command is also killed
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel = context.WithCancel(ctx)
	defer cancel() //cancel heartbeat context if this function exits

	taskCommand, err := shellwords.Parse(task.Body)
	if err != nil {
		fmt.Println("failed parsing arguments")
	}

	var command string
	var args []string

	if len(taskCommand) > 0 {
		command = taskCommand[0]
	}
	if len(taskCommand) > 1 {
		args = taskCommand[1:]
	}

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdin = bytes.NewBuffer(stdin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	for _, e := range os.Environ() {
		cmd.Env = append(cmd.Env, e)
	}

	envs := map[string]string{
		"key": "value",
	}

	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	fmt.Printf("  command: %v; args: %v\n", command, args)
	fmt.Print("> ")
	err = cmd.Run()
	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", command, args)
			exitCode = defaultFailedCode
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	return exitCode, nil
}
