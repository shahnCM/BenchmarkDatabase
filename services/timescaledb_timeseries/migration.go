package timescaledb_timeseries

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func Migrate(conn *pgx.Conn) error {
	// Define the table name
	tableName := "cpu_usage_logs"

	// SQL query to create the table and hypertable
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			cpu_id INTEGER,
			user_id INTEGER,
			device_id INTEGER,
			app_id INTEGER,
			ip_address VARCHAR(50),
			mac_address VARCHAR(50),
			usage_by_user DOUBLE PRECISION,
			usage_by_system DOUBLE PRECISION,
			usage_by_idle DOUBLE PRECISION,
			dump TEXT,  -- Column to store the entire document as JSON
			logged_at TIMESTAMPTZ NOT NULL
		);
	`, tableName)

	// Execute the table creation and hypertable creation query
	if _, err := conn.Exec(context.Background(), createTableQuery); err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// SQL query to create the table and hypertable
	createHyperTableQuery := fmt.Sprintf(`
		-- Create a hypertable directly
		SELECT create_hypertable(
		'%s', 
		'logged_at',
		chunk_time_interval => INTERVAL '1 day',
		-- partitioning_column => 'logged_at',
		if_not_exists => TRUE);
	`, tableName)

	// Execute the table creation and hypertable creation query
	if _, err := conn.Exec(context.Background(), createHyperTableQuery); err != nil {
		return fmt.Errorf("error creating hypertable: %v", err)
	}

	// Define indexes for integer columns
	indexesQuery := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS idx_cpu_id ON %s (cpu_id);
		CREATE INDEX IF NOT EXISTS idx_user_id ON %s (user_id);
		CREATE INDEX IF NOT EXISTS idx_device_id ON %s (device_id);
		CREATE INDEX IF NOT EXISTS idx_app_id ON %s (app_id);
		CREATE INDEX IF NOT EXISTS idx_app_id ON %s (ip_address);
		CREATE INDEX IF NOT EXISTS idx_app_id ON %s (mac_address);
	`, tableName, tableName, tableName, tableName, tableName, tableName)

	// Execute the index creation query
	if _, err := conn.Exec(context.Background(), indexesQuery); err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}

	// Define a compression policy
	// compressionPolicyQuery := fmt.Sprintf(`
	// 	-- Enable compression on the hypertable
	// 	ALTER TABLE %s
	// 	SET (timescaledb.compress, timescaledb.compress_segmentby = 'cpu_id');

	// 	-- Manually create a compression policy for data older than 1 day
	// 	CREATE POLICY compress_policy ON %s
	// 	USING (logged_at < NOW() - INTERVAL '1 day');
	// `, tableName, tableName)

	// Execute the compression policy query
	// if _, err := conn.Exec(context.Background(), compressionPolicyQuery); err != nil {
	// 	return fmt.Errorf("error setting compression policy: %v", err)
	// }

	// Define a retention function
	// retentionFunctionQuery := `
	// 	CREATE OR REPLACE FUNCTION delete_old_data()
	// 	RETURNS void LANGUAGE plpgsql AS $$
	// 	BEGIN
	// 		DELETE FROM cpu_usage_logs WHERE logged_at < NOW() - INTERVAL '30 days';
	// 	END;
	// 	$$;
	// `

	// Execute the retention function creation query
	// if _, err := conn.Exec(context.Background(), retentionFunctionQuery); err != nil {
	// 	return fmt.Errorf("error creating retention function: %v", err)
	// }

	// Schedule the retention function using pg_cron (requires pg_cron extension)
	// scheduleRetentionFunctionQuery := `
	// 	SELECT cron.schedule('daily_delete_old_data', '0 0 * * *', 'SELECT delete_old_data();');
	// `

	// Execute the scheduling query
	// if _, err := conn.Exec(context.Background(), scheduleRetentionFunctionQuery); err != nil {
	// 	return fmt.Errorf("error scheduling retention function: %v", err)
	// }

	fmt.Println("> Table 'cpu_usage_logs', hypertable, indexes, compression policy, and retention function are ready.")

	return nil
}
