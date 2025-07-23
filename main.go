package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
<<<<<<< HEAD
	"webscraper/config"
	"webscraper/database"
	"webscraper/repository"
	"webscraper/web"
	"webscraper/usecase"
=======
	"webscraper/internal/config"
	"webscraper/internal/infrastructure/database"
	repository "webscraper/internal/infrastructure/persistence"
	"webscraper/internal/infrastructure/web"
	"webscraper/internal/usecase"
>>>>>>> master
)

const (
	configFile = "config.yaml"
	appName    = "WebScraper"
	version    = "1.0"
)

func main() {
	log.Printf("Starting %s v%s", appName, version)
<<<<<<< HEAD
=======
	// Cargar configuración
>>>>>>> master
	cfg, err := loadConfig()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
<<<<<<< HEAD
=======
	// Inicializar base de datos
>>>>>>> master
	db, err := initializeDatabase(cfg)

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
<<<<<<< HEAD
	repo := repository.NewScrapingRepository(db)
	uc := usecase.NewScrapingUseCase(repo, cfg)
	server := web.NewServer(cfg.Server.Port, uc)
=======
	// Inicializar repositorios
	scrapingRepo := repository.NewScrapingRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Inicializar casos de uso
	scrapingUC := usecase.NewScrapingUseCase(scrapingRepo, cfg)
	authUC := usecase.NewAuthUseCase(userRepo, cfg)

	// Inicializar servidor web
	server := web.NewServer(cfg.Server.Port, cfg, scrapingUC, authUC)

	// Configurar manejo de señales para shutdown graceful
>>>>>>> master
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
<<<<<<< HEAD

	if err != nil {
		return nil, fmt.Errorf("error loading config file '%s': %w", configFile, err)
	}
=======
	if err != nil {
		return nil, fmt.Errorf("error loading config file '%s': %w", configFile, err)
	}

>>>>>>> master
	log.Printf("Configuration loaded successfully")
	log.Printf("- Server port: %s", cfg.Server.Port)
	log.Printf("- Database path: %s", cfg.Database.Path)
	log.Printf("- User agent: %s", cfg.Scraping.UserAgent)
	log.Printf("- Request timeout: %ds", cfg.Scraping.Timeout)
	log.Printf("- Max links: %d", cfg.Scraping.MaxLinks)
	log.Printf("- Max images: %d", cfg.Scraping.MaxImages)
	log.Printf("- Analytics enabled: %t", cfg.Features.EnableAnalytics)
	log.Printf("- Caching enabled: %t", cfg.Features.EnableCaching)
<<<<<<< HEAD
=======
	log.Printf("- Authentication required: %t", cfg.Auth.RequireAuth)
	log.Printf("- Default user role: %s", cfg.Auth.DefaultRole)
	log.Printf("- Token duration: %d hours", cfg.Auth.TokenDuration)
>>>>>>> master

	return cfg, nil
}

func initializeDatabase(cfg *config.Config) (*database.SQLiteDB, error) {
	dataDir := filepath.Dir(cfg.Database.Path)
<<<<<<< HEAD

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory '%s': %w", dataDir, err)
	}
	db, err := database.NewSQLiteDB(cfg.Database.Path)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database: %w", err)
	}
	log.Printf("Database initialized successfully at: %s", cfg.Database.Path)
	
=======
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory '%s': %w", dataDir, err)
	}

	db, err := database.NewSQLiteDB(cfg.Database.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite database: %w", err)
	}

	log.Printf("Database initialized successfully at: %s", cfg.Database.Path)
>>>>>>> master
	return db, nil
}
