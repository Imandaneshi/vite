package model

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ObjectId  *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string              `bson:"username,omitempty" json:"username"`
	Email     string              `bson:"email" json:"email"`
	FirstName string              `bson:"firstName" json:"firstName"`
	LastName  string              `bson:"lastName" json:"lastName"`
	Password  string              `bson:"password,omitempty" json:"-"`
	Token     *Token              `bson:"-" json:"token"`
}

func (user *User) Create() error {
	logrus.Debug("adding user to database")

	if user.ObjectId != nil {

		logrus.WithFields(log.Fields{"user": user}).Info("user is already registered in database")

		return &errors.Error{Code: "alreadyCreated",
			Message: "This user object is already in database"}
	}

	// check if user with such username exists in database or not
	users := m.Collection(mongoUsersCollection)
	err := users.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&User{})
	if err == nil {
		log.WithFields(log.Fields{"user": user}).Info("user with such username already exists in database")
		return errors.AlreadyExistsError("User with this username already exists", nil)
	}

	hashPassword, hashingError := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashingError != nil {
		log.WithFields(log.Fields{"user": user}).Info("failed in hashing password")
		return &errors.Error{Code: "hashingFailed",
			Message: "Failed hashing user password"}
	}
	user.Password = string(hashPassword)

	res, insertErr := users.InsertOne(context.Background(), user)
	if insertErr != nil {
		log.WithFields(log.Fields{"user": user}).Info("failed inserting new user into mongo db")
		return errors.New("insertingUserFailed", "Failed inserting new user", insertErr)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ObjectId = &oid
	}
	log.WithFields(log.Fields{"user": user}).Info("successfully inserted new user into mongo db")
	return nil
}


func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func getUser(filters *bson.M) (*User, error) {
	log.WithFields(log.Fields{"filters": filters}).Debug("looking if user with such filters exists")

	var user User
	users := m.Collection(mongoUsersCollection)

	err := users.FindOne(context.Background(), filters).Decode(&user)
	if err != nil {
		log.WithFields(log.Fields{"filters": filters}).Info("user with such filters doesn't exists in database")
		return nil, errors.NotFoundError("user with this id doesn't exists")
	}

	log.WithFields(log.Fields{"user": user, "filters": filters}).Info("successfully found user with this ID")
	return &user, nil
}

func GetUserById(userId *primitive.ObjectID) (user *User, err error) {
	user, err = getUser(&bson.M{"_id": userId})
	return
}

func GetUserByUsername(username string) (user *User, err error) {
	user, err = getUser(&bson.M{"username": username})
	return
}

func GetUserFromGinContext(c *gin.Context) (*User, bool) {
	_user, exists := c.Get("user")

	if exists {
		user, ok  := _user.(*User)
		if ok {
			return user, true
		}
	}

	return nil, false
}

const (
	mongoUsersCollection    string = "users"
	mongoUsersUsernameIndex string = "uniqueUsernameIndex"
)