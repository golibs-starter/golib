package pubsub

import "fmt"

type DebugLog func(msgFormat string, args ...interface{})

var defaultDebugLog = func(msgFormat string, args ...interface{}) {
	_, _ = fmt.Printf(msgFormat+"\n", args...)
}
