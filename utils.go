package main

import (
	"fmt"
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
			audioLink := fmt.Sprintf(linkTpl, audioTitle)

			songLink := SongLink{Title: audioTitle, Link: audioLink}
			audioFiles = append(audioFiles, songLink)
		}
	}
	return audioFiles
}

func findSong(song string) (title string, songType string, songPath string) {
	if song == "" { // The song parameter isn't present
		title = "No song playing"
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
		}
	}
	return title, songType, songPath
}
