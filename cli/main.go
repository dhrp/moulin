package main

import (
	"fmt"
	"os"

	"github.com/dhrp/moulin/cli/command"
	"github.com/mitchellh/cli"
)

func main() {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("cliexample", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{

		"load": func() (cli.Command, error) {
			return &command.LoadCommand{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"work": func() (cli.Command, error) {
			return &command.Work{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"create": func() (cli.Command, error) {
			return &command.CreateTask{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"progress": func() (cli.Command, error) {
			return &command.Progress{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
		"peek": func() (cli.Command, error) {
			return &command.Peek{
				UI: &cli.ColoredUi{
					Ui:          ui,
					OutputColor: cli.UiColorGreen,
				},
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
