package project

import (
	"context"
	"errors"
	"testing"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/azisuazusa/todo-cli/internal/domain/project/mocks"
	"github.com/stretchr/testify/suite"
)

type UseCaseTestSuite struct {
	suite.Suite
	projectRepo     *mocks.ProjectRepository
	taskRepo        *mocks.TaskRepository
	integrationRepo map[entity.IntegrationType]IntegrationRepository
	useCase         UseCase
}

func (t *UseCaseTestSuite) SetupTest() {
	t.projectRepo = &mocks.ProjectRepository{}
	t.taskRepo = &mocks.TaskRepository{}
	t.integrationRepo = map[entity.IntegrationType]IntegrationRepository{
		entity.IntegrationTypeJIRA:   &mocks.IntegrationRepository{},
		entity.IntegrationTypeGitHub: &mocks.IntegrationRepository{},
	}
	t.useCase = New(t.projectRepo, t.integrationRepo, t.taskRepo)
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}

func (t *UseCaseTestSuite) TestGetAll() {
	tests := []struct {
		name        string
		expectedRes entity.Projects
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get projects",
			expectedRes: nil,
			expectedErr: errors.New("any-error"),
			mockFunc: func() {
				t.projectRepo.On("GetAll", context.Background()).Return(nil, errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			expectedRes: entity.Projects{},
			expectedErr: nil,
			mockFunc: func() {
				t.projectRepo.On("GetAll", context.Background()).Return(entity.Projects{}, nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc()
			actual, err := t.useCase.GetAll(context.Background())
			t.Equal(test.expectedRes, actual)
			t.Equal(test.expectedErr, err)
		})
	}

}

func (t *UseCaseTestSuite) TestInsert() {
	tests := []struct {
		name        string
		project     entity.Project
		expectedErr error
		mockFunc    func(project entity.Project)
	}{
		{
			name: "failed to insert project",
			project: entity.Project{
				Name: "any-project",
			},
			expectedErr: errors.New("any-error"),
			mockFunc: func(project entity.Project) {
				t.projectRepo.On("Insert", context.Background(), project).Return(errors.New("any-error")).Once()
			},
		},
		{
			name: "success",
			project: entity.Project{
				Name: "any-project",
			},
			expectedErr: nil,
			mockFunc: func(project entity.Project) {
				t.projectRepo.On("Insert", context.Background(), project).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.project)
			err := t.useCase.Insert(context.Background(), test.project)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestDelete() {
	tests := []struct {
		name        string
		projectID   string
		expectedErr error
		mockFunc    func(projectID string)
	}{
		{
			name:        "failed to delete project",
			projectID:   "any-project-id",
			expectedErr: errors.New("any-error"),
			mockFunc: func(projectID string) {
				t.projectRepo.On("Delete", context.Background(), projectID).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			projectID:   "any-project-id",
			expectedErr: nil,
			mockFunc: func(projectID string) {
				t.projectRepo.On("Delete", context.Background(), projectID).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.projectID)
			err := t.useCase.Delete(context.Background(), test.projectID)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestSelect() {
	tests := []struct {
		name        string
		projectID   string
		expectedErr error
		mockFunc    func(projectID string)
	}{
		{
			name:        "failed to get project by id",
			projectID:   "any-project-id",
			expectedErr: errors.New("any-error"),
			mockFunc: func(projectID string) {
				t.projectRepo.On("GetByID", context.Background(), projectID).Return(entity.Project{}, errors.New("any-error")).Once()
			},
		},
		{
			name:        "failed to update project",
			projectID:   "any-project-id",
			expectedErr: errors.New("any-error"),
			mockFunc: func(projectID string) {
				project := entity.Project{
					ID: projectID,
				}
				t.projectRepo.On("GetByID", context.Background(), projectID).Return(project, nil).Once()
				project.IsSelected = true
				t.projectRepo.On("Update", context.Background(), project).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			projectID:   "any-project-id",
			expectedErr: nil,
			mockFunc: func(projectID string) {
				project := entity.Project{
					ID: projectID,
				}
				t.projectRepo.On("GetByID", context.Background(), projectID).Return(project, nil).Once()
				project.IsSelected = true
				t.projectRepo.On("Update", context.Background(), project).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.projectID)
			err := t.useCase.Select(context.Background(), test.projectID)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}
