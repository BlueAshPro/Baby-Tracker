package gestion

import (
	"net/http"
)

// Home sert la page d'accueil statique
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/acceuil.html")
}
