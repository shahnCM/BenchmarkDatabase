package mongodb_timeseries

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchAll(client *mongo.Client) error {

	// Define your database and collection
	collection := client.Database("benchmark").Collection("cpu_usage_logs")

	// Fetch the environment variable
	fetchLimitStr := os.Getenv("FETCH_LIMIT")

	// Convert it to an integer
	fetchLimit, err := strconv.Atoi(fetchLimitStr)
	if err != nil {
		log.Fatalf("Failed to parse FETCH_LIMIT: %v", err)
	}

	// Start measuring time
	startTime := time.Now()
	// Execute the query
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetLimit(int64(fetchLimit))) // Example limit
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	// Calculate the duration
	duration := time.Since(startTime)
	defer cursor.Close(context.Background())

	// Decode the result
	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		return fmt.Errorf("error decoding results: %v", err)
	}

	if len(result) == 0 {
		fmt.Println("> No results found.")
		return nil
	}

	// Print documents
	fmt.Printf("\n> Requested %d rows\n", fetchLimit)
	fmt.Printf("> Fetched %d documents\n", len(result))
	fmt.Println("> Fetch Duration: ", duration)

	return nil

}

// Define the aggregation pipeline using $facet
// pipeline := mongo.Pipeline{
// 	// Stage 1: $facet stage to run multiple pipelines
// 	{
// 		{Key: "$facet", Value: bson.D{
// 			// Pipeline to count documents
// 			{Key: "total", Value: bson.A{
// 				bson.D{{Key: "$count", Value: "total"}},
// 			}},
// 			// Pipeline to fetch documents with a limit
// 			{Key: "documents", Value: bson.A{
// 				bson.D{{Key: "$limit", Value: int64(1000)}},
// 			}},
// 		}},
// 	},
// }
// cursor, err := collection.Aggregate(context.TODO(), pipeline)

// Extract and print the count and documents
// counts, ok := result[0]["total"].(bson.A)
// if !ok {
// 	return fmt.Errorf("failed to convert count to bson.A")
// }
// documents, ok := result[0]["documents"].(bson.A)
// if !ok {
// 	return fmt.Errorf("failed to convert documents to bson.A")
// }

// Print count
// if len(counts) > 0 {
// 	count := counts[0].(bson.M)["total"]
// 	fmt.Printf("> Total count of documents: %v\n", count)
// } else {
// 	fmt.Println("> Count not found.")
// }

// for _, doc := range documents {
// 	fmt.Println(doc)
// }
