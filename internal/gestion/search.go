package gestion

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

func SearchArtists(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "Paramètre de recherche manquant", http.StatusBadRequest)
		return
	}

	query = strings.ToLower(query)

	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	results := make([]models.Artist, 0)
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, artist)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Erreur JSON: %v", err)
		http.Error(w, "Erreur encodage", http.StatusInternalServerError)
		return
	}

	log.Printf("Recherche '%s' - %d résultats", query, len(results))
}
