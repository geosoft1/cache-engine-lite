package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func apiUpdateKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if _, ok := cache[key]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var data data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m.Lock()
	cache[key] = data
	m.Unlock()
}

func apiUpdateKeyQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	query.Del("key")
	// NOTE for more accuracy use RFC3339Nano
	query.Set("_time_", time.Now().UTC().Format(time.RFC3339))
	if _, ok := cache[key]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// map[string][]string -> map[string]interface{}
	var data = map[string]data{}
	for k, v := range query {
		if len(v) > 1 {
			data[k] = v
		} else {
			data[k] = v[0]
		}
	}
	m.Lock()
	cache[key] = data
	m.Unlock()
}
