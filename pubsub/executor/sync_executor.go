package executor

type SyncExecutor struct {
}

func NewSyncExecutor() *SyncExecutor {
	return &SyncExecutor{}
}

func (s SyncExecutor) Execute(fn func()) {
	fn()
}
