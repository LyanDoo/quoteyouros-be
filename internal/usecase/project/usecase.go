package project

import (
	"context"

	"github.com/quoteyouros/backend/internal/domain"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
)

type ProjectUseCase struct {
	projectRepo domain.ProjectRepository
}

// New creates a new project use case
func New(projectRepo domain.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo: projectRepo,
	}
}

// CreateProject creates a new project
func (u *ProjectUseCase) CreateProject(ctx context.Context, req *domain.CreateProjectRequest) (*domain.Project, error) {
	logger.Debug("createProject: validating request", "name", req.Name)

	// Create project
	project := domain.NewProject(req.Name, req.Icon, req.Desc, req.Tech, req.URL)

	logger.Info("createProject: creating project in database", "project_id", project.ID, "name", project.Name)
	if err := u.projectRepo.CreateProject(ctx, project); err != nil {
		logger.Error("createProject: failed to create project", "project_id", project.ID, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to create project: " + err.Error())
	}

	logger.Info("createProject: project created successfully", "project_id", project.ID, "name", project.Name)
	return project, nil
}

// GetProject retrieves a single project
func (u *ProjectUseCase) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	logger.Debug("getProject: retrieving project", "project_id", id)

	project, err := u.projectRepo.GetProject(ctx, id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "no rows in result set" {
			logger.Warn("getProject: project not found", "project_id", id)
			return nil, apperrors.NotFound("project not found")
		}
		logger.Error("getProject: failed to retrieve project", "project_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve project: " + err.Error())
	}

	logger.Debug("getProject: project retrieved successfully", "project_id", id, "name", project.Name)
	return project, nil
}

// GetAllProjects retrieves all projects
func (u *ProjectUseCase) GetAllProjects(ctx context.Context) ([]*domain.Project, error) {
	logger.Debug("getAllProjects: retrieving all projects")

	projects, err := u.projectRepo.GetAllProjects(ctx)
	if err != nil {
		logger.Error("getAllProjects: failed to retrieve projects", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve projects: " + err.Error())
	}

	logger.Info("getAllProjects: projects retrieved successfully", "count", len(projects))
	return projects, nil
}

// UpdateProject updates an existing project
func (u *ProjectUseCase) UpdateProject(ctx context.Context, id string, req *domain.UpdateProjectRequest) (*domain.Project, error) {
	logger.Debug("updateProject: validating request", "project_id", id, "name", req.Name)

	// Check if project exists
	existingProject, err := u.projectRepo.GetProject(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("updateProject: project not found", "project_id", id)
			return nil, apperrors.NotFound("project not found")
		}
		logger.Error("updateProject: failed to check existing project", "project_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to update project: " + err.Error())
	}

	// Update fields
	existingProject.Name = req.Name
	existingProject.Icon = req.Icon
	existingProject.Desc = req.Desc
	existingProject.Tech = req.Tech
	existingProject.URL = req.URL

	logger.Info("updateProject: updating project", "project_id", id, "name", req.Name)
	if err := u.projectRepo.UpdateProject(ctx, id, existingProject); err != nil {
		logger.Error("updateProject: failed to update project", "project_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to update project: " + err.Error())
	}

	logger.Info("updateProject: project updated successfully", "project_id", id, "name", req.Name)
	return existingProject, nil
}

// DeleteProject deletes a project
func (u *ProjectUseCase) DeleteProject(ctx context.Context, id string) error {
	logger.Debug("deleteProject: validating project existence", "project_id", id)

	// Check if project exists
	_, err := u.projectRepo.GetProject(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("deleteProject: project not found", "project_id", id)
			return apperrors.NotFound("project not found")
		}
		logger.Error("deleteProject: failed to check project", "project_id", id, "error", err.Error())
		return apperrors.InternalServerError("failed to delete project: " + err.Error())
	}

	logger.Info("deleteProject: deleting project", "project_id", id)
	if err := u.projectRepo.DeleteProject(ctx, id); err != nil {
		logger.Error("deleteProject: failed to delete project", "project_id", id, "error", err.Error())
		return apperrors.InternalServerError("failed to delete project: " + err.Error())
	}

	logger.Info("deleteProject: project deleted successfully", "project_id", id)
	return nil
}
