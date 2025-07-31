package server

import (
	"log"
	"net/http"
	"strings"
	"webscraper-v2/internal/config"
	"webscraper-v2/internal/presentation/handlers"
	"webscraper-v2/internal/presentation/middleware"
	"webscraper-v2/internal/presentation/routes"
	"webscraper-v2/internal/usecase"

	"github.com/gorilla/mux"
)

type Server struct {
	port       string
	config     *config.Config
	router     *mux.Router
	scheduleUC *usecase.ScheduleUseCase
}

func NewServer(
	port string,
	cfg *config.Config,
	scrapingUC *usecase.ScrapingUseCase,
	authUC *usecase.AuthUseCase,
	scheduleUC *usecase.ScheduleUseCase,
) *Server {
	// Initialize middleware
	jwtMiddleware := middleware.NewJWTMiddleware(authUC)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUC)
	scrapingHandler := handlers.NewScrapingHandler(scrapingUC)
	scheduleHandler := handlers.NewScheduleHandler(scheduleUC)
	commonHandler := handlers.NewCommonHandler(cfg)

	// Initialize router
	routerManager := routes.NewRouter(
		cfg,
		jwtMiddleware,
		authHandler,
		scrapingHandler,
		scheduleHandler,
		commonHandler,
	)

	return &Server{
		port:       port,
		config:     cfg,
		router:     routerManager.SetupRoutes(),
		scheduleUC: scheduleUC,
	}
}

func (s *Server) Start() error {
	// Start the scheduler
	s.scheduleUC.StartScheduler()

	// Log available endpoints
	s.logEndpoints()

	log.Printf("Server listening on port %s", s.port)
	log.Printf("Server URL: http://localhost:%s", s.port)

	return http.ListenAndServe(":"+s.port, s.router)
}

func (s *Server) logEndpoints() {
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
		"POST /api/schedules - Create schedule",
		"GET  /api/schedules - Get user schedules",
		"GET  /api/schedules/{id} - Get specific schedule",
		"PUT  /api/schedules/{id} - Update schedule",
		"DELETE /api/schedules/{id} - Delete schedule",
		"GET  /api/admin/users - Get all users (admin only)",
		"GET  /api/health - Health check",
	}

	authStatus := "disabled"
	if s.config.Auth.RequireAuth {
		authStatus = "required"
	}

	log.Printf("Server configuration (Auth: %s)\nAvailable endpoints:\n  %s",
		authStatus, strings.Join(endpoints, "\n  "))
}
