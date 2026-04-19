package service

import (
	"errors"
	"practice-8/repository"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService_TableDriven(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	t.Run("RegisterUser", func(t *testing.T) {
		user := &repository.User{ID: 7, Name: "Admin"}
		email := "admin@domain.com"

		tests := []struct {
			name    string
			setup   func()
			wantErr string
		}{
			{
				name: "user_already_exists",
				setup: func() {
					mockRepo.EXPECT().GetByEmail(email).Return(user, nil)
				},
				wantErr: "already exists",
			},
			{
				name: "success",
				setup: func() {
					mockRepo.EXPECT().GetByEmail(email).Return(nil, nil)
					mockRepo.EXPECT().CreateUser(user).Return(nil)
				},
				wantErr: "",
			},
			{
				name: "repo_error_on_create",
				setup: func() {
					mockRepo.EXPECT().GetByEmail(email).Return(nil, nil)
					mockRepo.EXPECT().CreateUser(user).Return(errors.New("connection lost"))
				},
				wantErr: "connection lost",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setup()
				err := svc.RegisterUser(user, email)
				if tt.wantErr != "" {
					require.ErrorContains(t, err, tt.wantErr)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})

	t.Run("UpdateUserName", func(t *testing.T) {
		targetID := 42
		user := &repository.User{ID: targetID, Name: "OldName"}

		tests := []struct {
			name    string
			newName string
			setup   func()
			wantErr string
		}{
			{
				name:    "empty_name",
				newName: "",
				setup:   func() {},
				wantErr: "name cannot be empty",
			},
			{
				name:    "user_not_found",
				newName: "NewName",
				setup: func() {
					mockRepo.EXPECT().GetUserByID(targetID).Return(nil, errors.New("not found"))
				},
				wantErr: "not found",
			},
			{
				name:    "success",
				newName: "FreshName",
				setup: func() {
					mockRepo.EXPECT().GetUserByID(targetID).Return(user, nil)
					mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
				},
				wantErr: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setup()
				err := svc.UpdateUserName(targetID, tt.newName)
				if tt.wantErr != "" {
					require.ErrorContains(t, err, tt.wantErr)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		tests := []struct {
			name   string
			id     int
			setup  func()
			errMsg string
		}{
			{
				name:   "protect_admin",
				id:     1,
				setup:  func() {},
				errMsg: "not allowed to delete admin user",
			},
			{
				name: "success",
				id:   99,
				setup: func() {
					mockRepo.EXPECT().DeleteUser(99).Return(nil)
				},
				errMsg: "",
			},
			{
				name: "repo_error",
				id:   88,
				setup: func() {
					mockRepo.EXPECT().DeleteUser(88).Return(errors.New("db locked"))
				},
				errMsg: "db locked",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setup()
				err := svc.DeleteUser(tt.id)
				if tt.errMsg != "" {
					require.ErrorContains(t, err, tt.errMsg)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})
}