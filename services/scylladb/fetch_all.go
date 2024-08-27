package scylladb

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

func FetchAll(session *gocql.Session) error {

	// Fetch the environment variable
	fetchLimitStr := os.Getenv("FETCH_LIMIT")

	// Convert it to an integer
	fetchLimit, err := strconv.Atoi(fetchLimitStr)
	if err != nil {
		return fmt.Errorf("error to parse FETCH_LIMIT: %v", err)
	}

	keySpace := "benchmark"
	tableName := "cpu_usage_logs"
	count := 0
	queryCount := 0

	// First query, start without any paging state
	var pagingState []byte
	var iter *gocql.Iter

	// Record the start time
	startTime := time.Now()

	for {
		// Create your query
		query := session.Query(fmt.Sprintf(`SELECT * FROM %s LIMIT %d`, keySpace+"."+tableName, fetchLimit)).PageState(pagingState)

		// Execute the query
		iter = query.Iter()
		// Process rows here

		// Iterate over the results
		queryCount++

		// Iterate over the results
		count += iter.NumRows()

		// Get the paging state for the next iteration
		pagingState = iter.PageState()

		// If pagingState is empty, you've reached the end
		if len(pagingState) == 0 {
			break
		}
	}

	// Measure the duration
	duration := time.Since(startTime)

	// Check for errors during iteration
	if err := iter.Close(); err != nil {
		return fmt.Errorf("error closing iterator: %v", err)
	}

	// Print documents
	fmt.Printf("\n> Requested %d rows\n", fetchLimit)
	fmt.Printf("> Fetched %d documents\n", count)
	fmt.Printf("> Total query counts %d \n", queryCount)
	fmt.Println("> Fetch Duration: ", duration)

	return nil
}
