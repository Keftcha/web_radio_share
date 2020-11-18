package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func stream(w http.ResponseWriter, r *http.Request) {
	// Song absolute path
	song := r.URL.Query().Get("song")
	strmSvc.Start(song)

	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	// Check if the user is authentificated
	if !checkCredentials(
		username,
		password,
	) {
		// Redicect the user if he isn't authentificated
		http.Redirect(w, r, "/signin", 401)
		return
	}

	// Get the current directory of the user is in
	crntDir := r.URL.Query().Get("crntDir")
	if crntDir == "" {
		crntDir = musicRootDir
	}

	// Absolut path of the current directory
	absCrntDir := filepath.Clean(crntDir) + "/"

	// Load directoies and musics
	dirs, musics, err := loadDirectoryContent(absCrntDir)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/hoster/?crntDir=%s&username=%s&password=%s", musicRootDir, username, password), 500)
		fmt.Printf("Failed to load the %s directory content: %s\n", crntDir, err)
		return
	}

	// Our markers
	markers := make(map[string]interface{})
	// Directories and musics
	markers["Dirs"] = dirs
	markers["Musics"] = musics
	// Current directory where the user is
	markers["CrntDir"] = crntDir
	// User authentification
	markers["Authenticated"] = true
	// User credentials
	markers["Username"] = username
	markers["Passwd"] = password
	// Currently playing song markers
	if song != "" {
		if mime, err := mimetype.DetectFile(song); err == nil {
			markers["Title"] = filepath.Base(song)
			markers["Path"] = song
			markers["MimeType"] = mime
		} else {
			markers["Title"] = "Failed to find the mime type"
			fmt.Printf("Failed to finde `%s` file mimetype: %s\n", song, err)
		}
	} else {
		markers["Title"] = "No song playing"
	}

	// Load and execute the template
	tpl, _ := template.ParseFiles("page/player.html")
	if err := tpl.Execute(w, markers); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/hoster/?crntDir=%s&username=%s&password=%s", musicRootDir, username, password), 500)
		fmt.Printf("Failed to execute `page/player.html` template: %s\n", err)
		return
	}

	fmt.Println(fmt.Sprintf("You now streaming: \033[33m%s\033[0m", markers["Title"]))
}

func signin(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("page/sign_in.html")
	tpl.Execute(w, musicRootDir)

	fmt.Println(
		fmt.Sprintf(
			"\033[31m/!\\\033[0m The host \033[34m%s\033[0m is on the sign in page !",
			r.RemoteAddr,
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
	var files []string = nil //loadDirectoryTree(musicRootDir)

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

func broadcasting(w http.ResponseWriter, r *http.Request) {
}
