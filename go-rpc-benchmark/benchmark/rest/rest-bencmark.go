package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {

	start := time.Now()

	requests := 10000

	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			http.Get("http://localhost:8080/user")

		}()
	}

	wg.Wait()

	duration := time.Since(start)

	fmt.Println("REST requests:", requests)
	fmt.Println("Time:", duration)
	fmt.Println("RPS:", float64(requests)/duration.Seconds())
}
