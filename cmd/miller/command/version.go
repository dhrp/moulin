package command

import (
	"github.com/mitchellh/cli"
)

// Version is for displaying the version
type Version struct {
	UI cli.Ui
}

// Run (LoadCommand) executes the actual action
func (c *Version) Run(args []string) int {
	c.UI.Output("Version ")
	return 0
}

// Help (LoadCommand) shows help
func (c *Version) Help() string {
	return "Get the progress of a queue, this shows the lenghts of the various parts."
}

// Synopsis is the short description
func (c *Version) Synopsis() string {
	return "Get the progress of a queue"
}
