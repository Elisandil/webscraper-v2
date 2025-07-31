package routes

import (
	"net/http"
	"webscraper-v2/internal/config"
	"webscraper-v2/internal/presentation/handlers"
	"webscraper-v2/internal/presentation/middleware"

	"github.com/gorilla/mux"
)

type Router struct {
	router          *mux.Router
	config          *config.Config
	jwtMiddleware   *middleware.JWTMiddleware
	authHandler     *handlers.AuthHandler
	scrapingHandler *handlers.ScrapingHandler
	scheduleHandler *handlers.ScheduleHandler
	commonHandler   *handlers.CommonHandler
}

func NewRouter(
	config *config.Config,
	jwtMiddleware *middleware.JWTMiddleware,
	authHandler *handlers.AuthHandler,
	scrapingHandler *handlers.ScrapingHandler,
	scheduleHandler *handlers.ScheduleHandler,
	commonHandler *handlers.CommonHandler,
) *Router {
	return &Router{
		router:          mux.NewRouter(),
		config:          config,
		jwtMiddleware:   jwtMiddleware,
		authHandler:     authHandler,
		scrapingHandler: scrapingHandler,
		scheduleHandler: scheduleHandler,
		commonHandler:   commonHandler,
	}
}

func (rt *Router) SetupRoutes() *mux.Router {
	// Apply global middleware
	rt.router.Use(
		middleware.LoggingMiddleware,
		middleware.CORSMiddleware,
		middleware.ContentTypeMiddleware,
	)

	// Static files
	rt.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))),
	)

	// Authentication routes (public)
	rt.setupAuthRoutes()

	// API routes
	rt.setupAPIRoutes()

	// Admin routes
	rt.setupAdminRoutes()

	// Main page
	rt.router.HandleFunc("/", rt.commonHandler.Index).Methods("GET")

	return rt.router
}

func (rt *Router) setupAuthRoutes() {
	auth := rt.router.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/register", rt.authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", rt.authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", rt.authHandler.RefreshToken).Methods("POST")
}

func (rt *Router) setupAPIRoutes() {
	api := rt.router.PathPrefix("/api").Subrouter()

	if rt.config.Auth.RequireAuth {
		api.Use(rt.jwtMiddleware.RequireAuth)
	} else {
		api.Use(rt.jwtMiddleware.OptionalAuth)
	}

	// Scraping routes
	api.HandleFunc("/scrape", rt.scrapingHandler.Scrape).Methods("POST")
	api.HandleFunc("/results", rt.scrapingHandler.GetResults).Methods("GET")
	api.HandleFunc("/results/{id:[0-9]+}", rt.scrapingHandler.GetResult).Methods("GET")
	api.HandleFunc("/results/{id:[0-9]+}", rt.scrapingHandler.DeleteResult).Methods("DELETE")

	// Schedule routes
	api.HandleFunc("/schedules", rt.scheduleHandler.Create).Methods("POST")
	api.HandleFunc("/schedules", rt.scheduleHandler.GetAll).Methods("GET")
	api.HandleFunc("/schedules/{id:[0-9]+}", rt.scheduleHandler.GetByID).Methods("GET")
	api.HandleFunc("/schedules/{id:[0-9]+}", rt.scheduleHandler.Update).Methods("PUT")
	api.HandleFunc("/schedules/{id:[0-9]+}", rt.scheduleHandler.Delete).Methods("DELETE")

	// Common routes
	api.HandleFunc("/profile", rt.authHandler.Profile).Methods("GET")
	api.HandleFunc("/health", rt.commonHandler.Health).Methods("GET")

	// Handle not found routes in API
	api.PathPrefix("/").HandlerFunc(rt.commonHandler.NotFound)
}

func (rt *Router) setupAdminRoutes() {
	adminAPI := rt.router.PathPrefix("/api/admin").Subrouter()
	adminAPI.Use(rt.jwtMiddleware.RequireAuth)
	adminAPI.Use(rt.jwtMiddleware.RequireRole("admin"))
}
