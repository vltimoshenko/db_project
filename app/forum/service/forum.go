package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) CreateForum(body io.ReadCloser) (Forum, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		//return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var forum NewForum
	err = json.Unmarshal(bytes, &forum)
	if err != nil {
		//return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	_, err = s.Repository.GetUserByNickname(forum.User)
	if err != nil {
		return Forum{}, fmt.Errorf(messages.UserNotFound)
	}
	returnForum, err := s.Repository.GetForumBySlug(forum.Slug)
	if err == nil {
		return returnForum, fmt.Errorf(messages.ForumAlreadyExists)
	}

	err = s.Repository.CreateForum(forum)

	returnForum = Forum{
		Posts:  0,
		Slug:   forum.Slug,
		Thread: 0,
		Title:  forum.Title,
		User:   forum.User,
	}
	return returnForum, err
}

func (s Service) CreateThread(thread NewThread, forumSlug string) (Thread, error) {

	_, err := s.Repository.GetForumBySlug(forumSlug)
	if err != nil {
		return Thread{}, fmt.Errorf(messages.ForumDoesNotExist)
	}

	returnThread, err := s.Repository.GetThreadBySlug(thread.Slug) //consider check order
	if err == nil {
		return returnThread, fmt.Errorf(messages.ThreadAlreadyExists)
	}

	id, err := s.Repository.CreateThread(thread, forumSlug)

	returnThread = Thread{
		Author:  thread.Author,
		Created: thread.Created,
		Forum:   forumSlug,
		ID:      id,
		Message: thread.Message,
		Slug:    thread.Slug,
		Title:   thread.Title,
		Votes:   0,
	}
	return returnThread, err
}
