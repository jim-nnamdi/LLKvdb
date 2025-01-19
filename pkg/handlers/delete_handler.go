package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jim-nnamdi/Lkvs/pkg/model"
)

var _ http.Handler = &deleteHandler{}

type deleteHandler struct {
	Fsys *model.Filesys
}

func NewDeleteHandler(Fsys *model.Filesys) *deleteHandler {
	return &deleteHandler{Fsys: Fsys}
}

func (handler *deleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		key         = r.FormValue("key")
		DataRemoved = "Data successfully deleted and removed from disk!"
	)
	keyval, _ := strconv.ParseInt(key, 10, 64)
	handler.Fsys.Delete(keyval)
	response := map[string]interface{}{}
	response["message"] = DataRemoved
	json.NewEncoder(w).Encode(response)
}
