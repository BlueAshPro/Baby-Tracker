package gestion

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"groupie-tracker/internal/api"
)

func ArtistePage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	artist, err := api.FetchArtistByID(id)
	if err != nil {
		log.Printf("Erreur fetch artiste %d: %v", id, err)
		http.Error(w, "Artiste non trouvé", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("static/artiste.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, artist); err != nil {
		log.Printf("Erreur execution: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
		return
	}

	log.Printf("Page artiste affichée: %s (ID: %d)", artist.Name, id)
}
