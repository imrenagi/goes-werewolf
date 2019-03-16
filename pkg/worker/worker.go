package worker

import (
	"fmt"
	"log"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerPool chan chan Job) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:         id,
		Work:       make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool)}

	return worker
}

type Worker struct {
	ID         int
	Work       chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker pool.
			w.WorkerPool <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				// todo run the job
				err := work.Payload.Fn()
				if err != nil {
					log.Println(err)
				}

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}