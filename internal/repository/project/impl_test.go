package project

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/stretchr/testify/suite"
)

type RepoImplTestSuite struct {
	suite.Suite
	db       sqlmock.Sqlmock
	repoImpl RepoImpl
}

func (t *RepoImplTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	t.db = mock

	t.repoImpl = RepoImpl{db: db}
}

func TestRepoImpl(t *testing.T) {
	suite.Run(t, new(RepoImplTestSuite))
}

func (t *RepoImplTestSuite) TestGetAll() {
	query := "SELECT * FROM projects"
	tests := []struct {
		name        string
		expectedRes entity.Projects
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get all projects",
			expectedRes: nil,
			expectedErr: errors.New("failed to get all projects"),
			mockFunc: func() {
				t.db.ExpectQuery(query).WillReturnError(errors.New("failed to get all projects"))
			},
		},
		{
			name:        "failed to scan project",
			expectedRes: nil,
			expectedErr: errors.New("sql: expected 2 destination arguments in Scan, not 5"),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("project-1", "Project 1")
				t.db.ExpectQuery(query).WillReturnRows(rows)
			},
		},
		{
			name:        "failed to convert project to entity",
			expectedRes: nil,
			expectedErr: errors.New("failed to unmarshal integrations: invalid character 'i' looking for beginning of value"),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow("project-1", "Project 1", "description", true, "integrations")
				t.db.ExpectQuery(query).WillReturnRows(rows)
			},
		},
		{
			name: "success",
			expectedRes: entity.Projects{
				{
					ID:          "project-1",
					Name:        "Project 1",
					Description: "description",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      "JIRA",
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				},
			},
			expectedErr: nil,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow("project-1", "Project 1", "description", true, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`)
				t.db.ExpectQuery(query).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc()
			res, err := t.repoImpl.GetAll(context.Background())
			t.Equal(tt.expectedRes, res)
			if err != nil {
				err = errors.Unwrap(err)
				t.Equal(tt.expectedErr.Error(), err.Error())
				return
			}
			t.Nil(err)
		})
	}
}

func (t *RepoImplTestSuite) TestInsert() {
	query := "INSERT INTO projects (id, name, description, is_selected, integrations) VALUES (?, ?, ?, ?, ?)"
	tests := []struct {
		name        string
		project     entity.Project
		expectedErr error
		mockFunc    func(project entity.Project)
	}{
		{
			name: "failed to insert project",
			project: entity.Project{
				ID:          "project-1",
				Name:        "Project 1",
				Description: "description",
				IsSelected:  true,
				Integrations: []entity.Integration{
					{
						IsEnabled: true,
						Type:      "JIRA",
						Details: map[string]string{
							"token": "token",
						},
					},
				},
			},
			expectedErr: sql.ErrConnDone,
			mockFunc: func(project entity.Project) {
				t.db.ExpectExec(query).WithArgs(project.ID, project.Name, project.Description, project.IsSelected, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`).WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name: "success",
			project: entity.Project{
				ID:          "project-1",
				Name:        "Project 1",
				Description: "description",
				IsSelected:  true,
				Integrations: []entity.Integration{
					{
						IsEnabled: true,
						Type:      "JIRA",
						Details: map[string]string{
							"token": "token",
						},
					},
				},
			},
			expectedErr: nil,
			mockFunc: func(project entity.Project) {
				t.db.ExpectExec(query).WithArgs(project.ID, project.Name, project.Description, project.IsSelected, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.project)
			err := t.repoImpl.Insert(context.Background(), tt.project)
			if err != nil {
				err = errors.Unwrap(err)
				t.Equal(tt.expectedErr.Error(), err.Error())
				return
			}
			t.Nil(err)
		})
	}
}

