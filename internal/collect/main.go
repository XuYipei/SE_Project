package collect

import (
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
)

func Init(reset bool, count int) error {
	var err error = nil
	if reset {
		err = tools.InitCollection(config.CollectTable)
		if err != nil {
			return err
		}
		err = tools.InitCollection(config.CollectContentTable)
		if err != nil {
			return err
		}
	}
	collectCount = count
	return err
}
