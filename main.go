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

	http.HandleFunc("/", gestion.Home)
	http.HandleFunc("/artiste", gestion.ArtistePage)
	http.HandleFunc("/api/search", gestion.SearchArtists)
	http.HandleFunc("/api/filter", gestion.FilterArtists)
	http.HandleFunc("/api/artist/", gestion.ArtistDetailsAPI)
	http.HandleFunc("/api/artists", gestion.GetAllArtists)

	printServerInfo()

	port := ":8080"
	log.Fatal(http.ListenAndServe(port, nil))
}

func printBanner() {
	fmt.Println("================================")
	fmt.Println("   GROUPIE TRACKER SERVER")
	fmt.Println("================================")
}

func printServerInfo() {
	fmt.Println("\nServeur démarré sur http://localhost:8080")
	fmt.Println("\nEndpoints:")
	fmt.Println("  GET  /")
	fmt.Println("  GET  /artiste?id=X")
	fmt.Println("  GET  /api/search?q=XXX")
	fmt.Println("  GET  /api/filter")
	fmt.Println("================================\n")
}
