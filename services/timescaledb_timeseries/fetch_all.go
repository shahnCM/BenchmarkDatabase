package timescaledb_timeseries

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

func FetchAll(conn *pgx.Conn) error {

	// Fetch the environment variable
	fetchLimitStr := os.Getenv("FETCH_LIMIT")

	// Convert it to an integer
	fetchLimit, err := strconv.Atoi(fetchLimitStr)
	if err != nil {
		log.Fatalf("Failed to parse FETCH_LIMIT: %v", err)
	}

	// set databasse name
	databaseTable := "cpu_usage_logs"

	query := fmt.Sprintf(`
		SELECT 
		cpu_id, user_id, device_id, app_id, ip_address, mac_address, 
		usage_by_user, usage_by_system, usage_by_idle,
		dump, logged_at
		FROM %s
		LIMIT %d;
	`, databaseTable, fetchLimit)

	// Start measuring time
	startTime := time.Now()
	// Execute the query
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	// Calculate the duration
	duration := time.Since(startTime)
	defer rows.Close()

	var rowsLen int
	for rows.Next() {
		rowsLen++
	}

	fmt.Printf("\n> Fetched %d rows\n", fetchLimit)
	fmt.Printf("> Fetched %d rows\n", rowsLen)
	fmt.Println("> Fetch Duration: ", duration)

	return nil

}

// Iterate over the rows and process the data
// for rows.Next() {
// 	var loggedAt time.Time
// 	var cpuID, deviceID, appID, userID, ipAddress, macAddress string
// 	var usageByUser, usageByIdle, usageBySystem float64
// 	var dump string // Use *string for nullable strings

// 	err := rows.Scan(&loggedAt, &cpuID, &usageByUser, &usageBySystem, &usageByIdle,
// 		&deviceID, &appID, &userID, &ipAddress, &macAddress, &dump)
// 	if err != nil {
// 		log.Fatalf("Row scan failed: %v\n", err)
// 	}
// }
