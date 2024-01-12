package main

import (
	"fmt"
	"os"

	"github.com/dhrp/moulin/cmd/miller/command"
	"github.com/mitchellh/cli"
)

var (
	// Version is set during build with ldflags
	Version string
	// Build is set during build with ldflags
	Build string
)

func main() {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("moulin-cli", Version)
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
		"list": func() (cli.Command, error) {
			return &command.List{
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
