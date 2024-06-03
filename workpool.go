package main

import (
	"fmt"
	"time"
)

type Job struct {
	ID       int
	Duration time.Duration
}

type Worker struct {
	ID          int
	Jobchannel  chan Job
	WorkerPool  chan chan Job
	Quitchannel chan bool
}

type Dispatcher struct {
	WorkerPool chan chan Job
	Maxworkers int
	Jobqueue   chan Job
}

func NewWorker(id int, workpool chan chan Job) Worker {
	return Worker{
		ID:          id,
		Jobchannel:  make(chan Job),
		WorkerPool:  workpool,
		Quitchannel: make(chan bool),
	}
}

func NewDispatcher(maxworkers int, jobqueue chan Job) *Dispatcher {
	workerpool := make(chan chan Job, maxworkers)
	return &Dispatcher{
		WorkerPool: workerpool,
		Maxworkers: maxworkers,
		Jobqueue:   jobqueue,
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.Jobchannel

			select {
			case job := <-w.Jobchannel:
				fmt.Printf("Worker %d: Started job %d\n", w.ID, job.ID)
				time.Sleep(job.Duration)
				fmt.Printf("Worker %d: Finished job %d\n", w.ID, job.ID)

			case <-w.Quitchannel:
				fmt.Printf("Worker %d: Stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w Worker) Quit() {
	go func() {
		w.Quitchannel <- true
	}()
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.Maxworkers; i++ {
		worker := NewWorker(i+1, d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.Jobqueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}
