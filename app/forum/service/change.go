package service

import (
	"fmt"
	"strconv"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) ChangeThread(threadUpdate ThreadUpdate, slugOrId string) (Thread, error) {

	threadID, err := strconv.ParseInt(slugOrId, 10, 64)

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

	if threadUpdate.Message == " " || len(threadUpdate.Message) == 0 {
		threadUpdate.Message = thread.Message
	} else {
		thread.Message = threadUpdate.Message
	}

	if threadUpdate.Title == " " || len(threadUpdate.Title) == 0 {
		threadUpdate.Title = thread.Title
	} else {
		thread.Title = threadUpdate.Title
	}

	err = s.Repository.ChangeThread(threadUpdate, thread.ID)
	if err != nil {
		return thread, fmt.Errorf(messages.ThreadDoesNotExist) //should be another error
	}

	return thread, nil
}

func (s Service) ChangeUser(newUser NewUser, nickname string) (User, error) {
	_, err := s.Repository.GetUserByNickname(nickname)
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
