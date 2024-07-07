package service

import (
	"errors"
	"olshop/features/users"
	"olshop/helpers/encrypt"
)

type userService struct {
	repo users.Repository
	hash encrypt.HashInterface
}

func New(repo users.Repository, hash encrypt.HashInterface) users.Service {
	return &userService{
		repo: repo,
		hash: hash,
	}
}

func (service *userService) Register(newUser users.User) error {
	if newUser.Name == "" {
		return errors.New("validate: name can't be empty")
	}
	if newUser.Email == "" {
		return errors.New("validate: email can't be empty")
	}
	if newUser.Password == "" {
		return errors.New("validate: password can't be empty")
	}

	encrypt, err := service.hash.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	newUser.Password = encrypt

	if err := service.repo.Register(newUser); err != nil {
		return err
	}

	return nil
}

func (service *userService) Login(email string, password string) (*users.User, error) {
	if email == "" {
		return nil, errors.New("validate: email can't be empty")
	}
	if password == "" {
		return nil, errors.New("validate: password can't be empty")
	}

	result, err := service.repo.Login(email)
	if err != nil {
		return nil, err
	}

	if err := service.hash.Compare(result.Password, password); err != nil {
		return nil, errors.New("validate: wrong password")
	}

	return result, nil
}

func (service *userService) Update(id uint, updateUser users.User) error {
	if id == 0 {
		return errors.New("validate: invalid id")
	}
	if updateUser.Name == "" {
		return errors.New("validate: name can't be empty")
	}
	if updateUser.Email == "" {
		return errors.New("validate: email can't be empty")
	}
	if updateUser.Password == "" {
		return errors.New("validate: password can't be empty")
	}
	if updateUser.Username == "" {
		return errors.New("validate: username can't be empty")
	}
	if updateUser.ImageRaw == nil {
		return errors.New("validate: image can't be empty")
	}

	encrypt, err := service.hash.HashPassword(updateUser.Password)
	if err != nil {
		return err
	}

	updateUser.Password = encrypt

	if err := service.repo.Update(id, updateUser); err != nil {
		return err
	}

	return nil
}

func (service *userService) Delete(id uint) error {
	if id == 0 {
		return errors.New("validate: invalid id")
	}

	if err := service.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (service *userService) GetById(id uint) (*users.User, error) {
	if id == 0 {
		return nil, errors.New("validate: invalid id")
	}

	result, err := service.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
