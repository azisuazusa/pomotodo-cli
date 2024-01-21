package setting

import (
	"context"
	"fmt"
)

type UseCase interface {
	SetIntegration(ctx context.Context, integration Integration) error
	Upload(ctx context.Context) error
	Download(ctx context.Context) error
}

type useCase struct {
	settingRepo     SettingRepository
	integrationRepo map[IntegrationType]IntegrationRepository
}

func New(settingRepo SettingRepository, integrationRepo map[IntegrationType]IntegrationRepository) UseCase {
	return &useCase{
		settingRepo:     settingRepo,
		integrationRepo: integrationRepo,
	}
}

func (u *useCase) SetIntegration(ctx context.Context, integration Integration) error {
	err := u.settingRepo.SetIntegration(ctx, integration)
	if err != nil {
		return fmt.Errorf("error while setting dropbox token: %w", err)
	}

	return nil
}

func (u *useCase) Upload(ctx context.Context) error {
	integration, err := u.settingRepo.GetIntegration(ctx)
	if err != nil {
		return fmt.Errorf("error while getting integration: %w", err)
	}

	if err = u.integrationRepo[integration.Type].Upload(ctx, integration); err != nil {
		return fmt.Errorf("error while uploading: %w", err)
	}

	return nil
}

func (u *useCase) Download(ctx context.Context) error {
	integration, err := u.settingRepo.GetIntegration(ctx)
	if err != nil {
		return fmt.Errorf("error while getting integration: %w", err)
	}

	if err = u.integrationRepo[integration.Type].Download(ctx, integration); err != nil {
		return fmt.Errorf("error while downloading: %w", err)
	}

	return nil
}
