package collect

import (
	"sync"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func FindCollectId(Id int) (CollectContent, error) {
	mg := tools.FindOne(
		config.Database,
		config.CollectContentTable,
		bson.M{"id": Id},
	)
	var data CollectContent
	err := mg.Decode(&data)
	return data, err
}

/*
 * fdas
 */
func FindCollects(userId string) ([]Collect, error) {
	mg := tools.FindOne(
		config.Database,
		config.CollectTable,
		bson.M{"userId": userId},
	)
	var data UserCollect
	err := mg.Decode(&data)
	if err != nil {
		return data.Collects, err
	}
	// result := make([]Collect, 0)
	// for _, x := range data.Collects {
	// 	if (!x.Privacy && !private) || (private) {
	// 		result = append(result, x)
	// 	}
	// }
	return data.Collects, err
}

func FindFavorite(userId string, probId int) (bool, int) {
	colls, err := FindCollects(userId)
	if err != nil {
		return false, -1
	}

	var wg sync.WaitGroup
	result := false
	cId := -1
	for _, collInfo := range colls {
		wg.Add(1)
		go func(wg *sync.WaitGroup, id int) {
			defer wg.Done()
			coll, errC := FindCollectId(id)
			if errC != nil {
				return
			}
			for i := 0; i < len(coll.Probs); i += 1 {
				if coll.Probs[i] == probId {
					result = true
					cId = 1
				}
			}
		}(&wg, collInfo.Id)
	}
	wg.Wait()
	return result, cId
}
