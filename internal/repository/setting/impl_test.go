package setting

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/stretchr/testify/suite"
)

type RepoImplTestSuite struct {
	suite.Suite
	db       sqlmock.Sqlmock
	repoImpl RepoImpl
}

func (s *RepoImplTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	s.db = mock

	s.repoImpl = RepoImpl{db: db}
}

func TestRepoImpl(t *testing.T) {
	suite.Run(t, new(RepoImplTestSuite))
}

func (s *RepoImplTestSuite) TestSetSyncIntegration() {
	tests := []struct {
		name             string
		paramIntegration syncintegration.SyncIntegration
		expectedErr      error
		mock             func(paramIntegration syncintegration.SyncIntegration)
	}{
		{
			name: "failed to set sync integration",
			paramIntegration: syncintegration.SyncIntegration{
				Type: syncintegration.Dropbox,
				Details: map[string]string{
					"token": "token-test",
				},
			},
			expectedErr: errors.New("any-error"),
			mock: func(paramIntegration syncintegration.SyncIntegration) {
				model, _ := CreateModelFromSyncIntegration(paramIntegration)
				query := "INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?"
				s.db.ExpectExec(query).
					WithArgs(model.Key, model.Value, model.Value).
					WillReturnError(errors.New("any-error"))
			},
		},
		{
			name: "success to set sync integration",
			paramIntegration: syncintegration.SyncIntegration{
				Type: syncintegration.Dropbox,
				Details: map[string]string{
					"token": "token-test",
				},
			},
			expectedErr: nil,
			mock: func(paramIntegration syncintegration.SyncIntegration) {
				model, _ := CreateModelFromSyncIntegration(paramIntegration)
				query := "INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?"
				s.db.ExpectExec(query).
					WithArgs(model.Key, model.Value, model.Value).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.paramIntegration)
			err := s.repoImpl.SetSyncIntegration(context.Background(), tt.paramIntegration)

			if err != nil {
				err = errors.Unwrap(err)
			}

			s.Equal(tt.expectedErr, err)
		})
	}
}

func (s *RepoImplTestSuite) TestGetSyncIntegration() {
	tests := []struct {
		name        string
		expected    syncintegration.SyncIntegration
		expectedErr error
		mock        func()
	}{
		{
			name:        "failed to get sync integration",
			expected:    syncintegration.SyncIntegration{},
			expectedErr: errors.New("any-error"),
			mock: func() {
				query := "SELECT key, value FROM settings WHERE key = ?"
				s.db.ExpectQuery(query).
					WithArgs("sync_integration").
					WillReturnError(errors.New("any-error"))
			},
		},
		{
			name: "success to get sync integration",
			expected: syncintegration.SyncIntegration{
				Type: syncintegration.Dropbox,
				Details: map[string]string{
					"token": "token-test",
				},
			},
			expectedErr: nil,
			mock: func() {
				query := "SELECT key, value FROM settings WHERE key = ?"
				s.db.ExpectQuery(query).
					WithArgs("sync_integration").
					WillReturnRows(sqlmock.NewRows([]string{"key", "value"}).AddRow("sync_integration", "{\"type\":\"dropbox\",\"details\":{\"token\":\"token-test\"}}"))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			actual, err := s.repoImpl.GetSyncIntegration(context.Background())

			if err != nil {
				err = errors.Unwrap(err)
			}

			s.Equal(tt.expectedErr, err)
			s.Equal(tt.expected, actual)
		})
	}
}
