package task

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/google/uuid"
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
	Description  sql.NullString
	IsStarted    bool
	CompletedAt  sql.NullTime
	ParentTaskID sql.NullString
	Integration  sql.NullString
	Histories    sql.NullString
}

func (tm TaskModel) ToEntity() (entity.Task, error) {
	task := entity.Task{
		ID:           tm.ID,
		ProjectID:    tm.ProjectID,
		Name:         tm.Name,
		Description:  tm.Description.String,
		IsStarted:    tm.IsStarted,
		CompletedAt:  tm.CompletedAt.Time,
		ParentTaskID: tm.ParentTaskID.String,
	}

	if tm.Integration.Valid {
		var integrationModel TaskIntegrationModel
		err := json.Unmarshal([]byte(tm.Integration.String), &integrationModel)
		if err != nil {
			return entity.Task{}, fmt.Errorf("failed to unmarshal integration: %w", err)
		}

		task.Integration = entity.TaskIntegration{
			ID:   integrationModel.ID,
			Type: entity.IntegrationType(integrationModel.Type),
		}
	}

	if tm.Histories.Valid {
		var historyModels []TaskHistoryModel
		err := json.Unmarshal([]byte(tm.Histories.String), &historyModels)
		if err != nil {
			return entity.Task{}, fmt.Errorf("failed to unmarshal histories: %w", err)
		}

		for _, historyModel := range historyModels {
			task.Histories = append(task.Histories, entity.TaskHistory{
				StartedAt: historyModel.StartedAt,
				StoppedAt: historyModel.StoppedAt,
			})
		}
	}

	return task, nil
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

	if task.ID == "" {
		task.ID = uuid.NewString()
	}

	return TaskModel{
		ID:           task.ID,
		ProjectID:    task.ProjectID,
		Name:         task.Name,
		Description:  sql.NullString{String: task.Description, Valid: task.Description != ""},
		IsStarted:    task.IsStarted,
		CompletedAt:  sql.NullTime{Time: task.CompletedAt, Valid: !task.CompletedAt.IsZero()},
		ParentTaskID: sql.NullString{String: task.ParentTaskID, Valid: task.ParentTaskID != ""},
		Integration:  sql.NullString{String: string(integrationBytes), Valid: len(integrationBytes) > 0},
		Histories:    sql.NullString{String: string(historiesBytes), Valid: len(historiesBytes) > 0},
	}, nil
}
