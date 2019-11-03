package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestGetLink(t *testing.T) {

	sampleAddress := "https://google.com"
	sampleCode := "onvceqn"

	err := prepareTestDatabase()
	if err != nil {
		log.Fatal("failed loading test database", err)
	}

	err = m.Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {

		if err := sessionContext.StartTransaction(); err != nil {
			log.Fatal(err)
		}

		links := m.Collection(mongoLinksCollection)
		res, err := links.InsertOne(context.Background(), &Link{Code: sampleCode, Address: sampleAddress})
		if err != nil {
			log.Fatal("failed inserting test data in mongodb", err)
		}
		var objectId *primitive.ObjectID
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			objectId = &oid
		}

		link, err := GetLink(sampleCode)

		assert.Equal(t, err, nil)
		assert.NotEqual(t, link.ObjectId, nil)
		assert.Equal(t, link.Address, sampleAddress)
		assert.Equal(t, link.Code, sampleCode)
		assert.Equal(t, link.ObjectId, objectId)

		if err := sessionContext.AbortTransaction(context.Background()); err != nil {
			log.Fatal(err)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
