package jira

import (
	"context"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type JiraRepository interface {
	AddWorklog(ctx context.Context, issueID, taskName string, timeSpent time.Duration, integrationEntity entity.Integration) error
}

type ProjectRepository interface {
	GetSelectedProject(ctx context.Context) (entity.Project, error)
}

type TaskRepository interface {
	GetByID(ctx context.Context, id string) (entity.Task, error)
}
