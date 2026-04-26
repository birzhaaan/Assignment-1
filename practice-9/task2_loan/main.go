package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*CachedResponse
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]*CachedResponse)}
}


func (m *MemoryStore) Get(key string) (*CachedResponse, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	resp, exists := m.data[key]
	return resp, exists
}

func (m *MemoryStore) StartProcessing(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[key]; exists {
		return false
	}
	m.data[key] = &CachedResponse{Completed: false}
	return true
}

func (m *MemoryStore) Finish(key string, status int, body []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if resp, exists := m.data[key]; exists {
		resp.StatusCode = status
		resp.Body = body
		resp.Completed = true
	}
}


var store = NewMemoryStore()

func IdempotencyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		if cached, exists := store.Get(key); exists {
			if cached.Completed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
			} else {
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
			}
			return
		}

		if !store.StartProcessing(key) {
			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}

		next(w, r)
	}
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	
	fmt.Println("Processing started...")
	time.Sleep(2 * time.Second) 

	response := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": uuid.New().String(),
	}
	body, _ := json.Marshal(response)

	store.Finish(key, http.StatusOK, body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func main() {
	handler := IdempotencyMiddleware(PaymentHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	fmt.Println("Simulating Double-Click Attack...")
	key := uuid.New().String()
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			req, _ := http.NewRequest("POST", server.URL, nil)
			req.Header.Set("Idempotency-Key", key)
			
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Request %d error: %v\n", id, err)
				return
			}
			fmt.Printf("Request %d: Status %d\n", id, resp.StatusCode)
		}(i)
	}

	wg.Wait()
	
	fmt.Println("\nRequest after completion:")
	req, _ := http.NewRequest("POST", server.URL, nil)
	req.Header.Set("Idempotency-Key", key)
	resp, _ := http.DefaultClient.Do(req)
	fmt.Printf("Final Request: Status %d\n", resp.StatusCode)
}