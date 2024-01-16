package task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type TaskHistoryModel struct {
	StartedAt time.Time `json:"started_at"`
	StoppedAt time.Time `json:"stopped_at"`
}

type TaskIntegrationModel struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type TaskModel struct {
	ID           string
	ProjectID    string
	Name         string
	Description  string
	IsStarted    bool
	CompletedAt  time.Time
	ParentTaskID string
	Integration  string
	Histories    string
}

func (tm TaskModel) ToEntity() (entity.Task, error) {
	var integrationModel TaskIntegrationModel
	err := json.Unmarshal([]byte(tm.Integration), &integrationModel)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to unmarshal integration: %w", err)
	}

	var historyModels []TaskHistoryModel
	err = json.Unmarshal([]byte(tm.Histories), &historyModels)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to unmarshal histories: %w", err)
	}

	var histories []entity.TaskHistory
	for _, historyModel := range historyModels {
		histories = append(histories, entity.TaskHistory{
			StartedAt: historyModel.StartedAt,
			StoppedAt: historyModel.StoppedAt,
		})
	}

	return entity.Task{
		ID:           tm.ID,
		ProjectID:    tm.ProjectID,
		Name:         tm.Name,
		Description:  tm.Description,
		IsStarted:    tm.IsStarted,
		CompletedAt:  tm.CompletedAt,
		ParentTaskID: tm.ParentTaskID,
		Integration: entity.TaskIntegration{
			ID:   integrationModel.ID,
			Type: entity.IntegrationType(integrationModel.Type),
		},
		Histories: histories,
	}, nil
}

func CreateModel(task entity.Task) (TaskModel, error) {
	integrationBytes, err := json.Marshal(TaskIntegrationModel{
		ID:   task.Integration.ID,
		Type: string(task.Integration.Type),
	})
	if err != nil {
		return TaskModel{}, fmt.Errorf("failed to marshal integration: %w", err)
	}

	var historyModels []TaskHistoryModel
	for _, history := range task.Histories {
		historyModels = append(historyModels, TaskHistoryModel{
			StartedAt: history.StartedAt,
			StoppedAt: history.StoppedAt,
		})
	}
	historiesBytes, err := json.Marshal(historyModels)
	if err != nil {
		return TaskModel{}, fmt.Errorf("failed to marshal histories: %w", err)
	}

	return TaskModel{
		ID:           task.ID,
		ProjectID:    task.ProjectID,
		Name:         task.Name,
		Description:  task.Description,
		IsStarted:    task.IsStarted,
		CompletedAt:  task.CompletedAt,
		ParentTaskID: task.ParentTaskID,
		Integration:  string(integrationBytes),
		Histories:    string(historiesBytes),
	}, nil
}
