package main

import (
	"fmt"
	"io/ioutil"
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

// File representing a file
// - absPath → absolute path of the file
// - name → name of the file
// - mimeType → mimeType of the file
type File struct {
	AbsPath  string
	Name     string
	MimeType string
}

func checkCredentials(username, password string) bool {
	return (username == os.Getenv("username") && password == os.Getenv("password"))
}

// Load directories and music files in an absolute directory
func loadDirectoryContent(dirPath string) ([]File, []File, error) {
	// Read dir content
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}

	// Separate directories and musics
	dirs := make([]File, 0)
	// Add the parent and rood directory to the directory list
	if dirPath != musicRootDir {
		dirs = append(dirs, File{musicRootDir, "/", "text/directory"})
		dirs = append(
			dirs,
			File{
				strings.TrimSuffix(dirPath, filepath.Base(dirPath)+"/"),
				"/..",
				"text/directory",
			},
		)
	}
	musics := make([]File, 0)
	for _, file := range files {
		if file.IsDir() {
			dir := File{filepath.Clean(dirPath + file.Name()), file.Name(), "text/directory"}
			dirs = append(dirs, dir)
		} else {
			absFlePath := filepath.Clean(dirPath + file.Name())
			// Detect the mimetype of the file
			mimeType, err := mimetype.DetectFile(absFlePath)
			generalMimeType := strings.Split(mimeType.String(), "/")[0]
			if err != nil {
				fmt.Printf("MimeType detection failed on file `%s`: %s\n", absFlePath, err)
			} else if generalMimeType == "audio" { // Check if it's an audio file
				music := File{absFlePath, file.Name(), mimeType.String()}
				musics = append(musics, music)
			}
		}
	}

	return dirs, musics, nil
}

func makeSongsLink(files []string, linkTpl string) []SongLink {
	// Remove the `musicRootDir` part of files path and remove all non audio files
	audioFiles := make([]SongLink, 0)
	for _, path := range files {
		mt, err := mimetype.DetectFile(path)
		mime := strings.Split(mt.String(), "/")[0] // General type of file
		// There is no error and it's an audio file
		if err == nil && mime == "audio" {
			// Build the song url link
			audioTitle := strings.TrimPrefix(path, musicRootDir)
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
		songPath = fmt.Sprintf("%s%s", musicRootDir, song)

		// Check if we can acces to the file
		if _, err := os.Stat(songPath); err != nil { // We can't acces the file
			title = "Error finding file"
			fmt.Printf("Error accessing the file %s: %s\n", songPath, err)
		} else { // We can acces the file
			// Get the mime type of the song
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
