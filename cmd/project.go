package cmd

import (
	"github.com/azisuazusa/todo-cli/internal/presenter/project"
	"github.com/urfave/cli/v2"
)

func projectCLI(presenter *project.Presenter) *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "Manage projects",
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a new project",
				Action: func(c *cli.Context) error {
					return presenter.Add(c.Context)
				},
			},
			{
				Name:  "remove",
				Usage: "Remove a project",
				Action: func(c *cli.Context) error {
					return presenter.Remove(c.Context)
				},
			},
			{
				Name:  "list",
				Usage: "List all projects",
				Action: func(c *cli.Context) error {
					return presenter.GetProjects(c.Context)
				},
			},
			{
				Name:  "sync-task",
				Usage: "Sync tasks for a project",
				Action: func(c *cli.Context) error {
					return presenter.SyncTasks(c.Context)
				},
			},
			{
				Name:  "select",
				Usage: "Select a project",
				Action: func(c *cli.Context) error {
					return presenter.Select(c.Context)
				},
			},
			{
				Name:  "add-integration",
				Usage: "Add an integration to a project",
				Action: func(c *cli.Context) error {
					return presenter.AddIntegration(c.Context)
				},
			},
			{
				Name:  "remove-integration",
				Usage: "Remove an integration from a project",
				Action: func(c *cli.Context) error {
					return presenter.RemoveIntegration(c.Context)
				},
			},
			{
				Name:  "enable-integration",
				Usage: "Enable an integration for a project",
				Action: func(c *cli.Context) error {
					return presenter.EnableIntegration(c.Context)
				},
			},
			{
				Name:  "disable-integration",
				Usage: "Disable an integration for a project",
				Action: func(c *cli.Context) error {
					return presenter.DisableIntegration(c.Context)
				},
			},
		},
	}
}
