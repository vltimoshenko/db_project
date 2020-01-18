package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/db_project/pkg/messages"
	// . "github.com/db_project/pkg/models"

	"github.com/gorilla/mux"
)

func (h *Handler) GetThread(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	slugOrID, ok := mux.Vars(r)["slug_or_id"]
	if !ok {

	}

	forum, err := h.Service.GetThread(slugOrID)

	if err != nil {
		if err.Error() == messages.ThreadDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	w.WriteHeader(code)
	answer, _ := json.Marshal(forum)
	w.Write(answer)
}

func (h *Handler) GetForum(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	forumSlug, ok := mux.Vars(r)["slug"]
	if !ok {

	}

	forum, err := h.Service.GetForum(forumSlug)

	if err != nil {
		if err.Error() == messages.ThreadAlreadyExists {
			code = 409
		}

		if err.Error() == messages.UserNotFound {
			SetError(w, 404, err.Error())
			return
		}

		if err.Error() == messages.ForumDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	w.WriteHeader(code)
	answer, _ := json.Marshal(forum)
	w.Write(answer)
}

func (h *Handler) GetThreads(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	forumSlug, ok := mux.Vars(r)["slug"]

	if !ok {
		SetError(w, 404, fmt.Errorf(messages.ForumDoesNotExist).Error())
		return
	}

	params := h.ParseThreadsAndUsersQuery(r.URL.Query())
	params["forum"] = forumSlug

	threads, err := h.Service.GetThreads(params)

	if err != nil {
		if err.Error() == messages.ThreadAlreadyExists {
			code = 409
		}

		if err.Error() == messages.UserNotFound {
			SetError(w, 404, err.Error())
			return
		}

		if err.Error() == messages.ForumDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	w.WriteHeader(code)
	answer, _ := json.Marshal(threads)
	w.Write(answer)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	forumSlug, ok := mux.Vars(r)["slug"]

	if !ok {
		SetError(w, 404, fmt.Errorf(messages.ForumDoesNotExist).Error())
		return
	}

	params := h.ParseThreadsAndUsersQuery(r.URL.Query())
	params["forum"] = forumSlug

	users, err := h.Service.GetUsers(params)

	if err != nil {
		if err.Error() == messages.ForumDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	w.WriteHeader(code)
	answer, _ := json.Marshal(users)
	w.Write(answer)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	nickname, ok := mux.Vars(r)["nickname"]
	if !ok {

	}

	user, err := h.Service.GetUser(nickname)

	if err != nil {
		if err.Error() == messages.UserNotFound {
			SetError(w, 404, err.Error())
			return
		}
	}

	w.WriteHeader(code)
	answer, _ := json.Marshal(user)
	w.Write(answer)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		//
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		//
	}

	paramsMap, _ := url.ParseQuery(r.URL.Query().Encode())

	fmt.Printf("RELATED %d\n", len(paramsMap["related"]))
	var params []string
	for _, str := range paramsMap["related"] {
		params = append(params, str)
	}

	fmt.Printf("RELATED %s\n", params)

	var postInfo map[string]interface{}
	postInfo, err = h.Service.GetPost(id, params)

	if err != nil {
		if err.Error() == messages.PostDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	answer, _ := json.Marshal(postInfo)
	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	slugOrID, ok := mux.Vars(r)["slug_or_id"]

	if !ok {
		SetError(w, 404, fmt.Errorf(messages.ForumDoesNotExist).Error())
		return
	}

	params := h.ParsePostsQuery(r.URL.Query())

	posts, err := h.Service.GetPosts(slugOrID, params)

	if err != nil {
		if err.Error() == messages.ThreadDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	w.WriteHeader(code)
	answer, _ := json.Marshal(posts)
	w.Write(answer)
}