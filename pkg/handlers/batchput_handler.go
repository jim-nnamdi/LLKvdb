package handlers

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("keyList", keyList)
	valList := strings.Split(vals, ",")
	fmt.Println("valList", valList)
	if len(keyList) != len(valList) {
		http.Error(w, "keys and values must be of the same length", http.StatusBadRequest)
		return
	}

	ikeys := make([]int64, len(keyList))
	for i, v := range keyList {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			http.Error(w, "invalid conversion from string to int64", http.StatusBadRequest)
			return
		}
		ikeys[i] = val
	}
	handler.Fsys.BatchPut(ikeys, valList)
	response := map[string]interface{}{}
	response["message"] = DataPersisted
	json.NewEncoder(w).Encode(response)
}
