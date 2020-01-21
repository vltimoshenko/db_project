package repository

import (
	"fmt"
	"strconv"

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

// func (r *Repository) ChangeThreadRate(dif int, threadID int) error {
// 	_, err := r.DbConn.Exec(sql_queries.UpdateThreadRating, dif, threadID)
// 	return err
// }

func (r *Repository) ChangeVote(updateVote Vote, slugOrID string) error {
	threadID, err := strconv.Atoi(slugOrID)
	if err != nil {
		_, err = r.DbConn.Exec(sql_queries.UpdateVoteByThreadSlug, updateVote.Voice, updateVote.Nickname,
			slugOrID)
	} else {
		_, err = r.DbConn.Exec(sql_queries.UpdateVoteByThreadID, updateVote.Voice, updateVote.Nickname,
			threadID)
	}
	return err

}

func (r *Repository) ChangeUser(user NewUser, nickname string) error {
	_, err := r.DbConn.Exec(sql_queries.UpdateUserByNickname, user.About,
		user.Email, user.Fullname, nickname) //should read id?

	if err != nil {
		// fmt.Println(err.Error())
		return fmt.Errorf(messages.UserAlreadyExists)
	}
	return nil
}

func (r *Repository) ChangePost(postUpdate PostUpdate, postID int64, isEdited bool) error {
	_, err := r.DbConn.Exec(sql_queries.UpdatePost, postUpdate.Message,
		isEdited, postID)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return fmt.Errorf(messages.ThreadDoesNotExist)
	// }

	return err
}
