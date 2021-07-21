package main

import (
	"net/http"
)

func admCreateKey(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()
	key := r.FormValue("key")
	if _, ok := cache[key]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}
	cache[key] = nil
	w.Header().Set("Location", "/keys/"+key)
	w.WriteHeader(http.StatusCreated)
}
