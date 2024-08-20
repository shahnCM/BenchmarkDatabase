package cratedb_timeseries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func Migrate(conn *pgx.Conn) error {
	// Define table name
	tableName := "cpu_usage_logs"

	// Create table query with indexing
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			cpu_id INT,
			user_id INT,
			device_id INT,
			app_id INT,
			ip_address STRING,
			mac_address STRING,
			usage_by_user DOUBLE,
			usage_by_system DOUBLE,
			usage_by_idle DOUBLE,
			dump STRING,
			logged_at TIMESTAMP
		) CLUSTERED BY (cpu_id) PARTITIONED BY (logged_at) WITH (number_of_replicas = 1);
	`, tableName)

	// Execute the table creation query
	if _, err := conn.Exec(context.Background(), createTableQuery); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	fmt.Println("> Table 'cpu_usage_logs' with indexes is ready.")
	return nil
}
