package scylladb

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

func FilterQuery(session *gocql.Session) error {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the IP address from the environment variable
	ipAddress := os.Getenv("FILTER_QUERY_MONGODB_IP_ADDRESS")
	if ipAddress == "" {
		log.Fatal("FILTER_QUERY_MONGODB_IP_ADDRESS is not set in .env")
	}

	keySpace := "benchmark"
	tableName := "cpu_usage_logs"
	// Define the query with a specific MAC address
	query := fmt.Sprintf("SELECT * FROM %s WHERE ip_address = ?", keySpace+"."+tableName)

	// Start the timer
	start := time.Now()

	// Execute the query
	iter := session.Query(query, ipAddress).Iter()

	// Measure the duration
	duration := time.Since(start)

	// Ensure that the iterator is drained and closed to avoid resource leaks
	if err := iter.Close(); err != nil {
		return fmt.Errorf("error closing iterator: %v", err)
	}

	// Print documents
	fmt.Println("> Filter Query on 'ip_address' Duration: ", duration)

	return nil
}
