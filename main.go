package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
)

var VERSION = "1.1.0-20191116"

var (
	cacheFile    = flag.String("cache-file", "cache.json", "cache file name")
	httpAddress  = flag.String("http", ":8080", "http address")
	httpsAddress = flag.String("https", ":8090", "https address")
)

type data = interface{}

var cache = map[string]data{}
var m = sync.Mutex{}
var api = mux.NewRouter()
var adm = api.PathPrefix("/admin").Subrouter()

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	*cacheFile = filepath.Join(path, *cacheFile)
	file, err := os.Open(*cacheFile)
	if err != nil {
		log.Println(err)
	}
	if err := json.NewDecoder(file).Decode(&cache); err != nil {
		log.Println(err)
	}
	// logging middleware
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr, r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
	// authorization middleware
	adm.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Auth-Token") != os.Getenv("XAUTHTOKEN") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	// core API
	adm.HandleFunc("/keys", admCreateKey).Methods("POST")
	adm.HandleFunc("/keys/{key}", admDeleteKey).Methods("DELETE")
	adm.HandleFunc("/keys", admGetKeys).Methods("GET")
	adm.HandleFunc("/keys", admSaveKeys).Methods("PUT")
	api.HandleFunc("/keys", apiGetKeys).Methods("GET")
	api.HandleFunc("/keys/{key}", apiGetKey).Methods("GET")
	api.HandleFunc("/keys/{key}", apiUpdateKey).Methods("PUT")
	api.HandleFunc("/update", apiUpdateKeyQuery).Methods("GET")
	api.HandleFunc("/version", apiGetVersion).Methods("GET")
	// core servers
	log.Printf("starting http services")
	go func() {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		if err := http.ListenAndServeTLS(*httpsAddress, filepath.Join(path, "server.crt"), filepath.Join(path, "server.key"), api); err != nil {
			log.Println(err)
		}
	}()
	if err := http.ListenAndServe(*httpAddress, api); err != nil {
		log.Fatalln(err)
	}
}
