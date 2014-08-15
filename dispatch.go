package wrker

// Dispatch connects and starts the workers in the WorkPool to an incoming chan
// of jobs and begins routing. Returns a error chan
func Dispatch(jobs chan Job, pool *WorkPool) chan error {
	var n = len(pool.Workers)
	workerq := make(chan chan Job, n)
	er := make(chan error)

	for _, wkr := range pool.Workers {
		wkr.Start(workerq, er)
	}

	go route(jobs, workerq)

	return er
}

// route pushes received jobs then sends them to an available worker
func route(jobs <-chan Job, workerq <-chan chan Job) {
	for {
		select {
		case job := <-jobs:
			go func() {
				worker := <-workerq // receive available worker from queue
				worker <- job       // send the worker the job
			}()
		}
	}
}
