package jira

import (
	"context"
	"fmt"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type UseCase interface {
	AddWorklog(ctx context.Context, task entity.Task, timeSpent time.Duration) error
}

type useCase struct {
	jiraRepo    JiraRepository
	projectRepo ProjectRepository
	taskRepo    TaskRepository
}

func New(jiraRepo JiraRepository, projectRepo ProjectRepository, taskRepo TaskRepository) UseCase {
	return &useCase{
		jiraRepo:    jiraRepo,
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

func (u *useCase) AddWorklog(ctx context.Context, task entity.Task, timeSpent time.Duration) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return err
	}

	if len(project.Integrations) == 0 {
		return fmt.Errorf("no integration found")
	}

	integrationEntity := entity.Integration{}
	for _, integration := range project.Integrations {
		if integration.Type == entity.IntegrationTypeJIRA {
			integrationEntity = integration
			break
		}

		return fmt.Errorf("no jira integration found")
	}

	if task.ParentTaskID != "" {
		task, err = u.taskRepo.GetByID(ctx, task.ParentTaskID)
		if err != nil {
			return fmt.Errorf("error while getting parent task: %w", err)
		}
	}

	err = u.jiraRepo.AddWorklog(ctx, task.ID, task.Name, timeSpent, integrationEntity)
	if err != nil {
		return fmt.Errorf("error while adding worklog: %w", err)
	}

	return nil

}
