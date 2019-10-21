package model

import (
	"context"
	"github.com/imandaneshi/vite/pkg/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	m *mongo.Database // cached mongo database for internal use
)

// SetupMongo connects to mongo and fills m variable with a mongo database
// so that other models can use it easily
func SetupMongo() error {
	log.Debugf("connecting to mongodb %s", config.Database.Uri)
	clientOptions := options.Client().ApplyURI(config.Database.Uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error("failed connecting to mongodb:", err)
		return err
	}
	log.Info("successfully connected to mongodb")
	m = client.Database(config.Database.DatabaseName)
	return nil
}