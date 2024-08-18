package questdb_timeseries

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

func FetchAll(conn *pgx.Conn) error {

	// set databasse name
	databaseTable := "cpu_usage_logs"

	query := fmt.Sprintf(`
		SELECT logged_at, cpu_id, usage_by_user, usage_by_system, usage_by_idle, 
		       device_id, app_id, user_id, ip_address, mac_address, dump
		FROM %s
		LIMIT 100000;
	`, databaseTable)

	// Start measuring time
	startTime := time.Now()

	// Execute the query
	rows, err := conn.Query(context.TODO(), query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	rows.Close()

	// Calculate the duration
	duration := time.Since(startTime)

	fmt.Println("> Total count of rows: 1000000")
	fmt.Println("> Fetched 100000 rows")
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
