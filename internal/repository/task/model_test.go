package task

import (
	"database/sql"
	"testing"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestTaskModelToEntity(t *testing.T) {
	tests := []struct {
		name    string
		task    TaskModel
		want    entity.Task
		wantErr bool
	}{
		{
			name: "unmarshal integration failed",
			task: TaskModel{
				ID:        "1",
				ProjectID: "1",
				Name:      "test",
				Description: sql.NullString{
					String: "description",
					Valid:  true,
				},
				IsStarted: false,
				CompletedAt: sql.NullTime{
					Time:  time.Time{},
					Valid: false,
				},
				ParentTaskID: sql.NullString{
					String: "",
					Valid:  false,
				},
				Integration: sql.NullString{
					String: "{invalid integration",
					Valid:  true,
				},
				Histories: sql.NullString{
					String: "",
					Valid:  false,
				},
			},
			want:    entity.Task{},
			wantErr: true,
		},
		{
			name: "unmarshal histories success",
			task: TaskModel{
				ID:        "1",
				ProjectID: "1",
				Name:      "test",
				Description: sql.NullString{
					String: "description",
					Valid:  true,
				},
				IsStarted: false,
				CompletedAt: sql.NullTime{
					Time:  time.Time{},
					Valid: false,
				},
				ParentTaskID: sql.NullString{
					String: "",
					Valid:  false,
				},
				Integration: sql.NullString{
					String: "",
					Valid:  false,
				},
				Histories: sql.NullString{
					String: "{invalid histories",
					Valid:  true,
				},
			},
			want:    entity.Task{},
			wantErr: true,
		},
		{
			name: "success",
			task: TaskModel{
				ID:        "1",
				ProjectID: "1",
				Name:      "test",
				Description: sql.NullString{
					String: "description",
					Valid:  true,
				},
				IsStarted: false,
				CompletedAt: sql.NullTime{
					Time:  time.Time{},
					Valid: false,
				},
				ParentTaskID: sql.NullString{
					String: "",
					Valid:  false,
				},
				Integration: sql.NullString{
					String: `{"id": "1", "type": "test"}`,
					Valid:  true,
				},
				Histories: sql.NullString{
					String: `[{"started_at": "2021-01-01T00:00:00Z", "stopped_at": "2021-01-01T00:00:00Z"}]`,
					Valid:  true,
				},
			},
			want: entity.Task{
				ID:           "1",
				ProjectID:    "1",
				Name:         "test",
				Description:  "description",
				IsStarted:    false,
				CompletedAt:  time.Time{},
				ParentTaskID: "",
				Integration: entity.TaskIntegration{
					ID:   "1",
					Type: "test",
				},
				Histories: []entity.TaskHistory{
					{
						StartedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						StoppedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.task.ToEntity()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestCreateModel(t *testing.T) {
	paramTask := entity.Task{
		ID:           "1",
		ProjectID:    "1",
		Name:         "test",
		Description:  "description",
		IsStarted:    false,
		CompletedAt:  time.Time{},
		ParentTaskID: "",
		Integration: entity.TaskIntegration{
			ID:   "1",
			Type: "test",
		},
		Histories: []entity.TaskHistory{
			{
				StartedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				StoppedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	expectedTask := TaskModel{
		ID:        "1",
		ProjectID: "1",
		Name:      "test",
		Description: sql.NullString{
			String: "description",
			Valid:  true,
		},
		IsStarted: false,
		CompletedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		ParentTaskID: sql.NullString{
			String: "",
			Valid:  false,
		},
		Integration: sql.NullString{
			String: `{"id":"1","type":"test"}`,
			Valid:  true,
		},
		Histories: sql.NullString{
			String: `[{"started_at":"2021-01-01T00:00:00Z","stopped_at":"2021-01-01T00:00:00Z"}]`,
			Valid:  true,
		},
	}

	got, err := CreateModel(paramTask)

	assert.Nil(t, err)
	assert.Equal(t, expectedTask, got)
}
