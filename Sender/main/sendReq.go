package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8080"
	var wg sync.WaitGroup
	mu := sync.Mutex{} // Mutex to synchronize output

	// Create a channel to receive results
	results := make(chan string, 25)

	// Create a goroutine for each request
	for i := 1; i <= 25; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			startTime := time.Now() // Record the start time
			resp, err := http.Get(url)
			duration := time.Since(startTime) // Calculate the duration

			if err != nil {
				results <- fmt.Sprintf("Request %d: Error: %v", i, err)
			} else {
				defer resp.Body.Close()
				results <- fmt.Sprintf("Request %d: Status code %d, Time taken: %v", i, resp.StatusCode, duration)
			}
		}(i)
	}

	// Close results channel after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print out results as they come in
	for result := range results {
		mu.Lock()
		fmt.Println(result)
		mu.Unlock()
	}
}
