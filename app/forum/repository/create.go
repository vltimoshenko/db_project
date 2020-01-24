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

const PacketSize = 30

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

var mutexMapMutex sync.Mutex

func (r *Repository) CreatePosts(posts []Post, threadID int64, forum string) ([]Post, error) {
	userList := make(map[string]bool)
	if len(posts) == 0 {
		return posts, nil
	}
	created := time.Now().Format(time.RFC3339)

	for i := 0; i < len(posts); i += PacketSize {
		currentPacket := posts[i:int(math.Min(float64(i+PacketSize), float64(len(posts))))]
		currentPacket, err := r.createPostsByPacket(threadID, forum, currentPacket, created)
		if err != nil {
			// fmt.Printf("Rep CreatePosts: %s\n", err.Error())
			return posts, err
		}

		for j, post := range currentPacket {
			posts[i+j] = post
			userList[post.Author] = true
		}
	}

	query := createPacketQuery(sql_queries.InsertForumUsers, 2, len(userList), `ON CONFLICT DO NOTHING`)
	params := make([]interface{}, 0, len(userList))
	for key := range userList {
		params = append(params, key, forum)
	}

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

	query := createPacketQuery(sql_queries.InsertPosts, 6, len(posts), ` RETURNING id;`)

	rows, err := r.DbConn.Query(query, params...)

	if err != nil || (rows != nil && rows.Err() != nil) {
		// fmt.Printf("createPostsByPacket: %s\n", err.Error())
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
			// fmt.Printf("createPostsByPacket: %s\n", err.Error())
			return posts, err
		}
		posts[i].Forum = forumSLug
		posts[i].Created = created
		posts[i].IsEdited = false
		posts[i].Thread = threadId
		i++
	}

	if i == 0 && len(posts) > 0 {
		_, err := r.GetThreadByID(threadId)
		if err != nil {
			// fmt.Println("createPostsByPacket: i == 0 && len(posts) > 0")
			return posts, fmt.Errorf(messages.ThreadDoesNotExist)
		}

		_, err = r.GetUserByNickname(posts[0].Author)
		if err != nil {
			// fmt.Printf("createPostsByPacket: %s\n", err.Error())
			return posts, fmt.Errorf(messages.UserNotFound)
		}
		// fmt.Println("createPostsByPacket: messages.ParentInAnotherThread")
		return posts, fmt.Errorf(messages.ParentInAnotherThread)
	}

	return posts, nil
}
