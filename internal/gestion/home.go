package gestion

import (
	"html/template"
	"log"
	"net/http"

	"groupie-tracker/internal/api"
)

type HomeData struct {
	Artists []interface{}
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Récupérer les artistes depuis l'API
	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch artistes: %v", err)
		http.Error(w, "Erreur de chargement des artistes", http.StatusInternalServerError)
		return
	}

	// Convertir en interface{} pour le template
	artistData := make([]interface{}, len(artists))
	for i, a := range artists {
		artistData[i] = a
	}

	data := HomeData{
		Artists: artistData,
	}

	tmpl, err := template.ParseFiles("static/accueil.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Erreur execution: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
		return
	}

	log.Printf("Page accueil affichée: %d artistes", len(artists))
}
