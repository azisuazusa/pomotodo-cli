package setting

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/setting"
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

func (r *RepoImpl) SetIntegration(ctx context.Context, integration setting.Integration) error {
	model, err := CreateModelFromIntegration(integration)
	if err != nil {
		return fmt.Errorf("failed to create model from integration: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO setting (id, type, details) VALUES (?, ?, ?)", model.ID, model.SettingType, model.SettingDetails)
	if err != nil {
		return fmt.Errorf("failed to insert setting: %w", err)
	}

	return nil
}

func (r *RepoImpl) GetIntegration(ctx context.Context) (setting.Integration, error) {
	var model SettingModel
	err := r.db.QueryRowContext(ctx, "SELECT id, type, details FROM setting").Scan(&model.ID, &model.SettingType, &model.SettingDetails)
	if err != nil {
		return setting.Integration{}, fmt.Errorf("failed to get setting: %w", err)
	}

	return model.ToIntegrationEntity()
}
