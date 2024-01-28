package project

import (
	"context"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type ProjectRepository interface {
	GetAll(ctx context.Context) (entity.Projects, error)
	Insert(ctx context.Context, project entity.Project) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (entity.Project, error)
	Update(ctx context.Context, project entity.Project) error
	SetSelectedProject(ctx context.Context, id string) error
	GetSelectedProject(ctx context.Context) (entity.Project, error)
}

type TaskRepository interface {
	Upsert(ctx context.Context, task entity.Task) error
	GetStartedTask(ctx context.Context) (entity.Task, error)
	Update(ctx context.Context, task entity.Task) error
}

type IntegrationRepository interface {
	GetTasks(ctx context.Context, projectID string, details map[string]string) (entity.Tasks, error)
}
