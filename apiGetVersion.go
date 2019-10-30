package main

import (
	"fmt"
	"net/http"
	"runtime"
)

func apiGetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Server", fmt.Sprintf("%s (%s)", runtime.Version(), runtime.GOOS))
	w.Write([]byte(VERSION))
}
