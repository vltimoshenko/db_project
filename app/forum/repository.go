package forum

import (
	. "github.com/db_project/pkg/models"
)

type RepositoryInterface interface {
	CreateForum(NewForum) error
	CreateThread(NewThread, string) (int, error)
	CreateUser(NewUser, string) error
	CreatePosts([]NewPost, int, string) ([]Post, error)
	CreateVote(Vote, string) error

	GetThreadBySlug(string) (Thread, error)
	GetThreadByID(int) (Thread, error)
	GetPostByID(int) (Post, error)

	GetForumBySlug(string) (Forum, error)
	GetThreads(params map[string]interface{}) ([]Thread, error)
	GetUsers(params map[string]interface{}) ([]User, error)
	GetUserByNickname(string) (User, error)
	GetUserByEmail(string) (User, error)

	GetThreadsBySlug(string) ([]Thread, error)
	GetThreadsByID(int) ([]Thread, error)
	GetPosts(int, string, string, string, string) ([]Post, error)

	// GetPostsFlat(int, map[string]interface{}) ([]Post, error)
	// GetPostsTree(int, map[string]interface{}) ([]Post, error)
	// GetPostsParentTree(int, map[string]interface{}) ([]Post, error)

	GetVoteByThreadID(nickname string, thread int) (Vote, error)
	GetVoteByThreadSlug(nickname string, threadSlug string) (Vote, error)

	ChangeUser(NewUser, string) error
	ChangeThread(ThreadUpdate, int) error
	ChangeThreadRate(int, int) error
	ChangeVote(Vote, string) error
	ChangePost(PostUpdate, int, bool) error

	ClearDB() error
	GetStatus() (Status, error)
}
