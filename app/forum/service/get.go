package service

import (
	"fmt"
	"strconv"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
)

func (s Service) GetThread(slugOrID string) (Thread, error) {
	threadID, err := strconv.Atoi(slugOrID)

	var thread Thread
	if err != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrID)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}

	return thread, err
}

func (s Service) GetForum(forumSlug string) (Forum, error) {
	forum, err := s.Repository.GetForumBySlug(forumSlug)
	if err != nil {
		return forum, fmt.Errorf(messages.ForumDoesNotExist)
	}
	return forum, nil
}

func (s Service) GetThreads(params map[string]interface{}) ([]Thread, error) {
	_, err := s.Repository.GetForumBySlug(params["forum"].(string))
	if err != nil {
		return []Thread{}, fmt.Errorf(messages.ForumDoesNotExist)
	}

	threads, err := s.Repository.GetThreads(params)
	if err != nil {
		return threads, fmt.Errorf(messages.ForumDoesNotExist) //specific error
	}
	return threads, nil
}

func (s Service) GetUsers(params map[string]interface{}) ([]User, error) {
	_, err := s.Repository.GetForumBySlug(params["forum"].(string))
	if err != nil {
		return []User{}, fmt.Errorf(messages.ForumDoesNotExist)
	}

	users, err := s.Repository.GetUsers(params)
	if err != nil {
		return users, fmt.Errorf(messages.ForumDoesNotExist)
	}
	return users, nil
}

func (s Service) GetPost(postID int64, params []string) (map[string]interface{}, error) {
	postInfo := make(map[string]interface{})
	post, err := s.Repository.GetPostByID(postID)
	if err != nil {
		return postInfo, fmt.Errorf(messages.PostDoesNotExist)
	}

	postInfo["post"] = post
	for _, obj := range params {
		switch obj {
		case "user":
			postInfo["author"], _ = s.Repository.GetUserByNickname(post.Author)
		case "forum":
			postInfo["forum"], _ = s.Repository.GetForumBySlug(post.Forum)
		case "thread":
			postInfo["thread"], _ = s.Repository.GetThreadByID(post.Thread)
			//errors could be eliminated
		}
	}

	return postInfo, err
}

func (s Service) GetPosts(slugOrID string, limit int64, since string, sort string, desc bool) ([]Post, error) {
	threadID, err := strconv.Atoi(slugOrID)

	var thread Thread
	if err != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrID)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}

	if err != nil {
		fmt.Println(err)
		return []Post{}, fmt.Errorf(messages.ThreadDoesNotExist)
	}

	posts, err := s.Repository.GetPosts(thread.ID, limit, since, sort, desc)

	return posts, err
}

func (s Service) GetUser(nickname string) (User, error) {
	user, err := s.Repository.GetUserByNickname(nickname)
	if err != nil {
		return user, fmt.Errorf(messages.UserNotFound)
	}
	return user, nil
}
