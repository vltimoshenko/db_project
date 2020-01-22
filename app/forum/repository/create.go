package repository

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/db_project/pkg/sql_queries"
)

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

func (r *Repository) CreateThread(thread NewThread, forum string) (int64, error) {
	var id int64
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

func init() {
	mutexMapMutex = sync.Mutex{}
}

var postCounter int64

var mutexMapMutex sync.Mutex

func (r *Repository) CreatePosts(posts []Post, threadID int64, forum string) ([]Post, error) {
	// postCounter++
	// threadId, err := r.getThreadId(thread)
	// if err != nil {
	// 	return posts, structs.InternalError{E: structs.ErrorNoThread}
	// }

	// var cnt int64
	// if row := r.DB.QueryRow(context.Background(), `SELECT count(id) from Thread WHERE id=$1;`, threadId); row.Scan(&cnt) != nil || cnt == 0 {
	// 	return posts, structs.InternalError{E: structs.ErrorNoThread}
	// }
	// if len(posts) == 0 {
	// 	return posts, nil
	// }
	// var forumSlug string
	// err = r.DB.QueryRow(context.Background(), `SELECT forum FROM Thread WHERE Thread.id=$1`, threadId).Scan(&forumSlug)
	// if err != nil {
	// 	return posts, structs.InternalError{E: structs.ErrorNoThread, Explain: err.Error()}
	// }

	userList := make(map[string]bool)
	if len(posts) == 0 {
		return posts, nil
	}
	postPacketSize := 30
	created := time.Now().Format(time.RFC3339)

	for i := 0; i < len(posts); i += postPacketSize {
		currentPacket := posts[i:int(math.Min(float64(i+postPacketSize), float64(len(posts))))]
		currentPacket, err := r.createPostsByPacket(threadID, forum, currentPacket, created)
		if err != nil {
			fmt.Printf("Rep CreatePosts: %s\n", err.Error())
			return posts, err
		}

		for j, post := range currentPacket {
			posts[i+j] = post
			userList[post.Author] = true
		}
	}

	// query := `UPDATE ForumPosts SET posts=posts+$2 WHERE forum=$1;`
	// _, err := r.DbConn.Exec(query, forumSlug, len(posts))
	// if err != nil {
	// 	return posts, structs.InternalError{E: structs.ErrorNoThread, Explain: err.Error()}
	// }
	// atomic.AddInt32(&forumPostsAccess.hasNewUpdates, 1)

	prefix := `INSERT INTO forum_users(person, forum) VALUES `
	postfix := `ON CONFLICT DO NOTHING`
	query := CreatePacketQuery(prefix, 2, len(userList), postfix)
	params := make([]interface{}, 0, len(userList))
	for key := range userList {
		params = append(params, key, forum)
	}
	// GetMutex(forumSlug)
	// defer FreeMutex(forumSlug)
	mutexMapMutex.Lock()
	defer mutexMapMutex.Unlock()
	_, err := r.DbConn.Exec(query, params...)
	if err != nil {
		return posts, fmt.Errorf(messages.ThreadDoesNotExist)
	}
	return posts, nil
}

func (r *Repository) createPostsByPacket(threadId int64, forumSLug string, posts []Post, created string) ([]Post, error) {
	var params []interface{}

	for _, post := range posts {
		var parent sql.NullInt64
		parent.Int64 = post.Parent
		if post.Parent != 0 {
			parent.Valid = true
		}

		params = append(params, post.Author, post.Message, parent, threadId, created, forumSLug)
	}

	query := sql_queries.InsertPosts
	postfix := ` RETURNING id;`

	query = CreatePacketQuery(query, 6, len(posts), postfix) //10

	rows, err := r.DbConn.Query(query, params...)
	fmt.Println("createPostsByPacket: QueryComplete")

	if err != nil || (rows != nil && rows.Err() != nil) {
		fmt.Printf("createPostsByPacket: %s\n", err.Error())
		if strings.Contains(err.Error(), "post_parent_constraint") {
			return posts, fmt.Errorf(messages.ParentInAnotherThread)
		} else {
			return posts, fmt.Errorf(messages.UserNotFound)
		}
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(&(posts[i].ID))
		if err != nil {
			fmt.Printf("createPostsByPacket: %s\n", err.Error())
			return posts, err //undef error
		}
		posts[i].Forum = forumSLug
		posts[i].Created = created
		posts[i].IsEdited = false
		posts[i].Thread = threadId
		i++
	}

	// if i == 0 && len(posts) > 0 {
	// 	_, err := r.GetThreadByID(threadId) //conv
	// 	if err != nil {
	// 		fmt.Println("createPostsByPacket: i == 0 && len(posts) > 0")
	// 		return posts, fmt.Errorf(messages.ThreadDoesNotExist)
	// 	}

	// 	_, err = r.GetUserByNickname(posts[0].Author)
	// 	if err != nil {
	// 		fmt.Printf("createPostsByPacket: %s\n", err.Error())
	// 		return posts, fmt.Errorf(messages.UserNotFound)
	// 	}
	// 	fmt.Println("createPostsByPacket: messages.ParentInAnotherThread")
	// 	return posts, fmt.Errorf(messages.ParentInAnotherThread)
	// }

	var cnt int64
	if i == 0 && len(posts) > 0 {
		fmt.Println("createPostsByPacket: i == 0 && len(posts) > 0")
		if row := r.DbConn.QueryRow(`SELECT count(id) from threads WHERE id=$1;`, threadId); row.Scan(&cnt) != nil || cnt == 0 {
			// fmt.Printf("createPostsByPacket: %s\n", err.Error())
			return posts, fmt.Errorf(messages.UserNotFound)
		} else if row := r.DbConn.QueryRow(`SELECT COUNT(nickname) FROM persons WHERE nickname=$1`, posts[0].Author); row.Scan(&cnt) != nil || cnt == 0 {
			// fmt.Printf("createPostsByPacket: %s\n", err.Error())
			return posts, fmt.Errorf(messages.UserNotFound)
		} else {
			fmt.Println("createPostsByPacket: messages.ParentInAnotherThread")
			return posts, fmt.Errorf(messages.ParentInAnotherThread)
		}
	}
	fmt.Printf("createPostsByPacket: i = %d, len(posts) = %d", i, len(posts))
	return posts, nil
}
