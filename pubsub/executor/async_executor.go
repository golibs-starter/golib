package executor

type AsyncExecutor struct {
}

func NewAsyncExecutor() *AsyncExecutor {
	return &AsyncExecutor{}
}

func (a AsyncExecutor) Execute(fn func()) {
	go fn()
}
