package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/odink789/project-management/models"
	"github.com/odink789/project-management/repositories"
)

//service ini adalah logika bisnis nya

type UserService interface {
	Register(user *models.User) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user *models.User) error {
	//kita harus mengecek email yang terdaftar apakah sudah di pakai atau blm
	//hasing password
	//set role
	//simpan user

	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}

	hased, err := utils.HashPassword(user.password)
	if err != nil {
		return err
	}

	user.Password = hased
	user.Role = "user"
	user.PublicId = uuid.New()
	return s.repo.Create(user)

}
