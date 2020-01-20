package repository

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetThreadsBySlug(slug string) ([]Thread, error) {
	threads := []Thread{}

	rows, err := r.DbConn.Queryx(sql_queries.SelectThreadBySlug, slug)
	if err != nil {
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		thread := Thread{}
		var timetz time.Time
		err := rows.Scan(&thread)
		if err != nil {
			return threads, err
		}
		thread.Created = timetz.Format(time.RFC3339Nano)
		threads = append(threads, thread)
	}
	return threads, nil
}

func (r *Repository) GetThreadsByID(id int) ([]Thread, error) {
	threads := []Thread{}

	rows, err := r.DbConn.Queryx(sql_queries.SelectThreadByID, id)
	if err != nil {
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		thread := Thread{}
		// var timetz time.Time
		err := rows.Scan(&thread)
		if err != nil {
			return threads, err
		}
		// thread.Created = timetz.Format(time.RFC3339Nano)
		threads = append(threads, thread)
	}
	return threads, nil
}

func (r *Repository) GetThreadByID(id int) (Thread, error) {
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

func (r *Repository) GetPostByID(ID int) (Post, error) {

	// row := r.DbConn.QueryRowx(sql_queries.SelectPostByIDThreadID, ID, threadID)
	row := r.DbConn.QueryRowx(sql_queries.SelectPostByID, ID)

	// var timetz time.Time
	var post Post
	err := row.StructScan(&post)
	if err != nil {
		fmt.Println(err)
		return post, fmt.Errorf(messages.PostDoesNotExist)
	}
	// scanPost.Created = timetz.Format(time.RFC3339Nano)

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

func (r *Repository) GetUsers(params map[string]interface{}) ([]User, error) { //TODO
	users := []User{}

	queryStr := paramsGetUsers(params)

	query, args, err := sqlx.Named(queryStr, params)
	if err != nil {
		log.Print(err)
		return users, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		log.Print(err)
		return users, err
	}

	query = r.DbConn.Rebind(query)

	rows, err := r.DbConn.Queryx(query, args...)

	if err != nil {
		log.Print(err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.StructScan(&user)
		if err != nil {
			log.Printf("GetUsers: %s\n", err)
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

func (r *Repository) GetVoteByThreadID(nickname string, thread int) (Vote, error) {
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

func (r *Repository) GetPosts(threadID int, limit, since, sort, desc string) (Posts []Post, Err error) {
	posts := []Post{}

	var rows *sqlx.Rows
	var err error
	if sort == "flat" {
		if desc == "false" {
			rows, err = r.DbConn.Queryx(sql_queries.SelectPostsFlat, threadID, limit, since)
		} else {
			rows, err = r.DbConn.Queryx(sql_queries.SelectPostsFlatDesc, threadID, limit, since)
		}

	} else if sort == "tree" {
		if desc == "false" {
			if since != "0" && since != "999999999" {
				rows, err = r.DbConn.Queryx(sql_queries.SelectPostsTree, threadID, 100000)
			} else {
				rows, err = r.DbConn.Queryx(sql_queries.SelectPostsTree, threadID, limit)
			}
		} else {
			if since != "0" && since != "999999999" {
				rows, err = r.DbConn.Queryx(sql_queries.SelectPostsTreeSinceDesc, threadID)
			} else {
				rows, err = r.DbConn.Queryx(sql_queries.SelectPostsTreeDesc, threadID, limit, 1000000)
			}
		}
	} else if sort == "parent_tree" {
		if desc == "false" {
			rows, err = r.DbConn.Queryx(sql_queries.SelectPostsParentTree, threadID)
		} else {
			rows, err = r.DbConn.Queryx(sql_queries.SelectPostsParentTreeDesc, threadID)
		}
	}

	if err != nil {
		return posts, err
	}
	defer rows.Close()

	if sort != "parent_tree" {
		for rows.Next() {
			scanPost := Post{}
			// var timetz time.Time
			err := rows.StructScan(&scanPost)
			if err != nil {
				return posts, err
			}
			// scanPost.Created = timetz.Format(time.RFC3339Nano)
			posts = append(posts, scanPost)
		}
	} else {
		count := 0
		limitDigit, _ := strconv.Atoi(limit)

		for rows.Next() {
			scanPost := Post{}
			// var timetz time.Time
			err := rows.StructScan(&scanPost)

			if err != nil {
				return posts, err
			}

			if scanPost.Parent == 0 {
				count = count + 1
			}
			if count > limitDigit && (since == "0" || since == "999999999") {
				break
			} else {
				// scanPost.Created = timetz.Format(time.RFC3339Nano)
				posts = append(posts, scanPost)
			}

		}
	}

	if since != "0" && since != "999999999" && sort == "tree" {
		limitDigit, _ := strconv.Atoi(limit)
		sinceDigit, _ := strconv.Atoi(since)
		sincePosts := []Post{}
		counter := 0

		if desc == "false" {
			startIndex := 1000000000
			minValue := 100000000000
			for i := 0; i < len(posts); i++ {
				if posts[i].ID == sinceDigit {
					startIndex = i + 1
					break
				}
				if (posts[i].ID > sinceDigit) && (posts[i].ID < minValue) {
					startIndex = i
					minValue = posts[i].ID
				}
			}
			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter < len(posts) {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter++
			}
		} else {
			startIndex := -1000000000
			maxValue := 0
			for i := len(posts) - 1; i >= 0; i-- {
				if posts[i].ID == sinceDigit {
					startIndex = i - 1
					break
				}
				if (posts[i].ID < sinceDigit) && (posts[i].ID > maxValue) {
					startIndex = i
					maxValue = posts[i].ID
				}
			}

			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter >= 0 {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter--
			}
		}
		return sincePosts, nil
	}

	if since != "0" && since != "999999999" && sort == "parent_tree" {
		limitDigit, _ := strconv.Atoi(limit)
		sinceDigit, _ := strconv.Atoi(since)
		sincePosts := []Post{}
		counter := 0
		if desc == "false" {
			startIndex := 1000000000
			minValue := 100000000000
			for i := 0; i < len(posts); i++ {
				if posts[i].ID == sinceDigit {
					startIndex = i + 1
					break
				}
				if (posts[i].ID > sinceDigit) && (posts[i].ID < minValue) {
					startIndex = i
					minValue = posts[i].ID
				}
			}
			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter < len(posts) {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				sincePostsCount++
				counter++
			}
		} else {
			startIndex := -1000000000
			maxValue := 100000000000
			for i := len(posts) - 1; i >= 0; i-- {
				if posts[i].ID == sinceDigit {
					startIndex = i + 1
					break
				}
				if (posts[i].ID < sinceDigit) && (posts[i].ID < maxValue) {
					startIndex = i
					maxValue = posts[i].ID
				}
			}

			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter < len(posts) {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter++
			}
		}
		return sincePosts, nil
	}

	return posts, nil
}

func (r *Repository) GetUserByNickname(nickname string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByNickname, nickname)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		fmt.Printf("Repository GetUserByNickname: %s", err)
		return user, err
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (User, error) {

	row := r.DbConn.QueryRowx(sql_queries.SelectUserByEmail, email)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		fmt.Printf("Repository GetUserByNickname: %s", err)
		return user, err
	}

	return user, nil
}
