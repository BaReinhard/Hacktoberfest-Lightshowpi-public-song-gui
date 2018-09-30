package main

import (
	"log"
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
	appengine.Main()
}
