package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/usecase"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type JWTMiddleware struct {
	authUseCase *usecase.AuthUseCase
}

func NewJWTMiddleware(authUseCase *usecase.AuthUseCase) *JWTMiddleware {
	return &JWTMiddleware{authUseCase: authUseCase}
}

func (m *JWTMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			m.sendUnauthorized(w, "Authorization header required")
			return
		}
		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			m.sendUnauthorized(w, "Invalid authorization header format")
			return
		}
		tokenString := tokenParts[1]
		claims, err := m.authUseCase.ValidateToken(tokenString)

		if err != nil {
			m.sendUnauthorized(w, "Invalid or expired token")
			return
		}
		user, err := m.authUseCase.GetUserByID(claims.UserID)

		if err != nil || user == nil {
			m.sendUnauthorized(w, "User not found")
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				m.sendForbidden(w, "User not found in context")
				return
			}
			hasRole := false

			for _, requiredRole := range roles {

				if user.Role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				m.sendForbidden(w, "Insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (m *JWTMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := tokenParts[1]
		claims, err := m.authUseCase.ValidateToken(tokenString)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := m.authUseCase.GetUserByID(claims.UserID)

		if err != nil || user == nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) sendUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	response := ErrorResponse{
		Error:   "Unauthorized",
		Message: message,
		Code:    http.StatusUnauthorized,
	}
	sendJSONResponse(w, response)
}

func (m *JWTMiddleware) sendForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	response := ErrorResponse{
		Error:   "Forbidden",
		Message: message,
		Code:    http.StatusForbidden,
	}
	sendJSONResponse(w, response)
}

func GetUserFromContext(ctx context.Context) *entity.User {

	if user, ok := ctx.Value(UserContextKey).(*entity.User); ok {
		return user
	}
	return nil
}

func sendJSONResponse(w http.ResponseWriter, response interface{}) {

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
