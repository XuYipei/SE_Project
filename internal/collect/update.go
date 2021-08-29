package collect

import (
	"strconv"
	"sync"
	"time"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"github.com/tuplz/tuplz-be/internal/prob"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadUser(userId string) error {
	err := tools.UpdateOne(
		config.Database,
		config.HyperTable,
		bson.M{"key": "collect"},
		bson.M{"$inc": bson.M{"value": 1}},
	)
	if err != nil {
		return err
	}

	col := CollectContent{
		UserId:  userId,
		Name:    "默认收藏",
		Id:      generateId(),
		UpdTime: time.Now(),
	}
	err = tools.InsertOne(config.Database, config.CollectContentTable, col)
	if err != nil {
		return err
	}

	data := UserCollect{
		UserId: userId,
		Collects: []Collect{
			Collect{
				Id:      col.Id,
				Name:    "默认收藏",
				UpdTime: time.Now(),
			},
		},
	}
	err = tools.InsertOne(config.Database, config.CollectTable, data)
	return err
}

func UploadCollect(id int, probId int) error {
	var wg sync.WaitGroup
	var errD, errC error
	var data CollectContent
	var probSt prob.ProbStruct

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		data, errD = findOneCollectContent(id)
	}(&wg)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		probSt, errC = prob.FindOneProb(bson.M{"id": strconv.Itoa(probId)})
	}(&wg)
	wg.Wait()

	if errD != nil {
		return errD
	}
	if errC != nil {
		return errD
	}

	var uData UserCollect
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		data.Probs = append(data.Probs, probId)
		data.ProbCs = append(data.ProbCs, probSt)
		errC = tools.UpdateOne(
			config.Database,
			config.CollectContentTable,
			bson.M{"id": id},
			bson.M{"$set": data},
		)
	}(&wg)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		uData, errD = findOneUserCollect(data.UserId)
		if errD != nil {
			return
		}
		for i := 0; i < len(uData.Collects); i += 1 {
			if uData.Collects[i].Id == id {
				uData.Collects[i].Count += 1
			}
		}
		errD = tools.UpdateOne(
			config.Database,
			config.CollectTable,
			bson.M{"userId": data.UserId},
			bson.M{"$set": uData},
		)
	}(&wg)
	wg.Wait()

	if errC != nil {
		return errC
	}
	return errD
}

func UploadCollects(userId string, name string, disc string) error {
	err := tools.UpdateOne(
		config.Database,
		config.HyperTable,
		bson.M{"key": "collect"},
		bson.M{"$inc": bson.M{"value": 1}},
	)
	if err != nil {
		return err
	}

	UpdTime := time.Now()
	col := CollectContent{
		UserId:  userId,
		Id:      generateId(),
		UpdTime: UpdTime,
		Name:    name,
		Disc:    disc,
	}
	err = tools.InsertOne(
		config.Database,
		config.CollectContentTable,
		col,
	)
	if err != nil {
		return err
	}

	uData, err := findOneUserCollect(userId)
	if err != nil {
		return err
	}

	data := Collect{
		Id:      col.Id,
		Name:    name,
		UpdTime: UpdTime,
	}
	uData.Collects = append(uData.Collects, data)

	err = tools.UpdateOne(
		config.Database,
		config.CollectTable,
		bson.M{"userId": userId},
		bson.M{"$set": uData},
	)
	return err
}

func UpdateCollects(userId string, id int, name string, disc string) error {
	data, err := findOneCollectContent(id)
	if err != nil {
		return err
	}
	err = tools.UpdateOne(
		config.Database,
		config.CollectContentTable,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"name": name, "discription": disc}},
	)
	if err != nil {
		return err
	}

	uData, err := findOneUserCollect(data.UserId)
	if err != nil {
		return err
	}
	// log.print(uData)

	for i, x := range uData.Collects {
		if x.Id == id {
			uData.Collects[i].Name = name
			uData.Collects[i].Disc = disc
		}
	}
	err = tools.UpdateOne(
		config.Database,
		config.CollectTable,
		bson.M{"userId": data.UserId},
		bson.M{"$set": uData},
	)
	return err
}

func DeleteCollects(id int) error {
	data, err := findOneCollectContent(id)
	if err != nil {
		return err
	}
	err = tools.DeleteOne(
		config.Database,
		config.CollectContentTable,
		bson.M{"id": id},
	)
	if err != nil {
		return err
	}

	uData, err := findOneUserCollect(data.UserId)
	if err != nil {
		return err
	}

	colls := uData.Collects
	uData.Collects = make([]Collect, 0)
	for _, x := range colls {
		if x.Id != id {
			uData.Collects = append(uData.Collects, x)
		}
	}
	err = tools.UpdateOne(
		config.Database,
		config.CollectTable,
		bson.M{"userId": data.UserId},
		bson.M{"$set": uData},
	)
	return err
}

func DeleteCollect(userId string, id int, probId int) error {
	data, err := findOneCollectContent(id)
	if err != nil {
		return err
	}
	pColls := data.Probs
	data.Probs = make([]int, 0)
	for _, x := range pColls {
		if x != probId {
			data.Probs = append(data.Probs, x)
		}
	}
	pCColls := data.ProbCs
	data.ProbCs = make([]prob.ProbStruct, 0)
	for _, x := range pCColls {
		if x.Index != probId {
			data.ProbCs = append(data.ProbCs, x)
		}
	}

	err = tools.UpdateOne(
		config.Database,
		config.CollectContentTable,
		bson.M{"id": id},
		bson.M{"$set": data},
	)
	if err != nil {
		return err
	}

	uData, err := findOneUserCollect(data.UserId)
	if err != nil {
		return err
	}

	for i := 0; i < len(uData.Collects); i += 1 {
		if uData.Collects[i].Id == id {
			uData.Collects[i].Count -= 1
		}
	}
	err = tools.UpdateOne(
		config.Database,
		config.CollectTable,
		bson.M{"userId": data.UserId},
		bson.M{"$set": uData},
	)
	return err
}
