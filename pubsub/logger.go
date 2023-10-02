package pubsub

import (
	"context"
	"fmt"
)

type DebugLog func(ctx context.Context, msgFormat string, args ...interface{})

var defaultDebugLog DebugLog = func(_ context.Context, msgFormat string, args ...interface{}) {
	_, _ = fmt.Printf(msgFormat+"\n", args...)
}
