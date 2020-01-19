package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

func (r *Repository) CreatePosts(posts []NewPost, threadID int, forum string) ([]Post, error) {
	tx, _ := r.DbConn.Begin()

	created := time.Now().Format(time.RFC3339)
	// created := time.Now().Format("1970-01-01 03:00:00+03")
	returnPosts := []Post{}

	for _, post := range posts {
		// if post.Parent != 0 {
		// 	// _, err := r.GetPostByIDThreadID(post.Parent, threadID)
		// 	_, err := r.GetPostByID(post.Parent)

		// 	if err != nil {
		// 		tx.Rollback()
		// 		return []Post{}, fmt.Errorf(messages.ParentPostDoesNotExist)
		// 	}
		// }
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

func (r *Repository) createPost(tx *sql.Tx, post NewPost, threadID int, forum string, created string) (int, error) {
	var lastID int

	if post.Parent != 0 {
		parentPost, err := r.GetPostByID(post.Parent)

		if err != nil {
			fmt.Printf("createPost %s", err.Error())
			return 0, fmt.Errorf(messages.ParentPostDoesNotExist) //
		} else if parentPost.Thread != threadID {
			return 0, fmt.Errorf(messages.ParentInAnotherThread) //
		}
	}

	stmt, err := tx.Prepare(sql_queries.InsertPost)
	if err != nil {
		fmt.Printf("createPost %s", err.Error())
		return lastID, fmt.Errorf(messages.UserNotFound)
	}

	defer stmt.Close()

	err = stmt.QueryRow(post.Author, post.Message, post.Parent, threadID, forum, created).Scan(&lastID)
	if err != nil {
		fmt.Printf("createPost %s", err.Error())
		return lastID, fmt.Errorf(messages.UserNotFound)
	}
	return lastID, nil
}

func (r *Repository) CreateVote(vote Vote, threadID int) error {
	_, err := r.DbConn.Exec(sql_queries.InsertVote, vote.Nickname, vote.Voice,
		threadID)
	return err
}

func (r *Repository) CreateForum(forum NewForum) error {
	_, err := r.DbConn.Exec(sql_queries.InsertForum, forum.Slug, forum.Title, forum.User)
	if err != nil {
		fmt.Printf("CreateForum error: %s", err.Error())
	}
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
	// var id int
	_, _ = r.DbConn.Exec(sql_queries.InsertUser, user.About,
		user.Email, user.Fullname, nickname)

	// fmt.Println(row)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return fmt.Errorf(messages.UserAlreadyExists)
	// }

	// fmt.Println(id)
	return nil
}
