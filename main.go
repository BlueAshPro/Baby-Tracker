package main

import (
	"fmt"
	"log"
	"net/http"

	"groupie-tracker/internal/gestion"
)

func main() {
	printBanner()

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", gestion.Home)
	http.HandleFunc("/artiste", gestion.ArtistePage)
	http.HandleFunc("/api/search", gestion.SearchArtists)
	http.HandleFunc("/api/filter", gestion.FilterArtists)

	port := ":8080"
	log.Fatal(http.ListenAndServe(port, nil))
}

func printBanner() {
	fmt.Println("Serveur lancé sur: http://localhost:8080")
}
