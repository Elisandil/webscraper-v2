package routes

import (
	"net/http"
	"webscraper-v2/internal/infrastructure/config"
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
	chatHandler     *handlers.ChatHandler
	commonHandler   *handlers.CommonHandler
	strictLimiter   *middleware.RateLimiter
	moderateLimiter *middleware.RateLimiter
	generalLimiter  *middleware.RateLimiter
	publicLimiter   *middleware.RateLimiter
}

func NewRouter(
	config *config.Config,
	jwtMiddleware *middleware.JWTMiddleware,
	authHandler *handlers.AuthHandler,
	scrapingHandler *handlers.ScrapingHandler,
	scheduleHandler *handlers.ScheduleHandler,
	chatHandler *handlers.ChatHandler,
	commonHandler *handlers.CommonHandler,
) *Router {
	return &Router{
		router:          mux.NewRouter(),
		config:          config,
		jwtMiddleware:   jwtMiddleware,
		authHandler:     authHandler,
		scrapingHandler: scrapingHandler,
		scheduleHandler: scheduleHandler,
		chatHandler:     chatHandler,
		commonHandler:   commonHandler,
		strictLimiter:   middleware.NewStrictRateLimiter(),
		moderateLimiter: middleware.NewModerateRateLimiter(),
		generalLimiter:  middleware.NewGeneralRateLimiter(),
		publicLimiter:   middleware.NewPublicRateLimiter(),
	}
}

func (rt *Router) Shutdown() {
	rt.strictLimiter.Shutdown()
	rt.moderateLimiter.Shutdown()
	rt.generalLimiter.Shutdown()
	rt.publicLimiter.Shutdown()
}

func (rt *Router) SetupRoutes() *mux.Router {
	rt.router.Use(
		middleware.LoggingMiddleware,
		middleware.CORSMiddleware,
		middleware.ContentTypeMiddleware,
		rt.generalLimiter.Limit,
	)

	rt.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))),
	)
	rt.setupAuthRoutes()

	rt.setupPublicRoutes()

	rt.setupAPIRoutes()

	rt.setupAdminRoutes()

	return rt.router
}

func (rt *Router) setupAuthRoutes() {
	auth := rt.router.PathPrefix("/api/auth").Subrouter()
	auth.Use(rt.strictLimiter.Limit)

	auth.HandleFunc("/register", rt.authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", rt.authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", rt.authHandler.RefreshToken).Methods("POST")
	auth.HandleFunc("/logout", rt.authHandler.Logout).Methods("POST")
}

func (rt *Router) setupPublicRoutes() {
	public := rt.router.PathPrefix("/api/public").Subrouter()
	public.Use(rt.publicLimiter.Limit)

	public.HandleFunc("/scrape", rt.scrapingHandler.PublicScrape).Methods("POST")
}

func (rt *Router) setupAPIRoutes() {
	api := rt.router.PathPrefix("/api").Subrouter()
	api.Use(rt.jwtMiddleware.RequireAuth)

	scraping := api.PathPrefix("/scrape").Subrouter()
	scraping.Use(rt.moderateLimiter.Limit)
	scraping.HandleFunc("", rt.scrapingHandler.Scrape).Methods("POST")

	api.HandleFunc("/results", rt.scrapingHandler.GetResults).Methods("GET")
	api.HandleFunc("/results/{id:[0-9]+}", rt.scrapingHandler.GetResult).Methods("GET")
	api.HandleFunc("/schedules", rt.scheduleHandler.Create).Methods("POST")
	api.HandleFunc("/schedules", rt.scheduleHandler.GetAll).Methods("GET")
	api.HandleFunc("/schedules/{id:[0-9]+}", rt.scheduleHandler.GetByID).Methods("GET")
	api.HandleFunc("/schedules/{id:[0-9]+}", rt.scheduleHandler.Update).Methods("PUT")
	api.HandleFunc("/schedules/{id:[0-9]+}", rt.scheduleHandler.Delete).Methods("DELETE")

	api.HandleFunc("/chat/parse", rt.chatHandler.ParseMessage).Methods("POST")
	api.HandleFunc("/chat/execute", rt.chatHandler.ExecuteAction).Methods("POST")

	api.HandleFunc("/profile", rt.authHandler.Profile).Methods("GET")
	api.HandleFunc("/health", rt.commonHandler.Health).Methods("GET")
	api.HandleFunc("/profile", rt.authHandler.Profile).Methods("GET")
	api.HandleFunc("/health", rt.commonHandler.Health).Methods("GET")

	api.PathPrefix("/").HandlerFunc(rt.commonHandler.NotFound)
}

func (rt *Router) setupAdminRoutes() {
	adminAPI := rt.router.PathPrefix("/api/admin").Subrouter()
	adminAPI.Use(rt.jwtMiddleware.RequireAuth)
	adminAPI.Use(rt.jwtMiddleware.RequireRole("admin"))
}
