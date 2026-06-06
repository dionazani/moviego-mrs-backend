package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN                 string
	ServerPort            string
	MasterUserRoleRegular string
	MasterUserRoleAdmin   string
	JwtSecret             string
	LoginLockoutDuration  time.Duration
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// Match PostgreSQL DSN format for GORM
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE") +
		" TimeZone=" + os.Getenv("DB_TIMEZONE")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Parse login lockout duration from env, defaulting to 10 minutes.
	lockoutDuration := 10 * time.Minute
	if raw := os.Getenv("USER_LOGIN_LOCKOUT_DURATION"); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil {
			lockoutDuration = d
		} else {
			log.Printf("Warning: invalid USER_LOGIN_LOCKOUT_DURATION value '%s', using default 10m", raw)
		}
	}

	return &Config{
		DBDSN:                 dsn,
		ServerPort:            port,
		MasterUserRoleRegular: os.Getenv("MASTER_USER_ROLE_REGULAR"),
		MasterUserRoleAdmin:   os.Getenv("MASTER_USER_ROLE_ADMIN"),
		JwtSecret:             os.Getenv("JWT_SECRET"),
		LoginLockoutDuration:  lockoutDuration,
	}
}
