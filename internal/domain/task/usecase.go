package task

import (
	"context"
	"fmt"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type UseCase interface {
	GetUncompleteTasks(ctx context.Context) (entity.Tasks, error)
	Add(ctx context.Context, task entity.Task) error
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Complete(ctx context.Context, id string) error
}

type useCase struct {
	taskRepo    TaskRepository
	syncRepo    SyncRepository
	projectRepo ProjectRepository
}

func New(taskRepo TaskRepository, syncRepo SyncRepository, projectRepo ProjectRepository) UseCase {
	return &useCase{
		taskRepo:    taskRepo,
		syncRepo:    syncRepo,
		projectRepo: projectRepo,
	}
}

func (u *useCase) GetUncompleteTasks(ctx context.Context) (entity.Tasks, error) {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting selected project: %w", err)
	}

	return u.taskRepo.GetUncompleteTasks(ctx, project.ID)
}

func (u *useCase) Add(ctx context.Context, task entity.Task) error {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return fmt.Errorf("error while getting selected project: %w", err)
	}

	task.ProjectID = project.ID
	err = u.taskRepo.Insert(ctx, task)
	if err != nil {
		return fmt.Errorf("error while inserting task: %w", err)
	}

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}

func (u *useCase) Start(ctx context.Context, id string) error {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error while getting task: %w", err)
	}

	task.IsStarted = true

	err = u.taskRepo.Update(ctx, task)
	if err != nil {
		return fmt.Errorf("error while updating task: %w", err)
	}

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}

func (u *useCase) Stop(ctx context.Context, id string) error {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error while getting task: %w", err)
	}

	task.IsStarted = false
	task.Histories[len(task.Histories)-1].StoppedAt = time.Now()

	err = u.taskRepo.Update(ctx, task)
	if err != nil {
		return fmt.Errorf("error while updating task: %w", err)
	}

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}

func (u *useCase) Remove(ctx context.Context, id string) error {
	err := u.taskRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error while deleting task: %w", err)
	}

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}

func (u *useCase) Complete(ctx context.Context, id string) error {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error while getting task: %w", err)
	}

	timeNow := time.Now()
	task.CompletedAt = timeNow
	task.IsStarted = false
	task.Histories[len(task.Histories)-1].StoppedAt = timeNow

	err = u.taskRepo.Update(ctx, task)
	if err != nil {
		return fmt.Errorf("error while updating task: %w", err)
	}

	go func() {
		err := u.syncRepo.Sync(ctx)
		if err != nil {
			fmt.Printf("error while syncing: %v\n", err)
		}
	}()

	return nil
}
