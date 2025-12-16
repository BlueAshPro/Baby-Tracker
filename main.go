package main

import (
	"groupie-tracker/internal/gestion"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/gestion/", gestion.Home)
	http.HandleFunc("/artists", gestion.Artists)
	http.HandleFunc("/search", gestion.Search) // Barre de recherche

	http.HandleFunc("/", gestion.Home)

	cssFS := http.FileServer(http.Dir("./css/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFS))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/static/accueil/", http.StripPrefix("/static/accueil/", http.FileServer(http.Dir("static/accueil"))))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
