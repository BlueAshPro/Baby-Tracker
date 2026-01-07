package gestion

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

func FilterArtists(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	creationMin := params.Get("creationDateMin")
	creationMax := params.Get("creationDateMax")
	albumYear := params.Get("firstAlbumYear")
	membersMin := params.Get("membersMin")
	membersMax := params.Get("membersMax")
	location := strings.ToLower(strings.TrimSpace(params.Get("location")))

	// Récupération de tous les artistes depuis l'API externe
	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	var relations *models.RelationData
	if location != "" {
		relations, err = api.FetchRelations()
		if err != nil {
			log.Printf("Erreur fetch relations: %v", err)
		}
	}

	var minYear, maxYear, minMembers, maxMembers int
	var hasMinYear, hasMaxYear, hasMinMembers, hasMaxMembers bool

	if creationMin != "" {
		if val, err := strconv.Atoi(creationMin); err == nil {
			minYear = val
			hasMinYear = true
		}
	}
	if creationMax != "" {
		if val, err := strconv.Atoi(creationMax); err == nil {
			maxYear = val
			hasMaxYear = true
		}
	}
	if membersMin != "" {
		if val, err := strconv.Atoi(membersMin); err == nil {
			minMembers = val
			hasMinMembers = true
		}
	}
	if membersMax != "" {
		if val, err := strconv.Atoi(membersMax); err == nil {
			maxMembers = val
			hasMaxMembers = true
		}
	}

	results := make([]models.Artist, 0, len(artists))

	for _, artist := range artists {
		if hasMinYear && artist.CreationDate < minYear {
			continue
		}
		if hasMaxYear && artist.CreationDate > maxYear {
			continue
		}

		if albumYear != "" && !strings.Contains(artist.FirstAlbum, albumYear) {
			continue
		}

		memberCount := len(artist.Members)
		if hasMinMembers && memberCount < minMembers {
			continue
		}
		if hasMaxMembers && memberCount > maxMembers {
			continue
		}

		// Filtre : Ville de concert
		if location != "" && relations != nil {
			locationMatch := false
			for _, rel := range relations.Index {
				if rel.ID == artist.ID {
					for loc := range rel.DatesLocations {
						if strings.Contains(strings.ToLower(loc), location) {
							locationMatch = true
							break
						}
					}
					break
				}
			}
			if !locationMatch {
				continue
			}
		}

		results = append(results, artist)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Erreur JSON: %v", err)
		http.Error(w, "Erreur encodage", http.StatusInternalServerError)
		return
	}

	log.Printf("Filtrage appliqué - %d résultats sur %d artistes", len(results), len(artists))
}
