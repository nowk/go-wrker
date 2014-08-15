package wrker

// WorkPool

type WorkPool struct {
	Workers []*Worker
}

// NewPool returns a WorkPool connected to a collection of workers
func NewPool(num int) (pool *WorkPool) {
	pool = &WorkPool{}

	for i := 0; i < num; i++ {
		w := NewWorker(i)
		pool.Add(w)
	}

	return
}

// Add adds a worker to the pool
func (p *WorkPool) Add(worker *Worker) {
	p.Workers = append(p.Workers, worker)
}
