package task

import (
	"context"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type TaskRepository interface {
	GetUncompleteTasks(ctx context.Context, projectID string) (entity.Tasks, error)
	Insert(ctx context.Context, task entity.Task) error
	Update(ctx context.Context, task entity.Task) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (entity.Task, error)
}

type SyncRepository interface {
	Sync(ctx context.Context) error
}

type ProjectRepository interface {
	GetSelectedProject(ctx context.Context) (entity.Project, error)
}
