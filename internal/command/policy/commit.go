package policy

import (
	"raind/internal/core/policy"

	"github.com/urfave/cli/v2"
)

func CommandCommit() *cli.Command {
	return &cli.Command{
		Name:   "commit",
		Usage:  "commit policy",
		Action: runCommit,
	}
}

func runCommit(ctx *cli.Context) error {
	service := policy.NewServicePolicyCommit()
	err := service.Commit()
	if err != nil {
		return err
	}

	return nil
}
