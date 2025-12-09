package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/gestion/", gestion.home)
	http.HandleFunc("/artists", gestion.artists)

	http.HandleFunc("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
