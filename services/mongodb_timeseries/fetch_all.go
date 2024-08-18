package mongodb_timeseries

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchAll(client *mongo.Client) error {

	// Define your database and collection
	collection := client.Database("benchmark").Collection("cpu_usage_logs")

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

	// Start measuring time
	startTime := time.Now()

	// Run the aggregation
	// cursor, err := collection.Aggregate(context.TODO(), pipeline)
	// Execute the query
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetLimit(1000000)) // Example limit
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer cursor.Close(context.Background())
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	// Calculate the duration
	duration := time.Since(startTime)

	// Decode the result
	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		return fmt.Errorf("error decoding results: %v", err)
	}

	if len(result) == 0 {
		fmt.Println("> No results found.")
		return nil
	}

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

	// Print documents
	fmt.Printf("> Fetched %d documents\n", len(result))
	// for _, doc := range documents {
	// 	fmt.Println(doc)
	// }

	fmt.Println("> Fetch Duration: ", duration)

	return nil

}
