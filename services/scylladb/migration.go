package scylladb

import (
	"fmt"
	"os"

	"github.com/gocql/gocql"
)

func Migrate(session *gocql.Session) error {

	scyllaKeySpace := os.Getenv("SCYLLADB_KEYSPACE")
	if scyllaKeySpace == "" {
		return fmt.Errorf("env variable scyllaKeySpace is not set")
	}

	tableName := "cpu_usage_logs"

	// CQL query to create the table
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			cpu_id INT,
			user_id INT,
			device_id INT,
			app_id INT,
			ip_address TEXT,
			mac_address TEXT,
			usage_by_user DOUBLE,
			usage_by_system DOUBLE,
			usage_by_idle DOUBLE,
			dump TEXT,  -- Column to store the entire document as JSON
			bucket_date DATE,
			logged_at TIMESTAMP,
			PRIMARY KEY ((bucket_date), logged_at)
		)  
			WITH compression = {'sstable_compression': ''};
	`, scyllaKeySpace+"."+tableName)

	// Compression Methods
	// org.apache.cassandra.io.compress.SnappyCompressor

	if err := session.Query(createTableQuery, nil).Exec(); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Execute each index creation query individually
	queries := []string{
		"CREATE INDEX IF NOT EXISTS idx_user_id ON benchmark.cpu_usage_logs (user_id);",
		"CREATE INDEX IF NOT EXISTS idx_cpu_id ON benchmark.cpu_usage_logs (cpu_id);",
		"CREATE INDEX IF NOT EXISTS idx_device_id ON benchmark.cpu_usage_logs (device_id);",
		"CREATE INDEX IF NOT EXISTS idx_app_id ON benchmark.cpu_usage_logs (app_id);",
		"CREATE INDEX IF NOT EXISTS idx_ip_address ON benchmark.cpu_usage_logs (ip_address);",
		"CREATE INDEX IF NOT EXISTS idx_mac_address ON benchmark.cpu_usage_logs (mac_address);",
	}

	for _, query := range queries {
		if err := session.Query(query, nil).Exec(); err != nil {
			return fmt.Errorf("error creating index: %v", err)
		}
	}

	fmt.Println("> Table 'cpu_usage_logs', hypertable, indexes, compression policy, and retention function are ready.")

	return nil
}
