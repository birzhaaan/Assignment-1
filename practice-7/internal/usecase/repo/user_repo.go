package repo

import (
	"practice-7/internal/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{Db: db}
}

func (r *UserRepo) CreateUser(user *entity.User) error {
	return r.Db.Create(user).Error
}

func (r *UserRepo) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.Db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.Db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepo) UpdateUserRole(id string, newRole string) error {
	return r.Db.Model(&entity.User{}).Where("id = ?", id).Update("role", newRole).Error
}