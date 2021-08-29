package prob

import (
	"context"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func extractCur(cur *mongo.Cursor, maxLength int) []ProbStruct {
	result := make([]ProbStruct, 0)
	len := 0
	for cur.Next(context.TODO()) && (len < maxLength || maxLength == 0) {
		var recmd ProbStruct
		cur.Decode(&recmd)
		result = append(result, recmd)
		len += 1
	}
	return result
}

func FindOneProb(filter bson.M) (ProbStruct, error) {
	resmongo := tools.FindOne(
		config.Database,
		config.ProblemTable,
		filter,
	)
	var result ProbStruct
	err := resmongo.Decode(&result)
	return result, err
}

func FindOneCont(filter bson.M) (Content, error) {
	resmongo := tools.FindOne(
		config.Database,
		config.ProbContentTable,
		filter,
	)
	var result Content
	err := resmongo.Decode(&result)

	return result, err
}

// type sortVisitList []ProbStruct

// func (s sortVisitList) Less(i, j int) bool {
// 	return s[i].Visit > s[j].Visit
// }
// func (s sortVisitList) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }
// func (s sortVisitList) Len() int {
// 	return len(s)
// }
// func sortVisit(a []ProbStruct) {
// 	var b sortVisitList = a
// 	sort.Sort(b)
// 	a = b
// }
