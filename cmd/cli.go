package cmd

import (
	"database/sql"
	"os"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	jiraDomain "github.com/azisuazusa/todo-cli/internal/domain/jira"
	projectDomain "github.com/azisuazusa/todo-cli/internal/domain/project"
	syncintegrationDomain "github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
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

	// Repositories
	taskRepo := taskRepository.New(db)
	projectRepo := projectRepository.New(db)
	jiraRepo := jira.New()
	settingIntegrationRepo := map[syncintegrationDomain.SyncIntegrationType]syncintegrationDomain.IntegrationRepository{
		syncintegrationDomain.Dropbox: dropbox.New(),
	}
	projectIntegrationRepo := map[entity.IntegrationType]projectDomain.IntegrationRepository{
		entity.IntegrationTypeJIRA: jiraRepo,
	}

	// UseCases
	taskUseCase := taskDomain.New(taskRepo, projectRepo)
	settingUseCase := syncintegrationDomain.New(settingRepository.New(db), settingIntegrationRepo)
	projectUseCase := projectDomain.New(projectRepo, projectIntegrationRepo, taskRepo)
	jiraUseCase := jiraDomain.New(jiraRepo, projectRepo, taskRepo)

	// Presenters
	taskPresenter := taskPresenter.New(taskUseCase, settingUseCase, jiraUseCase)
	settingPresenter := settingPresenter.New(settingUseCase)
	projectPresenter := projectPresenter.New(projectUseCase, settingUseCase, taskUseCase)

	commands := taskCLI(taskPresenter)
	commands = append(commands, projectCLI(projectPresenter), settingCLI(settingPresenter), setupCLI(db))
	return &cli.App{
		Name:     "todo",
		Usage:    "todo-cli is a CLI for managing your todo list",
		Commands: commands,
	}

}
