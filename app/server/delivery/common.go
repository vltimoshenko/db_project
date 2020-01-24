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
	w.Header().Set("Content-Type", "application/json;")

	_ = h.Service.ClearDB()
	w.WriteHeader(200)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;")

	status, _ := h.Service.GetStatus()
	// if err != nil {
	// 	SetError(w, 500, err.Error())
	// }

	answer, _ := status.MarshalJSON()
	w.WriteHeader(200)
	w.Write(answer)
}
