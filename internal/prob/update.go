package prob

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func Update(prob ProbStruct) error {
	err := tools.UpdateOne(
		config.Database,
		config.ProblemTable,
		bson.M{"id": prob.Id},
		prob,
	)
	return err
}

func Visit(id string) error {
	err := tools.UpdateOne(
		config.Database,
		config.ProblemTable,
		bson.M{"id": id},
		bson.M{"$inc": bson.M{"visit": 1}},
	)
	return err
}

func UpdateLike(id string, upd int) error {
	err := tools.UpdateOne(
		config.Database,
		config.ProblemTable,
		bson.M{"id": id},
		bson.M{"$inc": bson.M{"like": upd}},
	)
	return err
}
