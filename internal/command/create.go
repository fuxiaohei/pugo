package command

import "github.com/urfave/cli/v2"

// NewCreateCommand returns a new cli.Command for the create subcommand.
func NewCreateCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "create",
		Usage:       "create new post or page",
		Description: "create new post or page with default template",
		Flags:       append(GlobalFlags, buildFlags...),
		Action: func(c *cli.Context) error {

			initGlobalFlags(c)

			return nil
		},
	}
	return cmd
}
