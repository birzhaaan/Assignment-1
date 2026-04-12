package usecase

import "practice-7/internal/entity"

type UserInterface interface {
	RegisterUser(user *entity.User) (*entity.User, string, error)
	LoginUser(input *entity.LoginUserDTO) (string, error)
	GetMe(userID string) (*entity.User, error)
	PromoteUser(id string) error 
}