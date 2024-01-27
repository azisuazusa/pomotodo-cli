package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	_ "github.com/mattn/go-sqlite3"
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

func (s *RepoImplTestSuite) TestGetUncompleteParentTasks() {
	tests := []struct {
		name           string
		projectID      string
		expectedResult entity.Tasks
		expectedError  error
		mock           func(projectID string)
	}{
		{
			name:      "failed to get uncomplete parent tasks",
			projectID: "1",
			mock: func(projectID string) {
				query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND parent_task_id = '' OR parent_task_id IS NULL`
				s.db.ExpectQuery(query).WithArgs(projectID).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name:      "failed to scan task",
			projectID: "1",
			mock: func(projectID string) {
				query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND parent_task_id = '' OR parent_task_id IS NULL`
				s.db.ExpectQuery(query).WithArgs(projectID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
			},
			expectedError: errors.New("sql: expected 1 destination arguments in Scan, not 9"),
		},
		{
			name:      "success",
			projectID: "1",
			mock: func(projectID string) {
				query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND parent_task_id = '' OR parent_task_id IS NULL`
				s.db.ExpectQuery(query).WithArgs(projectID).WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "name", "description", "is_started", "completed_at", "parent_task_id", "integration", "histories"}).AddRow("1", "1", "name", "description", false, nil, nil, nil, nil))
			},
			expectedResult: entity.Tasks{
				{
					ID:           "1",
					ProjectID:    "1",
					Name:         "name",
					Description:  "description",
					IsStarted:    false,
					CompletedAt:  time.Time{},
					ParentTaskID: "",
					Integration:  entity.TaskIntegration{},
					Histories:    []entity.TaskHistory(nil),
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.projectID)

			result, err := s.repoImpl.GetUncompleteParentTasks(context.Background(), tt.projectID)

			s.Equal(tt.expectedResult, result)

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
		})
	}
}

func (s *RepoImplTestSuite) TestGetUncompleteSubTask() {
	tests := []struct {
		name           string
		projectID      string
		expectedResult map[string]entity.Tasks
		expectedError  error
		mock           func(projectID string)
	}{
		{
			name:      "failed to get uncomplete sub tasks",
			projectID: "1",
			mock: func(projectID string) {
				query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND (parent_task_id IS NOT NULL AND parent_task_id != '')`
				s.db.ExpectQuery(query).WithArgs(projectID).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name:      "failed to scan task",
			projectID: "1",
			mock: func(projectID string) {
				query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND (parent_task_id IS NOT NULL AND parent_task_id != '')`
				s.db.ExpectQuery(query).WithArgs(projectID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
			},
			expectedError: errors.New("sql: expected 1 destination arguments in Scan, not 9"),
		},
		{
			name:      "success",
			projectID: "1",
			mock: func(projectID string) {
				query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND (parent_task_id IS NOT NULL AND parent_task_id != '')`
				s.db.ExpectQuery(query).WithArgs(projectID).WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "name", "description", "is_started", "completed_at", "parent_task_id", "integration", "histories"}).AddRow("1", "1", "name", "description", false, nil, "1", nil, nil))
			},
			expectedResult: map[string]entity.Tasks{
				"1": {
					{
						ID:           "1",
						ProjectID:    "1",
						Name:         "name",
						Description:  "description",
						IsStarted:    false,
						CompletedAt:  time.Time{},
						ParentTaskID: "1",
						Integration:  entity.TaskIntegration{},
						Histories:    []entity.TaskHistory(nil),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.projectID)

			result, err := s.repoImpl.GetUncompleteSubTask(context.Background(), tt.projectID)

			s.Equal(tt.expectedResult, result)

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
		})
	}
}

func (s *RepoImplTestSuite) TestInsert() {
	tests := []struct {
		name          string
		task          entity.Task
		expectedError error
		mock          func(task entity.Task)
	}{
		{
			name: "failed to insert task",
			task: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "name",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentTaskID: "1",
				Integration: entity.TaskIntegration{
					ID:   "1",
					Type: entity.IntegrationType("JIRA"),
				},
				Histories: []entity.TaskHistory(nil),
			},
			mock: func(task entity.Task) {
				taskModel, _ := CreateModel(task)
				query := `INSERT INTO tasks (id, project_id, name, description, is_started, completed_at, parent_task_id, integration, histories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
				s.db.ExpectExec(query).WithArgs(taskModel.ID, taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration, taskModel.Histories).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name: "success",
			task: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "name",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentTaskID: "1",
				Integration: entity.TaskIntegration{
					ID:   "1",
					Type: entity.IntegrationType("JIRA"),
				},
				Histories: []entity.TaskHistory(nil),
			},
			mock: func(task entity.Task) {
				taskModel, _ := CreateModel(task)
				query := `INSERT INTO tasks (id, project_id, name, description, is_started, completed_at, parent_task_id, integration, histories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
				s.db.ExpectExec(query).WithArgs(taskModel.ID, taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration, taskModel.Histories).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.task)

			err := s.repoImpl.Insert(context.Background(), tt.task)

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
		})
	}
}

