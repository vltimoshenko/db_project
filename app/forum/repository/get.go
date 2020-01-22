package repository

import (
	"fmt"
"log"
	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetThreadByID(id int64) (Thread, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectThreadByID, id)

	// thread := Thread{}
	// // var timetz time.Time
	// err := row.StructScan(&thread)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return thread, fmt.Errorf(messages.ThreadDoesNotExist)
	// }
	// // thread.Created = timetz.Format(time.RFC3339Nano)
	// c := 'd'
	// fmt.Println(c)
	var thread Thread
	var slug pgtype.Text

	err := row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message,
		&slug, &thread.Title, &thread.Votes)
	thread.Slug = slug.String
	if err != nil {
		fmt.Println(err)
		return thread, fmt.Errorf(messages.ThreadDoesNotExist)
	}
	return thread, err
}

func (r *Repository) GetThreadBySlug(threadSlug string) (Thread, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectThreadBySlug, threadSlug)

	// var thread Thread
	// err := row.StructScan(&thread)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return thread, fmt.Errorf(messages.ThreadDoesNotExist)
	// }
	var thread Thread
	var slug pgtype.Text

	err := row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message,
		&slug, &thread.Title, &thread.Votes)
	thread.Slug = slug.String
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return post, fmt.Errorf(messages.PostDoesNotExist)
	}

	return post, nil
}

func (r *Repository) GetThreads(params map[string]interface{}) ([]Thread, error) {
	threads := []Thread{}

	queryStr := paramsThreadsToQuery(params)

	query, args, err := sqlx.Named(queryStr, params)
	if err != nil {
		log.Print(err)
		return threads, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		log.Print(err)
		return threads, err
	}

	query = r.DbConn.Rebind(query)

	rows, err := r.DbConn.Queryx(query, args...)

	if err != nil {
		log.Print(err)
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
			fmt.Println(err)
			return threads, fmt.Errorf(messages.ThreadDoesNotExist)
		}

		threads = append(threads, thread)
	}
	return threads, nil
}

// query := fmt.Sprintf(
// 	`select "user".* from "user"
// 			 join forum_user on nickname = forum_user.user
// 			where forum = $1 %s order by nickname %s %s`,
// 	sinceFilter, r.getOrder(desc), r.getLimit(limit),
// )

func (r *Repository) GetUsers(params map[string]interface{}) ([]User, error) {
	users := []User{}

	queryStr := paramsGetUsers(params)

	query, args, err := sqlx.Named(queryStr, params)
	if err != nil {
		fmt.Printf("GetUsers: %s\n", err.Error())
		return users, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		fmt.Printf("GetUsers: %s\n", err.Error())
		return users, err
	}

	query = r.DbConn.Rebind(query)

	rows, err := r.DbConn.Queryx(query, args...)

	if err != nil {
		fmt.Printf("GetUsers: %s\n", err.Error())
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.StructScan(&user)
		if err != nil {
			fmt.Printf("GetUsers: %s\n", err)
			return users, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) GetForumBySlug(slug string) (Forum, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectForumBySlug, slug)

	var forum Forum
	err := row.StructScan(&forum)
	if err != nil {
		fmt.Println(err)
		return forum, err //fmt.Errorf()
	}

	return forum, nil
}

func (r *Repository) GetVoteByThreadID(nickname string, thread int64) (Vote, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectVoteByThreadID, nickname, thread)

	var vote Vote
	err := row.StructScan(&vote)
	if err != nil {
		fmt.Println(err)
		return vote, err //fmt.Errorf()
	}

	return vote, nil
}

func (r *Repository) GetVoteByThreadSlug(nickname string, slug string) (Vote, error) {
	row := r.DbConn.QueryRowx(sql_queries.SelectVoteByThreadSlug, nickname, slug)

	var vote Vote
	err := row.StructScan(&vote)
	if err != nil {
		fmt.Println(err)
		return vote, err //fmt.Errorf()
	}

	return vote, nil
}

func (r *Repository) GetUserByNickname(nickname string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByNickname, nickname)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		fmt.Printf("Repository GetUserByNickname: %s\n", err)
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByEmail, email)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		fmt.Printf("Repository GetUserByEmail: %s\n", err)
		return user, err
	}

	return user, nil
}
