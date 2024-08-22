package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8080"

	for i := 1; i <= 25; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Request %d: Error: %v\n", i, err)
		} else {
			fmt.Printf("Request %d: Status code %d\n", i, resp.StatusCode)
			resp.Body.Close()
		}
		time.Sleep(time.Second / 2)
	}
}
