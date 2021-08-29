package user

import (
	"time"
)

type Verify struct {
	State bool      `bson:"state"`
	Code  string    `bson:"code"`
	Time  time.Time `bson:"time"`
}

type User struct {
	Id       string `bson:"id"`
	Password string `bson:"password"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Verify   Verify `bson:"verify"`
}

type Information struct {
	Id     string `bson:"id"`
	Photo  string `bson:"photo"`
	Disc   string `bson:"description"`
	Exp    int    `bson:"exp"`
	Gender bool   `bson:"gender"`
}
