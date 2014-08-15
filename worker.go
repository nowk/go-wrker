package wrker

import "sync"

// Worker

type Worker struct {
	id   interface{}
	Job  chan Job
	stop chan chan struct{}
	sync *sync.WaitGroup
}

// NewWorker returns a new Worker
func NewWorker(id interface{}) (w *Worker) {
	w = &Worker{
		id:   id,
		Job:  make(chan Job),
		stop: make(chan chan struct{}),
	}

	return
}

func (w Worker) GetID() interface{} {
	return w.id
}

// Start the worker on the queue and process incoming job
func (w Worker) Start(workerq chan<- chan Job, er chan<- error) {
	go func() {
		for {
			workerq <- w.Job

			select {
			case ch := <-w.stop:
				close(w.Job)
				ch <- struct{}{}

				return
			case job := <-w.Job:
				if err := job.Do(); err != nil {
					er <- err
				}
				w.sync.Done()
			}
		}
	}()
}

// Stop sends stop on the worker to close the worker
func (w *Worker) Stop() {
	ch := make(chan struct{})
	w.stop <- ch
	<-ch
}
