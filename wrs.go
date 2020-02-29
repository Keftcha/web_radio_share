package main

import (
	"fmt"
	"log"
	"net/http"
)

func listen(w http.ResponseWriter, r *http.Request) {
	song := r.URL.Query().Get("song")
	if song != "" {
		fmt.Fprintf(w, "Currently playing: %s !", song)
		fmt.Println("Listening: ", song)
	} else {
		fmt.Fprint(w, "No song is playing.")
		fmt.Println("Not listening")
	}
}

func main() {
	fmt.Println("The server is now running !")

	http.HandleFunc("/hoster", listen)
	http.HandleFunc("/hoster/", listen)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
