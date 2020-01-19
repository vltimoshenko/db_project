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

func (s Service) GetPost(postID int, params []string) (map[string]interface{}, error) {
	postInfo := make(map[string]interface{})
	post, err := s.Repository.GetPostByID(postID)
	if err != nil {
		return postInfo, fmt.Errorf(messages.PostDoesNotExist)
	}

	postInfo["post"] = post
	for _, obj := range params {
		switch obj {
		case "user":
			postInfo["author"], err = s.Repository.GetUserByNickname(post.Author)
		case "forum":
			postInfo["forum"], err = s.Repository.GetForumBySlug(post.Forum)
		case "thread":
			postInfo["thread"], err = s.Repository.GetThreadByID(post.Thread)
			//errors could be eliminated
		}
	}

	return postInfo, err
}

func (s Service) GetPosts(slugOrID string, params map[string]interface{}) ([]Post, error) {
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

	// var posts []Post
	// if params["sort"] == "tree" {
	// 	posts, err = s.Repository.GetPostsTree(thread.ID, params)
	// } else if params["sort"] == "parent_tree" {
	// 	posts, err = s.Repository.GetPostsParentTree(thread.ID, params)
	// } else {
	// 	posts, err = s.Repository.GetPostsFlat(thread.ID, params)
	// }
	var limit, since, sort, desc string
	if params["limit"] != nil {
		limit = params["limit"].(string)
	} else {
		limit = "100"
	}

	if params["sort"] != nil {
		sort = params["sort"].(string)
	} else {
		sort = "flat"
	}

	if params["desc"] != nil {
		desc = params["desc"].(string)
	} else {
		desc = "false"
	}

	if params["since"] != nil {
		since = params["since"].(string)
	} else {
		if desc == "false" {
			since = "0"
		} else {
			since = "999999999"
		}
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
