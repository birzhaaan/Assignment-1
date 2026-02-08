package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetExternalTodos(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "external api unavailable"})
		return
	}
	defer resp.Body.Close()

	var data any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "failed to parse external json"})
		return
	}

	writeJSON(w, http.StatusOK, data)
}
