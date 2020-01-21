package repository

import (
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
