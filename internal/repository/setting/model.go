package setting

import (
	"encoding/json"

	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
)

const KeySyncIntegration = "sync_integration"

type IntegrationModel struct {
	Type    string            `json:"type"`
	Details map[string]string `json:"details"`
}

type SettingModel struct {
	Key   string
	Value string
}

func CreateModelFromSyncIntegration(integration syncintegration.SyncIntegration) (SettingModel, error) {
	integrationModel := IntegrationModel{
		Type:    string(integration.Type),
		Details: integration.Details,
	}

	detailBytes, err := json.Marshal(integrationModel)
	if err != nil {
		return SettingModel{}, err
	}

	return SettingModel{
		Key:   KeySyncIntegration,
		Value: string(detailBytes),
	}, nil
}

func (sm SettingModel) ToSyncIntegration() (syncintegration.SyncIntegration, error) {
	var integrationModel IntegrationModel
	err := json.Unmarshal([]byte(sm.Value), &integrationModel)
	if err != nil {
		return syncintegration.SyncIntegration{}, err
	}

	return syncintegration.SyncIntegration{
		Type:    syncintegration.SyncIntegrationType(integrationModel.Type),
		Details: integrationModel.Details,
	}, nil
}
