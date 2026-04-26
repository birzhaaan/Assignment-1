package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func doSomethingUnreliable() error {
	return errors.New("network congestion")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	var err error
	const maxRetries = 5
	baseDelay := 100 * time.Millisecond
	maxDelay := 5 * time.Second

	fmt.Println("=== Exponential Backoff + Jitter ===")
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			fmt.Println("Success!")
			break
		}

		if attempt == maxRetries-1 {
			break
		}

		backoffTime := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
		if backoffTime > maxDelay {
			backoffTime = maxDelay
		}

		jitter := time.Duration(rand.Int63n(int64(backoffTime)))
		fmt.Printf("Attempt %d failed, waiting %v (backoff %v + jitter) before next retry...\n", attempt+1, jitter, backoffTime)
		
		time.Sleep(jitter)
	}
	
	if err != nil {
		fmt.Println("Operation completely failed after all jitter attempts.")
	}
}