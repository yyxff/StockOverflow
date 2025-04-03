package server

import (
	"StockOverflow/internal/database"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

type ServerEntry struct{}

func (serverEntry *ServerEntry) Enter() {
	// Setup logger
	logger := log.New(os.Stdout, "EXCHANGE: ", log.LstdFlags|log.Lshortfile)

	// Initialize database connection
	dbm := database.DatabaseMaster{
		ConnStr: getDBConnStr(),
		DbName:  getEnvOrDefault("DB_NAME", "stockoverflow"),
	}

	// Connect to database
	logger.Println("Connecting to database...")
	dbm.Connect()
	dbm.CreateDB()
	dbm.Init()

	// Create and start the server
	server := NewServer(logger)
	server.SetDB(dbm.Db)
	// Start server in a goroutine
	go func() {
		logger.Println("Starting exchange server on port 12345...")
		if err := server.Start(":12345"); err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for terminate signal
	sig := <-sigCh
	fmt.Printf("Received signal %v, shutting down...\n", sig)

	// Stop the server
	if err := server.Stop(); err != nil {
		logger.Fatalf("Error during shutdown: %v", err)
	}

	logger.Println("Server shutdown complete")
}

// getDBConnStr returns the database connection string from environment
// variables or uses default values
func getDBConnStr() string {
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "passw0rd")

	return fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
		user, password, host, port)
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
