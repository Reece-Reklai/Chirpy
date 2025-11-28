package main

import (
	"log"
	"net/http"
)

// type fileHandler struct {
// }
//
// func (fileHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
//
// }

func main() {
	mux := http.NewServeMux()
	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(server.ListenAndServe())
}
