package repository

import (
	"fmt"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

func (r *Repository) ChangeThread(threadUpdate ThreadUpdate, threadID int) error {
	// var id int
	//should do two methods by slug and by id
	_ = r.DbConn.QueryRow(sql_queries.UpdateThreadByID, threadUpdate.Message,
		threadUpdate.Title, threadID) //should read id?

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return fmt.Errorf(messages.ThreadDoesNotExist)
	// }

	// fmt.Println(id)
	return nil
}

func (r *Repository) ChangeThreadRate(dif int, threadID int) error {
	_, err := r.DbConn.Exec(sql_queries.UpdateThreadRating, dif, threadID)
	return err
}

func (r *Repository) ChangeVote(updateVote Vote, threadID int) error {
	_, err := r.DbConn.Exec(sql_queries.UpdateVote, updateVote.Voice, updateVote.Nickname,
		threadID)
	return err
}

func (r *Repository) ChangeUser(user NewUser, nickname string) error {
	var id int
	err := r.DbConn.QueryRow(sql_queries.UpdateUserByNickname, user.About,
		user.Email, user.Fullname, nickname).Scan(&id) //should read id?

	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf(messages.UserAlreadyExists)
	}

	fmt.Println(id)
	return nil
}

func (r *Repository) ChangePost(postUpdate PostUpdate, postID int) error {
	// var id int
	//should do two methods by slug and by id
	_, err := r.DbConn.Exec(sql_queries.UpdatePost, postUpdate.Message,
		true, postID) //should read id?

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return fmt.Errorf(messages.ThreadDoesNotExist)
	// }

	// fmt.Println(id)
	return err
}