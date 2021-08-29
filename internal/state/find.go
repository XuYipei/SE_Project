package state

import (
	"context"
	"time"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func find(id string, key string, upd bool) (bool, error) {
	cur, err := tools.FindMany(
		config.Database,
		config.StateTable,
		bson.M{"id": id},
	)
	if err != nil {
		return false, err
	}
	defer cur.Close(context.Background())

	nowTime := time.Now()
	result := extractCur(cur)

	// log.print(result)

	for i := 0; i < len(result); i++ {
		if result[i].Act.Add(TimeLimit).Before(nowTime) {
			delete(result[i].Id, result[i].Key)
		}
	}
	for i := 0; i < len(result); i++ {
		if result[i].Act.Add(TimeLimit).After(nowTime) && result[i].Key == key {
			if upd {
				update(id, key, nowTime)
			}
			return true, nil
		}
	}
	return false, nil
}

func Find(userId string, userKey string) (bool, error) {
	result, err := find(userId, userKey, false)
	return result, err
}
func FindUpdate(userId string, userKey string) (bool, error) {
	result, err := find(userId, userKey, true)
	return result, err
}
