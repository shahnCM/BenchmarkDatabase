package timescaledb_timeseries

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v4"

	"DbBenchmark/utils"
)

func Seed(conn *pgx.Conn, seedAmount int) error {
	// WorkingDir path
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	// Parse the json file
	docPath := filepath.Join(wd, "doc.json")
	file, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("failed to read doc.json: %v", err)
	}

	// Unmarshal to doc
	var doc map[string]any
	err = json.Unmarshal(file, &doc)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	tableName := "cpu_usage_logs"

	// Prepare data
	loggedAt := time.Now().Format(time.RFC3339) // Format for timestamp
	cpuID := 1                                  // Example value, adjust as needed
	deviceID := 1                               // Example value, adjust as needed
	appID := 1                                  // Example value, adjust as needed
	userID := 1                                 // Example value, adjust as needed
	usageUser := 0.5                            // Example value, adjust as needed
	usageIdle := 0.3                            // Example value, adjust as needed
	usageSystem := 0.2                          // Example value, adjust as needed
	ipAddress := "192.168.1.1"                  // Example value, adjust as needed
	macAddress := "00:00:00:00:00:00"           // Example value, adjust as needed

	// Prepare dump data from JSON
	dump, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal dump data: %v", err)
	}
	dumpString := string(dump)

	// Create insert query
	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (
			cpu_id, user_id, device_id, app_id, ip_address, mac_address, 
			usage_by_user, usage_by_system, usage_by_idle,
			dump, logged_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		);
	`, tableName)

	// Start measuring time
	startTime := time.Now()

	// Seed the data
	for i := 1; i <= seedAmount; i++ {
		if _, err := conn.Exec(context.TODO(), insertQuery,
			cpuID, userID, deviceID, appID,
			ipAddress, macAddress,
			usageUser, usageSystem, usageIdle,
			dumpString, loggedAt,
		); err != nil {
			return fmt.Errorf("failed to insert data %d: %v", i, err)
		}

		// Print progress
		if i%1000 == 0 { // Print progress every 1000 rows
			go utils.PrintProgressBar(i, seedAmount)
		}
	}

	// Calculate the duration
	duration := time.Since(startTime)

	fmt.Printf("\n> Rows inserted successfully\n")
	fmt.Println("> Duration: ", duration)

	return nil
}
