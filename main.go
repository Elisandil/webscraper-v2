package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"webscraper-v2/internal/config"
	"webscraper-v2/internal/infrastructure/database"
	repository "webscraper-v2/internal/infrastructure/persistence"
	"webscraper-v2/internal/infrastructure/web"
	"webscraper-v2/internal/usecase"
)

const (
	configFile = "config.yaml"
	appName    = "WebScraper"
	version    = "1.0"
)

func main() {
	log.Printf("Starting %s v%s", appName, version)
	// Cargar configuración
	cfg, err := loadConfig()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	// Inicializar base de datos
	db, err := initializeDatabase(cfg)

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	// Inicializar repositorios
	scrapingRepo := repository.NewScrapingRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Inicializar casos de uso
	scrapingUC := usecase.NewScrapingUseCase(scrapingRepo, cfg)
	authUC := usecase.NewAuthUseCase(userRepo, cfg)

	// Inicializar servidor web
	server := web.NewServer(cfg.Server.Port, cfg, scrapingUC, authUC)

	// Configurar manejo de señales para shutdown graceful
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		log.Printf("Server URL: http://localhost:%s", cfg.Server.Port)

		if err := server.Start(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()
	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	log.Printf("%s v%s shutdown complete", appName, version)
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("error loading config file '%s': %w", configFile, err)
	}

	log.Printf("Configuration loaded successfully")
	log.Printf("- Server port: %s", cfg.Server.Port)
	log.Printf("- Database path: %s", cfg.Database.Path)
	log.Printf("- User agent: %s", cfg.Scraping.UserAgent)
	log.Printf("- Request timeout: %ds", cfg.Scraping.Timeout)
	log.Printf("- Max links: %d", cfg.Scraping.MaxLinks)
	log.Printf("- Max images: %d", cfg.Scraping.MaxImages)
	log.Printf("- Analytics enabled: %t", cfg.Features.EnableAnalytics)
	log.Printf("- Caching enabled: %t", cfg.Features.EnableCaching)
	log.Printf("- Authentication required: %t", cfg.Auth.RequireAuth)
	log.Printf("- Default user role: %s", cfg.Auth.DefaultRole)
	log.Printf("- Token duration: %d hours", cfg.Auth.TokenDuration)

	return cfg, nil
}

func initializeDatabase(cfg *config.Config) (*database.SQLiteDB, error) {
	dataDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory '%s': %w", dataDir, err)
	}

	db, err := database.NewSQLiteDB(cfg.Database.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database: %w", err)
	}

	log.Printf("Database initialized successfully at: %s", cfg.Database.Path)
	return db, nil
}
