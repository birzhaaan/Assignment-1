package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func doSomethingUnreliable() error {
	if rand.Intn(10) < 8 { 
		return errors.New("service unavailable")
	}
	return nil
}

func main() {
	var err error
	const maxRetries = 5
	baseDelay := 100 * time.Millisecond
	maxDelay := 5 * time.Second

	fmt.Println("=== Exponential Backoff ===")
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			fmt.Println("Operation succeeded!")
			break
		}

		if attempt == maxRetries-1 {
			break
		}

		backoffTime := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
		
		if backoffTime > maxDelay {
			backoffTime = maxDelay
		}

		fmt.Printf("Attempt %d failed, waiting %v before next retry...\n", attempt+1, backoffTime)
		time.Sleep(backoffTime)
	}

	if err != nil {
		fmt.Printf("Finally failed: %v\n", err)
	}
}