package entity

import "time"

type TaskHistory struct {
	StartedAt time.Time
	StoppedAt time.Time
}

type IntegrationType string

const (
	IntegrationTypeJIRA   IntegrationType = "JIRA"
	IntegrationTypeGitHub IntegrationType = "GitHub"
)

type TaskIntegration struct {
	ID   string
	Type IntegrationType
}

type Task struct {
	ID           string
	ProjectID    string
	Name         string
	Description  string
	IsStarted    bool
	CompletedAt  time.Time
	ParentTaskID string
	Integration  TaskIntegration
	Histories    []TaskHistory
}

func (t Task) TimeSpent() time.Duration {
	var timeSpent time.Duration
	for _, history := range t.Histories {
		timeSpent += history.StoppedAt.Sub(history.StartedAt)
	}

	return timeSpent
}

func (t *Task) Stop() {
	t.IsStarted = false
	t.Histories[len(t.Histories)-1].StoppedAt = time.Now()
}

func (t *Task) Start() {
	t.IsStarted = true
	t.Histories = append(t.Histories, TaskHistory{
		StartedAt: time.Now(),
	})
}

func (t *Task) Complete() {
	timeNow := time.Now()
	if t.IsStarted {
		t.IsStarted = false
		t.Histories[len(t.Histories)-1].StoppedAt = timeNow
	}
	t.CompletedAt = timeNow
}

type Tasks []Task
