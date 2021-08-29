package monitor

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
)

func FindCollStatus(collName string) (CollInfo, error) {
	var result CollInfo
	mg, err := tools.CollStatus(config.Database, collName)
	if err != nil {
		return CollInfo{}, err
	}
	// var doc bson.M
	// mg.Decode(&doc)
	err = mg.Decode(&result)
	return result, err
}

func FindDbStatus() ([]CollInfo, []string, error) {
	result := make([]CollInfo, 6)
	colls := []string{
		config.UserTable,
		config.UserContentTable,
		config.CommentTable,
		config.CommentContentTable,
		config.ProblemTable,
		config.ProbContentTable,
	}

	for i, collName := range colls {
		collInfo, err := FindCollStatus(collName)
		if err != nil {
			return result, colls, err
		}
		result[i] = collInfo
	}
	return result, colls, nil
}
