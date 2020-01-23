package delivery

import (
	"net/http"

	"github.com/db_project/app/forum"
	. "github.com/db_project/pkg/models"
)

type Handler struct {
	Service forum.ServiceInterface
}

func SetError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errJSON, _ := Error{msg}.MarshalJSON()
	w.Write(errJSON)
	return
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	_ = h.Service.ClearDB()
	// if err != nil {
	// 	if err.Error() == messages.UserAlreadyExists {
	// 		code = 409
	// 	}
	// 	// if err.Error() == messages.UserNotFound {
	// 	// 	SetError(w, 404, err.Error())
	// 	// 	return
	// 	// }
	// 	log.Println(err.Error())
	// }

	w.WriteHeader(code)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	code := 200

	status, _ := h.Service.GetStatus()
	// if err != nil {
	// 	SetError(w, 500, err.Error())
	// }
	answer, _ := status.MarshalJSON()

	w.WriteHeader(code)
	w.Write(answer)
}
