package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func process(job Job) {
	fmt.Println("Processing Job:", job.ID)
	time.Sleep(time.Second)
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("worker %d started job %d\n", id, job.ID)
		process(job)
		fmt.Printf("worker %d finished job %d\n", id, job.ID)
	}
}

func main() {
	numWorkers := 3
	queueSize := 5

	jobs := make(chan Job, queueSize)

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	for j := 1; j <= 10; j++ {
		fmt.Println("enqueue job", j)
		jobs <- Job{ID: j}
	}

	close(jobs)
	wg.Wait()
}
