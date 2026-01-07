package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"groupie-tracker/internal/models"
)

const (
	BaseURL = "https://groupietrackers.herokuapp.com/api"
	Timeout = 10 * time.Second
)

var httpClient = &http.Client{
	Timeout: Timeout,
}

func FetchArtists() ([]models.Artist, error) {
	resp, err := httpClient.Get(BaseURL + "/artists")
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture: %w", err)
	}

	var artists []models.Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, fmt.Errorf("erreur JSON: %w", err)
	}

	return artists, nil
}

func FetchArtistByID(id int) (*models.Artist, error) {
	url := fmt.Sprintf("%s/artists/%d", BaseURL, id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur connexion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("artiste non trouvé")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API: status %d", resp.StatusCode)
	}

	var artist models.Artist
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		return nil, fmt.Errorf("erreur parsing: %w", err)
	}

	return &artist, nil
}

func FetchLocations() (*models.DateLocation, error) {
	resp, err := httpClient.Get(BaseURL + "/locations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locations models.DateLocation
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, err
	}

	return &locations, nil
}

func FetchRelations() (*models.RelationData, error) {
	resp, err := httpClient.Get(BaseURL + "/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relations models.RelationData
	if err := json.NewDecoder(resp.Body).Decode(&relations); err != nil {
		return nil, err
	}

	return &relations, nil
}

func FetchRelationByID(id int) (*models.Relation, error) {
	url := fmt.Sprintf("%s/relation/%d", BaseURL, id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API: status %d", resp.StatusCode)
	}

	var relation models.Relation
	if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
		return nil, err
	}

	return &relation, nil
}
