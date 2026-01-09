package gestion

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"groupie-tracker/internal/api"
)

type ConcertData struct {
	Location string `json:"location"`
	Date     string `json:"date"`
}

type ArtistPageData struct {
	ID           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Concerts     []ConcertData
}

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

	// Récupérer les concerts
	relations, err := api.FetchRelationByID(id)
	var concerts []ConcertData
	if err == nil {
		for location, dates := range relations.DatesLocations {
			for _, date := range dates {
				concerts = append(concerts, ConcertData{
					Location: location,
					Date:     date,
				})
			}
		}
	}

	data := ArtistPageData{
		ID:           artist.ID,
		Image:        artist.Image,
		Name:         artist.Name,
		Members:      artist.Members,
		CreationDate: artist.CreationDate,
		FirstAlbum:   artist.FirstAlbum,
		Concerts:     concerts,
	}

	tmpl, err := template.ParseFiles("static/artiste.html")
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

	log.Printf("Page artiste affichée: %s (ID: %d)", artist.Name, id)
}
