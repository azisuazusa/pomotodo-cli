package cmd

import (
	"database/sql"
	"os"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	projectDomain "github.com/azisuazusa/todo-cli/internal/domain/project"
	settingDomain "github.com/azisuazusa/todo-cli/internal/domain/setting"
	taskDomain "github.com/azisuazusa/todo-cli/internal/domain/task"
	projectPresenter "github.com/azisuazusa/todo-cli/internal/presenter/project"
	settingPresenter "github.com/azisuazusa/todo-cli/internal/presenter/setting"
	taskPresenter "github.com/azisuazusa/todo-cli/internal/presenter/task"
	"github.com/azisuazusa/todo-cli/internal/repository/dropbox"
	"github.com/azisuazusa/todo-cli/internal/repository/jira"
	projectRepository "github.com/azisuazusa/todo-cli/internal/repository/project"
	settingRepository "github.com/azisuazusa/todo-cli/internal/repository/setting"
	taskRepository "github.com/azisuazusa/todo-cli/internal/repository/task"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func TodoCLI() *cli.App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", homeDir+"/.todo-cli.db")
	if err != nil {
		panic(err)
	}

	taskRepo := taskRepository.New(db)
	projectRepo := projectRepository.New(db)
	taskUseCase := taskDomain.New(taskRepo, projectRepo)

	settingIntegrationRepo := map[settingDomain.IntegrationType]settingDomain.IntegrationRepository{
		settingDomain.Dropbox: dropbox.New(),
	}
	settingUseCase := settingDomain.New(settingRepository.New(db), settingIntegrationRepo)

	projectIntegrationRepo := map[entity.IntegrationType]projectDomain.IntegrationRepository{
		entity.IntegrationTypeJIRA: jira.New(),
	}
	projectUseCase := projectDomain.New(projectRepo, projectIntegrationRepo, taskRepo)

	taskPresenter := taskPresenter.New(taskUseCase, settingUseCase)
	settingPresenter := settingPresenter.New(settingUseCase)
	projectPresenter := projectPresenter.New(projectUseCase, settingUseCase)
	commands := taskCLI(taskPresenter)
	return &cli.App{
		Name:     "todo",
		Usage:    "todo-cli is a CLI for managing your todo list",
		Commands: append(commands, projectCLI(projectPresenter), settingCLI(settingPresenter)),
	}

}
