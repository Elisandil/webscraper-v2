package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/presentation/middleware"
	"webscraper-v2/internal/presentation/response"
	"webscraper-v2/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entity.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Username) == "" {
		response.SendErrorResponse(w, "Username is required", http.StatusBadRequest, "")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		response.SendErrorResponse(w, "Email is required", http.StatusBadRequest, "")
		return
	}
	if len(req.Password) < 6 {
		response.SendErrorResponse(w, "Password must be at least 6 characters", http.StatusBadRequest, "")
		return
	}
	resp, err := h.authUseCase.Register(&req)

	if err != nil {
		log.Printf("Registration error for user %s: %v", req.Username, err)
		response.SendErrorResponse(w, "Registration failed", http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("User registered successfully: %s", req.Username)
	response.SendSuccessResponse(w, "User registered successfully", resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entity.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		response.SendErrorResponse(w, "Username and password are required", http.StatusBadRequest, "")
		return
	}
	resp, err := h.authUseCase.Login(&req)

	if err != nil {
		log.Printf("Login error for user %s: %v", req.Username, err)
		response.SendErrorResponse(w, "Login failed", http.StatusUnauthorized, "Invalid credentials")
		return
	}
	log.Printf("User logged in successfully: %s", req.Username)
	response.SendSuccessResponse(w, "Login successful", resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(req.Token) == "" {
		response.SendErrorResponse(w, "Token is required", http.StatusBadRequest, "")
		return
	}
	resp, err := h.authUseCase.RefreshToken(req.Token)

	if err != nil {
		log.Printf("Token refresh error: %v", err)
		response.SendErrorResponse(w, "Token refresh failed", http.StatusUnauthorized, err.Error())
		return
	}
	response.SendSuccessResponse(w, "Token refreshed successfully", resp)
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	response.SendSuccessResponse(w, "Profile retrieved successfully", user)
}
