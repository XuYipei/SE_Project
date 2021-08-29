package user

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func FindOneUser(filter bson.M) (User, error) {
	resmongo := tools.FindOne(
		config.Database,
		config.UserTable,
		filter,
	)
	var result User
	err := resmongo.Decode(&result)
	return result, err
}

func FindOneUserInfo(filter bson.M) (Information, error) {
	resmongo := tools.FindOne(
		config.Database,
		config.UserContentTable,
		filter,
	)
	var result Information
	err := resmongo.Decode(&result)
	return result, err
}
