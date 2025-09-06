package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
