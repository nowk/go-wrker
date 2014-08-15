package wrker

// WorkPool

type WorkPool struct {
	Workers []*Worker
}

// NewPool returns a WorkPool connected to a collection of workers
func NewPool(num int) (pool *WorkPool) {
	pool = &WorkPool{}

	for i := 0; i < num; i++ {
		wkr := NewWorker(i)
		pool.Workers = append(pool.Workers, wkr)
	}

	return
}
