package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"webscraper/config"
	"webscraper/domain/entity"
	"webscraper/usecase"

	"github.com/gorilla/mux"
)

type Server struct {
	port        string
	config      *config.Config
	usecase     *usecase.ScrapingUseCase
	authUseCase *usecase.AuthUseCase
	router      *mux.Router
	jwtMw       *JWTMiddleware
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewServer(port string, cfg *config.Config, uc *usecase.ScrapingUseCase, authUC *usecase.AuthUseCase) *Server {
	s := &Server{
		port:        port,
		config:      cfg,
		usecase:     uc,
		authUseCase: authUC,
		router:      mux.NewRouter(),
		jwtMw:       NewJWTMiddleware(authUC),
	}
	s.setupRoutes()
	s.setupMiddleware()
	return s
}

func (s *Server) setupRoutes() {
	// Archivos estáticos
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Rutas de autenticación (públicas)
	auth := s.router.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/register", s.registerHandler).Methods("POST")
	auth.HandleFunc("/login", s.loginHandler).Methods("POST")
	auth.HandleFunc("/refresh", s.refreshTokenHandler).Methods("POST")

	// API general
	api := s.router.PathPrefix("/api").Subrouter()

	// Rutas públicas o con autenticación opcional
	if s.config.Auth.RequireAuth {
		// Si RequireAuth está habilitado, todas las rutas requieren autenticación
		api.Use(s.jwtMw.RequireAuth)
	} else {
		// Si no, usar autenticación opcional
		api.Use(s.jwtMw.OptionalAuth)
	}

	api.HandleFunc("/scrape", s.scrapeHandler).Methods("POST")
	api.HandleFunc("/results", s.resultsHandler).Methods("GET")
	api.HandleFunc("/results/{id:[0-9]+}", s.resultHandler).Methods("GET")
	api.HandleFunc("/results/{id:[0-9]+}", s.deleteResultHandler).Methods("DELETE")
	api.HandleFunc("/health", s.healthHandler).Methods("GET")
	api.HandleFunc("/profile", s.profileHandler).Methods("GET")

	// Rutas que requieren roles específicos
	adminAPI := s.router.PathPrefix("/api/admin").Subrouter()
	adminAPI.Use(s.jwtMw.RequireAuth)
	adminAPI.Use(s.jwtMw.RequireRole("admin"))
	adminAPI.HandleFunc("/users", s.getUsersHandler).Methods("GET")

	// Manejar rutas no encontradas en API
	api.PathPrefix("/").HandlerFunc(s.notFoundHandler)

	// Página principal
	s.router.HandleFunc("/", s.indexHandler).Methods("GET")
}

func (s *Server) setupMiddleware() {
	s.router.Use(s.loggingMiddleware, s.corsMiddleware, s.contentTypeMiddleware)
}

func (s *Server) Start() error {
	endpoints := []string{
		"GET  / - Web interface",
		"POST /api/auth/register - Register user",
		"POST /api/auth/login - Login user",
		"POST /api/auth/refresh - Refresh token",
		"GET  /api/profile - Get user profile",
		"POST /api/scrape - Scrape URL",
		"GET  /api/results - Get all results",
		"GET  /api/results/{id} - Get specific result",
		"DELETE /api/results/{id} - Delete result",
		"GET  /api/admin/users - Get all users (admin only)",
		"GET  /api/health - Health check",
	}
	authStatus := "disabled"

	if s.config.Auth.RequireAuth {
		authStatus = "required"
	}
	log.Printf("Server listening on port %s (Auth: %s)\nAvailable endpoints:\n  %s",
		s.port, authStatus, strings.Join(endpoints, "\n  "))
	return http.ListenAndServe(":"+s.port, s.router)
}

// Handlers de autenticación
func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(req.Username) == "" {
		s.sendErrorResponse(w, "Username is required", http.StatusBadRequest, "")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		s.sendErrorResponse(w, "Email is required", http.StatusBadRequest, "")
		return
	}
	if len(req.Password) < 6 {
		s.sendErrorResponse(w, "Password must be at least 6 characters", http.StatusBadRequest, "")
		return
	}
	response, err := s.authUseCase.Register(&req)

	if err != nil {
		log.Printf("Registration error for user %s: %v", req.Username, err)
		s.sendErrorResponse(w, "Registration failed", http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("User registered successfully: %s", req.Username)
	s.sendSuccessResponse(w, "User registered successfully", response)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		s.sendErrorResponse(w, "Username and password are required", http.StatusBadRequest, "")
		return
	}
	response, err := s.authUseCase.Login(&req)

	if err != nil {
		log.Printf("Login error for user %s: %v", req.Username, err)
		s.sendErrorResponse(w, "Login failed", http.StatusUnauthorized, "Invalid credentials")
		return
	}
	log.Printf("User logged in successfully: %s", req.Username)
	s.sendSuccessResponse(w, "Login successful", response)
}

func (s *Server) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Token) == "" {
		s.sendErrorResponse(w, "Token is required", http.StatusBadRequest, "")
		return
	}
	response, err := s.authUseCase.RefreshToken(req.Token)

	if err != nil {
		log.Printf("Token refresh error: %v", err)
		s.sendErrorResponse(w, "Token refresh failed", http.StatusUnauthorized, err.Error())
		return
	}
	s.sendSuccessResponse(w, "Token refreshed successfully", response)
}

