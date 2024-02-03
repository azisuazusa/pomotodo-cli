package project

import (
	"context"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type useCase struct {
	projectRepo      ProjectRepository
	taskRepo         TaskRepository
	integrationRepos map[entity.IntegrationType]IntegrationRepository
}

type UseCase interface {
	GetAll(ctx context.Context) (entity.Projects, error)
	Insert(ctx context.Context, project entity.Project) error
	Delete(ctx context.Context, id string) error
	Select(ctx context.Context, id string) error
	SyncTasks(ctx context.Context) error
	AddIntegration(ctx context.Context, integration entity.Integration) error
	RemoveIntegration(ctx context.Context, integrationType entity.IntegrationType) error
	GetSelected(ctx context.Context) (entity.Project, error)
	EnableIntegration(ctx context.Context, integrationType entity.IntegrationType) error
	DisableIntegration(ctx context.Context, integrationType entity.IntegrationType) error
}

func New(projectRepo ProjectRepository, integrationRepos map[entity.IntegrationType]IntegrationRepository, taskRepo TaskRepository) UseCase {
	return &useCase{
		projectRepo:      projectRepo,
		taskRepo:         taskRepo,
		integrationRepos: integrationRepos,
	}
}

func (u *useCase) GetAll(ctx context.Context) (entity.Projects, error) {
	return u.projectRepo.GetAll(ctx)
}

func (u *useCase) Insert(ctx context.Context, project entity.Project) error {
	err := u.projectRepo.Insert(ctx, project)
	if err != nil {
		return fmt.Errorf("error while inserting project: %w", err)
	}

	return nil
}

func (u *useCase) Delete(ctx context.Context, id string) error {
	err := u.projectRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error while deleting project: %w", err)
	}

	return nil
}

func (u *useCase) Select(ctx context.Context, id string) error {
	err := u.projectRepo.SetSelectedProject(ctx, id)
	if err != nil {
		return fmt.Errorf("error selecting project %w", err)
	}

	return nil
}

func (u *useCase) SyncTasks(ctx context.Context) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return fmt.Errorf("error while getting selected project: %w", err)
	}

	if len(project.Integrations) == 0 {
		fmt.Println("no integration found")
		return nil
	}

	var tasks entity.Tasks
	for _, integration := range project.Integrations {
		if integration.IsEnabled {
			integrationRepo := u.integrationRepos[integration.Type]
			integrationTasks, err := integrationRepo.GetTasks(ctx, project.ID, integration.Details)
			if err != nil {
				return fmt.Errorf("error while syncing tasks: %w", err)
			}

			tasks = append(tasks, integrationTasks...)
		}
	}

	for _, task := range tasks {
		err := u.taskRepo.Upsert(ctx, task)
		if err != nil {
			return fmt.Errorf("error while inserting task: %w", err)
		}
	}

	return nil
}

func (u *useCase) AddIntegration(ctx context.Context, integration entity.Integration) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return fmt.Errorf("error while getting selected project: %w", err)
	}

	project.Integrations = append(project.Integrations, integration)
	err = u.projectRepo.Update(ctx, project)
	if err != nil {
		return fmt.Errorf("error while updating project: %w", err)
	}

	return nil
}

func (u *useCase) RemoveIntegration(ctx context.Context, integrationType entity.IntegrationType) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return fmt.Errorf("error while getting selected project: %w", err)
	}

	for i, integration := range project.Integrations {
		if integration.Type == integrationType {
			project.Integrations = append(project.Integrations[:i], project.Integrations[i+1:]...)
			break
		}
	}

	err = u.projectRepo.Update(ctx, project)
	if err != nil {
		return fmt.Errorf("error while updating project: %w", err)
	}

	return nil
}

func (u *useCase) GetSelected(ctx context.Context) (entity.Project, error) {
	return u.projectRepo.GetSelectedProject(ctx)
}

func (u *useCase) EnableIntegration(ctx context.Context, integrationType entity.IntegrationType) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return fmt.Errorf("error while getting selected project: %w", err)
	}

	for i, integration := range project.Integrations {
		if integration.Type == integrationType {
			project.Integrations[i].IsEnabled = true
			break
		}
	}

	err = u.projectRepo.Update(ctx, project)
	if err != nil {
		return fmt.Errorf("error while updating project: %w", err)
	}

	return nil
}

func (u *useCase) DisableIntegration(ctx context.Context, integrationType entity.IntegrationType) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return fmt.Errorf("error while getting selected project: %w", err)
	}

	for i, integration := range project.Integrations {
		if integration.Type == integrationType {
			project.Integrations[i].IsEnabled = false
			break
		}
	}

	err = u.projectRepo.Update(ctx, project)
	if err != nil {
		return fmt.Errorf("error while updating project: %w", err)
	}

	return nil
}
