package state

import (
	"time"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
)

var (
	TimeLimit time.Duration
)

type StateStruct struct {
	Id  string    `bson:"id"`
	Key string    `bson:"key"`
	Act time.Time `bson:"act"`
}

var stateCount = 1

func Init() error {
	TimeLimit = config.Limits.UOnline
	err := tools.InitCollection(config.StateTable)
	return err
}
