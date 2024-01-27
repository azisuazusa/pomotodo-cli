package task

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	_ "github.com/mattn/go-sqlite3"
)

type RepoImpl struct {
	db *sql.DB
}

func New(db *sql.DB) *RepoImpl {
	return &RepoImpl{db: db}
}

func (ri *RepoImpl) GetUncompleteParentTasks(ctx context.Context, projectID string) (entity.Tasks, error) {
	var tasks []entity.Task
	query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND parent_task_id = '' OR parent_task_id IS NULL`
	rows, err := ri.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task TaskModel
		err := rows.Scan(&task.ID, &task.ProjectID, &task.Name, &task.Description, &task.IsStarted, &task.CompletedAt, &task.ParentTaskID, &task.Integration, &task.Histories)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		taskEntity, err := task.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert task to entity: %w", err)
		}

		tasks = append(tasks, taskEntity)
	}

	return tasks, nil
}

func (ri *RepoImpl) GetUncompleteSubTask(ctx context.Context, projectID string) (map[string]entity.Tasks, error) {
	query := `SELECT * FROM tasks WHERE completed_at IS NULL AND project_id = ? AND (parent_task_id IS NOT NULL AND parent_task_id != '')`
	rows, err := ri.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := map[string]entity.Tasks{}
	for rows.Next() {
		var task TaskModel
		err := rows.Scan(&task.ID, &task.ProjectID, &task.Name, &task.Description, &task.IsStarted, &task.CompletedAt, &task.ParentTaskID, &task.Integration, &task.Histories)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		taskEntity, err := task.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert task to entity: %w", err)
		}

		tasks[task.ParentTaskID.String] = append(tasks[task.ParentTaskID.String], taskEntity)
	}

	return tasks, nil
}

func (ri *RepoImpl) Insert(ctx context.Context, taskEntity entity.Task) error {
	task, err := CreateModel(taskEntity)
	if err != nil {
		return fmt.Errorf("failed to create task model: %w", err)
	}

	query := `INSERT INTO tasks (id, project_id, name, description, is_started, completed_at, parent_task_id, integration, histories) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = ri.db.ExecContext(ctx, query, task.ID, task.ProjectID, task.Name, task.Description, task.IsStarted, task.CompletedAt, task.ParentTaskID, task.Integration, task.Histories)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

func (ri *RepoImpl) Update(ctx context.Context, taskEntity entity.Task) error {
	task, err := CreateModel(taskEntity)
	if err != nil {
		return fmt.Errorf("failed to create task model: %w", err)
	}

	query := `UPDATE tasks SET project_id = ?, name = ?, description = ?, is_started = ?, completed_at = ?, parent_task_id = ?, integration = ?, histories = ? WHERE id = ?`
	_, err = ri.db.ExecContext(ctx, query, task.ProjectID, task.Name, task.Description, task.IsStarted, task.CompletedAt, task.ParentTaskID, task.Integration, task.Histories, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (ri *RepoImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := ri.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (ri *RepoImpl) GetByID(ctx context.Context, id string) (entity.Task, error) {
	query := `SELECT * FROM tasks WHERE id = ?`
	row := ri.db.QueryRowContext(ctx, query, id)

	var task TaskModel
	err := row.Scan(&task.ID, &task.ProjectID, &task.Name, &task.Description, &task.IsStarted, &task.CompletedAt, &task.ParentTaskID, &task.Integration, &task.Histories)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to scan task: %w", err)
	}

	taskEntity, err := task.ToEntity()
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to convert task to entity: %w", err)
	}

	return taskEntity, nil
}

func (ri *RepoImpl) Upsert(ctx context.Context, taskEntity entity.Task) error {
	task, err := CreateModel(taskEntity)
	if err != nil {
		return fmt.Errorf("failed to create task model: %w", err)
	}

	query := `INSERT INTO tasks (id, project_id, name, description, is_started, completed_at, parent_task_id, integration) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET project_id = ?, name = ?, description = ?, is_started = ?, completed_at = ?, parent_task_id = ?, integration = ?`
	_, err = ri.db.ExecContext(ctx, query, task.ID, task.ProjectID, task.Name, task.Description, task.IsStarted, task.CompletedAt, task.ParentTaskID, task.Integration, task.ProjectID, task.Name, task.Description, task.IsStarted, task.CompletedAt, task.ParentTaskID, task.Integration)
	if err != nil {
		return fmt.Errorf("failed to upsert task: %w", err)
	}

	return nil
}

func (ri *RepoImpl) GetStartedTask(ctx context.Context) (entity.Task, error) {
	query := `SELECT * FROM tasks WHERE is_started = ? LIMIT 1`
	row := ri.db.QueryRowContext(ctx, query, true)

	var task TaskModel
	err := row.Scan(&task.ID, &task.ProjectID, &task.Name, &task.Description, &task.IsStarted, &task.CompletedAt, &task.ParentTaskID, &task.Integration, &task.Histories)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to scan task: %w", err)
	}

	taskEntity, err := task.ToEntity()
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to convert task to entity: %w", err)
	}

	return taskEntity, nil
}
