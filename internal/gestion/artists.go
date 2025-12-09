package gestion

import (
	"encoding/json"
	"net/http"

	"groupie-tracker-gui/internal/api"
)

// Artists récupère la liste des artistes via l'API interne et renvoie du JSON
func Artists(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(artists)
}
