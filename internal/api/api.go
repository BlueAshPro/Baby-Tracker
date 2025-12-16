package api

import (
	"encoding/json"
	"net/http"
)

var (
	ArtistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	LocationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	DatesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationURL  = "https://groupietrackers.herokuapp.com/api/relation"
)

type Artist struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Location struct {
	Id   int    `json:"id"`
	City string `json:"city"`
}

type Date struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// SearchArtist peut interroger l’API et retourner tous les artistes
func GetArtists() ([]Artist, error) {
	resp, err := http.Get(ArtistsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}

// Pareil pour Locations, Dates, Relation
func GetLocations() ([]Location, error) {
	resp, err := http.Get(LocationsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locations []Location
	err = json.NewDecoder(resp.Body).Decode(&locations)
	return locations, err
}

func GetDates() ([]Date, error) {
	resp, err := http.Get(DatesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dates []Date
	err = json.NewDecoder(resp.Body).Decode(&dates)
	return dates, err
}

// Relation est spéciale → map dynamique
func GetRelations() ([]Date, error) {
	resp, err := http.Get(RelationURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relations []Date
	err = json.NewDecoder(resp.Body).Decode(&relations)
	return relations, err
}
