package project

import (
	"testing"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateProjectView(t *testing.T) {
	param := entity.Project{
		ID:          "project-1",
		Name:        "Project 1",
		Description: "Description 1",
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

	expected := ProjectView{
		Name:        "Project 1 [Integrated with: JIRA] (Selected)",
		Description: "Description 1",
	}

	res := CreateProjectView(param)

	assert.Equal(t, expected, res)
}
