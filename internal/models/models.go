package models

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Members      []string `json:"members"`
	Image        string   `json:"image"`
}
