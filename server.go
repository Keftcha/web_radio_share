package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("\033[35mThe server is now running !\033[0m")

	// Manage the music
	http.HandleFunc("/hoster", stream)
	http.HandleFunc("/hoster/", stream)

	// Sign in to manage the music
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signin/", signin)

	// Listen all music of the host
	http.HandleFunc("/listen", listen)
	http.HandleFunc("/listen/", listen)

	// Listen what the hoster listen
	http.HandleFunc("/listen_stream", listenStream)
	http.HandleFunc("/listen_stream/", listenStream)

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
		"/public/",
		http.StripPrefix(
			"/public/",
			http.FileServer(http.Dir("public/")),
		),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
