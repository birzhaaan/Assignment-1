package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

const baseDelay = 500 * time.Millisecond

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true 
	}
	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	backoff := float64(baseDelay) * math.Pow(2, float64(attempt))
	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return jitter
}

func ExecutePayment(ctx context.Context, url string) error {
	const maxRetries = 5
	client := &http.Client{}

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Attempt %d: Success!\n", attempt+1)
			return nil
		}

		if !IsRetryable(resp, err) {
			return fmt.Errorf("non-retryable error")
		}

		if attempt == maxRetries-1 {
			return fmt.Errorf("max retries exhausted")
		}

		wait := CalculateBackoff(attempt)
		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, wait)

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable) // 503
			return
		}
		w.WriteHeader(http.StatusOK) // 200 OK
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ExecutePayment(ctx, server.URL)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
	}
}