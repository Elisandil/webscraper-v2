package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"webscraper-v2/internal/infrastructure/config"
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
	routerMgr  *routes.Router
	scheduleUC *usecase.ScheduleUseCase
	httpServer *http.Server
}

func NewServer(
	port string,
	cfg *config.Config,
	scrapingUC *usecase.ScrapingUseCase,
	authUC *usecase.AuthUseCase,
	scheduleUC *usecase.ScheduleUseCase,
	chatUC *usecase.ChatUseCase,
) *Server {
	jwtMiddleware := middleware.NewJWTMiddleware(authUC)

	authHandler := handlers.NewAuthHandler(authUC)
	scrapingHandler := handlers.NewScrapingHandler(scrapingUC)
	scheduleHandler := handlers.NewScheduleHandler(scheduleUC)
	chatHandler := handlers.NewChatHandler(chatUC, scrapingUC, scheduleUC)
	commonHandler := handlers.NewCommonHandler(cfg)

	routerManager := routes.NewRouter(
		cfg,
		jwtMiddleware,
		authHandler,
		scrapingHandler,
		scheduleHandler,
		chatHandler,
		commonHandler,
	)

	return &Server{
		port:       port,
		config:     cfg,
		router:     routerManager.SetupRoutes(),
		routerMgr:  routerManager,
		scheduleUC: scheduleUC,
	}
}

func (s *Server) Start() error {
	s.scheduleUC.StartScheduler()

	s.logEndpoints()

	s.httpServer = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	log.Printf("Server listening on port %s", s.port)
	log.Printf("Server URL: http://localhost:%s", s.port)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("🛑 Shutting down HTTP server...")
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
	}
	log.Println("✅ HTTP server stopped")

	s.routerMgr.Shutdown()
	log.Println("✅ Rate limiters stopped")

	return nil
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
