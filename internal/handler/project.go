package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	projectusecase "github.com/quoteyouros/backend/internal/usecase/project"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type ProjectHandler struct {
	usecase *projectusecase.ProjectUseCase
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(usecase *projectusecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{usecase: usecase}
}

// GetAllProjects retrieves all projects
// GET /api/projects
func (h *ProjectHandler) GetAllProjects(c *fiber.Ctx) error {
	logger.Debug("getAllProjects: retrieving all projects")
	projects, err := h.usecase.GetAllProjects(c.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getAllProjects: failed to retrieve projects", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getAllProjects: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve projects", fiber.StatusInternalServerError)
	}

	logger.Info("getAllProjects: projects retrieved successfully", "count", len(projects))
	return response.SuccessResponse(c, fiber.StatusOK, projects, "Projects retrieved successfully")
}

// GetProject retrieves a single project by ID
// GET /api/projects/:id
func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("getProject: missing project ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "project ID is required", fiber.StatusBadRequest)
	}

	logger.Debug("getProject: retrieving project", "project_id", id)
	project, err := h.usecase.GetProject(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getProject: failed to retrieve project", "project_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getProject: unexpected error", "project_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve project", fiber.StatusInternalServerError)
	}

	logger.Debug("getProject: project retrieved successfully", "project_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, project, "Project retrieved successfully")
}

// CreateProject creates a new project
// POST /api/projects
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var req domain.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("createProject: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("createProject: attempting to create project", "name", req.Name)
	project, err := h.usecase.CreateProject(c.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("createProject: failed to create project", "name", req.Name, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("createProject: unexpected error", "name", req.Name, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to create project", fiber.StatusInternalServerError)
	}

	logger.Info("createProject: project created successfully", "project_id", project.ID, "name", project.Name)
	return response.SuccessResponse(c, fiber.StatusCreated, project, "Project created successfully")
}

// UpdateProject updates an existing project
// PUT /api/projects/:id
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("updateProject: missing project ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "project ID is required", fiber.StatusBadRequest)
	}

	var req domain.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("updateProject: failed to parse request body", "project_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("updateProject: attempting to update project", "project_id", id, "name", req.Name)
	project, err := h.usecase.UpdateProject(c.Context(), id, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("updateProject: failed to update project", "project_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("updateProject: unexpected error", "project_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to update project", fiber.StatusInternalServerError)
	}

	logger.Info("updateProject: project updated successfully", "project_id", id, "name", project.Name)
	return response.SuccessResponse(c, fiber.StatusOK, project, "Project updated successfully")
}

// DeleteProject deletes a project
// DELETE /api/projects/:id
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("deleteProject: missing project ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "project ID is required", fiber.StatusBadRequest)
	}

	logger.Info("deleteProject: attempting to delete project", "project_id", id)
	err := h.usecase.DeleteProject(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("deleteProject: failed to delete project", "project_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("deleteProject: unexpected error", "project_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to delete project", fiber.StatusInternalServerError)
	}

	logger.Info("deleteProject: project deleted successfully", "project_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{}, "Project deleted successfully")
}
