package api

import (
	"groupie-tracker-gui/internal/models"
	"net/http"
)

var artists []models.Artist

func GetArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist

	return artists, nil
}
