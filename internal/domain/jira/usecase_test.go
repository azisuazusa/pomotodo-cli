package jira

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/azisuazusa/todo-cli/internal/domain/jira/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UseCaseTestSuite struct {
	suite.Suite
	jiraRepo    *mocks.JiraRepository
	projectRepo *mocks.ProjectRepository
	taskRepo    *mocks.TaskRepository
	useCase     UseCase
}

func (t *UseCaseTestSuite) SetupTest() {
	t.jiraRepo = new(mocks.JiraRepository)
	t.projectRepo = new(mocks.ProjectRepository)
	t.taskRepo = new(mocks.TaskRepository)
	t.useCase = New(t.jiraRepo, t.projectRepo, t.taskRepo)
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}

func (t *UseCaseTestSuite) TestAddWorklog() {
	tests := []struct {
		name          string
		task          entity.Task
		timeSpent     time.Duration
		expectedError error
		mockFunc      func(task entity.Task, timeSpent time.Duration)
	}{
		{
			name: "failed to get selected project",
			task: entity.Task{
				ProjectID: "project-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: entity.ErrNoProjectSelected,
			mockFunc: func(_ entity.Task, _ time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{}, entity.ErrNoProjectSelected).Once()
			},
		},
		{
			name: "no integration found",
			task: entity.Task{
				ProjectID: "project-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: fmt.Errorf("no integration found"),
			mockFunc: func(_ entity.Task, _ time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID:           "project-1",
					Name:         "Project 1",
					Description:  "",
					IsSelected:   true,
					Integrations: []entity.Integration{},
				}, nil).Once()
			},
		},
		{
			name: "no jira integration found",
			task: entity.Task{
				ProjectID: "project-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: fmt.Errorf("no jira integration found"),
			mockFunc: func(_ entity.Task, _ time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
					Description: "",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      entity.IntegrationTypeGitHub,
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}, nil).Once()
			},
		},
		{
			name: "task is child and failed to get parent task",
			task: entity.Task{
				ProjectID:    "project-1",
				ParentTaskID: "parent-task-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: fmt.Errorf("error while getting parent task: %w", errors.New("any-error")),
			mockFunc: func(task entity.Task, _ time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
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
				t.taskRepo.On("GetByID", mock.Anything, task.ParentTaskID).Return(entity.Task{}, errors.New("any-error")).Once()
			},
		},
		{
			name: "failed to add worklog to parent task",
			task: entity.Task{
				ProjectID:    "project-1",
				ParentTaskID: "parent-task-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: fmt.Errorf("error while adding worklog: %w", errors.New("any-error")),
			mockFunc: func(task entity.Task, timeSpent time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
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
				t.taskRepo.On("GetByID", mock.Anything, task.ParentTaskID).Return(entity.Task{
					ID:   task.ParentTaskID,
					Name: "Parent Task 1",
				}, nil).Once()
				t.jiraRepo.On("AddWorklog", mock.Anything, "parent-task-1", "Parent Task 1", timeSpent, mock.Anything).Return(errors.New("any-error")).Once()
			},
		},
		{
			name: "failed to add worklog to task",
			task: entity.Task{
				ProjectID: "project-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: fmt.Errorf("error while adding worklog: %w", errors.New("any-error")),
			mockFunc: func(task entity.Task, timeSpent time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
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
				t.jiraRepo.On("AddWorklog", mock.Anything, task.ID, task.Name, timeSpent, mock.Anything).Return(errors.New("any-error")).Once()
			},
		},
		{
			name: "success add worklog to parent task",
			task: entity.Task{
				ProjectID:    "project-1",
				ParentTaskID: "parent-task-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: nil,
			mockFunc: func(task entity.Task, timeSpent time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
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
				t.taskRepo.On("GetByID", mock.Anything, task.ParentTaskID).Return(entity.Task{
					ID:   task.ParentTaskID,
					Name: "Parent Task 1",
				}, nil).Once()
				t.jiraRepo.On("AddWorklog", mock.Anything, "parent-task-1", "Parent Task 1", timeSpent, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "success add worklog to task",
			task: entity.Task{
				ProjectID: "project-1",
			},
			timeSpent:     time.Duration(0),
			expectedError: nil,
			mockFunc: func(task entity.Task, timeSpent time.Duration) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
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
				t.jiraRepo.On("AddWorklog", mock.Anything, task.ID, task.Name, timeSpent, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			if test.mockFunc != nil {
				test.mockFunc(test.task, test.timeSpent)
			}

			err := t.useCase.AddWorklog(context.Background(), test.task, test.timeSpent)
			t.Equal(test.expectedError, err)
		})
	}

}
