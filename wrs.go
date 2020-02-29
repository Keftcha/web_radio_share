package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Markers struct {
	Song  string
	Infos string
}

func listen(w http.ResponseWriter, r *http.Request) {
	song := r.URL.Query().Get("song")

	// Check the song parameter
	var infos string
	if song != "" {
		infos = fmt.Sprintf("Currently playing: %s !", song)
		fmt.Println("Listening: ", song)
	} else {
		infos = fmt.Sprint("No song is playing.")
		song = "No song playing"
		fmt.Println("Not listening")
	}

	// Set our markers
	markers := &Markers{
		Song:  song,
		Infos: infos,
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
