package http

import (
	"fmt"
	"runtime"
)

type Dispatcher struct {	
	WorkerPool chan chan Job
	MaxWorkers int
}

func NewDispatcher() *Dispatcher {
	
	maxWorkers := runtime.GOMAXPROCS(runtime.NumCPU())
	
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool , MaxWorkers : maxWorkers}
}

func (d *Dispatcher) Run() {    
	
	JobQueu = make(chan Job)
	
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()	
}

func (d *Dispatcher) dispatch() {	
	for {
		select {
		case job := <-JobQueu:				
			go func(job Job) {								
				jobChannel := <-d.WorkerPool
				fmt.Println("jobChannel <- job")		
				jobChannel <- job
				
			}(job)
		}
	}
}