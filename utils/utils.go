package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, message map[string]any) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(message)

	w.WriteHeader(code)
	_, _ = w.Write(j)

}
