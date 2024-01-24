package cmd

import (
	"github.com/azisuazusa/todo-cli/internal/presenter/task"
	"github.com/urfave/cli/v2"
)

func taskCLI(presenter *task.Presenter) []*cli.Command {

	return []*cli.Command{
		{
			Name:  "add",
			Usage: "Add a task",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "subtask",
					Usage: "Add a subtask",
				},
			},
			Action: func(c *cli.Context) error {
				if !c.Bool("subtask") {
					return presenter.Add(c.Context)
				}

				return presenter.AddSubTask(c.Context)
			},
		},
		{
			Name:  "start",
			Usage: "Start a task",
			Action: func(c *cli.Context) error {
				return presenter.Start(c.Context)
			},
		},
		{
			Name:  "stop",
			Usage: "Stop a task",
			Action: func(c *cli.Context) error {
				return presenter.Stop(c.Context)
			},
		},
		{
			Name:  "remove",
			Usage: "Remove a task",
			Action: func(c *cli.Context) error {
				return presenter.Remove(c.Context)
			},
		},
		{
			Name:  "complete",
			Usage: "Complete a task",
			Action: func(c *cli.Context) error {
				return presenter.Complete(c.Context)
			},
		},
		{
			Name:  "list",
			Usage: "List tasks",
			Action: func(c *cli.Context) error {
				return presenter.GetUncompleteTasks(c.Context)
			},
		},
	}

}
