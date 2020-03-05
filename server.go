package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("The server is now running !")

	http.HandleFunc("/hoster", stream)
	http.HandleFunc("/hoster/", stream)

	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signin/", signin)

	http.HandleFunc("/listen", listen)
	http.HandleFunc("/listen/", listen)

	// Serve sound as static file when path start with `/music/`
	http.Handle(
		"/music/",
		http.StripPrefix(
			"/music/",
			http.FileServer(http.Dir("/music/")),
		),
	)

	// Serve images as staticfile when path start with `/image/`
	http.Handle(
		"/image/",
		http.StripPrefix(
			"/image/",
			http.FileServer(http.Dir("img/")),
		),
	)

	// Serve statics files
	http.Handle(
		"/",
		http.StripPrefix(
			"/",
			http.FileServer(http.Dir("public/")),
		),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
