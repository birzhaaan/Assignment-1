package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func openDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}
	defer db.Close()

	// health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, "ok")
	})

	// GET/POST users
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			rows, err := db.Query("SELECT id, name, email FROM users ORDER BY id")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer rows.Close()

			var users []User
			for rows.Next() {
				var u User
				if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				users = append(users, u)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)

		case http.MethodPost:
			var u User
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				http.Error(w, "bad json", 400)
				return
			}
			err := db.QueryRow(
				"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
				u.Name, u.Email,
			).Scan(&u.ID)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(u)

		default:
			http.Error(w, "method not allowed", 405)
		}
	})

	// PUT/DELETE by id: /users/{id}
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/users/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", 400)
			return
		}

		switch r.Method {
		case http.MethodPut:
			var u User
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				http.Error(w, "bad json", 400)
				return
			}
			_, err := db.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", u.Name, u.Email, id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(204)

		case http.MethodDelete:
			_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(204)

		default:
			http.Error(w, "method not allowed", 405)
		}
	})

	log.Println("Starting the Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
