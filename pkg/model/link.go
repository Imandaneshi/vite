package model

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/imandaneshi/vite/pkg/config"
	"github.com/imandaneshi/vite/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type webPage struct {
	Title       string
	Keywords    []string
	Author      string
	Description string
	ThemeColor  string
	ImageUrl    string
}

type Link struct {
	ObjectId *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Address  string              `bson:"address,omitempty" json:"address"`
	Code     string              `bson:"code,omitempty" json:"code"`
	User     *primitive.ObjectID `bson:"userId,omitempty" json:"-"`
	WebPage  *webPage            `bson:"webPage,omitempty" json:"webPage"`
}

func (link *Link) Delete() error {
	links := m.Collection(mongoLinksCollection)
	deleteResult, deleteError := links.DeleteOne(context.TODO(), bson.M{"_id": link.ObjectId})
	if deleteError != nil {
		return errors.New("failedDeletingLink", "failed while deleting link from database", deleteError)
	}
	if deleteResult.DeletedCount != 1 {
		return errors.New("failedDeletingLink", "failed while deleting link from database", nil)
	}
	return nil
}

func (link *Link) Scrape() error {

	log.Debugf("scraping page %s", link.Address)

	var webPage webPage
	res, err := http.Get(link.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		message := fmt.Sprintf("failed scraping link %s: status code error: %d %s",
			link.Address, res.StatusCode, res.Status)
		log.Infof(message)
		return errors.ConnectionError(message, nil)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("failed reading html body while scraping %s: %", link.Address, err)
		return err
	}

	webPage.Title = doc.Find("title").Last().Text()

	// set meta tags with names and properties
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {

		name, _ := s.Attr("name")
		switch name {
		case "description":
			webPage.Description, _ = s.Attr("content")
		case "keywords":
			keywordsString, _ := s.Attr("content")
			webPage.Keywords = strings.Split(keywordsString, ",")
		case "author":
			webPage.Author, _ = s.Attr("content")
		case "theme-color":
			webPage.ThemeColor, _ = s.Attr("content")
		}

		property, _ := s.Attr("property")
		switch property {
		case "og:image":
			webPage.ImageUrl, _ = s.Attr("content")
		}

	})

	link.WebPage = &webPage

	links := m.Collection(mongoLinksCollection)

	updateResult, updateError := links.UpdateOne(context.Background(),
		bson.M{"_id": link.ObjectId}, bson.M{"$set": bson.M{"webPage": webPage}})

	if updateError != nil {
		log.Error("failed updating link %s: %", link.Address, updateError)
		return updateError
	}

	if updateResult.ModifiedCount < 1 {
		return errors.New("FailedUpdatingLink", "failed updating link", nil)
	}

	log.Infof("successfully scraped page %s", webPage)

	return nil
}

// TODO: move to Create and refactor
func GenerateRandomShortenLink(address string, userId *primitive.ObjectID) (*Link, error) {
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

	var linkId *primitive.ObjectID
	logFields := log.Fields{
		"code":    randomCode,
		"address": address,
		"ID":      linkId,
	}
	log.WithFields(logFields).Debugf("inserting new link into mongo db: %s | %s", randomCode, address)
	link := &Link{
		Address: address,
		Code:    randomCode,
		User:    userId,
	}
	links := m.Collection(mongoLinksCollection)
	res, err := links.InsertOne(context.Background(), link)
	if err != nil {
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
	code := GenerateRandomString(n, "abcdefghijklmnopqrstuvwxyz1234567890")
	if validate {
		links := m.Collection(mongoLinksCollection)
		err := links.FindOne(context.Background(), bson.M{"code": code})
		if err == nil {
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
	mongoLinksCodeIndex  string = "uniqueCodeIndex"
)
