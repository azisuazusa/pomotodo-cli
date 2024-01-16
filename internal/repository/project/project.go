package project

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

func (ri *RepoImpl) GetAll(ctx context.Context) (entity.Project, error) {
	var project ProjectModel
	query := `SELECT * FROM projects`
	row := ri.db.QueryRowContext(ctx, query)
	err := row.Scan(&project.ID, &project.Name, &project.Description, &project.IsSelected, &project.Integrations)
	if err != nil {
		return entity.Project{}, fmt.Errorf("failed to scan project: %w", err)
	}

	projectEntity, err := project.ToEntity()
	if err != nil {
		return entity.Project{}, fmt.Errorf("failed to convert project to entity: %w", err)
	}

	return projectEntity, nil
}

func (ri *RepoImpl) Insert(ctx context.Context, projectEntity entity.Project) error {
	project, err := CreateModel(projectEntity)
	if err != nil {
		return fmt.Errorf("failed to create project model: %w", err)
	}

	query := `INSERT INTO projects (name, description, is_selected, integrations) VALUES (?, ?, ?, ?)`
	_, err = ri.db.ExecContext(ctx, query, project.Name, project.Description, project.IsSelected, project.Integrations)
	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}

	return nil
}

func (ri *RepoImpl) Update(ctx context.Context, projectEntity entity.Project) error {
	project, err := CreateModel(projectEntity)
	if err != nil {
		return fmt.Errorf("failed to create project model: %w", err)
	}

	query := `UPDATE projects SET name = ?, description = ?, is_selected = ?, integrations = ? WHERE id = ?`
	_, err = ri.db.ExecContext(ctx, query, project.Name, project.Description, project.IsSelected, project.Integrations, project.ID)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

func (ri *RepoImpl) Delete(ctx context.Context, projectEntity entity.Project) error {
	project, err := CreateModel(projectEntity)
	if err != nil {
		return fmt.Errorf("failed to create project model: %w", err)
	}

	query := `DELETE FROM projects WHERE id = ?`
	_, err = ri.db.ExecContext(ctx, query, project.ID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

func (ri *RepoImpl) GetByID(ctx context.Context, id int) (entity.Project, error) {
	var project ProjectModel
	query := `SELECT * FROM projects WHERE id = ?`
	row := ri.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&project.ID, &project.Name, &project.Description, &project.IsSelected, &project.Integrations)
	if err != nil {
		return entity.Project{}, fmt.Errorf("failed to scan project: %w", err)
	}

	projectEntity, err := project.ToEntity()
	if err != nil {
		return entity.Project{}, fmt.Errorf("failed to convert project to entity: %w", err)
	}

	return projectEntity, nil
}

func (ri *RepoImpl) GetSelectedProject(ctx context.Context) (entity.Project, error) {
	var project ProjectModel
	query := `SELECT * FROM projects WHERE is_selected = ? LIMIT 1`
	row := ri.db.QueryRowContext(ctx, query, true)
	err := row.Scan(&project.ID, &project.Name, &project.Description, &project.IsSelected, &project.Integrations)
	if err != nil {
		return entity.Project{}, fmt.Errorf("failed to scan project: %w", err)
	}

	projectEntity, err := project.ToEntity()
	if err != nil {
		return entity.Project{}, fmt.Errorf("failed to convert project to entity: %w", err)
	}

	return projectEntity, nil
}
