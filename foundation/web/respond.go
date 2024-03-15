package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, status int) error {
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return nil
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}
