package pubsub

type Event interface {
	GetName() string
	GetMessage() interface{}
}
