package project

import (
	"encoding/json"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type IntegrationModel struct {
	IsEnabled bool              `json:"is_enabled"`
	Type      string            `json:"type"`
	Details   map[string]string `json:"auth"`
}

type ProjectModel struct {
	ID           string
	Name         string
	Description  string
	IsSelected   bool
	Integrations string
}

func (pm ProjectModel) ToEntity() (entity.Project, error) {
	var integrationModels []IntegrationModel
	err := json.Unmarshal([]byte(pm.Integrations), &integrationModels)
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

	return entity.Project{
		ID:           pm.ID,
		Name:         pm.Name,
		Description:  pm.Description,
		IsSelected:   pm.IsSelected,
		Integrations: integrations,
	}, nil
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

	return ProjectModel{
		ID:           entity.ID,
		Name:         entity.Name,
		Description:  entity.Description,
		IsSelected:   entity.IsSelected,
		Integrations: string(jsonIntegrations),
	}, nil
}
