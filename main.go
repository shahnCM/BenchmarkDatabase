package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/joho/godotenv"

	"DbBenchmark/services/mongodb_timeseries"
)

func main() {

	// Load environment variables from the specified .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file")
	}

	// Define choices for database selection
	dbChoices := []string{"MongoDb Timeseries", "QuestDb Timeseries"}

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
		insertIntoQuestDB(numRows)
	default:
		log.Fatalf("Unsupported database choice: %s", selectedDB)
	}
}

func insertIntoQuestDB(rows int) {
	fmt.Printf("Inserting %d rows into QuestDB Timeseries...\n", rows)
	// Add your QuestDB insertion logic here
}
