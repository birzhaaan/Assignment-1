package main

import (
	"fmt"
	"sync"
)

func main() {
	var safeMap sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(value int) {
			defer wg.Done()
			safeMap.Store("key", value)
		}(i)
	}

	wg.Wait()

	value, _ := safeMap.Load("key")
	fmt.Printf("Value: %v\n", value)
}
