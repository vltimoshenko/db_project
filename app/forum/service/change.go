package service

import (
	"fmt"
	"strconv"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) ChangeThread(threadUpdate ThreadUpdate, slugOrId string) (Thread, error) {

	threadID, err := strconv.Atoi(slugOrId)

	var thread Thread
	//could be done one query
	if err != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrId)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}

	if err != nil {
		fmt.Println(err)
		return thread, fmt.Errorf(messages.ThreadDoesNotExist)
	}

	if threadUpdate.Message == "" {
		threadUpdate.Message = thread.Message
	}
	if threadUpdate.Title == "" {
		threadUpdate.Title = thread.Title
	}

	err = s.Repository.ChangeThread(threadUpdate, thread.ID)
	if err != nil {
		return thread, fmt.Errorf(messages.ThreadDoesNotExist) //should be another error
	}

	thread.Message = threadUpdate.Message
	thread.Title = threadUpdate.Title

	return thread, nil
}

func (s Service) ChangeUser(newUser NewUser, nickname string) (User, error) {

	oldUser, err := s.Repository.GetUserByNickname(nickname)
	if err != nil {
		return User{}, fmt.Errorf(messages.UserNotFound)
	}

	// mutex
	userByEmail, _ := s.Repository.GetUserByEmail(newUser.Email)
	// if err == nil {
	// 	return User{}, fmt.Errorf(messages.UserAlreadyExists)
	// }
	//consider db error

	if userByEmail.Nickname != "" && userByEmail.Nickname != nickname {
		return User{}, fmt.Errorf(messages.UserAlreadyExists)
	}

	if newUser.About == "" {
		newUser.About = oldUser.About
	}
	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}
	if newUser.Fullname == "" {
		newUser.Fullname = oldUser.Fullname
	}

	err = s.Repository.ChangeUser(newUser, nickname)

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

func (s Service) ChangePost(updatePost PostUpdate, postID int) (Post, error) {
	post, err := s.Repository.GetPostByID(postID)
	if err != nil {
		fmt.Println(err.Error())
		return post, fmt.Errorf(messages.PostDoesNotExist)
	}

	post.IsEdited = true
	if updatePost.Message == "" || updatePost.Message == post.Message {
		updatePost.Message = post.Message
		post.IsEdited = false
	} else {
		post.Message = updatePost.Message
	}

	err = s.Repository.ChangePost(updatePost, post.ID, post.IsEdited)

	return post, nil
}
