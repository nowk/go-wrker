package wrker

// WorkPool

type WorkPool struct {
	Workers []*Worker
	Queue   chan chan Job
	Errors  chan error
}

// NewPool returns a WorkPool connected to a collection of workers
func NewPool(num int) (pool *WorkPool) {
	pool = &WorkPool{
		Queue:  make(chan chan Job, num),
		Errors: make(chan error),
	}

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

// Dispatch starts the workers and routes jobs to available workers
func (p *WorkPool) Dispatch(jobs chan Job) chan error {
	for _, w := range p.Workers {
		w.Start(p.Queue, p.Errors)
	}

	go p.route(jobs)

	return p.Errors
}

// route takes a channel of jobs and passes them to an available worker in the
// queue
func (p *WorkPool) route(jobs <-chan Job) {
	for {
		select {
		case job := <-jobs:
			go func() {
				worker := <-p.Queue // receive available worker from queue
				worker <- job       // send the worker the job
			}()
		}
	}
}
