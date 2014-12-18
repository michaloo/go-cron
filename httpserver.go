package gocron

import (
	"encoding/json"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.MarshalIndent(Current_state, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if Last_err.Exit_status != 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Write(js)
}

func Http_server() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":18080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