func (t *RepoImplTestSuite) TestUpdate() {
	query := "UPDATE projects SET name = ?, description = ?, is_selected = ?, integrations = ? WHERE id = ?"
	tests := []struct {
		name        string
		project     entity.Project
		expectedErr error
		mockFunc    func(project entity.Project)
	}{
		{
			name: "failed to update project",
			project: entity.Project{
				ID:          "project-1",
				Name:        "Project 1",
				Description: "description",
				IsSelected:  true,
				Integrations: []entity.Integration{
					{
						IsEnabled: true,
						Type:      "JIRA",
						Details: map[string]string{
							"token": "token",
						},
					},
				},
			},
			expectedErr: sql.ErrConnDone,
			mockFunc: func(project entity.Project) {
				t.db.ExpectExec(query).WithArgs(project.Name, project.Description, project.IsSelected, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`, project.ID).WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name: "success",
			project: entity.Project{
				ID:          "project-1",
				Name:        "Project 1",
				Description: "description",
				IsSelected:  true,
				Integrations: []entity.Integration{
					{
						IsEnabled: true,
						Type:      "JIRA",
						Details: map[string]string{
							"token": "token",
						},
					},
				},
			},
			expectedErr: nil,
			mockFunc: func(project entity.Project) {
				t.db.ExpectExec(query).WithArgs(project.Name, project.Description, project.IsSelected, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`, project.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.project)
			err := t.repoImpl.Update(context.Background(), tt.project)
			if err != nil {
				err = errors.Unwrap(err)
				t.Equal(tt.expectedErr.Error(), err.Error())
				return
			}
			t.Nil(err)
		})
	}
}

func (t *RepoImplTestSuite) TestDelete() {
	query := "DELETE FROM projects WHERE id = ?"
	tests := []struct {
		name        string
		projectID   string
		expectedErr error
		mockFunc    func(projectID string)
	}{
		{
			name:        "failed to delete project",
			projectID:   "project-1",
			expectedErr: sql.ErrConnDone,
			mockFunc: func(projectID string) {
				t.db.ExpectExec(query).WithArgs(projectID).WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "success",
			projectID:   "project-1",
			expectedErr: nil,
			mockFunc: func(projectID string) {
				t.db.ExpectExec(query).WithArgs(projectID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.projectID)
			err := t.repoImpl.Delete(context.Background(), tt.projectID)
			if err != nil {
				err = errors.Unwrap(err)
				t.Equal(tt.expectedErr.Error(), err.Error())
				return
			}
			t.Nil(err)
		})
	}
}

func (t *RepoImplTestSuite) TestGetByID() {
	query := "SELECT * FROM projects WHERE id = ? LIMIT 1"
	tests := []struct {
		name        string
		projectID   string
		expected    entity.Project
		expectedErr error
		mockFunc    func(projectID string)
	}{
		{
			name:        "failed to get project by id",
			projectID:   "project-1",
			expected:    entity.Project{},
			expectedErr: sql.ErrNoRows,
			mockFunc: func(projectID string) {
				t.db.ExpectQuery(query).WithArgs(projectID).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "failed to scan project",
			projectID:   "project-1",
			expected:    entity.Project{},
			expectedErr: errors.New("sql: Scan error on column index 1, name \"name\": converting NULL to string is unsupported"),
			mockFunc: func(projectID string) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow(projectID, nil, "description", true, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`)
				t.db.ExpectQuery(query).WithArgs(projectID).WillReturnRows(rows)
			},
		},
		{
			name:        "failed to unmarshal integrations",
			projectID:   "project-1",
			expected:    entity.Project{},
			expectedErr: errors.New("failed to unmarshal integrations: invalid character 'i' looking for beginning of value"),
			mockFunc: func(projectID string) {
				project := entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
					Description: "description",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      "JIRA",
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow(project.ID, project.Name, project.Description, project.IsSelected, `invalid_json`)
				t.db.ExpectQuery(query).WithArgs(projectID).WillReturnRows(rows)
			},
		},
		{
			name:      "success",
			projectID: "project-1",
			expected: entity.Project{
				ID:          "project-1",
				Name:        "Project 1",
				Description: "description",
				IsSelected:  true,
				Integrations: []entity.Integration{
					{
						IsEnabled: true,
						Type:      "JIRA",
						Details: map[string]string{
							"token": "token",
						},
					},
				},
			},
			expectedErr: nil,
			mockFunc: func(projectID string) {
				project := entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
					Description: "description",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      "JIRA",
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}
				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow(project.ID, project.Name, project.Description, project.IsSelected, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`)
				t.db.ExpectQuery(query).WithArgs(projectID).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.projectID)
			project, err := t.repoImpl.GetByID(context.Background(), tt.projectID)
			if err != nil {
				err = errors.Unwrap(err)
				t.Equal(tt.expectedErr.Error(), err.Error())
				return
			}
			t.Equal(tt.expected, project)
		})
	}
}

