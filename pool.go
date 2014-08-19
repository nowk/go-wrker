package wrker

import "sync"

// WorkPool

type WorkPool struct {
	Workers []*Worker
	Queue   chan chan Job
	Errors  chan error
	stop    chan chan struct{}
	sync    *sync.WaitGroup
}

// NewPool returns a WorkPool connected to a collection of workers
func NewPool(num int) (pool *WorkPool) {
	pool = &WorkPool{
		Queue:  make(chan chan Job, num),
		Errors: make(chan error),
		stop:   make(chan chan struct{}),
		sync:   &sync.WaitGroup{},
	}

	for i := 0; i < num; i++ {
		w := NewWorker(i)
		w.sync = pool.sync
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
		case ch := <-p.stop:
			for _, w := range p.Workers {
				w.Stop()
			}
			close(p.Queue)
			close(p.Errors)
			ch <- struct{}{}

			return
		case job := <-jobs:
			p.sync.Add(1)

			go func() {
				defer func() {
					if err := recover(); err != nil {
						p.Errors <- err.(error)
					}
				}()

				worker, ok := <-p.Queue // receive available worker from queue
				if !ok {
					return
				}
				worker <- job // send the worker the job
			}()
		}
	}
}

// Drain provides a graceful close of the pool allowing any workers in the queue
// to finish.
func (p *WorkPool) Drain() {
	p.sync.Wait()
	p.Flush()
}

// Flush is an immediate method to close the channels of the pool
func (p *WorkPool) Flush() {
	ch := make(chan struct{})
	p.stop <- ch
	<-ch
}
