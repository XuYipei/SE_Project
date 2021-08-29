package user

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func Update(user User, info Information) error {
	err := tools.UpdateOne(
		config.Database,
		config.UserTable,
		bson.M{"id": user.Id},
		bson.M{"$set": bson.M{"name": user.Name, "email": user.Email}},
	)
	if err != nil {
		return err
	}
	return tools.UpdateOne(
		config.Database,
		config.UserContentTable,
		bson.M{"id": user.Id},
		bson.M{"$set": bson.M{
			"photo":       info.Photo,
			"description": info.Disc,
			"exp":         info.Exp,
			"gender":      info.Gender}},
	)
}

func Act(userId string) error {
	return tools.UpdateOne(
		config.Database,
		config.UserContentTable,
		bson.M{"id": userId},
		bson.M{"$inc": bson.M{"exp": 1}},
	)
}
