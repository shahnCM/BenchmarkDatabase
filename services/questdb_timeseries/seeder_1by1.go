package questdb_timeseries

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/questdb/go-questdb-client/v3"

	"DbBenchmark/utils"
)

func Seed(conn *pgx.Conn, seedAmount int) error {

	// WorkingDir path
	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Error getting working directory: %v", err)
	// }

	// Parse the json file
	// docPath := filepath.Join(wd, "doc.json")
	// file, err := os.ReadFile(docPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to read doc.json: %v", err)
	// }

	// Unmarshal to doc
	// var doc map[string]any
	// err = json.Unmarshal(file, &doc)
	// if err != nil {
	// 	return fmt.Errorf("failed to parse JSON: %v", err)
	// }

	// Random Json Data
	doc := utils.RandomizeJSON()

	tableName := "cpu_usage_logs"

	// Prepare data
	loggedAt := time.Now().Format(time.RFC3339)                // Format for timestamp
	cpuID := utils.RandomString(5)                             // Example value, adjust as needed
	deviceID := utils.RandomString(5)                          // Example value, adjust as needed
	appID := utils.RandomString(5)                             // Example value, adjust as needed
	userID := utils.RandomString(5)                            // Example value, adjust as needed
	usageUser := utils.RandomDouble(0, 100)                    // Example value, adjust as needed
	usageIdle := utils.RandomDouble(0, 100)                    // Example value, adjust as needed
	usageSystem := utils.RandomDouble(0, 100)                  // Example value, adjust as needed
	ipAddress := utils.RandomString(10) + "192.168.1.1"        // Example value, adjust as needed
	macAddress := utils.RandomString(10) + "00:00:00:00:00:00" // Example value, adjust as needed

	// Prepare dump data from JSON
	dump, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal dump data: %v", err)
	}
	dumpString := string(dump)

	// Get the environment variable value
	envVarValue := os.Getenv("QUESTDB_INFLUX_LINE_PROTO_FOR_INSERT")
	useInfluxLine, err := strconv.ParseBool(envVarValue)
	if err != nil {
		log.Fatalf("Failed to parse 'QUESTDB_INFLUX_LINE_PROTO_FOR_INSERT': %v", err)
	}

	if useInfluxLine {
		ctx := context.Background()
		client, err := questdb.LineSenderFromEnv(ctx)
		if err != nil {
			log.Fatalf("Failed to initiate INFLUX LINE PROTOCOL: %v", err)
		}
		defer client.Close(ctx)

		// Start measuring time
		startTime := time.Now()

		// Seed
		for i := 1; i <= seedAmount; i++ {
			cpuUsageTs, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			if err != nil {
				log.Fatal(err)
			}

			err = client.
				Table(tableName).
				Symbol("cpu_id", cpuID).
				Symbol("user_id", userID).
				Symbol("device_id", deviceID).
				Symbol("app_id", appID).
				StringColumn("mac_address", macAddress).
				StringColumn("ip_address", ipAddress).
				Float64Column("usage_by_user", usageUser).
				Float64Column("usage_by_system", usageSystem).
				Float64Column("usage_by_idle", usageIdle).
				StringColumn("dump", dumpString).
				At(ctx, cpuUsageTs)
			if err != nil {
				log.Println(err)
			}

			if err = client.Flush(ctx); err != nil {
				log.Println(err)
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

	// Create insert query
	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (
			logged_at, cpu_id, usage_by_user, usage_by_idle, usage_by_system,
			device_id, app_id, user_id, ip_address, mac_address, dump
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11
		);
	`, tableName)

	// Start measuring time
	startTime := time.Now()

	// Seed the data
	for i := 1; i <= seedAmount; i++ {
		if _, err := conn.Exec(context.TODO(), insertQuery,
			loggedAt, cpuID, usageUser, usageIdle, usageSystem,
			deviceID, appID, userID, ipAddress, macAddress, dumpString,
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
