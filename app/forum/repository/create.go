package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

func (r *Repository) CreatePosts(posts []NewPost, threadID int, forum string) ([]Post, error) {
	tx, _ := r.DbConn.Begin()

	created := time.Now().Format(time.RFC3339)
	returnPosts := []Post{}

	_, err := r.GetThreadByID(threadID)
	if err != nil {
		return []Post{}, fmt.Errorf(messages.ThreadDoesNotExist)
	}

	for _, post := range posts {
		lastID, err := r.createPost(tx, post, threadID, forum, created)
		if err != nil {
			tx.Rollback()
			return []Post{}, err
		}

		returnPost := Post{
			Author:   post.Author,
			Created:  created,
			Forum:    forum,
			ID:       lastID,
			IsEdited: false,
			Message:  post.Message,
			Parent:   post.Parent,
			Thread:   threadID,
		}
		returnPosts = append(returnPosts, returnPost)
	}

	tx.Commit()
	return returnPosts, nil
}

func (r *Repository) createPost(tx *sql.Tx, post NewPost, threadID int, forum string, created string) (int64, error) {
	var lastID int64
	var err error
	if post.Parent == 0 {
		stmt, _ := tx.Prepare(sql_queries.InsertPostWithoutParent)
		defer stmt.Close()

		err = stmt.QueryRow(post.Author, post.Message, threadID, forum, created).Scan(&lastID)
	} else {
		stmt, _ := tx.Prepare(sql_queries.InsertPost)
		defer stmt.Close()

		err = stmt.QueryRow(post.Author, post.Message, post.Parent, threadID, forum, created).Scan(&lastID)
	}

	if err != nil {
		fmt.Printf("createPost %s", err.Error())
		_, err = r.GetUserByNickname(post.Author)
		if err != nil {
			return lastID, fmt.Errorf(messages.UserNotFound)
		} else {
			return lastID, fmt.Errorf(messages.ParentInAnotherThread)
		}
	}

	return lastID, nil
}

func (r *Repository) CreateVote(vote Vote, slugOrID string) error {
	threadID, err := strconv.Atoi(slugOrID)
	if err != nil {
		_, err = r.DbConn.Exec(sql_queries.InsertVoteByThreadSlug, vote.Nickname, vote.Voice,
			slugOrID)
	} else {
		_, err = r.DbConn.Exec(sql_queries.InsertVoteByThreadID, vote.Nickname, vote.Voice,
			threadID)
	}
	// if err != nil {
	// 	fmt.Printf("Rep CreateVote: %s\n", err.Error())
	// }

	return err
}

func (r *Repository) CreateForum(forum NewForum) error {
	_, err := r.DbConn.Exec(sql_queries.InsertForum, forum.Slug, forum.Title, forum.User)
	// if err != nil {
	// 	fmt.Printf("CreateForum error: %s", err.Error())
	// }
	return err
}

func (r *Repository) CreateThread(thread NewThread, forum string) (int, error) {
	var id int
	var err error
	if thread.Slug == "" {
		if thread.Created == "" {
			err = r.DbConn.QueryRow(sql_queries.InsertThreadWithoutCreated, thread.Author,
				thread.Message, thread.Title, forum).Scan(&id)
		} else {
			err = r.DbConn.QueryRow(sql_queries.InsertThread, thread.Author, thread.Created,
				thread.Message, thread.Title, forum).Scan(&id)
		}
	} else {
		if thread.Created == "" {
			err = r.DbConn.QueryRow(sql_queries.InsertThreadWithSlugWithoutCreated, thread.Author,
				thread.Message, thread.Title, forum, thread.Slug).Scan(&id)
		} else {
			err = r.DbConn.QueryRow(sql_queries.InsertThreadWithSlug, thread.Author, thread.Created,
				thread.Message, thread.Title, forum, thread.Slug).Scan(&id)
		}
	}

	return id, err
}

func (r *Repository) CreateUser(user NewUser, nickname string) error {
	_, err := r.DbConn.Exec(sql_queries.InsertUser, user.About,
		user.Email, user.Fullname, nickname)
	// if err != nil {
	// 	fmt.Printf("Rep CreateUser: %s\n", err)
	// }
	return err
}
