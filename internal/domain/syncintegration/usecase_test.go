package syncintegration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration/mocks"
	"github.com/stretchr/testify/suite"
)

type UseCaseTestSuite struct {
	suite.Suite
	settingRepo      *mocks.SettingRepository
	integrationRepos map[syncintegration.SyncIntegrationType]*mocks.IntegrationRepository
	useCase          syncintegration.UseCase
}

func (t *UseCaseTestSuite) SetupTest() {
	t.settingRepo = &mocks.SettingRepository{}
	dropboxIntegrationRepo := &mocks.IntegrationRepository{}
	integrationRepos := map[syncintegration.SyncIntegrationType]syncintegration.IntegrationRepository{
		syncintegration.Dropbox: dropboxIntegrationRepo,
	}
	t.integrationRepos = map[syncintegration.SyncIntegrationType]*mocks.IntegrationRepository{
		syncintegration.Dropbox: dropboxIntegrationRepo,
	}
	t.useCase = syncintegration.New(t.settingRepo, integrationRepos)
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}

func (t *UseCaseTestSuite) SetSyncIntegration() {
	tests := []struct {
		name          string
		integration   syncintegration.SyncIntegration
		expectedError error
		mockFunc      func(syncintegration.SyncIntegration)
	}{
		{
			name: "failed to set sync integration",
			integration: syncintegration.SyncIntegration{
				Type: syncintegration.Dropbox,
			},
			expectedError: errors.New("any-error"),
			mockFunc: func(integration syncintegration.SyncIntegration) {
				t.settingRepo.On("SetSyncIntegration", context.Background(), integration).Return(errors.New("any-error")).Once()
			},
		},
		{
			name: "success",
			integration: syncintegration.SyncIntegration{
				Type: syncintegration.Dropbox,
			},
			expectedError: nil,
			mockFunc: func(integration syncintegration.SyncIntegration) {
				t.settingRepo.On("SetSyncIntegration", context.Background(), integration).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.integration)
			err := t.useCase.SetSyncIntegration(context.Background(), test.integration)
			t.Equal(test.expectedError, err)
		})
	}
}

func (t *UseCaseTestSuite) TestUpload() {
	tests := []struct {
		name          string
		expectedError error
		mockFunc      func()
	}{
		{
			name:          "failed to get sync integration",
			expectedError: errors.New("any-error"),
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{}, errors.New("any-error")).Once()
			},
		},
		{
			name:          "no integration found",
			expectedError: nil,
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{}, syncintegration.ErrSyncIntegrationNotFound).Once()
			},
		},
		{
			name:          "failed to upload",
			expectedError: errors.New("any-error"),
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}, nil).Once()
				t.integrationRepos[syncintegration.Dropbox].On("Upload", context.Background(), syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:          "success",
			expectedError: nil,
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}, nil).Once()
				t.integrationRepos[syncintegration.Dropbox].On("Upload", context.Background(), syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc()
			err := t.useCase.Upload(context.Background())
			if err != nil {
				err = errors.Unwrap(err)
			}

			t.Equal(test.expectedError, err)
		})
	}
}

func (t *UseCaseTestSuite) TestDownload() {
	tests := []struct {
		name          string
		expectedError error
		mockFunc      func()
	}{
		{
			name:          "failed to get sync integration",
			expectedError: errors.New("any-error"),
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{}, errors.New("any-error")).Once()
			},
		},
		{
			name:          "no integration found",
			expectedError: nil,
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{}, syncintegration.ErrSyncIntegrationNotFound).Once()
			},
		},
		{
			name:          "failed to download",
			expectedError: errors.New("any-error"),
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}, nil).Once()
				t.integrationRepos[syncintegration.Dropbox].On("Download", context.Background(), syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:          "success",
			expectedError: nil,
			mockFunc: func() {
				t.settingRepo.On("GetSyncIntegration", context.Background()).Return(syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}, nil).Once()
				t.integrationRepos[syncintegration.Dropbox].On("Download", context.Background(), syncintegration.SyncIntegration{
					Type: syncintegration.Dropbox,
				}).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc()
			err := t.useCase.Download(context.Background())
			if err != nil {
				err = errors.Unwrap(err)
			}

			t.Equal(test.expectedError, err)
		})
	}
}
