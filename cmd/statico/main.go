package main

import (
	"flag"
	"log"
	"net/http"
)

func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s - %s", r.Method, r.RemoteAddr, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	listenAddr := flag.String("addr", ":80", "HTTP listen address")
	staticDir := flag.String("dir", "/static/", "Static directory to serve")
	flag.Parse()

	fileServer := http.FileServer(http.Dir(*staticDir))
	log.Fatal(http.ListenAndServe(*listenAddr, logHandler(fileServer)))
}
