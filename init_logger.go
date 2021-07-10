package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/log"
)

func InitLogger() {
	logger, err := log.NewLogger(&log.Options{
		Development:    true,
		JsonOutputMode: true,
		CallerSkip:     2,
	})
	if err != nil {
		panic(fmt.Sprintf("Error when init logger: [%v]", err))
	}
	log.RegisterGlobal(logger)
}
