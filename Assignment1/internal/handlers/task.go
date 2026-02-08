package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	tasks  = make(map[int]Task)
	nextID = 1
	mu     sync.Mutex
)

const maxTitleLen = 100

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	doneFilter := r.URL.Query().Get("done") // "", "true", "false"

	mu.Lock()
	defer mu.Unlock()

	result := []Task{}
	for _, task := range tasks {
		// no filter
		if doneFilter == "" {
			result = append(result, task)
			continue
		}

		// filter true/false
		if doneFilter == "true" && task.Done {
			result = append(result, task)
		}
		if doneFilter == "false" && !task.Done {
			result = append(result, task)
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	task, ok := tasks[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	if input.Title == "" || len(input.Title) > maxTitleLen {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title must be non-empty and <= 100 characters",
		})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	task := Task{
		ID:    nextID,
		Title: input.Title,
		Done:  false,
	}

	tasks[nextID] = task
	nextID++

	writeJSON(w, http.StatusCreated, task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	// decode into map to validate that "done" exists and is boolean
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	val, ok := payload["done"]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "done is required"})
		return
	}

	done, ok := val.(bool)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "done must be boolean"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	task, exists := tasks[id]
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	task.Done = done
	tasks[id] = task

	writeJSON(w, http.StatusOK, map[string]bool{"updated": true})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := tasks[id]; !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	delete(tasks, id)
	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}
