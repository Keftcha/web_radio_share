package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type SongLink struct {
	Link  string
	Title string
}

type Markers struct {
	Title         string
	SongPath      string
	SongType      string
	FilesLinks    []SongLink
	Username      string
	Password      string
	Authenticated bool
}

func checkCredentials(username, password string) bool {
	return (username == os.Getenv("username") && password == os.Getenv("password"))
}

func loadDirectoryTree(root string) []string {
	files := make([]string, 0, 1)
	filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			files = append(files, path)
			return nil
		},
	)
	return files
}

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

func signin(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("page/sign_in.html")
	tpl.Execute(w, struct{}{})
}

func makeSongsLink(files []string, linkTpl string) []SongLink {
	// Remove the `/music/` part of files path and remove all non audio files
	audioFiles := make([]SongLink, 0)
	for _, path := range files {
		mt, err := mimetype.DetectFile(path)
		mime := strings.Split(mt.String(), "/")[0] // General type of file
		// There is no error and it's an audio file
		if err == nil && mime == "audio" {
			// Build the song url link
			audioTitle := path[len("/music/"):]
			audioLink := fmt.Sprintf(
				linkTpl,
				audioTitle,
			)

			songLink := SongLink{Title: audioTitle, Link: audioLink}
			audioFiles = append(audioFiles, songLink)
		}
	}
	return audioFiles
}

func findSong(song string) (title string, songType string, songPath string) {
	if song == "" { // The song parameter isn't present
		title = "No song playing"
		fmt.Println("Not listening")
	} else { // The song parameter is present
		songPath = fmt.Sprintf("/music/%s", song)

		// Check if we can acces to the file
		if _, err := os.Stat(songPath); err != nil { // We can't acces the file
			title = "Error finding file"
			fmt.Println(err)
		} else { // We can acces the file
			// get the mime type of the song
			var mimeType string
			if mime, err := mimetype.DetectFile(songPath); err == nil {
				mimeType = mime.String()
			} else {
				fmt.Println(err)
				mimeType = "application/octet-stream"
			}

			title = song
			songType = mimeType
			fmt.Println("Listening: ", title)
		}
	}
	return title, songType, songPath
}

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
