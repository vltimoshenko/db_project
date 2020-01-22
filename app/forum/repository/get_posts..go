package repository

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	. "github.com/db_project/pkg/models"
	"github.com/jackc/pgtype"
)

const (
	queryTemplateGetPostsSorted = `SELECT author, forum, created, posts.id, is_edited, message, coalesce(parent, 0), thread 
				FROM posts
					{{.Condition}}
					ORDER BY {{.OrderBy}}
					{{.Limit}}`
	queryTemplateGetPostsParentTree = `JOIN (
						SELECT parents.id FROM posts AS parents
						WHERE parents.thread=$1 AND parents.parent IS NULL
							{{- if .Since}} AND {{.Since}}{{- end}}
						ORDER BY parents.path[1] {{.Desc}}
						{{.Limit}}
						) as p ON path[1]=p.id`
)

var getPostsTemplate *template.Template
var getPostsParentTreeTemplate *template.Template

func init() {
	var err error
	getPostsTemplate, err = template.New("getPosts").Parse(queryTemplateGetPostsSorted)
	if err != nil {
		fmt.Println("Error: cannot create getPostsTemplate template: ", err)
		panic(err)
	}

	getPostsParentTreeTemplate, err = template.New("parent_tree").Parse(queryTemplateGetPostsParentTree)
	if err != nil {
		fmt.Println("Error: cannot create getPostsParentTreeTemplate template: ", err)
		panic(err)
	}
}

func (r *Repository) GetPosts(threadID int64, limit int64, since string, sort string, desc bool) ([]Post, error) {
	mainTemplateArgs := struct {
		Condition string
		OrderBy   string
		Limit     string
	}{}

	posts := make([]Post, 0)

	params := make([]interface{}, 0, 2)
	params = append(params, threadID)
	var placeholderSince, placeholderDesc string
	if desc {
		placeholderDesc = "DESC"
	} else {
		placeholderDesc = "ASC"
	}

	if limit != 0 {
		params = append(params, limit)
		mainTemplateArgs.Limit = `LIMIT $` + strconv.Itoa(len(params))
	}

	if since != "" {
		params = append(params, since)
		var compareSign string
		if desc {
			compareSign = "<"
		} else {
			compareSign = ">"
		}
		paramNum := len(params)
		queryGetPath := `SELECT %s FROM posts AS since WHERE since.id=%s`
		switch sort {
		case "flat":
			placeholderSince = fmt.Sprintf(`AND id%s$%d`, compareSign, paramNum)
		case "tree":
			placeholderSince = fmt.Sprintf(
				`AND path%s(%s)`,
				compareSign,
				fmt.Sprintf(queryGetPath, `since.path`, fmt.Sprintf(`$%d`, paramNum)),
			)
		case "parent_tree":
			placeholderSince = fmt.Sprintf(
				`parents.path[1]%s(%s)`,
				compareSign,
				fmt.Sprintf(queryGetPath, `since.path[1]`, fmt.Sprintf(`$%d`, paramNum)),
			)
		}
	}

	var err error
	switch sort {
	case "flat":
		mainTemplateArgs.Condition = `WHERE thread=$1 ` + placeholderSince
		mainTemplateArgs.OrderBy = fmt.Sprintf(`(created, id) %s`, placeholderDesc)
	case "tree":
		mainTemplateArgs.OrderBy = fmt.Sprintf(`(path, created) %s`, placeholderDesc)
		mainTemplateArgs.Condition = `WHERE thread=$1 ` + placeholderSince
	case "parent_tree":
		conditionBuffer := &bytes.Buffer{}
		err = getPostsParentTreeTemplate.Execute(conditionBuffer, struct {
			Since string
			Desc  string
			Limit string
		}{Since: placeholderSince, Desc: placeholderDesc, Limit: mainTemplateArgs.Limit})
		if err != nil {
			return posts, err
		}
		mainTemplateArgs.Condition = conditionBuffer.String()
		mainTemplateArgs.OrderBy = fmt.Sprintf(`path[1] %s, path`, placeholderDesc)
		mainTemplateArgs.Limit = ""
	}

	queryBuffer := &bytes.Buffer{}
	err = getPostsTemplate.Execute(queryBuffer, mainTemplateArgs)
	if err != nil {
		return posts, err
	}
	query := queryBuffer.String()

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

	// if len(posts) == 0 {
	// 	var sl pgtype.Text
	// 	err = r.DbConn.QueryRow(`SELECT slug from threads WHERE id=$1`, threadID).Scan(&sl)
	// 	if err != nil {
	// 		return posts, err
	// 	}
	// }
	return posts, nil
}
