package setting

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	_ "github.com/mattn/go-sqlite3"
)

type RepoImpl struct {
	db *sql.DB
}

func New(db *sql.DB) *RepoImpl {
	return &RepoImpl{
		db: db,
	}
}

func (r *RepoImpl) SetSyncIntegration(ctx context.Context, integration syncintegration.SyncIntegration) error {
	model, err := CreateModelFromSyncIntegration(integration)
	if err != nil {
		return fmt.Errorf("failed to create model from integration: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO settings (key, value) VALUES (?, ?)", model.Key, model.Value)
	if err != nil {
		return fmt.Errorf("failed to insert setting: %w", err)
	}

	return nil
}

func (r *RepoImpl) GetSyncIntegration(ctx context.Context) (syncintegration.SyncIntegration, error) {
	var model SettingModel
	err := r.db.QueryRowContext(ctx, "SELECT key, value FROM settings").Scan(&model.Key, &model.Value)
	if err != nil && err != sql.ErrNoRows {
		return syncintegration.SyncIntegration{}, fmt.Errorf("failed to get setting: %w", err)
	}

	if err == sql.ErrNoRows {
		return syncintegration.SyncIntegration{}, syncintegration.ErrSyncIntegrationNotFound
	}

	return model.ToSyncIntegration()
}
