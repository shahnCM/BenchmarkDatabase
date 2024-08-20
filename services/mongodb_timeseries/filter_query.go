package mongodb_timeseries

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FilterQuery(client *mongo.Client) error {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the IP address from the environment variable
	ipAddress := os.Getenv("FILTER_QUERY_MONGODB_IP_ADDRESS")
	if ipAddress == "" {
		log.Fatal("FILTER_QUERY_MONGODB_IP_ADDRESS is not set in .env")
	}

	// Define your database and collection
	collection := client.Database("benchmark").Collection("cpu_usage_logs")

	// Build the filter query
	filter := bson.M{
		"metadata.network.ip_address": ipAddress,
	}

	// Start measuring time
	startTime := time.Now()

	// Execute the query
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal("Error executing query")
	}

	// Calculate the duration
	duration := time.Since(startTime)

	defer cursor.Close(context.Background())

	// Iterate through the result set
	var results []bson.M
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal("Error decoding document")
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal("Error during cursor iteration")
	}

	// Print the results
	if len(results) == 0 {
		fmt.Println("No documents found with the specified IP address.")
	} else {
		fmt.Println("Documents found:")
		for _, result := range results {
			fmt.Printf("Document: %v\n", result)
		}
	}

	fmt.Println("> Filter Query on 'ip_address' Duration: ", duration)

	return nil
}
