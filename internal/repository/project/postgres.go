package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project *domain.Project) error {
	query := `
		INSERT INTO projects (id, name, icon, "desc", tech, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query, project.ID, project.Name, project.Icon, project.Desc, project.Tech, project.URL, project.CreatedAt, project.UpdatedAt)
	return err
}

func (r *ProjectRepository) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	query := `
		SELECT id, name, icon, "desc", tech, url, created_at, updated_at
		FROM projects WHERE id = $1
	`
	project := &domain.Project{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&project.ID, &project.Name, &project.Icon, &project.Desc, &project.Tech, &project.URL, &project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepository) GetAllProjects(ctx context.Context) ([]*domain.Project, error) {
	query := `
		SELECT id, name, icon, "desc", tech, url, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*domain.Project, 0)
	for rows.Next() {
		project := &domain.Project{}
		err := rows.Scan(
			&project.ID, &project.Name, &project.Icon, &project.Desc, &project.Tech, &project.URL, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, id string, project *domain.Project) error {
	query := `
		UPDATE projects
		SET name = $1, icon = $2, "desc" = $3, tech = $4, url = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.pool.Exec(ctx, query, project.Name, project.Icon, project.Desc, project.Tech, project.URL, project.UpdatedAt, id)
	return err
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
