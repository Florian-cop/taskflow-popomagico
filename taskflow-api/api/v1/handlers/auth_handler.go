package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"taskflow-api/api/v1/dto"
	sharedDomain "taskflow-api/internal/shared/domain"
	"taskflow-api/internal/user/application"
)

type AuthHandler struct {
	service *application.UserService
}

func NewAuthHandler(service *application.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.Register(r.Context(), application.RegisterUserDTO{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		if errors.Is(err, sharedDomain.ErrConflict) {
			http.Error(w, "email already used", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toAuthResponse(result))
}

// Login POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Login(r.Context(), application.LoginDTO{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, application.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, toAuthResponse(result))
}

// Me GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := sharedDomain.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

// LookupByEmail GET /api/v1/users/by-email?email=...
// Permet à l'UI de résoudre email → user pour l'ajout de membres.
// Note sécurité : expose l'existence d'un compte ; acceptable pour ce pilote.
func (h *AuthHandler) LookupByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email query param required", http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(user))
}

// SearchUsers GET /api/v1/users?search=ab&limit=10
// Retourne les utilisateurs dont l'email contient la sous-chaîne (pour autocomplete).
// Exige au moins 2 caractères pour éviter un dump complet.
func (h *AuthHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	if len(query) < 2 {
		writeJSON(w, http.StatusOK, []dto.UserResponse{})
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	users, err := h.service.SearchUsersByEmail(r.Context(), query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]dto.UserResponse, len(users))
	for i, u := range users {
		resp[i] = toUserResponse(u)
	}
	writeJSON(w, http.StatusOK, resp)
}

func toAuthResponse(r *application.AuthResultDTO) dto.AuthResponse {
	return dto.AuthResponse{
		Token: r.Token,
		User:  toUserResponse(&r.User),
	}
}

func toUserResponse(u *application.UserDTO) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
