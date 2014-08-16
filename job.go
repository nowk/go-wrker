package wrker

// Job is an interface for objects to be passed to a worker
type Job interface {
	Do() error
}
