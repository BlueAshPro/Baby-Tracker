package gestion

import (
	"log"
	"net/http"
	"os"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data, err := os.ReadFile("static/accueil.html")
	if err != nil {
		log.Printf("Erreur lecture fichier: %v", err)
		http.Error(w, "Erreur de chargement", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}
