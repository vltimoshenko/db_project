package delivery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		w.WriteHeader(400)
		return
	}

	bytes, _ := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	//return uuid.UUID{}, errors.New(BadRequestMsg)
	// }

	var posts []Post
	_ = json.Unmarshal(bytes, &posts)
	// if err != nil {
	// 	//return uuid.UUID{}, errors.New(InvalidJSONMsg)
	// }

	forum, err := h.Service.CreatePosts(posts, slugOrId)

	if err != nil {
		if err.Error() == messages.ParentInAnotherThread || err.Error() == messages.ParentPostDoesNotExist {
			// fmt.Printf("Handler CreatePosts: %s", err.Error())
			SetError(w, 409, err.Error())
			return
		}

		if err.Error() == messages.ThreadDoesNotExist || err.Error() == messages.UserNotFound {
			// fmt.Printf("Handler CreatePosts: %s", err.Error())
			SetError(w, 404, err.Error())
			return
		}
		// fmt.Printf("Handler CreatePosts Unknown error: %s", err.Error())
	}

	answer, _ := json.Marshal(forum)
	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 201

	bytes, _ := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	//return uuid.UUID{}, errors.New(BadRequestMsg)
	// }

	var forum NewForum
	err := json.Unmarshal(bytes, &forum)
	// if err != nil {
	// 	//return uuid.UUID{}, errors.New(InvalidJSONMsg)
	// }

	//could remove to creation moment or after get

	retForum, err := h.Service.CreateForum(forum)

	var answer []byte
	if err != nil {
		if err.Error() == messages.ForumAlreadyExists {
			code = 409
			answer, _ = json.Marshal(retForum)
		}

		if err.Error() == messages.UserNotFound {
			SetError(w, 404, err.Error())
			return
		}
	}
	if code != 409 {
		answer, _ = json.Marshal(NewForum{
			Slug:  retForum.Slug,
			Title: retForum.Title,
			User:  retForum.User,
		}) //could be removed
	}

	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 201

	forumSlug, _ := mux.Vars(r)["slug"]
	// if !ok {

	// }

	bytes, _ := ioutil.ReadAll(r.Body)
	// if err != nil {
	// }

	var thread NewThread
	err := json.Unmarshal(bytes, &thread)
	// if err != nil {
	// }

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
	w.Header().Set("Content-Type", "application/json")
	code := 201

	nickname, _ := mux.Vars(r)["nickname"]
	// if !ok {
	// 	//
	// }

	bytes, _ := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	//
	// }

	var newUser NewUser
	_ = json.Unmarshal(bytes, &newUser)
	// if err != nil {

	// }

	users, err := h.Service.CreateUser(newUser, nickname)
	if err != nil {
		if err.Error() == messages.UserAlreadyExists {
			code = 409
		}
		// fmt.Println(err.Error())
	}

	var answer []byte
	if err != nil {
		answer, _ = json.Marshal(users)
	} else {
		answer, _ = json.Marshal(users[0])
	}

	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	slugOrId, _ := mux.Vars(r)["slug_or_id"]
	// if !ok {
	// }

	bytes, _ := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	//
	// }

	var vote Vote
	_ = json.Unmarshal(bytes, &vote)
	// if err != nil {

	// }

	thread, err := h.Service.Vote(vote, slugOrId)

	if err != nil {
		if err.Error() == messages.ThreadDoesNotExist || err.Error() == messages.UserNotFound {
			SetError(w, 404, err.Error())
			return
		}
	}

	answer, _ := json.Marshal(thread)
	w.WriteHeader(code)
	w.Write(answer)
}
