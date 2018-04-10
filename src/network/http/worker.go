package http

import (
	"net/http"
	"fmt"
	"sync"
)



type Payload struct{	
	Conn 		http.ResponseWriter
	Message		*http.Request
	Wait		*sync.WaitGroup
}


type Job struct{
	Payload Payload	
}

var JobQueu chan Job

type Worker struct{
	WorkerPool  chan chan Job
	JobChannel  chan Job
	quit    	chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w Worker) Start() {		
	go func() {		
		for {
			//fmt.Println("w.WorkerPool <- w.JobChannel")
			
			w.WorkerPool <- w.JobChannel
			
			//fmt.Println("worker working")			

			select {
			case job := <-w.JobChannel:		
				fmt.Println("job := <-w.JobChannel")																									
				ReSponseMessage(job.Payload.Conn , job.Payload.Message, job.Payload.Wait)
				
			case <-w.quit:
				fmt.Println("<-w.quit")
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
