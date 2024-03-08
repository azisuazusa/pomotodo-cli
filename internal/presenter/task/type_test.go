package task

import (
	"testing"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateTaskView(t *testing.T) {
	paramNumber := "1."
	paramTask := entity.Task{
		ID:           "task-1",
		ProjectID:    "project-1",
		Name:         "Task 1",
		Description:  "Description 1",
		IsStarted:    true,
		ParentTaskID: "parent-task-1",
		Integration: entity.TaskIntegration{
			ID:   "integration-1",
			Type: "JIRA",
		},
	}

	expected := TaskView{
		Name:         "1. Task 1 (Started) (integration-1)",
		Description:  "Description 1",
		ParentTaskID: "parent-task-1",
	}

	res := CreateTaskView(paramNumber, paramTask)

	assert.Equal(t, expected, res)
}
