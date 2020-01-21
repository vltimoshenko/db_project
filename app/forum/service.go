package forum

import (
	. "github.com/db_project/pkg/models"
	"io"
)

type ServiceInterface interface {
	CreateForum(body io.ReadCloser) (Forum, error)
	CreateThread(NewThread, string) (Thread, error)
	CreateUser(NewUser, string) ([]User, error)
	CreatePosts(body io.ReadCloser, slugOrId string) ([]Post, error)

	GetUser(string) (User, error)
	GetForum(string) (Forum, error)
	GetThread(string) (Thread, error)
	GetPost(int64, []string) (map[string]interface{}, error)

	GetThreads(map[string]interface{}) ([]Thread, error)
	GetUsers(map[string]interface{}) ([]User, error)
	// GetPosts(string, map[string]interface{}) ([]Post, error)
	GetPosts(slugOrID string, limit int64, since string, sort string, desc bool) ([]Post, error)

	ChangeUser(NewUser, string) (User, error)
	ChangeThread(ThreadUpdate, string) (Thread, error)
	ChangePost(PostUpdate, int64) (Post, error)

	Vote(Vote, string) (Thread, error)
	ClearDB() error
	GetStatus() (Status, error)
}
