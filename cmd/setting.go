package cmd

import (
	"github.com/azisuazusa/todo-cli/internal/presenter/setting"
	"github.com/urfave/cli/v2"
)

func settingCLI(presenter *setting.Presenter) *cli.Command {
	return &cli.Command{
		Name:  "setting",
		Usage: "Manage settings",
		Subcommands: []*cli.Command{
			{
				Name:  "sync",
				Usage: "Sync todo-cli",
				Action: func(c *cli.Context) error {
					return presenter.Sync(c.Context)
				},
			},
			{
				Name:  "set-integration",
				Usage: "Manage integrations",
				Action: func(c *cli.Context) error {
					return presenter.SetIntegration(c.Context)
				},
			},
		},
	}
}
