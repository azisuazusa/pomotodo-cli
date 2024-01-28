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

func (ri *RepoImpl) GetAll(ctx context.Context) (entity.Projects, error) {
	query := `SELECT * FROM projects`
	rows, err := ri.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	result := entity.Projects{}
	for rows.Next() {
		project := ProjectModel{}
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.IsSelected, &project.Integrations)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		entity, err := project.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert project to entity: %w", err)
		}

		result = append(result, entity)
	}

	return result, nil

}

func (ri *RepoImpl) Insert(ctx context.Context, projectEntity entity.Project) error {
	project, err := CreateModel(projectEntity)
	if err != nil {
		return fmt.Errorf("failed to create project model: %w", err)
	}

	query := `INSERT INTO projects (id, name, description, is_selected, integrations) VALUES (?, ?, ?, ?, ?)`
	_, err = ri.db.ExecContext(ctx, query, project.ID, project.Name, project.Description, project.IsSelected, project.Integrations)
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

func (ri *RepoImpl) Delete(ctx context.Context, projectID string) error {
	query := `DELETE FROM projects WHERE id = ?`
	_, err := ri.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

func (ri *RepoImpl) GetByID(ctx context.Context, id string) (entity.Project, error) {
	var project ProjectModel
	query := `SELECT * FROM projects WHERE id = ? LIMIT 1`
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

func (ri *RepoImpl) SetSelectedProject(ctx context.Context, id string) (err error) {
	query := `UPDATE projects SET is_selected = ?`
	_, err = ri.db.ExecContext(ctx, query, false)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	query = `UPDATE projects SET is_selected = ? WHERE id = ?`
	_, err = ri.db.ExecContext(ctx, query, true, id)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}
