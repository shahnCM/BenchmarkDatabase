package cratedb_timeseries

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func FilterQuery(conn *pgx.Conn) error {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the IP address from the environment variable
	ipAddress := os.Getenv("FILTER_QUERY_CRATEDB_IP_ADDRESS")
	if ipAddress == "" {
		log.Fatal("FILTER_QUERY_CRATEDB_IP_ADDRESS is not set in .env")
	}

	// set databasse name
	databaseTable := "cpu_usage_logs"

	// Define the query
	query := fmt.Sprintf("SELECT user_id, cpu_id, device_id, app_id, ip_address, mac_address, usage_by_user, usage_by_idle, usage_by_system, dump FROM %s WHERE ip_address = $1", databaseTable)

	// Start measuring time
	startTime := time.Now()

	// Execute the query
	rows, err := conn.Query(context.TODO(), query, ipAddress)
	if err != nil {
		log.Fatal("Query execution failed:", err)
	}

	// Calculate the duration
	duration := time.Since(startTime)

	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var (
			userID        int
			cpuID         int
			deviceID      int
			appID         int
			ipAddress     string
			macAddress    string
			usageByUser   float64
			usageByIdle   float64
			usageBySystem float64
			dump          string
		)

		// Scan the result into the variables
		err := rows.Scan(&userID, &cpuID, &deviceID, &appID, &ipAddress, &macAddress, &usageByUser, &usageByIdle, &usageBySystem, &dump)
		if err != nil {
			log.Fatal("Row scan failed:", err)
		}

		// Print or process the row
		fmt.Printf("user_id: %d, cpu_id: %d, device_id: %d, app_id: %d, ip_address: %s, mac_address: %s, usage_by_user: %.2f, usage_by_idle: %.2f, usage_by_system: %.2f, dump: %s\n",
			userID, cpuID, deviceID, appID, ipAddress, macAddress, usageByUser, usageByIdle, usageBySystem, dump)
	}

	// Check for errors after iteration
	if rows.Err() != nil {
		log.Fatal("Row iteration error:", rows.Err())
	}

	fmt.Println("> Filter Query on 'ip_address' Duration: ", duration)

	return nil
}
