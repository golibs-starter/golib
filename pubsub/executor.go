package pubsub

type Executor interface {

	// Execute a function
	Execute(fn func())
}
