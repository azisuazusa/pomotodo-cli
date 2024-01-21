package task

import (
	"context"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type TaskRepository interface {
	GetUncompleteParentTasks(ctx context.Context, projectID string) (entity.Tasks, error)
	GetUncompleteSubTask(ctx context.Context, projectID string) (map[string]entity.Tasks, error)
	Insert(ctx context.Context, task entity.Task) error
	Update(ctx context.Context, task entity.Task) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (entity.Task, error)
	GetStartedTask(ctx context.Context) (entity.Task, error)
}

type SettingRepository interface {
	Sync(ctx context.Context) error
}

type ProjectRepository interface {
	GetSelectedProject(ctx context.Context) (entity.Project, error)
}
