package command

import (
	containercommand "raind/internal/command/container"
	imagecommand "raind/internal/command/image"
	logscommand "raind/internal/command/logs"
	policycommand "raind/internal/command/policy"

	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	app := &cli.App{
		Name:  "raind",
		Usage: "raind container runtime",
		Commands: []*cli.Command{
			{
				Name:  "container",
				Usage: "container operation",
				Subcommands: []*cli.Command{
					containercommand.CommandCreate(),
					containercommand.CommandStart(),
					containercommand.CommandStop(),
					containercommand.CommandRemove(),
					containercommand.CommandList(),
					containercommand.CommandAttach(),
					containercommand.CommandRun(),
					containercommand.CommandExec(),
					containercommand.CommandLogs(),
				},
			},
			{
				Name:  "image",
				Usage: "image operation",
				Subcommands: []*cli.Command{
					imagecommand.CommandPull(),
					imagecommand.CommandList(),
					imagecommand.CommandRemove(),
				},
			},
			{
				Name:  "policy",
				Usage: "policy operation",
				Subcommands: []*cli.Command{
					policycommand.CommandCreate(),
					policycommand.CommandList(),
					policycommand.CommandCommit(),
					policycommand.CommandRemove(),
					policycommand.CommandRevert(),
					policycommand.CommandChangeMode(),
				},
			},
			{
				Name:  "logs",
				Usage: "log operation",
				Subcommands: []*cli.Command{
					logscommand.CommandLogs(),
				},
			},
		},
	}

	// disable slice flag separator
	app.DisableSliceFlagSeparator = true

	return app
}
