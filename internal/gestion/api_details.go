package gestion

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
)

type ConcertData struct {
	Location string `json:"location"`
	Date     string `json:"date"`
}

func ArtistDetailsAPI(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		http.Error(w, "Format invalide", http.StatusBadRequest)
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	if len(parts) >= 5 && parts[4] == "concerts" {
		GetArtistConcerts(w, id)
	} else {
		http.Error(w, "Endpoint invalide", http.StatusNotFound)
	}
}

func GetArtistConcerts(w http.ResponseWriter, artistID int) {
	// Récupération des relations (concerts) de cet artiste
	relations, err := api.FetchRelationByID(artistID)
	if err != nil {
		log.Printf("Erreur fetch concerts: %v", err)
		http.Error(w, "Concerts non trouvés", http.StatusNotFound)
		return
	}

	// Conversion des données en format simple pour le frontend
	var concerts []ConcertData
	for location, dates := range relations.DatesLocations {
		for _, date := range dates {
			concerts = append(concerts, ConcertData{
				Location: location,
				Date:     date,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(concerts)

	log.Printf("Concerts artiste %d - %d concerts", artistID, len(concerts))
}
