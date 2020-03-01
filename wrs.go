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

type Markers struct {
	Song     string
	SongPath string
	SongType string
	Infos    string
	Files    []string
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

func listen(w http.ResponseWriter, r *http.Request) {
	song := r.URL.Query().Get("song")

	// Our markers
	markers := new(Markers)

	// Fill markers values
	if song == "" { // The song parameter isn't present
		markers.Infos = "No song is playing."
		markers.Song = "No song playing"
		fmt.Println("Not listening")
	} else { // The song parameter is present
		markers.SongPath = fmt.Sprintf("/music/%s", song)

		// Check if we can acces to the file
		if _, err := os.Stat(markers.SongPath); err != nil { // We can't acces the file
			markers.Infos = "An error occured finding the file."
			markers.Song = "Error finding file"
			fmt.Println(err)
		} else { // We can acces the file
			// get the mime type of the song
			var mimeType string
			if mime, err := mimetype.DetectFile(markers.SongPath); err == nil {
				mimeType = mime.String()
			} else {
				fmt.Println(err)
				mimeType = "application/octet-stream"
			}

			markers.Infos = fmt.Sprintf("Currently playing: %s !", song)
			markers.Song = song
			markers.SongType = mimeType
			fmt.Println("Listening: ", markers.Song)
		}
	}

	// Load directoies and files
	var files []string = loadDirectoryTree("/music")

	// Remove the `/music/` part of files path and remove all non audio files
	audioFiles := make([]string, 0)
	for _, path := range files {
		mt, err := mimetype.DetectFile(path)
		mime := strings.Split(mt.String(), "/")[0] // General type of file
		// There is no error and it's an audio file
		if err == nil && mime == "audio" {
			audioFiles = append(audioFiles, path[len("/music/"):])
		}
	}

	markers.Files = audioFiles

	// Load and execute the template
	tpl, _ := template.ParseFiles("player.html")
	tpl.Execute(w, markers)
}

func main() {
	fmt.Println("The server is now running !")

	http.HandleFunc("/hoster", listen)
	http.HandleFunc("/hoster/", listen)

	// Serve sound as static file when path start with `/music/`
	http.Handle(
		"/music/",
		http.StripPrefix(
			"/music/",
			http.FileServer(http.Dir("/music")),
		),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
