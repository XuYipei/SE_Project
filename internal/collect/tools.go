package collect

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func generateId() int {
	collectCount += 1
	// return tools.GenerateMd5(recmdCount)
	return collectCount
}

func findOneCollectContent(id int) (CollectContent, error) {
	mg := tools.FindOne(
		config.Database,
		config.CollectContentTable,
		bson.M{"id": id},
	)
	var data CollectContent
	err := mg.Decode(&data)
	return data, err
}

func findOneUserCollect(userId string) (UserCollect, error) {
	mg := tools.FindOne(
		config.Database,
		config.CollectTable,
		bson.M{"userId": userId},
	)
	var data UserCollect
	err := mg.Decode(&data)
	return data, err
}
