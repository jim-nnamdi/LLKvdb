package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jim-nnamdi/Lkvs/pkg/model"
)

var _ http.Handler = &batchPutHandler{}

type batchPutHandler struct {
	Fsys *model.Filesys
}

func NewBatchPutHandler(Fsys *model.Filesys) *batchPutHandler {
	return &batchPutHandler{Fsys: Fsys}
}

func (handler *batchPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		keys          = r.FormValue("keys")
		vals          = r.FormValue("values")
		DataPersisted = "Data successfully added and persisted to disk!"
	)
	if keys == "" || vals == "" {
		http.Error(w, "keys and values are required", http.StatusBadRequest)
		return
	}
	keyList := strings.Split(keys, ",")
	valList := strings.Split(vals, ",")
	if len(keyList) != len(valList) {
		http.Error(w, "keys and values must be of the same length", http.StatusBadRequest)
		return
	}
	data := make([]map[int64]string, 0, len(keyList))
	for i, key := range keyList {
		keyval, err := strconv.ParseInt(strings.TrimSpace(key), 10, 64)
		if err != nil {
			http.Error(w, "error:invalid key format, must be an integer", http.StatusBadRequest)
			return
		}
		data = append(data, map[int64]string{
			keyval: strings.TrimSpace(valList[i]),
		})
	}
	handler.Fsys.BatchPut(data)
	response := map[string]interface{}{}
	response["message"] = DataPersisted
	json.NewEncoder(w).Encode(response)
}
