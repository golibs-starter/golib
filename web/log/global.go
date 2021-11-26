package log

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/log"
)

var global log.Logger

func init() {
	var err error
	if global, err = log.NewLogger(&log.Options{}); err != nil {
		panic(fmt.Errorf("init global web logger error [%v]", err))
	}
}

// ReplaceGlobal Register a logger instance as global
func ReplaceGlobal(logger log.Logger) {
	global = logger
}
