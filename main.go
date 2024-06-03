package main

import (
	"fmt"
	"time"
)

func main() {

	jobQueue := make(chan Job, 10)

	dispatcher := NewDispatcher(3, jobQueue)
	dispatcher.Run()

	for i := 1; i <= 5; i++ {
		job := Job{
			ID:       i,
			Duration: time.Second * time.Duration(i),
		}
		jobQueue <- job
		fmt.Printf("Enqueued job %d\n", job.ID)
	}

	time.Sleep(10 * time.Second)
}
