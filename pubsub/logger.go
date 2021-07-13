package pubsub

type Logger interface {
	Debugf(msgFormat string, args ...interface{})
}
