package prob

import (
	"context"
	"sync"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindId(id string) (ProbStruct, error) {
	var err0, err1 error
	var wg sync.WaitGroup
	var prob ProbStruct
	var cont Content

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		prob, err0 = FindOneProb(
			bson.M{"id": id},
		)
	}(&wg)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		cont, err1 = FindOneCont(
			bson.M{"id": id},
		)
	}(&wg)

	wg.Wait()
	prob.Content = cont
	if err0 != nil {
		return prob, err0
	}
	return prob, err1
}

// todo: api update
func FindIdGt(index int, maxLength int) ([]ProbStruct, error) {
	result := make([]ProbStruct, 0)

	idxSort := options.Find()
	idxSort.SetSort(map[string]int{"index": 1})
	idxSort.SetLimit(int64(maxLength))
	cur, errFind := tools.FindManyOpt(
		config.Database,
		config.ProblemTable,
		bson.M{"index": bson.M{"$gt": index}},
		idxSort,
	)
	if errFind != nil {
		return result, errFind
	}
	defer cur.Close(context.TODO())

	result = extractCur(cur, maxLength)
	// sortVisit(result)
	return result, nil
}

func FindProbs(maxLength int) ([]ProbStruct, error) {
	return FindIdGt(-1, maxLength)
}
