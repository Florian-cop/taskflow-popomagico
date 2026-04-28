package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	v1dto "taskflow-api/api/v1/dto"
	"taskflow-api/internal/project/application"
	sharedDomain "taskflow-api/internal/shared/domain"
)

type ProjectHandler struct {
	service *application.ProjectService
}

func NewProjectHandler(s *application.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// GET /api/v2/projects
func (h *ProjectHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]v1dto.ProjectResponse, len(projects))
	for i, p := range projects {
		resp[i] = toProjectResponse(p)
	}
	writeListEnveloped(w, http.StatusOK, resp, len(resp))
}

// GET /api/v2/projects/{id}
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	project, err := h.service.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeEnveloped(w, http.StatusOK, toProjectResponse(project))
}

// POST /api/v2/projects
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req v1dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	ownerID, ok := sharedDomain.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	project, err := h.service.CreateProject(r.Context(), application.CreateProjectDTO{
		Name: req.Name, Description: req.Description, OwnerID: ownerID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeEnveloped(w, http.StatusCreated, toProjectResponse(project))
}

// POST /api/v2/projects/{id}/members
func (h *ProjectHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req v1dto.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	project, err := h.service.AddMember(r.Context(), application.AddMemberDTO{
		ProjectID: id, UserID: req.UserID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	writeEnveloped(w, http.StatusCreated, toProjectResponse(project))
}

func toProjectResponse(p *application.ProjectDTO) v1dto.ProjectResponse {
	members := make([]v1dto.MemberResponse, len(p.Members))
	for i, m := range p.Members {
		members[i] = v1dto.MemberResponse{
			UserID:   m.UserID,
			Role:     m.Role,
			JoinedAt: roundTime(m.JoinedAt),
		}
	}
	return v1dto.ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Members:     members,
		CreatedAt:   roundTime(p.CreatedAt),
		UpdatedAt:   roundTime(p.UpdatedAt),
	}
}

func roundTime(t time.Time) time.Time {
	return t.UTC()
}
