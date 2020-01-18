package delivery

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"

	"github.com/gorilla/mux"
)

func (h *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 201

	slugOrId, ok := mux.Vars(r)["slug_or_id"]
	if !ok {
	}

	forum, err := h.Service.CreatePosts(r.Body, slugOrId)

	if err != nil {
		if err.Error() == messages.UserNotFound {
			SetError(w, 409, err.Error())
			return
		}

		if err.Error() == messages.ThreadDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}

		if err.Error() == messages.ParentPostDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	answer, _ := json.Marshal(forum)
	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 201

	forum, err := h.Service.CreateForum(r.Body)
	if err != nil {
		if err.Error() == messages.ForumAlreadyExists {
			// code = 409
			SetError(w, 409, err.Error())
			return
		}

		if err.Error() == messages.UserNotFound {
			SetError(w, 404, err.Error())
			return
		}
	}

	answer, _ := json.Marshal(forum)
	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 201

	forumSlug, ok := mux.Vars(r)["slug"]
	if !ok {

	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
	}

	var thread NewThread
	err = json.Unmarshal(bytes, &thread)
	if err != nil {
	}

	forum, err := h.Service.CreateThread(thread, forumSlug)

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

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 201

	nickname, ok := mux.Vars(r)["nickname"]
	if !ok {
		//
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//
	}

	var newUser NewUser
	err = json.Unmarshal(bytes, &newUser)
	if err != nil {

	}

	user, err := h.Service.CreateUser(newUser, nickname)
	if err != nil {
		if err.Error() == messages.UserAlreadyExists {
			code = 409
		}
		// if err.Error() == messages.UserNotFound {
		// 	SetError(w, 404, err.Error())
		// 	return
		// }
		log.Println(err.Error())
	}

	answer, _ := json.Marshal(user)

	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	slugOrId, ok := mux.Vars(r)["slug_or_id"]
	if !ok {
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//
	}

	var vote Vote
	err = json.Unmarshal(bytes, &vote)
	if err != nil {

	}

	thread, err := h.Service.Vote(vote, slugOrId)

	if err != nil {
		if err.Error() == messages.ThreadDoesNotExist {
			SetError(w, 404, err.Error())
			return
		}
	}

	answer, _ := json.Marshal(thread)
	w.WriteHeader(code)
	w.Write(answer)
}
