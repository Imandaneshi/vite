package model

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func prepareTestDatabase() error {
	testMongoURI := os.Getenv("VITE_TEST_MONGO_URI")
	if testMongoURI == "" {
		testMongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(testMongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error("failed connecting to test mongodb:", err)
		return err
	}
	m = client.Database(testMongoDatabase)
	return nil
}

const (
	testMongoDatabase = "viteTestDatabase"
)