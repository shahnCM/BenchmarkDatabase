package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/joho/godotenv"

	"DbBenchmark/services/cratedb_timeseries"
	"DbBenchmark/services/mongodb_timeseries"
	"DbBenchmark/services/questdb_timeseries"
	"DbBenchmark/services/scylladb"
	"DbBenchmark/services/timescaledb_timeseries"
)

func main() {

	// Load environment variables from the specified .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file")
	}

	// Define choices for database selection
	dbChoices := []string{"MongoDb Timeseries", "QuestDb Timeseries", "TimescaleDb Timeseries", "CrateDb Timeseries", "ScyllaDb"}

	// Prompt user to select a database for benchmarking
	var selectedDB string
	err := survey.AskOne(&survey.Select{
		Message: "Select DB for Benchmark (Single Insert):",
		Options: dbChoices,
	}, &selectedDB)
	if err != nil {
		log.Fatalf("Error selecting database: %v", err)
	}

	// Prompt user to enter the number of rows to insert
	var numRowsStr string
	err = survey.AskOne(&survey.Input{
		Message: "Number of rows to insert:",
	}, &numRowsStr)
	if err != nil {
		log.Fatalf("Error reading number of rows: %v", err)
	}

	// Convert number of rows input to integer
	numRows, err := strconv.Atoi(numRowsStr)
	if err != nil {
		log.Fatalf("Error parsing number of rows: %v", err)
	}

	// Call the appropriate function based on the selected database
	switch selectedDB {

	case "MongoDb Timeseries":
		result, err := mongodb_timeseries.Invoke(numRows)
		if err != nil {
			log.Println(err)
		}
		if result != nil {
			log.Println(result)
		}
	case "QuestDb Timeseries":
		result, err := questdb_timeseries.Invoke(numRows)
		if err != nil {
			log.Println(err)
		}
		if result != nil {
			log.Println(result)
		}
	case "TimescaleDb Timeseries":
		result, err := timescaledb_timeseries.Invoke(numRows)
		if err != nil {
			log.Println(err)
		}
		if result != nil {
			log.Println(result)
		}
	case "CrateDb Timeseries":
		result, err := cratedb_timeseries.Invoke(numRows)
		if err != nil {
			log.Println(err)
		}
		if result != nil {
			log.Println(result)
		}
	case "ScyllaDb":
		result, err := scylladb.Invoke(numRows)
		if err != nil {
			log.Println(err)
		}
		if result != nil {
			log.Println(result)
		}
	default:
		log.Fatalf("Unsupported database choice: %s", selectedDB)
	}
}
