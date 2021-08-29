package comment

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func Upload(comment Comment) (Comment, error) {
	err := tools.UpdateOne(
		config.Database,
		config.HyperTable,
		bson.M{"key": "comment"},
		bson.M{"$inc": bson.M{"value": 1}},
	)
	if err != nil {
		return comment, err
	}

	comment.Id = generateId()
	comment.Content.Id = comment.Id
	err = tools.InsertOne(
		config.Database,
		config.CommentContentTable,
		comment.Content,
	)
	if err != nil {
		return comment, err
	}

	if len(comment.Content.Content) > config.Limits.CommentContent {
		comment.Content.Content = comment.Content.Content[:config.Limits.CommentContent]
	}
	err = tools.InsertOne(
		config.Database,
		config.CommentTable,
		comment,
	)
	return comment, err
}

func Update(comment Comment) (Comment, error) {
	err := tools.UpdateOne(
		config.Database,
		config.CommentContentTable,
		bson.M{"id": comment.Id},
		bson.M{"$set": bson.M{"content": comment.Content}},
	)
	if err != nil {
		return comment, err
	}

	if len(comment.Content.Content) > config.Limits.CommentContent {
		comment.Content.Content = comment.Content.Content[:config.Limits.CommentContent]
	}
	err = tools.UpdateOne(
		config.Database,
		config.CommentTable,
		bson.M{"id": comment.Id},
		bson.M{"$set": bson.M{"content.content": comment.Content.Content}},
	)
	return comment, err
}

func AddResponses(recommendId int, commentId int) error {
	par, err := findOneComment(bson.M{"id": recommendId})
	if err != nil {
		return err
	}

	par.Responses = append(par.Responses, commentId)
	err = tools.UpdateOne(
		config.Database,
		config.CommentTable,
		bson.M{"id": recommendId},
		bson.M{"$set": bson.M{"responses": par.Responses}},
	)
	return err
}

// func Score(recommend RecmdStruct) (RecmdStruct, error) {
// 	err := tools.UpdateOne(
// 		config.Database,
// 		config.RecommendTable,
// 		bson.M{"recmd_id": recommend.Id},
// 		bson.M{"$set": bson.M{"score": recommend.Score}},
// 	)
// 	return recommend, err
// }
