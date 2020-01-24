package repository

import (
	"fmt"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetThreadByID(id int64) (Thread, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectThreadByID, id)

	var thread Thread
	var slug pgtype.Text

	err := row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message,
		&slug, &thread.Title, &thread.Votes)
	thread.Slug = slug.String
	if err != nil {
		// fmt.Println(err)
		return thread, fmt.Errorf(messages.ThreadDoesNotExist)
	}
	return thread, err
}

func (r *Repository) GetThreadBySlug(threadSlug string) (Thread, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectThreadBySlug, threadSlug)

	var thread Thread
	var slug pgtype.Text

	err := row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message,
		&slug, &thread.Title, &thread.Votes)
	thread.Slug = slug.String
	if err != nil {
		// fmt.Println(err)
		return thread, fmt.Errorf(messages.ThreadDoesNotExist)
	}
	return thread, err
}

func (r *Repository) GetPostByID(ID int64) (Post, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectPostByID, ID)
	var parent pgtype.Int8

	var post Post
	err := row.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited,
		&post.Message, &parent, &post.Thread)
	post.Parent = parent.Int
	if err != nil {
		// fmt.Println(err)
		return post, fmt.Errorf(messages.PostDoesNotExist)
	}

	return post, nil
}

func (r *Repository) GetThreads(params map[string]interface{}) ([]Thread, error) {
	threads := []Thread{}

	queryStr := paramsThreadsToQuery(params)

	query, args, err := sqlx.Named(queryStr, params)
	if err != nil {
		// fmt.Printf("Rep GetThreads: %s", err.Error())
		return threads, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		// fmt.Printf("Rep GetThreads: %s", err.Error())
		return threads, err
	}

	query = r.DbConn.Rebind(query)

	rows, err := r.DbConn.Queryx(query, args...)

	if err != nil {
		// fmt.Printf("Rep GetThreads: %s", err.Error())
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		var thread Thread
		var slug pgtype.Text

		err := rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message,
			&slug, &thread.Title, &thread.Votes)
		thread.Slug = slug.String
		if err != nil {
			// fmt.Printf("Rep GetThreads: %s", err.Error())
			return threads, fmt.Errorf(messages.ThreadDoesNotExist)
		}

		threads = append(threads, thread)
	}

	if len(threads) == 0 {
		_, err = r.GetForumBySlug(params["forum"].(string))
	}
	return threads, err
}

func (r *Repository) GetUsers(params map[string]interface{}) ([]User, error) {
	users := []User{}

	queryStr := paramsGetUsers(params)

	query, args, err := sqlx.Named(queryStr, params)
	if err != nil {
		// fmt.Printf("GetUsers: %s\n", err.Error())
		return users, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		// fmt.Printf("GetUsers: %s\n", err.Error())
		return users, err
	}

	query = r.DbConn.Rebind(query)

	rows, err := r.DbConn.Queryx(query, args...)

	if err != nil {
		// fmt.Printf("GetUsers: %s\n", err.Error())
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.StructScan(&user)
		if err != nil {
			// fmt.Printf("GetUsers: %s\n", err)
			return users, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		_, err = r.GetForumBySlug(params["forum"].(string))
	}

	return users, err
}

func (r *Repository) GetForumBySlug(slug string) (Forum, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectForumBySlug, slug)

	var forum Forum
	err := row.StructScan(&forum)
	if err != nil {
		// fmt.Println(err)
		return forum, err //fmt.Errorf()
	}

	return forum, nil
}

func (r *Repository) GetVoteByThreadID(nickname string, thread int64) (Vote, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectVoteByThreadID, nickname, thread)

	var vote Vote
	err := row.StructScan(&vote)
	if err != nil {
		// fmt.Println(err)
		return vote, err //fmt.Errorf()
	}

	return vote, nil
}

func (r *Repository) GetVoteByThreadSlug(nickname string, slug string) (Vote, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectVoteByThreadSlug, nickname, slug)

	var vote Vote
	err := row.StructScan(&vote)
	if err != nil {
		// fmt.Println(err)
		return vote, err //fmt.Errorf()
	}

	return vote, nil
}

func (r *Repository) GetUserByNickname(nickname string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByNickname, nickname)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		// fmt.Printf("Repository GetUserByNickname: %s\n", err)
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByEmail, email)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		// fmt.Printf("Repository GetUserByEmail: %s\n", err)
		return user, err
	}

	return user, nil
}

func (r *Repository) GetPosts(threadID int64, limit int64, since string, sort string, desc bool) ([]Post, error) {
	posts := make([]Post, 0)

	query, params, err := createPostsQuery(threadID, limit, since, sort, desc)
	if err != nil {
		return posts, err
	}

	rows, err := r.DbConn.Query(query, params...)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var parent pgtype.Int8

		err := rows.Scan(&post.Author, &post.Forum, &post.Created, &post.ID, &post.IsEdited,
			&post.Message, &parent, &post.Thread)
		post.Parent = parent.Int
		if err != nil {
			return posts, err
		}

		if post.Parent == post.ID {
			post.Parent = 0
		}

		posts = append(posts, post)
	}

	if len(posts) == 0 {
		_, err = r.GetThreadByID(threadID)
	}

	return posts, err
}
