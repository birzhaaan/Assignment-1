package users

import (
	"database/sql"
	"fmt"
	"time"

	"golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User

	err := r.db.DB.Get(&user, "SELECT id, name FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateUser(name string) (*modules.User, error) {
	var id int

	err := r.db.DB.QueryRow(
		"INSERT INTO users (name) VALUES ($1) RETURNING id",
		name,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &modules.User{ID: id, Name: name}, nil
}

func (r *Repository) UpdateUser(id int, name string) (*modules.User, error) {
	res, err := r.db.DB.Exec("UPDATE users SET name=$1 WHERE id=$2", name, id)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &modules.User{ID: id, Name: name}, nil
}

// ✅ NEW: Delete with RowsAffected
func (r *Repository) DeleteUser(id int) (int64, error) {
	res, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected == 0 {
		return 0, fmt.Errorf("user with id %d not found", id)
	}

	return affected, nil
}
