package xhttp

import (
	"encoding/json"
	"net/http"

	"github.com/todennus/x/mime"
)

func WriteResponseJSON(w http.ResponseWriter, code int, obj any) error {
	jsonString, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", mime.ApplicationJSON)
	w.WriteHeader(code)
	_, err = w.Write(jsonString)
	return err
}
