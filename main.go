package main

import (
	"log"
	"os"

	"github.com/azisuazusa/todo-cli/cmd"
)

func main() {
	app := cmd.TodoCLI()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