func (s *RepoImplTestSuite) TestUpdate() {
	tests := []struct {
		name          string
		task          entity.Task
		expectedError error
		mock          func(task entity.Task)
	}{
		{
			name: "failed to update task",
			task: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "name",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentTaskID: "1",
				Integration: entity.TaskIntegration{
					ID:   "1",
					Type: entity.IntegrationType("JIRA"),
				},
				Histories: []entity.TaskHistory(nil),
			},
			mock: func(task entity.Task) {
				taskModel, _ := CreateModel(task)
				query := `UPDATE tasks SET project_id = ?, name = ?, description = ?, is_started = ?, completed_at = ?, parent_task_id = ?, integration = ?, histories = ? WHERE id = ?`
				s.db.ExpectExec(query).WithArgs(taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration, taskModel.Histories, taskModel.ID).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name: "success",
			task: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "name",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentTaskID: "1",
				Integration: entity.TaskIntegration{
					ID:   "1",
					Type: entity.IntegrationType("JIRA"),
				},
				Histories: []entity.TaskHistory(nil),
			},
			mock: func(task entity.Task) {
				taskModel, _ := CreateModel(task)
				query := `UPDATE tasks SET project_id = ?, name = ?, description = ?, is_started = ?, completed_at = ?, parent_task_id = ?, integration = ?, histories = ? WHERE id = ?`
				s.db.ExpectExec(query).WithArgs(taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration, taskModel.Histories, taskModel.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.task)

			err := s.repoImpl.Update(context.Background(), tt.task)

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
		})
	}
}

func (s *RepoImplTestSuite) TestDelete() {
	tests := []struct {
		name          string
		taskID        string
		expectedError error
		mock          func(taskID string)
	}{
		{
			name:   "failed to delete task",
			taskID: "1",
			mock: func(taskID string) {
				query := `DELETE FROM tasks WHERE id = ?`
				s.db.ExpectExec(query).WithArgs(taskID).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name:   "success",
			taskID: "1",
			mock: func(taskID string) {
				query := `DELETE FROM tasks WHERE id = ?`
				s.db.ExpectExec(query).WithArgs(taskID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.taskID)

			err := s.repoImpl.Delete(context.Background(), tt.taskID)

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
		})
	}
}

func (s *RepoImplTestSuite) TestGetByID() {
	tests := []struct {
		name          string
		taskID        string
		expectedTask  entity.Task
		expectedError error
		mock          func(taskID string)
	}{
		{
			name:   "failed to get task by id",
			taskID: "1",
			mock: func(taskID string) {
				query := `SELECT * FROM tasks WHERE id = ?`
				s.db.ExpectQuery(query).WithArgs(taskID).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name:   "success",
			taskID: "1",
			mock: func(taskID string) {
				query := `SELECT * FROM tasks WHERE id = ?`
				rows := sqlmock.NewRows([]string{"id", "project_id", "name", "description", "is_started", "completed_at", "parent_task_id", "integration", "histories"}).
					AddRow(taskID, "1", "name", "description", false, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), "1", "{}", "[]")
				s.db.ExpectQuery(query).WithArgs(taskID).WillReturnRows(rows)
			},
			expectedTask: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "name",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentTaskID: "1",
				Integration:  entity.TaskIntegration{},
				Histories:    []entity.TaskHistory(nil),
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.taskID)

			task, err := s.repoImpl.GetByID(context.Background(), tt.taskID)

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
			s.Equal(tt.expectedTask, task)
		})
	}
}

func (s *RepoImplTestSuite) TestUpsert() {
	tests := []struct {
		name          string
		task          entity.Task
		expectedError error
		mock          func(task entity.Task)
	}{
		{
			name: "failed to upsert task",
			task: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "name",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentTaskID: "1",
				Integration: entity.TaskIntegration{
					ID:   "1",
					Type: entity.IntegrationType("JIRA"),
				},
				Histories: []entity.TaskHistory(nil),
			},
			mock: func(task entity.Task) {
				taskModel, _ := CreateModel(task)
				query := `INSERT INTO tasks (id, project_id, name, description, is_started, completed_at, parent_task_id, integration) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET project_id = ?, name = ?, description = ?, is_started = ?, completed_at = ?, parent_task_id = ?, integration = ?`
				s.db.ExpectExec(query).WithArgs(taskModel.ID, taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration, taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name: "success",
			task: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "name",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentTaskID: "1",
				Integration: entity.TaskIntegration{
					ID:   "1",
					Type: entity.IntegrationType("JIRA"),
				},
				Histories: []entity.TaskHistory(nil),
			},
			mock: func(task entity.Task) {
				taskModel, _ := CreateModel(task)
				query := `INSERT INTO tasks (id, project_id, name, description, is_started, completed_at, parent_task_id, integration) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET project_id = ?, name = ?, description = ?, is_started = ?, completed_at = ?, parent_task_id = ?, integration = ?`
				s.db.ExpectExec(query).WithArgs(taskModel.ID, taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration, taskModel.ProjectID, taskModel.Name, taskModel.Description, taskModel.IsStarted, taskModel.CompletedAt, taskModel.ParentTaskID, taskModel.Integration).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock(tt.task)

			err := s.repoImpl.Upsert(context.Background(), tt.task)

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
		})
	}
}

func (s *RepoImplTestSuite) TestGetStartedTask() {
	tests := []struct {
		name          string
		expectedError error
		mock          func()
	}{
		{
			name: "failed to get started task",
			mock: func() {
				query := `SELECT * FROM tasks WHERE is_started = ? LIMIT 1`
				s.db.ExpectQuery(query).WithArgs(true).WillReturnError(errors.New("any-error"))
			},
			expectedError: errors.New("any-error"),
		},
		{
			name: "success",
			mock: func() {
				query := `SELECT * FROM tasks WHERE is_started = ? LIMIT 1`
				rows := sqlmock.NewRows([]string{"id", "project_id", "name", "description", "is_started", "completed_at", "parent_task_id", "integration", "histories"}).
					AddRow("1", "1", "name", "description", true, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), "1", "{}", "[]")
				s.db.ExpectQuery(query).WithArgs(true).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()

			_, err := s.repoImpl.GetStartedTask(context.Background())

			if err != nil {
				err = errors.Unwrap(err)
			}
			s.Equal(tt.expectedError, err)
		})
	}
}
