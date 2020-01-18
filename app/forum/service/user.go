package service

import (
	"fmt"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) CreateUser(newUser NewUser, nickname string) (User, error) {

	//see a bottle neck - could be done by one query
	_, err := s.Repository.GetUserByNickname(nickname)
	if err == nil {
		return User{}, fmt.Errorf(messages.UserAlreadyExists)
	}

	_, err = s.Repository.GetUserByEmail(nickname)
	if err == nil {
		return User{}, fmt.Errorf(messages.UserAlreadyExists)
	}

	err = s.Repository.CreateUser(newUser, nickname)

	if err != nil {
		return User{}, err
	}

	user := User{
		About:    newUser.About,
		Email:    newUser.Email,
		Fullname: newUser.Fullname,
		Nickname: nickname,
	}

	return user, nil
}

func (s Service) GetUser(nickname string) (User, error) {
	user, err := s.Repository.GetUserByNickname(nickname)
	if err != nil {
		return user, fmt.Errorf(messages.UserNotFound)
	}
	return user, nil
}
