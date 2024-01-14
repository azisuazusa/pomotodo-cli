package main

import "github.com/urfave/cli/v2"

func createApp() *cli.App {
	return &cli.App{
		Name:  "todo",
		Usage: "A todo list with timeboxing for increase productivity",
		Action: func(c *cli.Context) error {
			println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "task",
				Usage: "Manage tasks",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add a new task to the list",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Aliases:  []string{"n"},
								Usage:    "Add a name to the task",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "description",
								Aliases:  []string{"d"},
								Usage:    "Add a description to the task",
								Required: false,
							},
							&cli.BoolFlag{
								Name:     "subtask",
								Aliases:  []string{"s"},
								Usage:    "Add a subtask to the task",
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							task := Task{
								Name:        c.String("name"),
								Description: c.String("description"),
							}

							if c.Bool("subtask") {
								return task.AddSubTask()
							}

							return task.Add()
						},
					},
					{
						Name:  "complete",
						Usage: "Complete a task on the list",
						Action: func(c *cli.Context) error {
							return TaskComplete()
						},
					},
					{
						Name:  "list",
						Usage: "List all tasks on the list",
						Action: func(c *cli.Context) error {
							return TaskList()
						},
					},
					{
						Name:  "start",
						Usage: "Start a task on the list",
						Action: func(c *cli.Context) error {
							return TaskStart()
						},
					},
					{
						Name:  "stop",
						Usage: "Stop a task on the list",
						Action: func(c *cli.Context) error {
							return TaskStop()
						},
					},
					{
						Name:  "remove",
						Usage: "Remove a task on the list",
						Action: func(c *cli.Context) error {
							return TaskRemove()
						},
					},
				},
			},
			{
				Name:  "project",
				Usage: "Manage projects",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add a new project",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Aliases:  []string{"n"},
								Usage:    "Add a name to the project",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "description",
								Aliases:  []string{"d"},
								Usage:    "Add a description to the project",
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							project := Project{
								Name:        c.String("name"),
								Description: c.String("description"),
							}
							return project.Add()
						},
					},
					{
						Name:  "remove",
						Usage: "Remove a project",
						Action: func(c *cli.Context) error {
							return ProjectRemove()
						},
					},
					{
						Name:  "list",
						Usage: "List all projects",
						Action: func(c *cli.Context) error {
							return ProjectList()
						},
					},
					{
						Name:  "select",
						Usage: "Select a project",
						Action: func(c *cli.Context) error {
							return ProjectSelect()
						},
					},
					{
						Name:  "jira-setting",
						Usage: "Set Jira settings",
						Action: func(c *cli.Context) error {
							return JIRASetting()
						},
					},
					{
						Name:  "dropbox-setting",
						Usage: "Set Dropbox settings",
						Action: func(c *cli.Context) error {
							return DropboxSetting()
						},
					},
					{
						Name:  "sync",
						Usage: "Syncronize tasks with 3rd party services",
						Subcommands: []*cli.Command{
							{
								Name:  "jira",
								Usage: "Syncronize tasks with Jira",
								Action: func(c *cli.Context) error {
									return SyncJIRAIssues()
								},
							},
							{
								Name:  "dropbox",
								Usage: "Syncronize tasks with Dropbox",
								Action: func(c *cli.Context) error {
									return SyncDropbox()
								},
							},
						},
					},
				},
			},
		},
	}
}
