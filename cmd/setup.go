package cmd

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func setupCLI(db *sql.DB) *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "Setup todo-cli",
		Action: func(c *cli.Context) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("failed to get home dir: %v\n", err)
				return err
			}

			dbPath := fmt.Sprintf("%s/.todo-cli.db", homeDir)
			if _, err = os.Stat(dbPath); err == nil {
				fmt.Printf("db already exists: %s\n", dbPath)
				return nil
			}

			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				fmt.Printf("failed to open db: %v\n", err)
				return err
			}

			_, err = db.Exec(`
				CREATE TABLE IF NOT EXISTS projects (
					id VARCHAR PRIMARY KEY,
					name TEXT NOT NULL,
					description TEXT,
					is_selected BOOLEAN NOT NULL DEFAULT FALSE,
					integrations TEXT
				);

				CREATE TABLE IF NOT EXISTS tasks (
					id VARCHAR PRIMARY KEY,
					project_id VARCHAR NOT NULL,
					name TEXT NOT NULL,
					description TEXT,
					is_started BOOLEAN NOT NULL DEFAULT FALSE,
					completed_at DATETIME,
					parent_task_id VARCHAR,
					integration TEXT,
					histories TEXT,
					FOREIGN KEY (project_id) REFERENCES projects(id),
					FOREIGN KEY (parent_task_id) REFERENCES tasks(id)
				);

				CREATE TABLE IF NOT EXISTS settings (
					key VARCHAR PRIMARY KEY,
					value TEXT NOT NULL
				);`)
			if err != nil {
				fmt.Printf("failed to create tables: %v\n", err)
				return err
			}

			return nil
		},
	}
}
