package repository

import (
	"SessionCookie/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (rep *UserRepository) FindAllUsers() (*[]models.User, error) {
	users := &[]models.User{}
	err := rep.DB.Find(users).Error
	return users, err
}

func (rep *UserRepository) FindUserByName(name string) (*models.User, error) {

	user := &models.User{}
	err := rep.DB.First(user, "user_name = ?", name).Error
	return user, err
}
