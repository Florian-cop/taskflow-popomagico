package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	v1dto "taskflow-api/api/v1/dto"
	"taskflow-api/internal/task/application"
)

type TaskHandler struct {
	service *application.TaskService
}

func NewTaskHandler(s *application.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// GET /api/v2/projects/{id}/tasks
func (h *TaskHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tasks, err := h.service.GetTasksByProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]v1dto.TaskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toTaskResponse(t)
	}
	writeListEnveloped(w, http.StatusOK, resp, len(resp))
}

// POST /api/v2/projects/{id}/tasks
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	var req v1dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	task, err := h.service.CreateTask(r.Context(), application.CreateTaskDTO{
		Title: req.Title, Description: req.Description, ProjectID: projectID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeEnveloped(w, http.StatusCreated, toTaskResponse(task))
}

// PUT /api/v2/tasks/{id}/move
func (h *TaskHandler) Move(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	var req v1dto.MoveTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	task, err := h.service.MoveTask(r.Context(), application.MoveTaskDTO{
		TaskID: taskID, NewStatus: req.Status,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	writeEnveloped(w, http.StatusOK, toTaskResponse(task))
}

func toTaskResponse(t *application.TaskDTO) v1dto.TaskResponse {
	return v1dto.TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		AssigneeID:  t.AssigneeID,
		ProjectID:   t.ProjectID,
		CreatedAt:   t.CreatedAt.UTC(),
		UpdatedAt:   t.UpdatedAt.UTC(),
	}
}
