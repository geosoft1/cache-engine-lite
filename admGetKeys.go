package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func admGetKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Size", strconv.Itoa(len(cache)))
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(cache); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}
