package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func doSomethingUnreliable() error {
	if rand.Intn(10) < 7 {
		fmt.Println("Operation failed, retrying...")
		return errors.New("temporary failure")
	}
	fmt.Println("Operation succeeded!")
	return nil
}

func main() {
	var err error
	const maxRetries = 5
	const delay = 1 * time.Second

	fmt.Println("=== Fixed Delay ===")
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = doSomethingUnreliable()
		if err == nil {
			break
		}
		
		fmt.Printf("Attempt %d failed, waiting %v before next retry...\n", attempt+1, delay)
		if attempt < maxRetries-1 {
			time.Sleep(delay)
		}
	}

	if err != nil {
		fmt.Printf("Failed after %d attempts: %v\n", maxRetries, err)
	} else {
		fmt.Println("Succeeded within retry limit.")
	}
}