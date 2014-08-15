package wrker

// Worker

type Worker struct {
	id  interface{}
	Job chan Job
}

// NewWorker returns a new Worker
func NewWorker(id interface{}) (w *Worker) {
	w = &Worker{
		id:  id,
		Job: make(chan Job),
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
			case job := <-w.Job:
				if err := job.Do(); err != nil {
					er <- err
				}
			}
		}
	}()
}
