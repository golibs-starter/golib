package log

import (
	"gitlab.id.vin/vincart/golib/log"
)

var global log.Logger

// ReplaceGlobal Register a logger instance as global
func ReplaceGlobal(logger log.Logger) {
	global = logger
}
