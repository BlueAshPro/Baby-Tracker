package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", gestion.home)
	http.HandleFunc("/artists", gestion.artists)

	http.HandleFunc("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
