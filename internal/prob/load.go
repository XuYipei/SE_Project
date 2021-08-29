package prob

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func LoadFile(path string, contentOnly bool) error {
	clientp, err := tools.ConnectDb()
	if err != nil {
		return err
	}
	collectionp := clientp.Database(config.Database).Collection(config.ProblemTable)
	defer clientp.Disconnect(context.Background())

	clientc, err := tools.ConnectDb()
	if err != nil {
		return err
	}
	collectionc := clientc.Database(config.Database).Collection(config.ProbContentTable)
	defer clientc.Disconnect(context.Background())

	_, filename := filepath.Split(path)
	bytes, errFile := os.ReadFile(path)
	id := filename[0 : len(filename)-len(filepath.Ext(filename))]
	if errFile != nil {
		return errFile
	}

	var cont Content
	json.Unmarshal(bytes, &cont)
	cont.Id = id
	if contentOnly {
		_, errUpd := collectionc.UpdateOne(context.Background(), bson.M{"id": cont.Id}, bson.M{"$set": cont})
		if errUpd != nil {
			return errUpd
		}
	} else {
		_, errInsert := collectionc.InsertOne(context.Background(), cont)
		if errInsert != nil {
			return errInsert
		}
	}

	if len(cont.Sections) > config.Limits.ProbSections {
		cont.Sections = cont.Sections[:config.Limits.ProbSections]
	}
	for idx, v := range cont.Sections {
		if len(v.Content) > config.Limits.ProbContent {
			v.Content = v.Content[:config.Limits.ProbContent]
			cont.Sections[idx] = v
		}
	}
	cont.Samples = cont.Samples[:0]
	index, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}

	probInfo := ProbStruct{
		Id:      id,
		Content: cont,
		Index:   index,
		Like:    (rand.Int() % 40) + 20,
		Dislike: rand.Int() % 10,
		Visit:   0,
		UpdTime: time.Now(),
		File:    path,
	}

	if contentOnly {
		_, errUpd := collectionp.UpdateOne(
			context.Background(),
			bson.M{"id": cont.Id},
			bson.M{"$set": bson.M{"content": cont, "updTime": probInfo.UpdTime}},
		)
		if errUpd != nil {
			return errUpd
		}
	} else {
		_, errInsert := collectionp.InsertOne(context.Background(), probInfo)
		if errInsert != nil {
			return errInsert
		}
	}

	return nil
}

func LoadDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("load problem start: directory is %s\n", dir)

	var loads int32
	var wg sync.WaitGroup

	for offset := 0; offset < config.LoadProblemWorkers; offset += 1 {
		wg.Add(1)
		go func(offset int) {
			for i := offset; i < len(files); i += config.LoadProblemWorkers {
				atomic.AddInt32(&loads, 1)
				if filepath.Ext(files[i].Name()) != ".json" {
					continue
				}
				errOne := LoadFile(filepath.Join(dir, files[i].Name()), false)
				if errOne != nil {
					err = errOne
				}
			}
			wg.Done()
		}(offset)
	}

	wg.Wait()

	log.Printf("\nload problem done.\n")
	return nil
}
