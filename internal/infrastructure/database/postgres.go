package infrastructuredatabase

import (
	"context"
	"log"
	"time"

	"github.com/dionazani/moviego-mrs-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	// Open connection to Postgres using GORM
	// We set logger to Info to see executed SQL queries during development
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database interface: %v", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection successfully configured.")
	return db
}

type contextKey struct{}

var txKey = contextKey{}

// WithTransaction returns a new context containing the GORM transaction db instance.
func WithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// GetTx retrieves the GORM transaction db instance from the context if it exists.
func GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}

