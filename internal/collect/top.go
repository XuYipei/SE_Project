package collect

import (
	"time"

	"github.com/tuplz/tuplz-be/internal/prob"
)

type CollectContent struct {
	Id      int               `bson:"id" json:"id" form:"id"`
	Disc    string            `bson:"discription" json:"discription" form:"discription"`
	Name    string            `bson:"name" json:"name" form:"name"`
	Probs   []int             `bson:"probs" json:"probs" form:"probs"`
	ProbCs  []prob.ProbStruct `bson:"probCs" json:"probCs" form:"probCs"`
	Privacy bool              `bson:"privacy" json:"privacy" form:"privacy"`
	UpdTime time.Time         `bson:"updTime" json:"updTime" form:"updTime"`
	UserId  string            `bson:"usreId" json:"userId" form:"userId"`
}

type Collect struct {
	Id      int       `bson:"id"`
	Name    string    `bson:"name"`
	UpdTime time.Time `bson:"updTime"`
	Count   int       `bson:"count"`
	Disc    string    `bson:"discription"`
	Privacy bool      `bson:"privacy" json:"privacy" form:"privacy"`
}
type UserCollect struct {
	UserId   string    `bson:"userId" json:"userId" form:"userId"`
	Collects []Collect `bson:"collects" json:"collects" form:"collects"`
	Follows  []Collect `bson:"follows" json:"follows" form:"follows"`
}

var collectCount int
