package usecase

import (
	"fmt"
	"strings"

	"golang/internal/repository"
	"golang/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: r}
}

func (u *UserUsecase) GetUsers() ([]modules.User, error) {
	return u.repo.GetUsers()
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) CreateUser(name string) (*modules.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return u.repo.CreateUser(name)
}

func (u *UserUsecase) UpdateUser(id int, name string) (*modules.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return u.repo.UpdateUser(id, name)
}

// ✅ NEW: Delete
func (u *UserUsecase) DeleteUser(id int) (int64, error) {
	return u.repo.DeleteUser(id)
}
