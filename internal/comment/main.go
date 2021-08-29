package comment

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
)

func Init(reset bool, count int) error {
	var err error = nil
	if reset {
		err = tools.InitCollection(config.CommentTable)
		if err != nil {
			return err
		}
		err = tools.InitCollection(config.CommentContentTable)
		if err != nil {
			return err
		}
	}
	recmdCount = count
	return err
}
