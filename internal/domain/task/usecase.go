package task

import (
	"context"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
)

type UseCase interface {
	GetUncompleteTasks(ctx context.Context) (entity.Tasks, error)
	GetUncompleteParentTasks(ctx context.Context) (entity.Tasks, error)
	Add(ctx context.Context, task entity.Task) error
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context) error
	Remove(ctx context.Context, id string) error
	Complete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (entity.Task, error)
}

type useCase struct {
	taskRepo    TaskRepository
	projectRepo ProjectRepository
}

func New(taskRepo TaskRepository, projectRepo ProjectRepository) UseCase {
	return &useCase{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (u *useCase) GetUncompleteTasks(ctx context.Context) (entity.Tasks, error) {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting selected project: %w", err)
	}

	parentTasks, err := u.taskRepo.GetUncompleteParentTasks(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("error while getting uncomplete parent tasks: %w", err)
	}

	subTasks, err := u.taskRepo.GetUncompleteSubTask(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("error while getting uncomplete sub tasks: %w", err)
	}

	var tasks entity.Tasks
	for _, task := range parentTasks {
		tasks = append(tasks, task)
		tasks = append(tasks, subTasks[task.ID]...)
	}

	return tasks, nil

}

func (u *useCase) GetUncompleteParentTasks(ctx context.Context) (entity.Tasks, error) {
	project, err := u.projectRepo.GetSelectedProject(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting selected project: %w", err)
	}

	tasks, err := u.taskRepo.GetUncompleteParentTasks(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("error while getting uncomplete parent tasks: %w", err)
	}

	return tasks, nil
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

	return nil
}

func (u *useCase) Start(ctx context.Context, id string) error {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error while getting task: %w", err)
	}

	task.Start()
	err = u.taskRepo.SetStartedTask(ctx, task)
	if err != nil {
		return fmt.Errorf("error while setting started task: %w", err)
	}

	return nil
}

func (u *useCase) Stop(ctx context.Context) (err error) {
	task, err := u.taskRepo.GetStartedTask(ctx)
	if err != nil {
		return fmt.Errorf("error while getting started task: %w", err)
	}

	task.Stop()
	err = u.taskRepo.Update(ctx, task)
	if err != nil {
		return fmt.Errorf("error while updating task: %w", err)
	}

	return nil
}

func (u *useCase) Remove(ctx context.Context, id string) error {
	err := u.taskRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error while deleting task: %w", err)
	}

	return nil
}

func (u *useCase) Complete(ctx context.Context, id string) error {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error while getting task: %w", err)
	}

	task.Complete()
	err = u.taskRepo.Update(ctx, task)
	if err != nil {
		return fmt.Errorf("error while updating task: %w", err)
	}

	return nil
}

func (u *useCase) GetByID(ctx context.Context, id string) (entity.Task, error) {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("error while getting task: %w", err)
	}

	return task, nil
}
