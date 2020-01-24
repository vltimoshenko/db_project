package delivery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/db_project/pkg/messages"
	. "github.com/db_project/pkg/models"
	"github.com/gorilla/mux"
)

func (h *Handler) ChangeThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;")

	slugOrID, _ := mux.Vars(r)["slug_or_id"]
	// if !ok {
	// 	//
	// }

	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// if err != nil {
	// 	//
	// }

	var threadUpdate ThreadUpdate
	err := json.Unmarshal(bytes, &threadUpdate)
	// if err != nil {

	// }

	thread, err := h.Service.ChangeThread(threadUpdate, slugOrID)

	if err != nil {
		// log.Println(err.Error())
		SetError(w, 404, messages.ThreadDoesNotExist)
		return
	}

	answer, _ := json.Marshal(thread)

	w.WriteHeader(200)
	w.Write(answer)
}

func (h *Handler) ChangeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;")

	nickname, _ := mux.Vars(r)["nickname"]
	// if !ok {
	// 	//
	// }

	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// if err != nil {
	// 	//
	// }

	var newUser NewUser
	_ = json.Unmarshal(bytes, &newUser)
	// if err != nil {

	// }

	code := 200
	user, err := h.Service.ChangeUser(newUser, nickname)
	if err != nil {
		if err.Error() == messages.UserAlreadyExists {
			code = 409
		}
		if err.Error() == messages.UserNotFound {
			SetError(w, 404, err.Error())
			return
		}
		// log.ln(err.Error())
	}

	answer, _ := json.Marshal(user)

	w.WriteHeader(code)
	w.Write(answer)
}

func (h *Handler) ChangePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;")

	idStr, _ := mux.Vars(r)["id"]
	// if !ok {
	// 	//
	// }

	id, err := strconv.ParseInt(idStr, 10, 64)
	// if err != nil {
	// 	//
	// }

	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// if err != nil {
	// 	//
	// }

	var postUpdate PostUpdate
	err = json.Unmarshal(bytes, &postUpdate)
	// if err != nil {

	// }

	post, err := h.Service.ChangePost(postUpdate, id)

	if err != nil {
		// log.Println(err.Error())
		SetError(w, 404, messages.ThreadDoesNotExist)
		return
	}

	answer, _ := json.Marshal(post)

	w.WriteHeader(200)
	w.Write(answer)
}
