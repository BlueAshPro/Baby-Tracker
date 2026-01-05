package gestion

import (
	"encoding/json"
	"net/http"

	"groupie-tracker-gui/internal/api"
)

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "Recherche non retrouvée", http.StatusBadRequest)
		return
	}

	results, err := api.SearchArtist(query)
	if err != nil {
		http.Error(w, "Échec lors de la récupération des résultats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
