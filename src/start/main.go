package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	// Указание папки template как источника статических файлов
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fs := http.FileServer(http.Dir("../template"))
	http.Handle("/", http.StripPrefix("/", fs))
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
