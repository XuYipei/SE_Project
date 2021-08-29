package comment

import (
	"context"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func curComments(cur *mongo.Cursor, maxLength int) []Comment {
	result := make([]Comment, 0)
	len := 0
	for cur.Next(context.TODO()) && (len < maxLength || maxLength == 0) {
		var recmd Comment
		cur.Decode(&recmd)
		result = append(result, recmd)
		len += 1
	}
	return result
}

func generateId() int {
	recmdCount += 1
	// return tools.GenerateMd5(recmdCount)
	return recmdCount
}

func findOne(table string, filter bson.M) *mongo.SingleResult {
	return tools.FindOne(
		config.Database,
		table,
		filter,
	)
}

func findOneComment(filter bson.M) (Comment, error) {
	resmongo := findOne(config.CommentTable, filter)
	var result Comment
	err := resmongo.Decode(&result)
	return result, err
}

func findOneContent(filter bson.M) (Content, error) {
	resmongo := findOne(config.CommentContentTable, filter)
	var result Content
	err := resmongo.Decode(&result)
	return result, err
}
