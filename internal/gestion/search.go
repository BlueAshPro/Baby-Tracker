package gestion

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
)

func SearchArtists(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	// Si pas de recherche, rediriger vers l'accueil
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	query = strings.ToLower(query)

	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Récupérer les relations pour chercher dans les emplacements
	relations, _ := api.FetchRelations()

	// Filtrer les artistes
	var results []interface{}
	for _, artist := range artists {
		matched := false

		// Recherche dans le nom de l'artiste
		if strings.Contains(strings.ToLower(artist.Name), query) {
			matched = true
		}

		// Recherche dans les membres
		if !matched {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), query) {
					matched = true
					break
				}
			}
		}

		// Recherche dans les emplacements de concerts
		if !matched && relations != nil {
			for _, rel := range relations.Index {
				if rel.ID == artist.ID {
					for location := range rel.DatesLocations {
						if strings.Contains(strings.ToLower(location), query) {
							matched = true
							break
						}
					}
					break
				}
			}
		}

		// Recherche dans la date de création
		if !matched {
			creationStr := strconv.Itoa(artist.CreationDate)
			if strings.Contains(creationStr, query) {
				matched = true
			}
		}

		// Recherche dans la date du premier album
		if !matched {
			if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
				matched = true
			}
		}

		if matched {
			results = append(results, artist)
		}
	}

	data := HomeData{
		Artists: results,
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

	log.Printf("Recherche '%s' - %d résultats", query, len(results))
}
