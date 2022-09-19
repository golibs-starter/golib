package pubsub

import "fmt"

type DebugLog func(e Event, msgFormat string, args ...interface{})

var defaultDebugLog DebugLog = func(_ Event, msgFormat string, args ...interface{}) {
	_, _ = fmt.Printf(msgFormat+"\n", args...)
}
