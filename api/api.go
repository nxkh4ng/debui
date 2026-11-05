package api

import (
	"encoding/json"
	"net/http"
)

func Heathz(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	buf, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "JSON encode failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
