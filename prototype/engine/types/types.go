package types

type TaskStatus int

const (
	Pending TaskStatus = iota
	Ready
	Running
	Completed
	Failed
)

type Task struct {
	ID        string
	Priority  int // lower = higher priority
	Status    TaskStatus
	Resources []string // paths this task needs exclusive access to
	Deps      []string // task IDs this depends on
}

type ResourceLock struct {
	Path   string
	HeldBy string
}
