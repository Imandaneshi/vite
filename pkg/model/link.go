package model

import (
	"context"
	"github.com/imandaneshi/vite/pkg/config"
	"github.com/imandaneshi/vite/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TODO: use camelcase database names
type Link struct {
	ObjectId *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Address  string              `bson:"address,omitempty" json:"address"`
	Code     string              `bson:"code,omitempty" json:"code"`
	User     string              `bson:"userId,omitempty" json:"-"`
}

// TODO: move to Create and refactor
func GenerateRandomShortenLink(address string, userId string) (*Link, error) {
	// get a unique code
	log.Debug("generating a random code and validating for a duplicate")
	randomCode, exists := getRandomCode(config.Server.RandomCodeLength, true)
	if exists {
		log.Debug("getting a random code for a second time")
		randomCode, exists = getRandomCode(config.Server.RandomCodeLength, true)
		if exists {
			log.Info("failed in getting a new random code for a second time")
			return nil, errors.AlreadyExistsError("link with this code already exists", nil)
		}
		log.Info("successfully got a new random code for a second time")
	}
	log.Info("successfully got a new random code")

	var  linkId *primitive.ObjectID
	logFields := log.Fields{
		"code": randomCode,
		"address": address,
		"ID": linkId,
	}
	log.WithFields(logFields).Debugf("inserting new link into mongo db: %s | %s", randomCode, address)
	link := &Link{
		Address: address,
		Code: randomCode,
		User: userId,
	}
	links := m.Collection(mongoLinksCollection)
	res, err := links.InsertOne(context.Background(), link)
	if err != nil{
		log.WithFields(logFields).Info("failed inserting new link into mongo db")
		return nil, errors.New("insertingLinkFailed", "Failed inserting new link", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		link.ObjectId = &oid
		linkId = &oid
	}
	log.WithFields(logFields).Info("successfully inserted new link into mongo db")
	return link, nil
}

func getRandomCode(n int, validate bool) (string, bool) {
	code := GenerateRandomString(n,"abcdefghijklmnopqrstuvwxyz1234567890")
	if validate {
		links := m.Collection(mongoLinksCollection)
		err := links.FindOne(context.Background(), bson.M{"code": code})
		if err == nil{
			return "", true
		}
	}
	return code, false
}

func GetLink(code string) (*Link, error) {
	var result Link

	links := m.Collection(mongoLinksCollection)
	err := links.FindOne(context.TODO(), bson.D{{"code", code}}).Decode(&result)
	if err != nil {
		return nil, errors.NotFoundError("Link not found")
	}
	return &result, nil
}

func GenerateRandomString(length int, characters string) string {
	letterRunes := []rune(characters)
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

const (
	mongoLinksCollection string = "links"
	mongoLinksCodeIndex string = "uniqueCodeIndex"
)