func (t *RepoImplTestSuite) TestGetSelectedProject() {
	query := "SELECT * FROM projects WHERE is_selected = ? LIMIT 1"
	tests := []struct {
		name        string
		expected    entity.Project
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get selected project",
			expected:    entity.Project{},
			expectedErr: sql.ErrNoRows,
			mockFunc: func() {
				t.db.ExpectQuery(query).WithArgs(true).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "failed to scan project",
			expected:    entity.Project{},
			expectedErr: errors.New("sql: Scan error on column index 1, name \"name\": converting NULL to string is unsupported"),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow("project-1", nil, "description", true, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`)
				t.db.ExpectQuery(query).WithArgs(true).WillReturnRows(rows)
			},
		},
		{
			name:        "failed to unmarshal integrations",
			expected:    entity.Project{},
			expectedErr: errors.New("failed to unmarshal integrations: invalid character 'i' looking for beginning of value"),
			mockFunc: func() {
				project := entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
					Description: "description",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      "JIRA",
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}

				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow(project.ID, project.Name, project.Description, project.IsSelected, `invalid_json`)
				t.db.ExpectQuery(query).WithArgs(true).WillReturnRows(rows)
			},
		},
		{
			name: "success",
			expected: entity.Project{
				ID:          "project-1",
				Name:        "Project 1",
				Description: "description",
				IsSelected:  true,
				Integrations: []entity.Integration{
					{
						IsEnabled: true,
						Type:      "JIRA",
						Details: map[string]string{
							"token": "token",
						},
					},
				},
			},
			expectedErr: nil,
			mockFunc: func() {
				project := entity.Project{
					ID:          "project-1",
					Name:        "Project 1",
					Description: "description",
					IsSelected:  true,
					Integrations: []entity.Integration{
						{
							IsEnabled: true,
							Type:      "JIRA",
							Details: map[string]string{
								"token": "token",
							},
						},
					},
				}
				rows := sqlmock.NewRows([]string{"id", "name", "description", "is_selected", "integrations"}).AddRow(project.ID, project.Name, project.Description, project.IsSelected, `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`)
				t.db.ExpectQuery(query).WithArgs(true).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc()
			project, err := t.repoImpl.GetSelectedProject(context.Background())
			if err != nil {
				err = errors.Unwrap(err)
				t.Equal(tt.expectedErr.Error(), err.Error())
				return
			}
			t.Equal(tt.expected, project)
		})
	}
}

func (t *RepoImplTestSuite) TestSetSelectedProject() {
	queryDeselected := "UPDATE projects SET is_selected = ?"
	querySelected := "UPDATE projects SET is_selected = ? WHERE id = ?"

	tests := []struct {
		name        string
		projectID   string
		expectedErr error
		mockFunc    func(projectID string)
	}{
		{
			name:        "failed to deselected project",
			projectID:   "project-1",
			expectedErr: sql.ErrConnDone,
			mockFunc: func(_ string) {
				t.db.ExpectExec(queryDeselected).WithArgs(false).WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "failed to selected project",
			projectID:   "project-1",
			expectedErr: sql.ErrConnDone,
			mockFunc: func(projectID string) {
				t.db.ExpectExec(queryDeselected).WithArgs(false).WillReturnResult(sqlmock.NewResult(1, 1))
				t.db.ExpectExec(querySelected).WithArgs(true, projectID).WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "success",
			projectID:   "project-1",
			expectedErr: nil,
			mockFunc: func(projectID string) {
				t.db.ExpectExec(queryDeselected).WithArgs(false).WillReturnResult(sqlmock.NewResult(1, 1))
				t.db.ExpectExec(querySelected).WithArgs(true, projectID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func() {
			tt.mockFunc(tt.projectID)
			err := t.repoImpl.SetSelectedProject(context.Background(), tt.projectID)
			if err != nil {
				err = errors.Unwrap(err)
				t.Equal(tt.expectedErr.Error(), err.Error())
				return
			}
		})
	}
}
