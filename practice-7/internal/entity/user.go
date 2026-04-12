package entity

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Username  string    `json:"username" gorm:"unique"`
    Email     string    `json:"email" gorm:"unique"`
    Password  string    `json:"password"`
    Role      string    `json:"role"`
    Verified  bool      `json:"verified"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}