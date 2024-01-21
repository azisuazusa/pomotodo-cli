package setting

import (
	"context"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type SettingRepository interface {
	SetIntegration(ctx context.Context, integration Integration) error
	GetIntegration(ctx context.Context) (Integration, error)
}

type TaskRepository interface {
	GetTasks(ctx context.Context) (entity.Tasks, error)
}

type IntegrationRepository interface {
	Upload(ctx context.Context, integration Integration) error
	Download(ctx context.Context, integration Integration) error
}
