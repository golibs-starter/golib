package example

type SampleEventExecutor struct {
	// Such as using Worker Pool.
	// pool *worker.Pool
}

func NewSampleEventExecutor() *SampleEventExecutor {
	return &SampleEventExecutor{
		// pool: pool
	}
}

func (s SampleEventExecutor) Execute(fn func()) {
	// pool.Submit(fn)
}
