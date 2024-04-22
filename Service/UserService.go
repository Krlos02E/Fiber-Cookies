package service

import (
	repository "SessionCookie/Repository"
	"SessionCookie/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	Rep *repository.UserRepository
}

func (s *UserService) ListAllUsers() (*[]models.User, error) {
	users, err := s.Rep.FindAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) AuthenticateUser(user *models.User) *fiber.Error {

	foundUser, err := s.Rep.FindUserByName(user.UserName)
	if err != nil {
		return fiber.ErrNotFound
	}

	if foundUser.Password != user.Password {
		return fiber.ErrUnauthorized
	}

	return nil
}

func (s *UserService) CreateJWT(userName string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userName,
		"exp":  time.Now().Add(time.Minute * 10).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
