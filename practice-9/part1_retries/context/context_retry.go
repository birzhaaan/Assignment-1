package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func doSomethingUnreliable() error {
	return errors.New("database connection refused")
}

type RetryConfig struct {
	maxRetries int
	baseDelay  time.Duration
	maxDelay   time.Duration
}

func Retry(ctx context.Context, cfg RetryConfig) error {
	var err error
	for attempt := 0; attempt < cfg.maxRetries; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = doSomethingUnreliable()
		if err == nil {
			return nil
		}

		if attempt == cfg.maxRetries-1 {
			return err
		}

		backoff := cfg.baseDelay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > cfg.maxDelay {
			backoff = cfg.maxDelay
		}

		jitter := time.Duration(rand.Int63n(int64(backoff)))
		fmt.Printf("Attempt %d failed, waiting %v (max backoff: %v)...\n", attempt+1, jitter, backoff)
		
		select {
		case <-time.After(jitter):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	cfg := RetryConfig{
		maxRetries: 10,
		baseDelay:  200 * time.Millisecond,
		maxDelay:   2 * time.Second,
	}

	fmt.Println("=== Context-Aware Retry ===")
	err := Retry(ctx, cfg)
	
	if err != nil {
		fmt.Printf("Operation terminated: %v\n", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}