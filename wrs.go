package main

import (
	"fmt"
	"log"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	fmt.Println("Hello, World!")
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
