package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"webscraper-v2/internal/config"
	"webscraper-v2/internal/infrastructure/database"
	"webscraper-v2/internal/infrastructure/persistence"
	"webscraper-v2/internal/presentation/server"
	"webscraper-v2/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	// Initialize database
	db, err := database.NewSQLiteDB(cfg.Database.Path)

	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("⚠️  Error closing database: %v", err)
		}
	}()

	// Initialize repositories
	scrapingRepo := persistence.NewScrapingRepository(db)
	userRepo := persistence.NewUserRepository(db)
	scheduleRepo := persistence.NewScheduleRepository(db)

	// Initialize token repository
	tokenRepo := persistence.NewSQLiteTokenRepository(db)

	log.Println("✅ Database and repositories initialized")

	// Initialize use cases
	scrapingUC := usecase.NewScrapingUseCase(scrapingRepo, cfg)
	authUC := usecase.NewAuthUseCase(userRepo, tokenRepo, cfg)
	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo, scrapingUC, cfg)

	log.Println("✅ Use cases initialized")

	// Initialize server
	srv := server.NewServer(cfg.Server.Port, cfg, scrapingUC, authUC, scheduleUC)

	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 Starting WebScraper v2...")
		if err := srv.Start(); err != nil {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		log.Fatalf("❌ Server error: %v", err)
	case sig := <-shutdownChan:
		log.Printf("\n📥 Received signal: %v", sig)
		log.Println("🔄 Shutting down gracefully...")
		authUC.Shutdown()
		log.Println("  ✅ Auth service stopped")

		scheduleUC.StopScheduler() // Detener scheduler
		log.Println("  ✅ Scheduler stopped")

		log.Println("✅ Shutdown complete")
	}
}
