package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"taskflow-api/api/v1/dto"
	"taskflow-api/internal/task/application"
)

type TaskHandler struct {
	service *application.TaskService
}

func NewTaskHandler(service *application.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// GetTasksByProject GET /api/v1/projects/{id}/tasks
func (h *TaskHandler) GetTasksByProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	tasks, err := h.service.GetTasksByProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.TaskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toTaskResponse(t)
	}

	writeJSON(w, http.StatusOK, resp)
}

// CreateTask POST /api/v1/projects/{id}/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(r.Context(), application.CreateTaskDTO{
		Title:       req.Title,
		Description: req.Description,
		ProjectID:   projectID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toTaskResponse(task))
}

// MoveTask PUT /api/v1/tasks/{id}/move
func (h *TaskHandler) MoveTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	var req dto.MoveTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.MoveTask(r.Context(), application.MoveTaskDTO{
		TaskID:    taskID,
		NewStatus: req.Status,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(task))
}

func toTaskResponse(t *application.TaskDTO) dto.TaskResponse {
	return dto.TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		AssigneeID:  t.AssigneeID,
		ProjectID:   t.ProjectID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
