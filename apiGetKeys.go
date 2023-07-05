package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func apiGetKeys(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()
	var grp []any
	for _, key := range r.URL.Query()["key"] {
		grp = append(grp, cache[key])
	}
	if len(grp) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Size", strconv.Itoa(len(grp)))
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(grp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}
