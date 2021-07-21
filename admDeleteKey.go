package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func admDeleteKey(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()
	key := mux.Vars(r)["key"]
	if _, ok := cache[key]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	delete(cache, key)
}
