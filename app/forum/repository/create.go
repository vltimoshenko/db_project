package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

// func (r *Repository) CreateThread(thread NewThread, forum string) (int, error) {
// 	var id int
// 	err := r.DbConn.QueryRow(sql_queries.InsertThread, thread.Author,
// 		thread.Message, thread.Title, forum, thread.Slug).Scan(&id)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return id, fmt.Errorf(messages.UserNotFound)
// 	}
// 	fmt.Println(id)
// 	return id, nil
// }

func (r *Repository) CreatePosts(posts []NewPost, threadID int, forum string) ([]Post, error) {
	tx, _ := r.DbConn.Begin()

	created := time.Now()
	// .Format(time.RFC3339Nano)
	returnPosts := []Post{}

	for _, post := range posts {
		if post.Parent != 0 {
			// _, err := r.GetPostByIDThreadID(post.Parent, threadID)
			_, err := r.GetPostByID(post.Parent)

			if err != nil {
				tx.Rollback()
				return []Post{}, fmt.Errorf(messages.ParentPostDoesNotExist)
			}
		}

		lastID, err := r.createPost(tx, post, threadID, forum, created)
		if err != nil {
			tx.Rollback()
			return []Post{}, err
		}

		returnPost := Post{
			Author:   post.Author,
			Created:  created.String(),
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

func (r *Repository) createPost(tx *sql.Tx, post NewPost, threadID int, forum string, created time.Time) (int, error) {
	var lastID int

	if post.Parent != 0 {
		parentPost, err := r.GetPostByID(post.Parent)

		if err != nil || parentPost.Thread != threadID {
			return 0, fmt.Errorf(messages.ParentPostDoesNotExist) //
		}
	}

	stmt, err := tx.Prepare(sql_queries.InsertPost)
	if err != nil {
		fmt.Println(err)
		return lastID, fmt.Errorf(messages.UserNotFound)
	}

	defer stmt.Close()

	err = stmt.QueryRow(post.Author, post.Message, post.Parent, threadID, forum, created).Scan(&lastID)
	if err != nil {
		fmt.Println(err)
		return lastID, fmt.Errorf(messages.UserNotFound)
	}
	return lastID, nil
}

func (r *Repository) CreateVote(vote Vote, threadID int) error {
	_, err := r.DbConn.Exec(sql_queries.InsertVote, vote.Nickname, vote.Voice,
		threadID)
	return err
}
