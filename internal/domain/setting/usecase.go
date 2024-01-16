package setting

import (
	"context"
	"fmt"
)

type UseCase interface {
	SetDropboxToken(ctx context.Context, token string) error
}

type useCase struct {
	settingRepo SettingRepository
}

func New(settingRepo SettingRepository) UseCase {
	return &useCase{
		settingRepo: settingRepo,
	}
}

func (u *useCase) SetDropboxToken(ctx context.Context, token string) error {
	err := u.settingRepo.SetDropboxToken(ctx, token)
	if err != nil {
		return fmt.Errorf("error while setting dropbox token: %w", err)
	}

	return nil
}
