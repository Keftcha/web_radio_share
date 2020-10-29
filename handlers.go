package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func stream(w http.ResponseWriter, r *http.Request) {
	song := r.URL.Query().Get("song")

	// Check if the user is authentificated
	if !checkCredentials(
		r.URL.Query().Get("username"),
		r.URL.Query().Get("password"),
	) {
		// Redicect the user if he isn't authentificated
		http.Redirect(w, r, "/signin", 401)
		return
	}

	// Our markers
	markers := new(Markers)
	markers.Authenticated = true

	// Fill markers values
	markers.Title, markers.SongType, markers.SongPath = findSong(song)

	// Load directoies and files
	var files []string = loadDirectoryTree("/music")

	// Add the authentigication parameters to song links
	authTpl := fmt.Sprintf("&username=%s&password=%s", os.Getenv("username"), os.Getenv("password"))
	// Make the list of song
	linkTpl := "/hoster/?song=%s" + authTpl
	markers.FilesLinks = makeSongsLink(files, linkTpl)

	// Load and execute the template
	tpl, _ := template.ParseFiles("page/player.html")
	tpl.Execute(w, markers)

	fmt.Println(fmt.Sprintf("You now streaming: \033[33m%s\033[0m", markers.Title))
}

func signin(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("page/sign_in.html")
	tpl.Execute(w, struct{}{})

	fmt.Println(
		fmt.Sprintf(
			"\033[31m/!\\\033[0m The host \033[34m%s\033[0m is on the sign in page !",
			r.Host,
		),
	)
}

func listen(w http.ResponseWriter, r *http.Request) {
	song := r.URL.Query().Get("song")

	// Our markers
	markers := new(Markers)
	markers.Authenticated = false

	// Fill markers values
	markers.Title, markers.SongType, markers.SongPath = findSong(song)

	// Load directoies and files
	var files []string = loadDirectoryTree("/music")

	// Make the list of song
	linkTpl := "/listen/?song=%s"
	markers.FilesLinks = makeSongsLink(files, linkTpl)

	// Load and execute the template
	tpl, _ := template.ParseFiles("page/player.html")
	tpl.Execute(w, markers)

	fmt.Println(
		fmt.Sprintf(
			"\033[34m%s\033[0m listening \033[33m%s\033[0m",
			r.Host,
			markers.Title,
		),
	)
}

func listenStream(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("page/listen_stream.html")
	tpl.Execute(w, nil)
}
