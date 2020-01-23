package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) CreatePosts(posts []Post, slugOrId string) ([]Post, error) {
	threadID, err := strconv.ParseInt(slugOrId, 10, 64)

	var thread Thread
	if err != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrId)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}
	if err != nil {
		// fmt.Printf("Service CreatePosts: %s\n", err.Error())
		return []Post{}, errors.New(messages.ThreadDoesNotExist)
	}

	returnPosts, err := s.Repository.CreatePosts(posts, thread.ID, thread.Forum) //turn into int64
	if err != nil {
		// fmt.Printf("Service CreatePosts: %s\n", err.Error())
		return returnPosts, err
	}
	return returnPosts, err
}

func (s Service) Vote(vote Vote, slugOrId string) (Thread, error) {
	//should remove
	// if vote.Voice != 1 && vote.Voice != -1 {
	// 	return Thread{}, fmt.Errorf("Invalid value")
	// }
	threadID, convErr := strconv.ParseInt(slugOrId, 10, 64)

	var err error
	if convErr != nil {
		_, err = s.Repository.GetVoteByThreadSlug(vote.Nickname, slugOrId)
	} else {
		_, err = s.Repository.GetVoteByThreadID(vote.Nickname, threadID)
	}

	if err != nil {
		err = s.Repository.CreateVote(vote, slugOrId)
		if err != nil {
			// fmt.Printf("Vote: %s\n", err.Error())
			return Thread{}, fmt.Errorf(messages.UserNotFound)
		}
	} else {
		err = s.Repository.ChangeVote(vote, slugOrId)
		if err != nil {
			// fmt.Printf("Vote: %s\n", err.Error())
			return Thread{}, fmt.Errorf(messages.UserNotFound)
		}
	}

	var thread Thread
	if convErr != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrId)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}

	if err != nil {
		return thread, errors.New(messages.ThreadDoesNotExist)
	}

	return thread, err
}

func (s Service) CreateUser(newUser NewUser, nickname string) ([]User, error) {
	//see a bottle neck - could be done by one query
	var users []User
	err := s.Repository.CreateUser(newUser, nickname)

	if err != nil {
		// fmt.Printf("Service CreateUser: %s\n", err)
		userByNickname, err := s.Repository.GetUserByNickname(nickname)
		if err == nil {
			users = append(users, userByNickname)
		}

		userByEmail, err := s.Repository.GetUserByEmail(newUser.Email)
		if err == nil && userByNickname != userByEmail {
			users = append(users, userByEmail)
		}

		// fmt.Printf("Len users: %d\n", len(users))
		if len(users) > 0 {
			return users, fmt.Errorf(messages.UserAlreadyExists)
		}
	}

	user := User{
		About:    newUser.About,
		Email:    newUser.Email,
		Fullname: newUser.Fullname,
		Nickname: nickname,
	}

	users = append(users, user)
	return users, err
}

func (s Service) CreateForum(forum NewForum) (Forum, error) {

	user, err := s.Repository.GetUserByNickname(forum.User)
	if err != nil {
		return Forum{}, fmt.Errorf(messages.UserNotFound)
	}

	returnForum, err := s.Repository.GetForumBySlug(forum.Slug)
	if err == nil {
		return returnForum, fmt.Errorf(messages.ForumAlreadyExists)
	}

	forum.User = user.Nickname
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
	forum, err := s.Repository.GetForumBySlug(forumSlug)
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
		Forum:   forum.Slug,
		ID:      id,
		Message: thread.Message,
		Slug:    thread.Slug,
		Title:   thread.Title,
		Votes:   0,
	}
	if err != nil {
		// fmt.Printf("CreateThread: %s", err.Error())
		err = fmt.Errorf(messages.UserNotFound)
	}

	return returnThread, err
}
