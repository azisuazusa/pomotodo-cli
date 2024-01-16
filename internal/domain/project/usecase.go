package project

import (
	"context"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type useCase struct {
	projectRepo      ProjectRepository
	syncRepo         SyncRepository
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
}

func New(projectRepo ProjectRepository, syncRepo SyncRepository, taskRepo TaskRepository) UseCase {
	return &useCase{
		projectRepo: projectRepo,
		syncRepo:    syncRepo,
		taskRepo:    taskRepo,
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

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}

func (u *useCase) Delete(ctx context.Context, id string) error {
	err := u.projectRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error while deleting project: %w", err)
	}

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}

func (u *useCase) Select(ctx context.Context, id string) error {
	project, err := u.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error while getting project: %w", err)
	}

	project.IsSelected = true
	err = u.projectRepo.Update(ctx, project)
	if err != nil {
		return fmt.Errorf("error while updating project: %w", err)
	}

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}

func (u *useCase) SyncTasks(ctx context.Context) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return fmt.Errorf("error while getting selected project: %w", err)
	}

	var tasks entity.Tasks
	for _, integration := range project.Integrations {
		if integration.IsEnabled {
			integrationRepo := u.integrationRepos[integration.Type]
			integrationTasks, err := integrationRepo.GetTasks(ctx, integration.Details)
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

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}
