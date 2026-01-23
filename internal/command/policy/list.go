package policy

import (
	"raind/internal/core/policy"

	"github.com/urfave/cli/v2"
)

func CommandList() *cli.Command {
	return &cli.Command{
		Name:  "ls",
		Usage: "list policies",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
				Usage:    "policy type (ew/ns-obs/ns-enf)",
				Required: true,
			},
		},
		Action: runList,
	}
}

func runList(ctx *cli.Context) error {
	// option
	opt_type, err := validateType(ctx.String("type"))
	if err != nil {
		return err
	}

	service := policy.NewServicePolicyList()
	err = service.List(
		policy.ListRequestModel{
			ChainName: opt_type,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
