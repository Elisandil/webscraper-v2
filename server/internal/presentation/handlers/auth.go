package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/presentation/middleware"
	"webscraper-v2/internal/presentation/response"
	"webscraper-v2/internal/usecase"
)

const authCookieName = "auth_token"

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) setAuthCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	secure := os.Getenv("ENV") != "development"
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
}

func (h *AuthHandler) clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") != "development",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entity.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
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

	resp, err := h.authUseCase.Login(&req)
	if err != nil {
		log.Printf("Login error for user %s: %v", req.Username, err)
		response.SendErrorResponse(w, "Login failed", http.StatusUnauthorized, "Invalid credentials")
		return
	}
	log.Printf("User logged in successfully: %s", req.Username)
	h.setAuthCookie(w, resp.Token, resp.ExpiresAt)
	resp.Token = ""
	response.SendSuccessResponse(w, "Login successful", resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		response.SendErrorResponse(w, "No active session", http.StatusUnauthorized, "")
		return
	}

	resp, err := h.authUseCase.RefreshToken(cookie.Value)
	if err != nil {
		log.Printf("Token refresh error: %v", err)
		h.clearAuthCookie(w)
		response.SendErrorResponse(w, "Token refresh failed", http.StatusUnauthorized, err.Error())
		return
	}
	h.setAuthCookie(w, resp.Token, resp.ExpiresAt)
	resp.Token = ""
	response.SendSuccessResponse(w, "Token refreshed successfully", resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		response.SendErrorResponse(w, "No active session", http.StatusBadRequest, "")
		return
	}

	if err := h.authUseCase.RevokeToken(cookie.Value); err != nil {
		response.SendErrorResponse(w, "Failed to revoke token", http.StatusBadRequest, err.Error())
		return
	}
	h.clearAuthCookie(w)
	response.SendSuccessResponse(w, "Logged out successfully", nil)
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if user == nil {
		response.SendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}
	response.SendSuccessResponse(w, "Profile retrieved successfully", user)
}
