package mongodb_timeseries

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"DbBenchmark/utils"
)

func Seed(client *mongo.Client, seedAmount int) error {

	// WorkingDir path
	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Error getting working directory: %v", err)
	// }

	// Parse the json file
	// docPath := filepath.Join(wd, "doc.json")
	// file, err := os.ReadFile(docPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to read doc.json: %v", err)
	// }

	// Unmarshal to doc
	// var doc map[string]any
	// err = json.Unmarshal(file, &doc)
	// if err != nil {
	// 	return fmt.Errorf("failed to parse JSON: %v", err)
	// }

	// Random Json Data
	doc := utils.RandomizeJSON()

	// Access the database and collection
	database := client.Database("benchmark") // Replace with your database name
	collection := database.Collection("cpu_usage_logs")

	// Prepare data
	data := bson.D{
		{Key: "logged_at", Value: time.Now()},
		{Key: "metadata", Value: doc["metadata"]},
	}

	// Start measuring time
	startTime := time.Now()

	// seed
	for i := 1; i < seedAmount; i++ {
		if _, err := collection.InsertOne(context.TODO(), data); err != nil {
			log.Fatalf("Failed to insert document %d: %v", i, err)
		}

		// Update the progress bar
		if i%1000 == 0 {
			go utils.PrintProgressBar(i, seedAmount)
		}
	}

	// Calculate the duration
	duration := time.Since(startTime)

	fmt.Printf("\n> Rows inserted successfully\n")
	fmt.Println("> Duration: ", duration)

	return nil
}
