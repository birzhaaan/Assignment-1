package main

import (
	"log"
	"net/http"

	"Assignment1/internal/handlers"
	"Assignment1/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	// Tasks endpoint
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				handlers.GetTaskByID(w, r)
			} else {
				handlers.GetTasks(w, r) // supports ?done=true/false
			}
		case http.MethodPost:
			handlers.CreateTask(w, r)
		case http.MethodPatch:
			handlers.UpdateTask(w, r)
		case http.MethodDelete:
			handlers.DeleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// External API endpoint (bonus)
	mux.HandleFunc("/external/todos", handlers.GetExternalTodos)

	// middleware chain: Logging -> RequestID -> APIKey -> handlers
	handler := middleware.Logging("request received")(
		middleware.RequestID(
			middleware.APIKey(mux),
		),
	)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
