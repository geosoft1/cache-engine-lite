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

var VERSION = "1.6.0-20220721"

var (
	cacheFile    = flag.String("cache-file", "cache.json", "cache file name")
	httpAddress  = flag.String("http", ":8080", "http address")
	httpsAddress = flag.String("https", ":8090", "https address")
)

type data = interface{}

var cache = map[string]data{}
var m = sync.Mutex{}
var api = mux.NewRouter().PathPrefix(os.Getenv("INSTANCEID")).Subrouter()
var adm = api.PathPrefix("/admin").Subrouter()

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Auth-Token, Cache-Control, Cors-Control")
	w.Header().Set("Access-Control-Expose-Headers", "*")
}

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
	// file server
	api.PathPrefix("/static/").Handler(http.StripPrefix(os.Getenv("INSTANCEID")+"/static/", http.FileServer(http.Dir(filepath.Join(path, "static")))))
	// logging middleware
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr, r.Method, r.RequestURI)
			if r.Header.Get("Cors-Control") != "no-cors" {
				// open api need CORS
				enableCORS(w)
			}
			if r.Method == "OPTIONS" {
				return
			}
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
	adm.HandleFunc("/keys", admCreateKey).Methods("POST", "OPTIONS")
	adm.HandleFunc("/keys/{key}", admDeleteKey).Methods("DELETE", "OPTIONS")
	adm.HandleFunc("/keys", admGetKeys).Methods("GET", "OPTIONS")
	adm.HandleFunc("/keys", admSaveKeys).Methods("PUT", "OPTIONS")
	api.HandleFunc("/keys", apiGetKeys).Methods("GET")
	api.HandleFunc("/keys/{key}", apiGetKey).Methods("GET", "OPTIONS")
	api.HandleFunc("/keys/{key}", apiUpdateKey).Methods("PUT", "POST")
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
