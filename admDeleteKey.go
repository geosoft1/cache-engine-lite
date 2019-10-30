package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func admDeleteKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if _, ok := cache[key]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	m.Lock()
	delete(cache, key)
	m.Unlock()
}
