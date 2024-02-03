package syncintegration

import (
	"context"
	"fmt"
)

type UseCase interface {
	SetSyncIntegration(ctx context.Context, integration SyncIntegration) error
	Upload(ctx context.Context) error
	Download(ctx context.Context) error
}

type useCase struct {
	settingRepo      SettingRepository
	integrationRepos map[SyncIntegrationType]IntegrationRepository
}

func New(settingRepo SettingRepository, integrationRepo map[SyncIntegrationType]IntegrationRepository) UseCase {
	return &useCase{
		settingRepo:      settingRepo,
		integrationRepos: integrationRepo,
	}
}

func (u *useCase) SetSyncIntegration(ctx context.Context, integration SyncIntegration) error {
	err := u.settingRepo.SetSyncIntegration(ctx, integration)
	if err != nil {
		return fmt.Errorf("error while setting dropbox token: %w", err)
	}

	return nil
}

func (u *useCase) Upload(ctx context.Context) error {
	integration, err := u.settingRepo.GetSyncIntegration(ctx)
	if err != nil && err != ErrSyncIntegrationNotFound {
		return fmt.Errorf("error while getting integration: %w", err)
	}

	if err == ErrSyncIntegrationNotFound {
		return nil
	}

	if err = u.integrationRepos[integration.Type].Upload(ctx, integration); err != nil {
		return fmt.Errorf("error while uploading: %w", err)
	}

	return nil
}

func (u *useCase) Download(ctx context.Context) error {
	integration, err := u.settingRepo.GetSyncIntegration(ctx)
	if err != nil && err != ErrSyncIntegrationNotFound {
		return fmt.Errorf("error while getting integration: %w", err)
	}

	if err == ErrSyncIntegrationNotFound {
		return nil
	}

	if err = u.integrationRepos[integration.Type].Download(ctx, integration); err != nil {
		return fmt.Errorf("error while downloading: %w", err)
	}

	return nil
}
