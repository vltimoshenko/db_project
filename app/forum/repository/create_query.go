package repository

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/db_project/pkg/sql_queries"
)

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
	// fmt.Println(query)
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
	// fmt.Println(query)
	return query
}

func CreatePacketQuery(prefix string, batchSize int, batchCount int, postfix ...string) string {
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
