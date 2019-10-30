package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func admSaveKeys(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create(*cacheFile)
	if err != nil {
		w.WriteHeader(http.StatusInsufficientStorage)
		log.Println(err)
		return
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	if err := enc.Encode(&cache); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}
