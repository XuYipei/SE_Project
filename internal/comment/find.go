package comment

import (
	"context"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func FindRecommendsUserId(userId string, maxLength int) ([]Comment, error) {
	result := make([]Comment, 0)
	cur, err := tools.FindMany(
		config.Database,
		config.CommentTable,
		bson.M{"userId": userId},
	)
	if err != nil {
		return result, err
	}
	defer cur.Close(context.Background())

	result = curComments(cur, maxLength)
	return result, nil
}

func FindRecommendsProbId(probId int, maxLength int) ([]Comment, error) {
	result := make([]Comment, 0)
	cur, err := tools.FindMany(
		config.Database,
		config.CommentTable,
		bson.M{"probId": probId},
	)
	if err != nil {
		return result, err
	}
	defer cur.Close(context.Background())

	result = curComments(cur, maxLength)
	return result, nil
}

func FindCommentId(Id int) (Comment, error) {
	result, err := findOneComment(bson.M{"id": Id})
	if err != nil {
		return result, err
	}
	cont, err := findOneContent(bson.M{"id": Id})
	result.Content = cont

	return result, err
}