func (s *Server) profileHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())

	if user == nil {
		s.sendErrorResponse(w, "Authentication required", http.StatusUnauthorized, "")
		return
	}

	s.sendSuccessResponse(w, "Profile retrieved successfully", user)
}

func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Este es un endpoint de ejemplo para administradores
	// Aquí podrías implementar la lógica para obtener todos los usuarios
	s.sendSuccessResponse(w, "Admin endpoint - users list", map[string]string{
		"message": "This would return all users",
		"note":    "Implement user listing logic here",
	})
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./interface/templates/index.html")
}

func (s *Server) scrapeHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user != nil {
		log.Printf("Scraping request from user: %s", user.Username)
	}
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest, err.Error())
		return
	}
	if req.URL = strings.TrimSpace(req.URL); req.URL == "" {
		s.sendErrorResponse(w, "URL is required", http.StatusBadRequest, "")
		return
	}
	parsedURL, err := url.ParseRequestURI(req.URL)

	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		s.sendErrorResponse(w, "Invalid URL format", http.StatusBadRequest, "URL must include protocol (http:// or https://) and valid domain")
		return
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		s.sendErrorResponse(w, "Invalid URL scheme", http.StatusBadRequest, "Only HTTP and HTTPS protocols are supported")
		return
	}
	log.Printf("Scraping URL: %s", req.URL)
	result, err := s.usecase.ScrapeURL(req.URL)

	if err != nil {
		log.Printf("Error scraping URL %s: %v", req.URL, err)
		s.sendErrorResponse(w, "Failed to scrape URL", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Successfully scraped URL: %s (Status: %d, Words: %d)", req.URL, result.StatusCode, result.WordCount)
	s.sendSuccessResponse(w, "URL scraped successfully", result)
}

func (s *Server) resultsHandler(w http.ResponseWriter, r *http.Request) {
	results, err := s.usecase.GetAllResults()

	if err != nil {
		log.Printf("Error getting results: %v", err)
		s.sendErrorResponse(w, "Failed to retrieve results", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Retrieved %d scraping results", len(results))
	s.sendSuccessResponse(w, fmt.Sprintf("Retrieved %d results", len(results)), results)
}

func (s *Server) resultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := s.parseID(r)

	if err != nil {
		s.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}
	result, err := s.usecase.GetResult(id)

	if err != nil {
		log.Printf("Error getting result %d: %v", id, err)
		s.sendErrorResponse(w, "Failed to retrieve result", http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		s.sendErrorResponse(w, "Result not found", http.StatusNotFound, fmt.Sprintf("No result found with ID %d", id))
		return
	}
	log.Printf("Retrieved result ID: %d (%s)", id, result.URL)
	s.sendSuccessResponse(w, "Result retrieved successfully", result)
}

func (s *Server) deleteResultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := s.parseID(r)

	if err != nil {
		s.sendErrorResponse(w, "Invalid ID format", http.StatusBadRequest, "ID must be a valid number")
		return
	}
	result, err := s.usecase.GetResult(id)

	if err != nil {
		log.Printf("Error checking result %d: %v", id, err)
		s.sendErrorResponse(w, "Failed to check result", http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		s.sendErrorResponse(w, "Result not found", http.StatusNotFound, fmt.Sprintf("No result found with ID %d", id))
		return
	}
	if err := s.usecase.DeleteResult(id); err != nil {
		log.Printf("Error deleting result %d: %v", id, err)
		s.sendErrorResponse(w, "Failed to delete result", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Deleted result ID: %d (%s)", id, result.URL)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "ok", "timestamp": time.Now().UTC().Format(time.RFC3339),
		"service": "webscraper", "version": "2.0",
		"auth_enabled": s.config.Auth.RequireAuth,
	}
	s.sendSuccessResponse(w, "Service is healthy", health)
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.sendErrorResponse(w, "Endpoint not found", http.StatusNotFound,
		fmt.Sprintf("The requested endpoint %s %s does not exist", r.Method, r.URL.Path))
}

func (s *Server) parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
}

func (s *Server) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, details string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Message: details, Code: statusCode})
}

func (s *Server) sendSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	json.NewEncoder(w).Encode(SuccessResponse{Message: message, Data: data})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)
		log.Printf("%s %s - %d - %v - %s", r.Method, r.URL.Path, ww.statusCode, time.Since(start), r.RemoteAddr)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
