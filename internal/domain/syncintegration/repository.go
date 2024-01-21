package syncintegration

import (
	"context"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type SettingRepository interface {
	SetSyncIntegration(ctx context.Context, integration SyncIntegration) error
	GetSyncIntegration(ctx context.Context) (SyncIntegration, error)
}

type TaskRepository interface {
	GetTasks(ctx context.Context) (entity.Tasks, error)
}

type IntegrationRepository interface {
	Upload(ctx context.Context, integration SyncIntegration) error
	Download(ctx context.Context, integration SyncIntegration) error
}
