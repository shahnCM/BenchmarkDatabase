package scylladb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/table"

	"DbBenchmark/utils"
)

func Seed(session *gocqlx.Session, seedAmount int) error {
	keySpace := "benchmark"
	tableName := "cpu_usage_logs"

	cpuUsageLogsMetadata := table.Metadata{
		Name: keySpace + "." + tableName,
		Columns: []string{
			"cpu_id",
			"user_id",
			"device_id",
			"app_id",
			"ip_address",
			"mac_address",
			"usage_by_user",
			"usage_by_system",
			"usage_by_idle",
			"dump",
			"bucket_date",
			"logged_at",
		},
		PartKey: []string{"bucket_date"}, // Composite partition key
		SortKey: []string{"logged_at"},   // Clustering key
	}

	// personTable allows for simple CRUD operations based on personMetadata.
	cpuUsageLogsTable := table.New(cpuUsageLogsMetadata)

	// Start measuring time
	startTime := time.Now()

	// Seed the data
	for i := 1; i <= seedAmount; i++ {

		// Prepare data
		bucketDate := time.Now().Format("2006-01-02")
		loggedAt := time.Now().UTC()
		cpuID := utils.RandomInt(1, 100)
		deviceID := utils.RandomInt(1, 100)
		appID := utils.RandomInt(1, 100)
		userID := utils.RandomInt(1, 100)
		usageUser := utils.RandomDouble(0, 100)
		usageIdle := utils.RandomDouble(0, 100)
		usageSystem := utils.RandomDouble(0, 100)
		ipAddress := utils.RandomString(10) + "192.168.1.1"
		macAddress := utils.RandomString(10) + "00:00:00:00:00:00"

		// Prepare dump data from JSON
		// Random JSON Data
		doc := utils.RandomizeJSON()
		dump, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal dump data: %v", err)
		}
		dumpString := string(dump)

		// Define and bind the data using an anonymous struct
		q := session.Query(cpuUsageLogsTable.Insert()).BindStruct(struct {
			CpuID         int
			UserID        int
			DeviceID      int
			AppID         int
			IPAddress     string
			MACAddress    string
			UsageByUser   float64
			UsageBySystem float64
			UsageByIdle   float64
			Dump          string
			LoggedAt      time.Time
			BucketDate    string
		}{
			CpuID:         cpuID,
			UserID:        userID,
			DeviceID:      deviceID,
			AppID:         appID,
			IPAddress:     ipAddress,
			MACAddress:    macAddress,
			UsageByUser:   usageUser,
			UsageBySystem: usageSystem,
			UsageByIdle:   usageIdle,
			Dump:          dumpString,
			LoggedAt:      loggedAt,
			BucketDate:    bucketDate,
		})

		if err := q.ExecRelease(); err != nil {
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
