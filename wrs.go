package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Markers struct {
	Song     string
	SongPath string
	SongType string
	Infos    string
}

func listen(w http.ResponseWriter, r *http.Request) {
	song := r.URL.Query().Get("song")

	// Our markers
	markers := new(Markers)

	// Fill markers values
	if song == "" {
		markers.Infos = fmt.Sprint("No song is playing.")
		markers.Song = "No song playing"
		fmt.Println("Not listening")
	} else {
		markers.Infos = fmt.Sprintf("Currently playing: %s !", song)
		markers.Song = song
		markers.SongPath = fmt.Sprintf("/music/%s", song)
		fmt.Println("Listening: ", song)
	}

	// Load and execute the template
	tpl, _ := template.ParseFiles("player.html")
	tpl.Execute(w, markers)
}

func main() {
	fmt.Println("The server is now running !")

	http.HandleFunc("/hoster", listen)
	http.HandleFunc("/hoster/", listen)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
