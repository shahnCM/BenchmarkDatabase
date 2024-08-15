package mongodb_timeseries

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migrate(client *mongo.Client) error {

	// Define database and collection names
	dbName := "benchmark"
	collectionName := "cpu_usage_logs"

	// Get a handle to the database
	db := client.Database(dbName)

	// Check if the collection exists, if not, create it
	if err := ensureCollectionExists(db, collectionName); err != nil {
		return fmt.Errorf("error ensuring collection exists: %v", err)
	}

	fmt.Printf("> Database '%s' and collection '%s' are ready.\n", dbName, collectionName)

	return nil
}

// Ensure the collection exists, create it if it doesn't
func ensureCollectionExists(db *mongo.Database, collectionName string) error {
	collections, err := db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	for _, name := range collections {
		if name == collectionName {
			return nil // Collection exists
		}
	}

	// Create the time-series collection if it does not exist
	timeSeriesOpts := options.TimeSeries().SetTimeField("logged_at").SetMetaField("metadata").SetGranularity("seconds")
	createOpts := options.CreateCollection().SetTimeSeriesOptions(timeSeriesOpts)
	err = db.CreateCollection(context.TODO(), collectionName, createOpts)
	if err != nil {
		return err
	}

	// Create the collection if it does not exist
	collection := db.Collection(collectionName)

	// List of fields to be individually indexed
	fields := []string{
		"metadata.cpu.cpu_id",
		"metadata.cpu.usage.user",
		"metadata.cpu.usage.idle",
		"metadata.cpu.usage.system",
		"metadata.device.device_id",
		"metadata.application.app_id",
		"metadata.user.user_id",
		"metadata.network.ip_address",
		"metadata.network.mac_address",
	}

	// Create individual indexes for each field
	for _, field := range fields {
		indexModel := mongo.IndexModel{
			Keys: bson.D{{Key: field, Value: 1}}, // 1 for ascending order
		}

		_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			return fmt.Errorf("error creating index on %s: %v", field, err)
		}
	}

	// Create a TTL index on the logged_at field to expire documents after 30 days
	ttlIndexOpts := options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60)
	if _, err = collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "logged_at", Value: 1}},
		Options: ttlIndexOpts.SetPartialFilterExpression(bson.M{"metadata.cpu.cpu_id": bson.M{"$exists": true}}),
	}); err != nil {
		return fmt.Errorf("error creating TTL index on logged_at: %v", err)
	}

	return nil
}
