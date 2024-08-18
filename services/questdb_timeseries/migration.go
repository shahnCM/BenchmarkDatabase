package questdb_timeseries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func Migrate(conn *pgx.Conn) error {

	// Define database and collection names
	// dbName := "benchmark"
	tableName := "cpu_usage_logs"

	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		cpu_id SYMBOL,
		user_id SYMBOL,
		device_id SYMBOL,
		ip_address STRING,
		mac_address STRING,
		app_id SYMBOL,
		usage_by_user DOUBLE,
		usage_by_system DOUBLE,
		usage_by_idle DOUBLE,
		dump STRING,  -- Column to store the entire document as JSON
		logged_at TIMESTAMP
	) timestamp(logged_at) PARTITION BY DAY;`, tableName)

	// Execute the table creation query
	if _, err := conn.Exec(context.Background(), createTableQuery); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	fmt.Println("> Table 'cpu_usage_logs' and indexes are ready.")

	return nil
}
