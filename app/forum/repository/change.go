package repository

import (
	"database/sql"
	"strconv"

	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
	"github.com/jackc/pgtype"
)

func (r *Repository) ChangeThread(threadUpdate ThreadUpdate, slugOrID string) (Thread, error) {
	threadID, err := strconv.ParseInt(slugOrID, 10, 64)

	var row *sql.Row
	if err == nil {
		row = r.DbConn.QueryRow(sql_queries.UpdateThreadByID, threadUpdate.Message,
			threadUpdate.Title, threadID)
	} else {
		row = r.DbConn.QueryRow(sql_queries.UpdateThreadBySlug, threadUpdate.Message,
			threadUpdate.Title, slugOrID)
	}

	var thread Thread
	var slug pgtype.Text

	err = row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message,
		&slug, &thread.Title, &thread.Votes)
	thread.Slug = slug.String

	return thread, err
}

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

func (r *Repository) ChangeUser(user NewUser, nickname string) (User, error) {
	var retUser User
	err := r.DbConn.QueryRow(sql_queries.UpdateUserByNickname, nickname,
		user.About, user.Email, user.Fullname).Scan(&retUser.Nickname, &retUser.Fullname, &retUser.Email, &retUser.About)

	return retUser, err
}

func (r *Repository) ChangePost(postUpdate PostUpdate, postID int64) (Post, error) {
	row := r.DbConn.QueryRowx(sql_queries.UpdatePost, postUpdate.Message, postID)
	// "RETURNING author::text, created, forum, is_edited, thread, message, parent"
	var post Post
	var parent pgtype.Int8
	err := row.Scan(&post.Author, &post.Forum, &post.Created, &post.ID, &post.IsEdited,
		&post.Message, &parent, &post.Thread)
	post.Parent = parent.Int
	return post, err
}
