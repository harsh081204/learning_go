// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type Job struct {
// 	ID int
// }

// func process(job Job) {
// 	fmt.Println("Processing Job:", job.ID)
// 	time.Sleep(time.Second)
// }

// func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	for job := range jobs {
// 		fmt.Printf("worker %d started job %d\n", id, job.ID)
// 		process(job)
// 		fmt.Printf("worker %d finished job %d\n", id, job.ID)
// 	}
// }

// func main() {
// 	numWorkers := 3
// 	queueSize := 5

// 	jobs := make(chan Job, queueSize)

// 	var wg sync.WaitGroup

// 	for w := 1; w <= numWorkers; w++ {
// 		wg.Add(1)
// 		go worker(w, jobs, &wg)
// 	}

// 	for j := 1; j <= 10; j++ {
// 		fmt.Println("enqueue job", j)
// 		jobs <- Job{ID: j}
// 	}

// 	close(jobs)
// 	wg.Wait()
// }

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID int
	Err   error
}

type WorkerPool struct {
	workers int
	jobs    chan Job
	results chan Result

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorkerPool(workers, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &WorkerPool{
		workers: workers,
		jobs:    make(chan Job, queueSize),
		results: make(chan Result, queueSize),
		ctx:     ctx,
		cancel:  cancel,
	}

	p.startWorkers()

	return p
}

func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case job, ok := <-p.jobs:

			if !ok {
				return
			}

			func() {

				defer func() {
					if r := recover(); r != nil {
						p.results <- Result{
							JobID: job.ID,
							Err:   fmt.Errorf("panic: %v", r),
						}
					}
				}()

				err := processJob(job)

				p.results <- Result{
					JobID: job.ID,
					Err:   err,
				}

			}()
		}
	}
}

func processJob(job Job) error {
	fmt.Println("Processing job", job.ID)
	time.Sleep(time.Second)
	return nil
}

func (p *WorkerPool) Submit(job Job) error {
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("pool shutting down")
	case p.jobs <- job:
		return nil
	}
}

func (p *WorkerPool) TrySubmit(job Job) bool {
	select {
	case p.jobs <- job:
		return true
	default:
		return false
	}
}

func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

func (p *WorkerPool) Shutdown() {
	p.cancel()
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

func main() {

	pool := NewWorkerPool(4, 10)

	go func() {
		for r := range pool.Results() {
			fmt.Println("result:", r.JobID, r.Err)
		}
	}()

	for i := 1; i <= 20; i++ {

		err := pool.Submit(Job{
			ID:      i,
			Payload: "data",
		})

		if err != nil {
			fmt.Println("submit error:", err)
		}
	}

	pool.Shutdown()
}
