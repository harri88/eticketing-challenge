package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the payment service
type Config struct {
	// Database Configuration
	Database DatabaseConfig

	// Server Configuration
	Server ServerConfig

	// External Services Configuration
	TicketService  TicketServiceConfig
	LedgerService  LedgerServiceConfig
	PaymentGateway PaymentGatewayConfig

	// Environment
	Environment string
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port            int
	Host            string
	ReadTimeoutSec  int
	WriteTimeoutSec int
	Environment     string // development or production
}

// LedgerServiceConfig holds ledger service configuration
type LedgerServiceConfig struct {
	URL string
}

// TicketServiceConfig holds ticket service configuration
type TicketServiceConfig struct {
	URL     string
	Timeout int // in seconds
}

// PaymentGatewayConfig holds payment gateway configuration
type PaymentGatewayConfig struct {
	DefaultProvider string // stripe, paypal, local
	StripeAPIKey    string
	PayPalClientID  string
	PayPalSecret    string
}

// Load loads configuration from .env file and environment variables
func Load() *Config {
	// Load .env file if it exists - handle "file not found" gracefully
	// Standard .env location is in the current working directory
	envPath := ".env"

	// Expand ~ to home directory if needed
	if expandedPath, err := expandPath(envPath); err == nil {
		envPath = expandedPath
	}

	if _, err := os.Stat(envPath); err == nil {
		// File exists, try to load it
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: Failed to load .env file: %v", err)
		}
	} else if !os.IsNotExist(err) {
		// Other errors (permission denied, etc)
		log.Printf("Warning: Error checking .env file: %v", err)
	}
	// If file doesn't exist (IsNotExist), silently use environment variables

	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "payment_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: getEnvInt("DB_MAX_CONNS", 25),
			MinConns: getEnvInt("DB_MIN_CONNS", 5),
		},
		Server: ServerConfig{
			Port:            getEnvInt("SERVER_PORT", 8081),
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeoutSec:  getEnvInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeoutSec: getEnvInt("SERVER_WRITE_TIMEOUT", 10),
			Environment:     getEnv("ENVIRONMENT", "development"),
		},
		TicketService: TicketServiceConfig{
			URL:     getEnv("TICKET_SERVICE_URL", "http://localhost:5020"),
			Timeout: getEnvInt("TICKET_SERVICE_TIMEOUT", 30),
		},
		LedgerService: LedgerServiceConfig{
			URL: getEnv("LEDGER_SERVICE_URL", "http://localhost:8000"),
		},
		PaymentGateway: PaymentGatewayConfig{
			DefaultProvider: getEnv("PAYMENT_GATEWAY_PROVIDER", "stripe"),
			StripeAPIKey:    getEnv("STRIPE_API_KEY", ""),
			PayPalClientID:  getEnv("PAYPAL_CLIENT_ID", ""),
			PayPalSecret:    getEnv("PAYPAL_SECRET", ""),
		},
	}

	return cfg
}

// GetDatabaseURL returns the formatted PostgreSQL connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// IsDevelopment checks if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction checks if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv retrieves an environment variable with a default fallback
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvInt retrieves an integer environment variable with a default fallback
func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid integer value for %s, using default: %d\n", key, defaultVal)
	}
	return defaultVal
}

// expandPath expands ~ to the user's home directory
func expandPath(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}
		return home + path[1:], nil
	}

	return path, nil
}
