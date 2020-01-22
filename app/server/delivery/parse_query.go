package delivery

import (
	"net/url"
	"strconv"
	// . "github.com/db_project/pkg/models"
)

func (h *Handler) ParseThreadsAndUsersQuery(query url.Values) map[string]interface{} {
	params := make(map[string]interface{})

	if query.Get("limit") != "" {
		params["limit"], _ = strconv.Atoi(query.Get("limit"))
	}

	if query.Get("desc") == "true" {
		params["desc"] = query.Get("desc")
	}
	if query.Get("since") != "" {
		params["since"] = query.Get("since")
	}

	return params
}

func (h *Handler) ParsePostsQuery(query url.Values) map[string]interface{} {
	params := make(map[string]interface{})

	if query.Get("limit") != "" {
		params["limit"] = query.Get("limit")
	}

	if query.Get("desc") == "true" {
		params["desc"] = query.Get("desc")
	}
	if query.Get("since") != "" {
		params["since"] = query.Get("since")
	}

	if query.Get("sort") != "" {
		params["sort"] = query.Get("sort")
	}

	return params
}
