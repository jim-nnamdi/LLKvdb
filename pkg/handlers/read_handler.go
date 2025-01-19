package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		vars            = mux.Vars(r)
		DataUnavailable = "Associated Data for key could not be found on disk!"
		DataAvailable   = "Data successfully read"
	)
	keyval, _ := strconv.ParseInt(vars["key"], 10, 64)
	vals, ok := handler.Fsys.Read(int64(keyval))
	if !ok {
		fmt.Println("err", DataUnavailable)
		json.NewEncoder(w).Encode(DataUnavailable)
	}
	response := map[string]interface{}{}
	response["message"] = DataAvailable
	response["value"] = vals
	json.NewEncoder(w).Encode(response)
}
