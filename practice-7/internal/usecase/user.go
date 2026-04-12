package usecase

import (
	"errors"
	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	err := u.repo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}
	token, err := utils.GenerateJWT(user.ID, user.Role)
	return user, token, err
}

func (u *UserUseCase) LoginUser(input *entity.LoginUserDTO) (string, error) {
	user, err := u.repo.GetByUsername(input.Username)
	if err != nil {
		return "", errors.New("user not found")
	}
	if !utils.CheckPassword(user.Password, input.Password) {
		return "", errors.New("invalid password")
	}
	return utils.GenerateJWT(user.ID, user.Role)
}

func (u *UserUseCase) GetMe(userID string) (*entity.User, error) {
	return u.repo.GetByID(userID)
}

// Новое для задачи №2
func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.UpdateUserRole(id, "admin")
}