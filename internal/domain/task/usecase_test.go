package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/azisuazusa/todo-cli/internal/domain/task/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UseCaseTestSuite struct {
	suite.Suite
	taskRepo    *mocks.TaskRepository
	projectRepo *mocks.ProjectRepository
	useCase     UseCase
}

func (t *UseCaseTestSuite) SetupTest() {
	t.taskRepo = &mocks.TaskRepository{}
	t.projectRepo = &mocks.ProjectRepository{}
	t.useCase = New(t.taskRepo, t.projectRepo)
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}

func (t *UseCaseTestSuite) TestGetUncompleteTasks() {
	tests := []struct {
		name        string
		expectedRes entity.Tasks
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get select project",
			expectedRes: nil,
			expectedErr: entity.ErrNoProjectSelected,
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{}, entity.ErrNoProjectSelected).Once()
			},
		},
		{
			name:        "failed to get uncomplete tasks",
			expectedRes: nil,
			expectedErr: errors.New("failed to get uncomplete tasks"),
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID: "project-1",
				}, nil).Once()
				t.taskRepo.On("GetUncompleteParentTasks", mock.Anything, "project-1").Return(nil, errors.New("failed to get uncomplete tasks")).Once()
			},
		},
		{
			name:        "failed to get uncomplete subtasks",
			expectedRes: nil,
			expectedErr: errors.New("failed to get uncomplete subtasks"),
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID: "project-1",
				}, nil).Once()
				t.taskRepo.On("GetUncompleteParentTasks", mock.Anything, "project-1").Return(entity.Tasks{
					{
						ID: "task-1",
					},
				}, nil).Once()
				t.taskRepo.On("GetUncompleteSubTask", mock.Anything, "project-1").Return(nil, errors.New("failed to get uncomplete subtasks")).Once()
			},
		},
		{
			name: "success",
			expectedRes: entity.Tasks{
				{
					ID: "task-1",
				},
				{
					ID: "subtask-1",
				},
			},
			expectedErr: nil,
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID: "project-1",
				}, nil).Once()
				t.taskRepo.On("GetUncompleteParentTasks", mock.Anything, "project-1").Return(entity.Tasks{
					{
						ID: "task-1",
					},
				}, nil).Once()
				t.taskRepo.On("GetUncompleteSubTask", mock.Anything, "project-1").Return(map[string]entity.Tasks{
					"task-1": {
						{
							ID: "subtask-1",
						},
					},
				}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc()
			res, err := t.useCase.GetUncompleteTasks(context.Background())
			t.Equal(tt.expectedRes, res)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestGetUncompleteParentTasks() {
	tests := []struct {
		name        string
		expectedRes entity.Tasks
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get selected project",
			expectedRes: nil,
			expectedErr: entity.ErrNoProjectSelected,
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{}, entity.ErrNoProjectSelected).Once()
			},
		},
		{
			name:        "failed to get uncomplete parent tasks",
			expectedRes: nil,
			expectedErr: errors.New("failed to get uncomplete parent tasks"),
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID: "project-1",
				}, nil).Once()
				t.taskRepo.On("GetUncompleteParentTasks", mock.Anything, "project-1").Return(nil, errors.New("failed to get uncomplete parent tasks")).Once()
			},
		},
		{
			name: "success",
			expectedRes: entity.Tasks{
				{
					ID: "task-1",
				},
			},
			expectedErr: nil,
			mockFunc: func() {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID: "project-1",
				}, nil).Once()
				t.taskRepo.On("GetUncompleteParentTasks", mock.Anything, "project-1").Return(entity.Tasks{
					{
						ID: "task-1",
					},
				}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc()
			res, err := t.useCase.GetUncompleteParentTasks(context.Background())
			t.Equal(tt.expectedRes, res)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestAdd() {
	tests := []struct {
		name        string
		task        entity.Task
		expectedErr error
		mockFunc    func(task entity.Task)
	}{
		{
			name:        "failed to get selected project",
			task:        entity.Task{},
			expectedErr: entity.ErrNoProjectSelected,
			mockFunc: func(_ entity.Task) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{}, entity.ErrNoProjectSelected).Once()
			},
		},
		{
			name: "failed to add task",
			task: entity.Task{
				ID: "task-1",
			},
			expectedErr: errors.New("failed to add task"),
			mockFunc: func(task entity.Task) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID: "project-1",
				}, nil).Once()
				task.ProjectID = "project-1"
				t.taskRepo.On("Insert", mock.Anything, task).Return(errors.New("failed to add task")).Once()
			},
		},
		{
			name: "success",
			task: entity.Task{
				ID: "task-1",
			},
			expectedErr: nil,
			mockFunc: func(task entity.Task) {
				t.projectRepo.On("GetSelectedProject", mock.Anything).Return(entity.Project{
					ID: "project-1",
				}, nil).Once()
				task.ProjectID = "project-1"
				t.taskRepo.On("Insert", mock.Anything, task).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.task)
			err := t.useCase.Add(context.Background(), tt.task)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestStart() {
	tests := []struct {
		name        string
		taskID      string
		expectedErr error
		mockFunc    func(taskID string)
	}{
		{
			name:        "failed to get task by ID",
			taskID:      "task-1",
			expectedErr: errors.New("failed to get task by ID"),
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{}, errors.New("failed to get task by ID")).Once()
			},
		},
		{
			name:        "failed to start task",
			taskID:      "task-1",
			expectedErr: errors.New("failed to start task"),
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{
					ID: "task-1",
				}, nil).Once()
				t.taskRepo.On("SetStartedTask", mock.Anything, mock.Anything).Return(errors.New("failed to start task")).Once()
			},
		},
		{
			name:        "success",
			taskID:      "task-1",
			expectedErr: nil,
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{
					ID: "task-1",
				}, nil).Once()
				t.taskRepo.On("SetStartedTask", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.taskID)
			err := t.useCase.Start(context.Background(), tt.taskID)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestStop() {
	tests := []struct {
		name        string
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get started task",
			expectedErr: errors.New("failed to get started task"),
			mockFunc: func() {
				t.taskRepo.On("GetStartedTask", mock.Anything).Return(entity.Task{}, errors.New("failed to get started task")).Once()
			},
		},
		{
			name:        "failed to stop task",
			expectedErr: errors.New("failed to stop task"),
			mockFunc: func() {
				t.taskRepo.On("GetStartedTask", mock.Anything).Return(entity.Task{
					ID: "task-1",
					Histories: []entity.TaskHistory{
						{
							StartedAt: time.Time{},
						},
					},
				}, nil).Once()
				t.taskRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed to stop task")).Once()
			},
		},
		{
			name:        "success",
			expectedErr: nil,
			mockFunc: func() {
				t.taskRepo.On("GetStartedTask", mock.Anything).Return(entity.Task{
					ID: "task-1",
					Histories: []entity.TaskHistory{
						{
							StartedAt: time.Time{},
						},
					},
				}, nil).Once()
				t.taskRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc()

			err := t.useCase.Stop(context.Background())
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestRemove() {
	tests := []struct {
		name        string
		taskID      string
		expectedErr error
		mockFunc    func(taskID string)
	}{
		{
			name:        "failed to remove task",
			taskID:      "task-1",
			expectedErr: errors.New("failed to remove task"),
			mockFunc: func(taskID string) {
				t.taskRepo.On("Delete", mock.Anything, taskID).Return(errors.New("failed to remove task")).Once()
			},
		},
		{
			name:        "success",
			taskID:      "task-1",
			expectedErr: nil,
			mockFunc: func(taskID string) {
				t.taskRepo.On("Delete", mock.Anything, taskID).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.taskID)
			err := t.useCase.Remove(context.Background(), tt.taskID)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestComplete() {
	tests := []struct {
		name        string
		taskID      string
		expectedErr error
		mockFunc    func(taskID string)
	}{
		{
			name:        "failed to get task by ID",
			taskID:      "task-1",
			expectedErr: errors.New("failed to get task by ID"),
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{}, errors.New("failed to get task by ID")).Once()
			},
		},
		{
			name:        "failed to complete task",
			taskID:      "task-1",
			expectedErr: errors.New("failed to complete task"),
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{
					ID: "task-1",
				}, nil).Once()
				t.taskRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("failed to complete task")).Once()
			},
		},
		{
			name:        "success",
			taskID:      "task-1",
			expectedErr: nil,
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{
					ID: "task-1",
				}, nil).Once()
				t.taskRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.taskID)
			err := t.useCase.Complete(context.Background(), tt.taskID)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expectedErr, err)
		})
	}
}

func (t *UseCaseTestSuite) TestGetByID() {
	tests := []struct {
		name        string
		taskID      string
		expected    entity.Task
		expectedErr error
		mockFunc    func(taskID string)
	}{
		{
			name:        "failed to get task by ID",
			taskID:      "task-1",
			expected:    entity.Task{},
			expectedErr: errors.New("failed to get task by ID"),
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{}, errors.New("failed to get task by ID")).Once()
			},
		},
		{
			name:        "success",
			taskID:      "task-1",
			expected:    entity.Task{ID: "task-1"},
			expectedErr: nil,
			mockFunc: func(taskID string) {
				t.taskRepo.On("GetByID", mock.Anything, taskID).Return(entity.Task{ID: "task-1"}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.taskID)
			task, err := t.useCase.GetByID(context.Background(), tt.taskID)
			if err != nil {
				err = errors.Unwrap(err)
			}
			t.Equal(tt.expected, task)
			t.Equal(tt.expectedErr, err)
		})
	}
}
