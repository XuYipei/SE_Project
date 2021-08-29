package user

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

var userCount = 0

func generateId() string {
	userCount += 1
	return tools.GenerateMd5(userCount)
}

func Remove(info User) error {
	err := tools.DeleteOne(
		config.Database,
		config.UserTable,
		bson.M{"id": info.Id},
	)
	return err
}

func Init(reset bool, count int) error {
	var err error = nil
	if reset {
		err = tools.InitCollection(config.UserTable)
		if err != nil {
			return err
		}
		err = tools.InitCollection(config.UserContentTable)
		if err != nil {
			return err
		}
	}
	userCount = count
	return err
}
