package main

import (
	"fmt"
	"os"
	"pugo/pkg/cmd"
	"pugo/pkg/constants"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	version                 = "dev"
	commands []*cli.Command = []*cli.Command{
		cmd.NewInit(),
		cmd.NewBuild(),
		cmd.NewCreate(),
		cmd.NewServer(),
		{
			Name:  "version",
			Usage: "print the version of PuGo",
			Action: func(c *cli.Context) error {
				fmt.Println(constants.AppName(), version)
				return nil
			},
		},
	}
)

func main() {
	constants.SetAppVersion(version)

	app := &cli.App{
		Name:     constants.AppName(),
		Usage:    "a simple static site generator with markdown support",
		Version:  constants.AppVersion(),
		Commands: commands,
		Flags:    cmd.GetGlobalFlags(),
	}
	args := movePostfixOptions(os.Args)
	app.Run(args)
}

// Function to reorder arguments in "correct" order for urfave/cli
// Copied from https://github.com/ipfs/ipget/blob/5397b0666d7e90d78c1566ecb90f289dad9d9ec1/main.go#L142
// And changed start index from 1 to 2.
func movePostfixOptions(args []string) []string {
	var endArgs []string
	for idx := 2; idx < len(args); idx++ {
		if args[idx][0] == '-' {
			if !strings.Contains(args[idx], "=") {
				idx++
			}
			continue
		}
		if endArgs == nil {
			// on first write, make copy of args
			newArgs := make([]string, len(args))
			copy(newArgs, args)
			args = newArgs
		}
		// add to args accumulator
		endArgs = append(endArgs, args[idx])
		// remove from real args list
		args = args[:idx+copy(args[idx:], args[idx+1:])]
		idx--
	}

	// append extracted arguments to the real args
	return append(args, endArgs...)
}
