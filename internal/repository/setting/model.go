package setting

import (
	"encoding/json"

	"github.com/azisuazusa/todo-cli/internal/domain/setting"
	"github.com/google/uuid"
)

type IntegrationModel struct {
	Type    string            `json:"type"`
	Details map[string]string `json:"details"`
}

type SettingModel struct {
	ID             string
	SettingType    string
	SettingDetails string
}

func CreateModelFromIntegration(integration setting.Integration) (SettingModel, error) {
	detailBytes, err := json.Marshal(integration.Details)
	if err != nil {
		return SettingModel{}, err
	}

	return SettingModel{
		ID:             uuid.NewString(),
		SettingType:    string(integration.Type),
		SettingDetails: string(detailBytes),
	}, nil
}

func (sm SettingModel) ToIntegrationEntity() (setting.Integration, error) {
	var details map[string]string
	err := json.Unmarshal([]byte(sm.SettingDetails), &details)
	if err != nil {
		return setting.Integration{}, err
	}

	return setting.Integration{
		Type:    setting.IntegrationType(sm.SettingType),
		Details: details,
	}, nil
}
