package project

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/azisuazusa/todo-cli/internal/domain/project/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UseCaseTestSuite struct {
	suite.Suite
	projectRepo      *mocks.ProjectRepository
	taskRepo         *mocks.TaskRepository
	integrationRepos map[entity.IntegrationType]*mocks.IntegrationRepository
	useCase          UseCase
}

func (t *UseCaseTestSuite) SetupTest() {
	t.projectRepo = &mocks.ProjectRepository{}
	t.taskRepo = &mocks.TaskRepository{}
	jiraIntegrationRepo := &mocks.IntegrationRepository{}
	githubIntegrationRepo := &mocks.IntegrationRepository{}
	integrationRepos := map[entity.IntegrationType]IntegrationRepository{
		entity.IntegrationTypeJIRA:   jiraIntegrationRepo,
		entity.IntegrationTypeGitHub: githubIntegrationRepo,
	}
	t.integrationRepos = map[entity.IntegrationType]*mocks.IntegrationRepository{
		entity.IntegrationTypeJIRA:   jiraIntegrationRepo,
		entity.IntegrationTypeGitHub: githubIntegrationRepo,
	}
	t.useCase = New(t.projectRepo, integrationRepos, t.taskRepo)
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
			name:        "failed to select project",
			projectID:   "any-project-id",
			expectedErr: errors.New("any-error"),
			mockFunc: func(projectID string) {
				t.projectRepo.On("SetSelectedProject", context.Background(), projectID).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			projectID:   "any-project-id",
			expectedErr: nil,
			mockFunc: func(projectID string) {
				t.projectRepo.On("SetSelectedProject", context.Background(), projectID).Return(nil).Once()
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

func (t *UseCaseTestSuite) TestSyncTasks() {
	tests := []struct {
		name        string
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get selected project",
			expectedErr: errors.New("any-error"),
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{}, errors.New("any-error")).Once()
			},
		},
		{
			name:        "no integration found",
			expectedErr: nil,
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{}, nil).Once()
			},
		},
		{
			name:        "failed to get tasks",
			expectedErr: errors.New("any-error"),
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}, nil).Once()
				t.integrationRepos[entity.IntegrationTypeJIRA].On("GetTasks", context.Background(), "any-id", mock.Anything).Return(entity.Tasks{}, errors.New("any-error")).Once()
			},
		},
		{
			name:        "failed to upsert tasks",
			expectedErr: errors.New("any-error"),
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}, nil).Once()
				t.integrationRepos[entity.IntegrationTypeJIRA].On("GetTasks", context.Background(), "any-id", mock.Anything).Return(entity.Tasks{
					{
						ID:           "any-id",
						ProjectID:    "any-id",
						Name:         "any-name",
						Description:  "",
						IsStarted:    false,
						CompletedAt:  time.Time{},
						ParentTaskID: "",
						Integration: entity.TaskIntegration{
							ID:   "",
							Type: "",
						},
						Histories: []entity.TaskHistory{},
					},
				}, nil).Once()
				t.taskRepo.On("Upsert", context.Background(), mock.Anything).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			expectedErr: nil,
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}, nil).Once()
				t.integrationRepos[entity.IntegrationTypeJIRA].On("GetTasks", context.Background(), "any-id", mock.Anything).Return(entity.Tasks{
					{
						ID:           "any-id",
						ProjectID:    "any-id",
						Name:         "any-name",
						Description:  "",
						IsStarted:    false,
						CompletedAt:  time.Time{},
						ParentTaskID: "",
						Integration: entity.TaskIntegration{
							ID:   "",
							Type: "",
						},
						Histories: []entity.TaskHistory{},
					},
				}, nil).Once()
				t.taskRepo.On("Upsert", context.Background(), mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc()
			err := t.useCase.SyncTasks(context.Background())
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestAddIntegration() {
	tests := []struct {
		name        string
		integration entity.Integration
		expectedErr error
		mockFunc    func(integration entity.Integration)
	}{
		{
			name:        "failed to get selected project",
			integration: entity.Integration{},
			expectedErr: errors.New("any-error"),
			mockFunc: func(_ entity.Integration) {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{}, errors.New("any-error")).Once()
			},
		},
		{
			name: "failed to add integration",
			integration: entity.Integration{
				IsEnabled: true,
				Type:      entity.IntegrationTypeJIRA,
				Details: map[string]string{
					"": "",
				},
			},
			expectedErr: errors.New("any-error"),
			mockFunc: func(integration entity.Integration) {
				project := entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				t.projectRepo.On("GetSelectedProject", context.Background()).Return(project, nil).Once()
				project.Integrations = append(project.Integrations, integration)
				t.projectRepo.On("Update", context.Background(), project).Return(errors.New("any-error")).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.integration)
			err := t.useCase.AddIntegration(context.Background(), test.integration)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestRemoveIntegration() {
	tests := []struct {
		name        string
		integration entity.IntegrationType
		expectedErr error
		mockFunc    func(integration entity.IntegrationType)
	}{
		{
			name:        "failed to get selected project",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: errors.New("any-error"),
			mockFunc: func(_ entity.IntegrationType) {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{}, errors.New("any-error")).Once()
			},
		},
		{
			name:        "failed to remove integration",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: errors.New("any-error"),
			mockFunc: func(_ entity.IntegrationType) {
				project := entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				t.projectRepo.On("GetSelectedProject", context.Background()).Return(project, nil).Once()
				project.Integrations = []entity.Integration{}
				t.projectRepo.On("Update", context.Background(), project).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: nil,
			mockFunc: func(_ entity.IntegrationType) {
				project := entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeGitHub,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				t.projectRepo.On("GetSelectedProject", context.Background()).Return(project, nil).Once()
				project.Integrations = []entity.Integration{
					{
						IsEnabled: true,
						Type:      entity.IntegrationTypeGitHub,
						Details: map[string]string{
							"token": "token",
						},
					},
				}
				t.projectRepo.On("Update", context.Background(), project).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.integration)
			err := t.useCase.RemoveIntegration(context.Background(), test.integration)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestGetSelected() {
	tests := []struct {
		name        string
		expected    entity.Project
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get selectd project",
			expected:    entity.Project{},
			expectedErr: errors.New("any-error"),
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{}, errors.New("any-error")).Once()
			},
		},
		{
			name: "success",
			expected: entity.Project{
				ID:          "any-id",
				Name:        "any-name",
				Description: "",
				IsSelected:  true,
				Integrations: []entity.Integration{
					{
						IsEnabled: true,
						Type:      entity.IntegrationTypeJIRA,
						Details: map[string]string{
							"token": "token",
						},
					},
				},
			},
			expectedErr: nil,
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}, nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc()
			project, err := t.useCase.GetSelected(context.Background())
			t.Equal(test.expected, project)
			t.Equal(test.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestEnableIntegration() {
	tests := []struct {
		name        string
		integration entity.IntegrationType
		expectedErr error
		mockFunc    func(integration entity.IntegrationType)
	}{
		{
			name:        "failed to get selected project",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: errors.New("any-error"),
			mockFunc: func(_ entity.IntegrationType) {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{}, errors.New("any-error")).Once()
			},
		},
		{
			name:        "failed to update project",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: errors.New("any-error"),
			mockFunc: func(integrationType entity.IntegrationType) {
				project := entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: false,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				t.projectRepo.On("GetSelectedProject", context.Background()).Return(project, nil).Once()
				project.Integrations = []entity.Integration{
					{
						IsEnabled: true,
						Type:      entity.IntegrationTypeJIRA,
						Details: map[string]string{
							"token": "token",
						},
					},
				}
				t.projectRepo.On("Update", context.Background(), project).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: nil,
			mockFunc: func(integrationType entity.IntegrationType) {
				project := entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: false,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				t.projectRepo.On("GetSelectedProject", context.Background()).Return(project, nil).Once()
				project.Integrations = []entity.Integration{
					{
						IsEnabled: true,
						Type:      entity.IntegrationTypeJIRA,
						Details: map[string]string{
							"token": "token",
						},
					},
				}
				t.projectRepo.On("Update", context.Background(), project).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.integration)
			err := t.useCase.EnableIntegration(context.Background(), test.integration)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestDisableIntegration() {
	tests := []struct {
		name        string
		integration entity.IntegrationType
		expectedErr error
		mockFunc    func(integration entity.IntegrationType)
	}{
		{
			name:        "failed to get selected project",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: errors.New("any-error"),
			mockFunc: func(_ entity.IntegrationType) {
				t.projectRepo.On("GetSelectedProject", context.Background()).Return(entity.Project{}, errors.New("any-error")).Once()
			},
		},
		{
			name:        "failed to update project",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: errors.New("any-error"),
			mockFunc: func(integrationType entity.IntegrationType) {
				project := entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				t.projectRepo.On("GetSelectedProject", context.Background()).Return(project, nil).Once()
				project.Integrations = []entity.Integration{
					{
						IsEnabled: false,
						Type:      entity.IntegrationTypeJIRA,
						Details: map[string]string{
							"token": "token",
						},
					},
				}
				t.projectRepo.On("Update", context.Background(), project).Return(errors.New("any-error")).Once()
			},
		},
		{
			name:        "success",
			integration: entity.IntegrationTypeJIRA,
			expectedErr: nil,
			mockFunc: func(integrationType entity.IntegrationType) {
				project := entity.Project{
					ID:          "any-id",
					Name:        "any-name",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeJIRA,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				t.projectRepo.On("GetSelectedProject", context.Background()).Return(project, nil).Once()
				project.Integrations = []entity.Integration{
					{
						IsEnabled: false,
						Type:      entity.IntegrationTypeJIRA,
						Details: map[string]string{
							"token": "token",
						},
					},
				}
				t.projectRepo.On("Update", context.Background(), project).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc(test.integration)
			err := t.useCase.DisableIntegration(context.Background(), test.integration)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(test.expectedErr, err)
		})
	}
}
