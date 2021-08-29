package state

import (
	"time"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func delete(id string, key string) error {
	err := tools.DeleteOne(
		config.Database,
		config.StateTable,
		bson.M{"id": id, "key": key},
	)
	return err
}

func update(userId string, userKey string, newTime time.Time) error {
	err := tools.UpdateOne(
		config.Database,
		config.StateTable,
		bson.M{"id": userId, "key": userKey},
		bson.M{"$set": bson.M{"act": newTime}},
	)
	return err
}

func Logout(id string, key string) error {
	_, err := tools.DeleteMany(
		config.Database,
		config.StateTable,
		bson.M{"id": id, "key": key},
	)
	return err
}

func Login(id string) (string, error) {
	key := generateKey()
	err := tools.InsertOne(
		config.Database,
		config.StateTable,
		StateStruct{
			Id:  id,
			Key: key,
			Act: time.Now(),
		},
	)
	return key, err
}
