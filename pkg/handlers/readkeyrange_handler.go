package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jim-nnamdi/Lkvs/pkg/model"
)

var _ http.Handler = &readKeyRangeHandler{}

type readKeyRangeHandler struct {
	Fsys *model.Filesys
}

func NewReadKeyRangeHandler(Fsys *model.Filesys) *readKeyRangeHandler {
	return &readKeyRangeHandler{Fsys: Fsys}
}

func (handler *readKeyRangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		startkey        = r.FormValue("startkey")
		endkey          = r.FormValue("endkey")
		DataUnavailable = "Associated Data for key could not be found on disk!"
		DataAvailable   = "Data successfully read"
	)
	startkeyval, _ := strconv.ParseInt(startkey, 10, 64)
	endkeyval, _ := strconv.ParseInt(endkey, 10, 64)
	vals, err := handler.Fsys.ReadKeyRange(startkeyval, endkeyval)
	if err != nil {
		fmt.Println("err", DataUnavailable)
		json.NewEncoder(w).Encode(DataUnavailable)
	}
	response := map[string]interface{}{}
	response["message"] = DataAvailable
	response["value"] = vals
	json.NewEncoder(w).Encode(response)
}
