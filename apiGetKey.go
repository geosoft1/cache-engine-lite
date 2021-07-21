package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func apiGetKey(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()
	key := mux.Vars(r)["key"]
	if _, ok := cache[key]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if cache[key] == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(cache[key]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if r.Header.Get("Cache-Control") == "no-store" {
		cache[key] = nil
	}
}
