package repository

import (
	"bytes"
	"fmt"
	"github.com/db_project/pkg/sql_queries"
	"strconv"
	"strings"
	"text/template"
)

var selectPostsTemplate *template.Template
var selectPostsParentTreeTemplate *template.Template

func init() {
	var err error
	selectPostsTemplate, err = template.New("getPosts").Parse(sql_queries.SelectPostsSorted)
	if err != nil {
		fmt.Println("Error: cannot create getPostsTemplate template: ", err)
		panic(err)
	}

	selectPostsParentTreeTemplate, err = template.New("parent_tree").Parse(sql_queries.SelectPostsParentTree)
	if err != nil {
		fmt.Println("Error: cannot create getPostsParentTreeTemplate template: ", err)
		panic(err)
	}
}

func paramsThreadsToQuery(params map[string]interface{}) string {
	query := sql_queries.SelectThreadsWithParams
	if params["since"] != nil {
		if params["desc"] == "true" {
			query += "AND t.created <= :since "
		} else {
			query += "AND t.created >= :since "
		}
	}

	query += `ORDER BY t.created `

	if params["desc"] == "true" {
		query += "DESC "
	}

	if params["limit"] != nil {
		query += "LIMIT :limit "
	}

	query += ";"
	return query
}

func paramsGetUsers(params map[string]interface{}) string {
	query := sql_queries.SelectUsersWithParams
	if params["since"] != nil {
		if params["desc"] == "true" {
			query += "AND lower(p.nickname) < lower(:since) "
		} else {
			query += "AND lower(p.nickname) > lower(:since) "
		}
	}

	query += `ORDER BY lower(p.nickname) `

	if params["desc"] != nil {
		query += "DESC "
	}

	if params["limit"] != nil {
		query += "LIMIT :limit "
	}

	query += ";"
	return query
}

func createPacketQuery(prefix string, batchSize int, batchCount int, postfix ...string) string {
	pack := make([]string, 0, batchCount)
	batch := make([]string, 0, batchSize)

	for i := 0; i < batchCount; i++ {
		for j := 1; j <= batchSize; j++ {
			batch = append(batch, "$"+strconv.Itoa(batchSize*i+j))
		}
		pack = append(pack, "("+strings.Join(batch, ", ")+")")
		batch = batch[:0]
	}

	var res bytes.Buffer
	res.WriteString(prefix)
	if prefix[len(prefix)-1] != ' ' {
		res.WriteString(" ")
	}
	res.WriteString(strings.Join(pack, ", "))
	if len(postfix) > 0 {
		res.WriteString(" ")
		res.WriteString(strings.Join(postfix, " "))
	}
	res.WriteString(";")
	return res.String()
}

func createPostsQuery(threadID int64, limit int64, since string, sort string, desc bool) (string, []interface{}, error) {
	mainTemplateArgs := struct {
		Condition string
		OrderBy   string
		Limit     string
	}{}

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
		err = selectPostsParentTreeTemplate.Execute(conditionBuffer, struct {
			Since string
			Desc  string
			Limit string
		}{Since: placeholderSince, Desc: placeholderDesc, Limit: mainTemplateArgs.Limit})
		if err != nil {
			return "", nil, err
		}
		mainTemplateArgs.Condition = conditionBuffer.String()
		mainTemplateArgs.OrderBy = fmt.Sprintf(`path[1] %s, path`, placeholderDesc)
		mainTemplateArgs.Limit = ""
	}

	queryBuffer := &bytes.Buffer{}
	err = selectPostsTemplate.Execute(queryBuffer, mainTemplateArgs)
	if err != nil {
		return "", nil, err
	}
	query := queryBuffer.String()

	return query, params, nil
}
