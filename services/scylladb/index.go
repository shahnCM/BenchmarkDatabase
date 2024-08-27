package scylladb

import (
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

func Invoke(numRows int) (any, error) {
	scyllaDBURL := os.Getenv("SCYLLADB_URL")
	if scyllaDBURL == "" {
		return nil, fmt.Errorf("env variable TIMESCALEDB_URL is not set")
	}

	scyllaKeySpace := os.Getenv("SCYLLADB_KEYSPACE")
	if scyllaKeySpace == "" {
		return nil, fmt.Errorf("env variable scyllaKeySpace is not set")
	}

	replication := "{'class': 'SimpleStrategy', 'replication_factor': 1}"

	// Create a new cluster configuration
	cluster := gocql.NewCluster(scyllaDBURL)

	// Connect to the ScyllaDB cluster
	gcqlSession, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to ScyllaDB: %v", err)
	}
	defer gcqlSession.Close()

	// Create gocqlx gcqlSession
	gocqlxSession := gocqlx.NewSession(gcqlSession)
	defer gocqlxSession.Close()

	// CQL query to create the keyspace if it doesn't exist
	createKeyspaceQuery := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH replication = %s
		AND durable_writes = false;
	`, scyllaKeySpace, replication)

	// Execute the keyspace creation query
	if err := gcqlSession.Query(createKeyspaceQuery).Exec(); err != nil {
		return nil, fmt.Errorf("error creating keyspace: %v", err)
	}

	cluster.ConnectTimeout = 60 * time.Second
	cluster.Timeout = 60 * time.Second
	cluster.Consistency = gocql.One
	// cluster.NumConns = 500

	// if err := Migrate(gcqlSession); err != nil {
	// 	return nil, err
	// }

	// if err := Seed(&gocqlxSession, numRows); err != nil {
	// 	return nil, err
	// }

	if err := FetchAll(gcqlSession); err != nil {
		return nil, err
	}

	if err := FilterQuery(gcqlSession); err != nil {
		return nil, err
	}

	return nil, nil
}
