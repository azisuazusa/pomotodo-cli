package project

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/google/uuid"
)

type IntegrationModel struct {
	IsEnabled bool              `json:"is_enabled"`
	Type      string            `json:"type"`
	Details   map[string]string `json:"details"`
}

type ProjectModel struct {
	ID           string
	Name         string
	Description  sql.NullString
	IsSelected   bool
	Integrations sql.NullString
}

func (pm ProjectModel) ToEntity() (entity.Project, error) {
	project := entity.Project{
		ID:          pm.ID,
		Name:        pm.Name,
		Description: pm.Description.String,
		IsSelected:  pm.IsSelected,
	}

	if pm.Integrations.Valid {
		var integrationModels []IntegrationModel
		err := json.Unmarshal([]byte(pm.Integrations.String), &integrationModels)
		if err != nil {
			return entity.Project{}, fmt.Errorf("failed to unmarshal integrations: %w", err)
		}

		var integrations []entity.Integration
		for _, integration := range integrationModels {
			integrations = append(integrations, entity.Integration{
				IsEnabled: integration.IsEnabled,
				Type:      entity.IntegrationType(integration.Type),
				Details:   integration.Details,
			})
		}

		project.Integrations = integrations
	}

	return project, nil
}

func CreateModel(entity entity.Project) (ProjectModel, error) {
	var integrationModels []IntegrationModel
	for _, integration := range entity.Integrations {
		integrationModels = append(integrationModels, IntegrationModel{
			IsEnabled: integration.IsEnabled,
			Type:      string(integration.Type),
			Details:   integration.Details,
		})
	}

	jsonIntegrations, err := json.Marshal(integrationModels)
	if err != nil {
		return ProjectModel{}, fmt.Errorf("failed to marshal integrations: %w", err)
	}

	if entity.ID == "" {
		entity.ID = uuid.NewString()
	}

	return ProjectModel{
		ID:           entity.ID,
		Name:         entity.Name,
		Description:  sql.NullString{String: entity.Description, Valid: entity.Description != ""},
		IsSelected:   entity.IsSelected,
		Integrations: sql.NullString{String: string(jsonIntegrations), Valid: len(integrationModels) > 0},
	}, nil
}
