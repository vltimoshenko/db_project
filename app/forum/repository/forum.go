package repository

import (
	"fmt"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

func (r *Repository) CreateForum(forum NewForum) error {
	var id int
	err := r.DbConn.QueryRow(sql_queries.InsertForum, forum.Slug, forum.Title, forum.User).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf(messages.UserNotFound)
	}
	fmt.Println(id)
	return nil
}

func (r *Repository) CreateThread(thread NewThread, forum string) (int, error) {
	var id int
	err := r.DbConn.QueryRow(sql_queries.InsertThread, thread.Author,
		thread.Message, thread.Title, forum, thread.Slug).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		return id, fmt.Errorf(messages.UserNotFound)
	}
	fmt.Println(id)
	return id, nil
}
