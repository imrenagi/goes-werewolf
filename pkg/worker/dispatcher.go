package worker

import "fmt"

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkQueue chan Job
	WorkerPool chan chan Job
	maxWorkers int
}

// NewDispatcher creates new dispatcher. You can use worker.MAX_WORKERS and worker.defaultWorkQueue as the default
// value for this factory method
func NewDispatcher(maxWorkers int, workQueue chan Job) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		WorkQueue: workQueue,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, d.WorkerPool)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case work := <-d.WorkQueue:
			fmt.Println("Received work requeust")
			go func() {
				worker := <-d.WorkerPool
				fmt.Println("Dispatching work request")
				worker <- work
			}()
		}
	}
}