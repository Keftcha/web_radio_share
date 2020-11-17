package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/keftcha/wrs/streamingservice"
)

// Our streaming service
var strmSvc *streamingservice.StreamingService

// The root directory to find musics
var musicRootDir string

func init() {
	strmSvc = streamingservice.New()
	musicRootDir = "/music/"
}

func main() {
	fmt.Println("\033[35mThe server is now running !\033[0m")

	// Page to manage the music
	http.HandleFunc("/hoster", stream)
	http.HandleFunc("/hoster/", stream)

	// Page to sign in to manage the music
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signin/", signin)

	// Pahe to listen all music of the host
	http.HandleFunc("/listen", listen)
	http.HandleFunc("/listen/", listen)

	// Page to listen what the hoster listen
	http.HandleFunc("/listen_stream", listenStream)
	http.HandleFunc("/listen_stream/", listenStream)

	// Broadcast the song that the host listen
	http.HandleFunc("/broadcasting", broadcasting)
	http.HandleFunc("/broadcasting/", broadcasting)

	// Serve sound as static file when path start with `/music/`
	http.Handle(
		"/music/",
		http.StripPrefix(
			"/music/",
			http.FileServer(http.Dir(musicRootDir)),
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
