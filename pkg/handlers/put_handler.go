package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jim-nnamdi/Lkvs/pkg/model"
)

var _ http.Handler = &putHandler{}

type putHandler struct {
	Fsys *model.Filesys
}

func NewPutHandler(Fsys *model.Filesys) *putHandler {
	return &putHandler{Fsys: Fsys}
}

func (handler *putHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		key           = r.FormValue("key")
		val           = r.FormValue("value")
		DataPersisted = "Data successfully added and persisted to disk!"
	)
	keyval, _ := strconv.ParseInt(key, 10, 64)
	handler.Fsys.Put(int64(keyval), val)
	json.NewEncoder(w).Encode(DataPersisted)
}
