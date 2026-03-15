package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"practice5/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetPaginatedUsers(filter models.UserFilter) (models.PaginatedResponse, error) {
	page := filter.Page
	pageSize := filter.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	orderBy := "id"
	allowedOrderBy := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"gender":     true,
		"birth_date": true,
	}

	if allowedOrderBy[filter.OrderBy] {
		orderBy = filter.OrderBy
	}

	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argPos))
		args = append(args, *filter.ID)
		argPos++
	}

	if filter.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Name+"%")
		argPos++
	}

	if filter.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Email+"%")
		argPos++
	}

	if filter.Gender != nil {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argPos))
		args = append(args, *filter.Gender)
		argPos++
	}

	if filter.BirthDate != nil {
		conditions = append(conditions, fmt.Sprintf("birth_date = $%d", argPos))
		args = append(args, *filter.BirthDate)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM users" + whereClause

	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	offset := (page - 1) * pageSize

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birth_date
		FROM users
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argPos, argPos+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Gender,
			&user.BirthDate,
		)
		if err != nil {
			return models.PaginatedResponse{}, err
		}

		users = append(users, user)
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(user1ID int, user2ID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM user_friends uf1
		JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
		JOIN users u ON u.id = uf1.friend_id
		WHERE uf1.user_id = $1 AND uf2.user_id = $2
	`

	rows, err := r.db.Query(query, user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Gender,
			&user.BirthDate,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
