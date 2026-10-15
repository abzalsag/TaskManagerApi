package service

import (
	"TaskManagerApi/auth"
	"TaskManagerApi/model"
	"TaskManagerApi/repository"
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (service *UserService) Register(
	ctx context.Context,
	email string,
	password string,
) error {
	email = strings.TrimSpace(email)

	if email == "" || password == "" {
		return errors.New("email and password are required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    email,
		Password: string(passwordHash),
	}

	return service.repo.CreateUser(user)
}

func (service *UserService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	email = strings.TrimSpace(email)

	user, err := service.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
