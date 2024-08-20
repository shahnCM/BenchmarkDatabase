package cratedb_timeseries

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func Invoke(numRows int) (any, error) {

	crateDBURL := os.Getenv("CRATEDB_URL")
	if crateDBURL == "" {
		return nil, fmt.Errorf("env variable CRATEDB_URL is not set")
	}

	// Connect to QuestDB using the PostgreSQL protocol
	conn, err := pgx.Connect(context.TODO(), crateDBURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to CRATE DB: %v", err)
	}

	if err := Migrate(conn); err != nil {
		return nil, err
	}

	if err := Seed(conn, numRows); err != nil {
		return nil, err
	}

	if err := FetchAll(conn); err != nil {
		return nil, err
	}

	return nil, nil
}
