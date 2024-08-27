package timescaledb_timeseries

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func Invoke(numRows int) (any, error) {

	questDBURL := os.Getenv("TIMESCALEDB_URL")
	if questDBURL == "" {
		return nil, fmt.Errorf("env variable TIMESCALEDB_URL is not set")
	}

	// Connect to QuestDB using the PostgreSQL protocol
	conn, err := pgx.Connect(context.TODO(), questDBURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to TimescaleDb: %v", err)
	}

	// if err := Migrate(conn); err != nil {
	// 	return nil, err
	// }

	// if err := Seed(conn, numRows); err != nil {
	// 	return nil, err
	// }

	if err := FetchAll(conn); err != nil {
		return nil, err
	}

	// if err := FilterQuery(conn); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
