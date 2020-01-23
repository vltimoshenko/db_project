package service

import (
	"fmt"
	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"strconv"
	"sync"
)

func (s Service) GetThread(slugOrID string) (Thread, error) {
	threadID, err := strconv.ParseInt(slugOrID, 10, 64)

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
	threads, err := s.Repository.GetThreads(params)
	if err != nil {
		// fmt.Printf("Ser GetThreads: %s", err.Error())
		return threads, fmt.Errorf(messages.ForumDoesNotExist)
	}
	return threads, nil
}

func (s Service) GetUsers(params map[string]interface{}) ([]User, error) {
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
	// could do by gorutines
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, obj := range params {
		wg.Add(1)
		// switch obj {
		// case "user":
		// 	postInfo["author"], _ = s.Repository.GetUserByNickname(post.Author)
		// case "forum":
		// 	postInfo["forum"], _ = s.Repository.GetForumBySlug(post.Forum)
		// case "thread":
		// 	postInfo["thread"], _ = s.Repository.GetThreadByID(post.Thread)
		// }
		go s.getPostWorker(postInfo, &post, obj, &wg, &mu)
	}

	wg.Wait()
	return postInfo, err
}

func (s Service) getPostWorker(postInfo map[string]interface{}, post *Post, param string,
	wg *sync.WaitGroup, mu *sync.Mutex) {

	defer wg.Done()
	switch param {
	case "user":
		user, _ := s.Repository.GetUserByNickname(post.Author)
		mu.Lock()
		postInfo["author"] = user
		mu.Unlock()
	case "forum":
		forum, _ := s.Repository.GetForumBySlug(post.Forum)
		mu.Lock()
		postInfo["forum"] = forum
		mu.Unlock()
	case "thread":
		thread, _ := s.Repository.GetThreadByID(post.Thread)
		mu.Lock()
		postInfo["thread"] = thread
		mu.Unlock()
	}
}

func (s Service) GetPosts(slugOrID string, limit int64, since string, sort string, desc bool) ([]Post, error) {
	threadID, err := strconv.ParseInt(slugOrID, 10, 64)

	var thread Thread
	if err != nil {
		thread, err = s.Repository.GetThreadBySlug(slugOrID)
	} else {
		thread, err = s.Repository.GetThreadByID(threadID)
	}

	if err != nil {
		// fmt.Println(err)
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
