package model

import (
	"context"
	"github.com/imandaneshi/vite/pkg/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
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
		log.Fatal("failed connecting to mongodb:", err)
		return err
	}
	pingError := client.Ping(context.TODO(), nil)
	if pingError != nil {
		log.Fatal("failed pinging mongodb:", err)
	}
	log.Info("successfully connected to mongodb")

	// cache our mongo database in m variable
	m = client.Database(config.Database.DatabaseName)

	log.Debug("started creating mongo db indexes")
	err = checkIndexes()
	if err != nil {
		log.Fatal("failed creating mongo db indexes", err)
	}
	log.Info("successfully created mongo db indexes")
	return nil
}

func checkIndexes() (err error){

	// represents collections with their indexes
	collections := map[string][]mongo.IndexModel{
		mongoLinksCollection: {
			{Keys: bsonx.Doc{{"code", bsonx.Int32(1)}},
				Options: options.Index().SetUnique(true).SetName(mongoLinksCodeIndex),
			},
		},
		mongoUsersCollection: {
			{Keys: bsonx.Doc{{"username", bsonx.Int32(1)}},
				Options: options.Index().SetUnique(true).SetName(mongoUsersUsernameIndex),
			},
		},
		mongoTokensCollection: {
			{Keys: bsonx.Doc{{"value", bsonx.Int32(1)}},
				Options: options.Index().SetUnique(true).SetName(mongoTokensValueIndex),
			},
		},
	}

	// loop over collections and create indexes for each mongo collection
	for k,v := range collections{
		logFields := log.Fields{"collection": k}
		log.WithFields(logFields).Debugf("started creating indexes for collection %s", k)
		links := m.Collection(k)
		_, err := links.Indexes().CreateMany(context.Background(), v)
		if err != nil {
			return err
		}
	}

	return
}