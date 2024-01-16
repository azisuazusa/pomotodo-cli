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

type Tasks []Task
