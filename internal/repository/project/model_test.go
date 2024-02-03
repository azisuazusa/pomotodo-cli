package project

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestProjectModelToEntity(t *testing.T) {
	tests := []struct {
		name        string
		model       ProjectModel
		expectedRes entity.Project
		expectedErr error
	}{
		{
			name: "failed to unmarshal integrations",
			model: ProjectModel{
				ID:           "project-1",
				Name:         "Project 1",
				IsSelected:   true,
				Integrations: sql.NullString{String: "invalid-json", Valid: true},
			},
			expectedRes: entity.Project{},
			expectedErr: errors.New("invalid character 'i' looking for beginning of value"),
		},
		{
			name: "success",
			model: ProjectModel{
				ID:           "project-1",
				Name:         "Project 1",
				IsSelected:   true,
				Integrations: sql.NullString{String: `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`, Valid: true},
			},
			expectedRes: entity.Project{
				ID:         "project-1",
				Name:       "Project 1",
				IsSelected: true,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.model.ToEntity()
			assert.Equal(t, tt.expectedRes, res)
			if err != nil {
				err = errors.Unwrap(err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestCreateModel(t *testing.T) {
	projectEntity := entity.Project{
		ID:         "project-1",
		Name:       "Project 1",
		IsSelected: true,
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

	expectedModel := ProjectModel{
		ID:         "project-1",
		Name:       "Project 1",
		IsSelected: true,
		Integrations: sql.NullString{
			String: `[{"is_enabled":true,"type":"JIRA","details":{"token":"token"}}]`,
			Valid:  true,
		},
	}

	model, err := CreateModel(projectEntity)
	assert.Nil(t, err)
	assert.Equal(t, expectedModel, model)
}
