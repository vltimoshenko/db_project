package forum

import (
	. "github.com/db_project/pkg/models"
)

type ServiceInterface interface {
	CreateForum(NewForum) (Forum, error)
	CreateThread(NewThread, string) (Thread, error)
	CreateUser(NewUser, string) ([]User, error)
	CreatePosts(posts []Post, slugOrId string) ([]Post, error)

	GetUser(string) (User, error)
	GetForum(string) (Forum, error)
	GetThread(string) (Thread, error)
	GetPost(int64, []string) (map[string]interface{}, error)

	GetThreads(map[string]interface{}) ([]Thread, error)
	GetUsers(map[string]interface{}) ([]User, error)
	GetPosts(slugOrID string, limit int64, since string, sort string, desc bool) ([]Post, error)

	ChangeUser(NewUser, string) (User, error)
	ChangeThread(ThreadUpdate, string) (Thread, error)
	ChangePost(PostUpdate, int64) (Post, error)

	Vote(Vote, string) (Thread, error)
	ClearDB() error
	GetStatus() (Status, error)
}
