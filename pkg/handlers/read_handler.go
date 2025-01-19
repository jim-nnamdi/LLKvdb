package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jim-nnamdi/Lkvs/pkg/model"
)

var _ http.Handler = &readHandler{}

type readHandler struct {
	Fsys *model.Filesys
}

func NewReadHandler(Fsys *model.Filesys) *readHandler {
	return &readHandler{Fsys: Fsys}
}

func (handler *readHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		key             = r.FormValue("key")
		DataUnavailable = "Associated Data for key could not be found on disk!"
	)
	keyval, _ := strconv.ParseInt(key, 10, 64)
	vals, ok := handler.Fsys.Read(int64(keyval))
	if !ok {
		fmt.Println("err", DataUnavailable)
		json.NewEncoder(w).Encode(DataUnavailable)
	}
	json.NewEncoder(w).Encode(vals)
}
