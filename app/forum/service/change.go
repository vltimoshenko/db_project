package service

import (
	"fmt"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) ChangeThread(threadUpdate ThreadUpdate, slugOrId string) (Thread, error) {
	thread, err := s.Repository.ChangeThread(threadUpdate, slugOrId)
	return thread, err
}

func (s Service) ChangeUser(newUser NewUser, nickname string) (User, error) {
	_, err := s.Repository.GetUserByNickname(nickname) //could be removed
	if err != nil {
		return User{}, fmt.Errorf(messages.UserNotFound)
	}

	user, err := s.Repository.ChangeUser(newUser, nickname)

	if err != nil {
		fmt.Printf("Serv ChangeUser: %s", err.Error())
		return User{}, fmt.Errorf(messages.UserAlreadyExists)
	}

	return user, nil
}

func (s Service) ChangePost(updatePost PostUpdate, postID int64) (Post, error) {
	post, err := s.Repository.ChangePost(updatePost, postID)
	return post, err
}
