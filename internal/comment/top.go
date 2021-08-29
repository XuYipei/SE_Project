package comment

import (
	"time"
)

type Content struct {
	Id       int       `bson:"id"`
	UserId   string    `bson:"userId" json:"userId" form:"userId"`
	UserName string    `bson:"userName" json:"userName" form:"userName"`
	Content  string    `bson:"content" json:"content" form:"content"`
	Like     int       `bson:"like" json:"like" form:"like"`
	Dislike  int       `bson:"dislike" json:"dislike" form:"dislike"`
	Time     time.Time `bson:"time" json:"time" form:"time"`
}

type Comment struct {
	Id        int     `bson:"id"`
	ProbId    int     `bson:"probId"`
	Content   Content `bson:"content"`
	Responses []int   `bson:"responses"`
	Reply     int     `bson:"reply"`
	Recommend int     `bson:"recommend"`
	// Score   float64   `bson:"score"`
}

var recmdCount int = 0
