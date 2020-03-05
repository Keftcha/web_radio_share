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

	// Make the list of song
	linkTpl := "/hoster/?song=%s"
	markers.FilesLinks = makeSongsLink(files, linkTpl)
	// Add the authentigication parameters to song links
	authTpl := "&username=%s&password=%s"
	for idx, songLink := range markers.FilesLinks {
		markers.FilesLinks[idx].Link = songLink.Link + fmt.Sprintf(authTpl, os.Getenv("username"), os.Getenv("password"))
	}

	// Load and execute the template
	tpl, _ := template.ParseFiles("page/player.html")
	tpl.Execute(w, markers)
}

func signin(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("page/sign_in.html")
	tpl.Execute(w, struct{}{})
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
}
