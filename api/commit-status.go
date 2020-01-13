package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (data *DataMiddleWare) SetCommitStatus(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	commitId := vars["commitId"]

	decoder := json.NewDecoder(req.Body)

	var commitStatus CommitStatus
	err := decoder.Decode(&commitStatus)
	if err != nil {
		panic(err)
	}
	data.CommitStatus[commitId] = commitStatus
	w.WriteHeader(http.StatusNoContent)
}

func (data *DataMiddleWare) GetCommitStatus(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	commitId := vars["commitId"]

	encoder := json.NewEncoder(w)
	if val, ok := data.CommitStatus[commitId]; ok {
		encoder.Encode(val)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
