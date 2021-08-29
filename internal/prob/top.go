package prob

import (
	"time"

	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
)

type Content struct {
	Id       string   `json:"id" form:"id"`
	Title    string   `json:"title" form:"title"`
	Tags     []string `json:"tags" form:"tags"`
	Type     string   `json:"type" form:"type"`
	Sections []struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Misc    string `json:"misc"`
	} `json:"sections"  form:"sections"`
	Samples []struct {
		Title       string `json:"title" form:"title"`
		Input       string `json:"input" form:"input"`
		Output      string `json:"output" form:"output"`
		Explanation string `json:"explanation" form:"explanation"`
		Misc        string `json:"misc" form:"misc"`
	} `json:"samples" form:"samples"`
	Rules struct {
		Runtime string `json:"runtime" form:"runtime"`
		Memory  string `json:"memory" form:"memory"`
		Stack   string `json:"stack" form:"stack"`
		Source  string `json:"source" form:"source"`
	} `json:"rules" form:"rules"`
	Meta struct {
		Created string `json:"created" form:"created"`
		Updated string `json:"updated" form:"updated"`
		Checked string `json:"checked" form:"checked"`
	} `json:"meta" form:"meta"`
}

type ProbStruct struct {
	Id      string    `json:"id" bson:"id" form:"_id"`
	Index   int       `json:"index" bson:"index" form:"id"`
	Content Content   `json:"content" bson:"content" form:"content"`
	Like    int       `json:"like" bson:"like" form:"like"`
	Dislike int       `json:"dislike" bson:"dislike" form:"dislike"`
	Visit   int       `json:"visit" bson:"visit" form:"visit"`
	UpdTime time.Time `json:"upd_time" bson:"updTime"`
	File    string    `json:"file" bson:"file"`
	Url     string    `json:"url" bson:"url"`
}

func Init(reset bool, count int) error {
	if reset {
		err := tools.InitCollection(config.ProblemTable)
		if err != nil {
			return err
		}
		err = tools.InitCollection(config.ProbContentTable)
		if err != nil {
			return err
		}
		LoadDir(config.ProblemDir)
	}
	return nil
}
