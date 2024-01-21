package setting

type IntegrationType string

const Dropbox IntegrationType = "dropbox"

type Integration struct {
	Type    IntegrationType
	Details map[string]string
}
