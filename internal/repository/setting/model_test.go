package setting

import (
	"testing"

	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/stretchr/testify/assert"
)

func TestCreateModelFromSyncIntegration(t *testing.T) {
	paramIntegration := syncintegration.SyncIntegration{
		Type: syncintegration.Dropbox,
		Details: map[string]string{
			"token": "test-token",
		},
	}

	expected := SettingModel{
		Key:   "sync_integration",
		Value: "{\"type\":\"dropbox\",\"details\":{\"token\":\"test-token\"}}",
	}

	actual, err := CreateModelFromSyncIntegration(paramIntegration)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestSettingModelToSyncIntegration(t *testing.T) {
	paramSettingModel := SettingModel{
		Key:   "sync_integration",
		Value: "{\"type\":\"dropbox\",\"details\":{\"token\":\"test-token\"}}",
	}

	expected := syncintegration.SyncIntegration{
		Type: syncintegration.Dropbox,
		Details: map[string]string{
			"token": "test-token",
		},
	}

	actual, err := paramSettingModel.ToSyncIntegration()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
