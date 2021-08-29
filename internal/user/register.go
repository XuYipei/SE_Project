package user

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(user User) (User, error) {
	_, err := FindOneUser(bson.M{"name": user.Name})
	if err == nil {
		return user, errors.New("username: name already exists")
	}

	err = tools.UpdateOne(
		config.Database,
		config.HyperTable,
		bson.M{"key": "user"},
		bson.M{"$inc": bson.M{"value": 1}},
	)
	if err != nil {
		return user, err
	}

	user.Id = generateId()
	user.Verify = Verify{
		State: false,
		Code:  fmt.Sprintf("%d", rand.Int()%10000),
		Time:  time.Now(),
	}

	err = tools.InsertOne(
		config.Database,
		config.UserTable,
		user,
	)
	if err != nil {
		return user, err
	}

	return user, tools.InsertOne(
		config.Database,
		config.UserContentTable,
		Information{Id: user.Id},
	)
}

func SendEmail(Id string) (User, error) {
	target, err := FindOneUser(
		bson.M{"id": Id},
	)
	if err != nil {
		return User{Id: Id}, err
	}

	rand.Seed(time.Now().UnixNano())
	target.Verify = Verify{
		State: false,
		Code:  fmt.Sprintf("%d", rand.Int()%10000),
		Time:  time.Now(),
	}

	tools.SendEmail(
		[]string{target.Email},
		"<h1>Your verification code is: "+target.Verify.Code+"</h1>",
	)
	err = tools.UpdateOne(
		config.Database,
		config.UserTable,
		bson.M{"id": target.Id},
		bson.M{"$set": bson.M{"verify": target.Verify}},
	)
	return target, err
}

func VerifyUser(Id string, Code string) (User, error) {
	target, err := FindOneUser(
		bson.M{"id": Id},
	)
	if err != nil {
		return target, err
	}

	if target.Verify.State {
		return target, errors.New("emailVerify: already verified")
	}
	if Code != target.Verify.Code {
		return target, errors.New("emailVerify: wrond verification code")
	}
	if target.Verify.Time.Add(config.Limits.VerifyTime).Before(time.Now()) {
		return target, errors.New("emailVerify: time limit exceed")
	}

	target.Verify.State = true
	err = tools.UpdateOne(
		config.Database,
		config.UserTable,
		bson.M{"id": target.Id},
		bson.M{"$set": bson.M{"verify": target.Verify}},
	)

	// log.print(err)

	return target, err
}
