package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"taskflow-api/api/v1/dto"
	"taskflow-api/internal/project/application"
)

type ProjectHandler struct {
	service *application.ProjectService
}

func NewProjectHandler(service *application.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// GetAllProjects GET /api/v1/projects
func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.ProjectResponse, len(projects))
	for i, p := range projects {
		resp[i] = toProjectResponse(p)
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetProject GET /api/v1/projects/{id}
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	project, err := h.service.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, toProjectResponse(project))
}

// CreateProject POST /api/v1/projects
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Phase 1 : utilisateur simulé via header X-User-Id
	ownerID := r.Header.Get("X-User-Id")
	if ownerID == "" {
		ownerID = "default-user"
	}

	project, err := h.service.CreateProject(r.Context(), application.CreateProjectDTO{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toProjectResponse(project))
}

// AddMember POST /api/v1/projects/{id}/members
func (h *ProjectHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	project, err := h.service.AddMember(r.Context(), application.AddMemberDTO{
		ProjectID: id,
		UserID:    req.UserID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	writeJSON(w, http.StatusCreated, toProjectResponse(project))
}

func toProjectResponse(p *application.ProjectDTO) dto.ProjectResponse {
	members := make([]dto.MemberResponse, len(p.Members))
	for i, m := range p.Members {
		members[i] = dto.MemberResponse{
			UserID:   m.UserID,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		}
	}
	return dto.ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Members:     members,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
