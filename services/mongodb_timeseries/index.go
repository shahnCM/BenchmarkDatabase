package mongodb_timeseries

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Invoke(numRows int) (any, error) {

	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		return nil, fmt.Errorf("env variable MONGODB_URI is not set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	if err := Migrate(client); err != nil {
		return nil, err
	}

	if err := Seed(client, numRows); err != nil {
		return nil, err
	}

	if err := FetchAll(client); err != nil {
		return nil, err
	}

	if err := FilterQuery(client); err != nil {
		return nil, err
	}

	return nil, nil
}
