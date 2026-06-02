package api

import (
	"encoding/json"
	"net/http"
)

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func writeError(w http.ResponseWriter, errorText string, code int) {
	var (
		response []byte
		err      error
	)
	if response, err = json.Marshal(&struct {
		Error string `json:"error"`
	}{
		Error: errorText,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setHeaders(w)
	w.WriteHeader(code)
	w.Write([]byte(response))
}

func writeJSON(w http.ResponseWriter, data any) {
	var (
		response []byte
		err      error
	)

	if response, err = json.Marshal(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
