package gestion

import (
	"encoding/json"
	"log"
	"net/http"

	"groupie-tracker/internal/api"
)

func GetAllArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch artistes: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}
