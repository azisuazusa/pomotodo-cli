package project

import (
	"fmt"
	"strings"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/manifoldco/promptui"
)

var projectSelectTemplate = &promptui.SelectTemplates{
	Active:   `▸ {{ .Name | cyan }}`,
	Inactive: `  {{ .Name }}`,
	Selected: `{{ "✔" | green }} {{ "Selected" | bold }}: {{ .Name | cyan }}`,
	Details:  `{{ "Description:" }} {{ .Description }}`,
}

type ProjectView struct {
	Name        string
	Description string
}

func CreateProjectView(p entity.Project) ProjectView {
	projectName := p.Name
	if len(p.Integrations) > 0 {
		integrationTypes := []string{}
		for _, integration := range p.Integrations {
			integrationTypes = append(integrationTypes, string(integration.Type))
		}
		projectName = fmt.Sprintf("%s [Integrated with: %s]", projectName, strings.Join(integrationTypes, ", "))
	}

	if p.IsSelected {
		projectName = fmt.Sprintf("%s (Started)", projectName)
	}

	return ProjectView{
		Name:        projectName,
		Description: p.Description,
	}
}
