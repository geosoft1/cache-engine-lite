package main

import (
	"net/http"
)

func admCreateKey(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if _, ok := cache[key]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}
	m.Lock()
	cache[key] = nil
	m.Unlock()
	w.Header().Set("Location", "/keys/"+key)
	w.WriteHeader(http.StatusCreated)
}
