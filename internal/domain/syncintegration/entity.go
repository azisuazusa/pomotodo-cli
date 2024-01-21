package syncintegration

type SyncIntegrationType string

const Dropbox SyncIntegrationType = "dropbox"

type SyncIntegration struct {
	Type    SyncIntegrationType
	Details map[string]string
}
