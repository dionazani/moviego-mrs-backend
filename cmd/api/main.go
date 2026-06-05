package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dionazani/moviego-mrs-backend/internal/config"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/database"
	"github.com/dionazani/moviego-mrs-backend/internal/router"
)

func main() {
	log.Println("Starting MovieGo MRS Backend with Gin Framework...")

	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database (GORM)
	db := infrastructuredatabase.InitDB(cfg)

	// Verify database connection using standard ping
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve SQL DB instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// 3. Initialize Gin Engine Router
	appRouter := router.NewRouter(cfg, db)

	// 4. Setup Custom HTTP Server Configuration
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      appRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Run Server asynchronously and handle graceful shutdowns
	go func() {
		log.Printf("HTTP Server is listening and serving on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Critical server error: %v", err)
		}
	}()

	// Wait for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	// Create a timeout context for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited successfully")
}
