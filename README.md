# Groupie-Tracker - Documentation Technique

## Architecture Générale

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (HTML/CSS/JS)                   │
│              (accueil.html, artiste.html)                    │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP Requests
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   main.go (Server Router)                     │
│  - http.HandleFunc("/", gestion.Home)                        │
│  - http.HandleFunc("/artiste", gestion.ArtistePage)          │
│  - http.HandleFunc("/api/artists", gestion.GetAllArtists)    │
│  - http.HandleFunc("/api/search", gestion.SearchArtists)     │
│  - http.HandleFunc("/api/filter", gestion.FilterArtists)     │
│  - http.HandleFunc("/api/artist/", gestion.ArtistDetailsAPI) │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐  ┌─────▼─────┐  ┌────▼─────┐
   │ gestion  │  │ gestion   │  │ internal  │
   │ handlers │  │ searches/ │  │ api       │
   │          │  │ filters   │  │           │
   └─────┬────┘  └─────┬─────┘  └────┬─────┘
         │             │             │
    ┌────▼─────────────▼─────────────▼────┐
    │   internal/api/api.go                 │
    │  (API Wrapper - Groupie Tracker)      │
    │  - FetchArtists()                     │
    │  - FetchArtistByID(id)                │
    │  - FetchRelationByID(id)              │
    │  - FetchRelations()                   │
    │  - FetchLocations()                   │
    └────┬──────────────────────────────────┘
         │
         │ HTTP GET
         │
    ┌────▼───────────────────────────────────────┐
    │  External API                               │
    │  https://groupietrackers.herokuapp.com/api  │
    │  /artists, /artists/{id}, /relation/{id}    │
    └─────────────────────────────────────────────┘
```

## Flux d'une Requête HTTP

### 1. Requête Homepage: GET /

```
Browser                  main.go            gestion.Home         os.ReadFile()
  │                        │                    │                    │
  ├─ GET / ───────────────>│                    │                    │
  │                        ├── route to Home ──>│                    │
  │                        │                    ├─ vérifier URL ────>│
  │                        │                    │    (doit être "/") │
  │                        │                    │                    │
  │                        │                    ├─ lire static/accueil.html
  │                        │                    │                    │
  │<─ HTML content ────────┼────────────────────┤                    │
  │                        │                    │                    │
```

**Code: `internal/gestion/home.go`**
```go
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
```

### 2. Requête API: GET /api/artists

```
Frontend JS           main.go          gestion.GetAllArtists    api.FetchArtists()
     │                  │                      │                       │
     ├─ fetch() ───────>│                      │                       │
     │                  ├─ route to /api/artists ──>                   │
     │                  │                      ├─ appel api ──────────>│
     │                  │                      │                       │
     │                  │                      │  ┌─────────────────┐  │
     │                  │                      │  │ Connexion HTTP  │  │
     │                  │                      │  │ vers API externe│  │
     │                  │                      │  └─────────────────┘  │
     │                  │                      │                       │
     │<─ JSON array ────┼──────────────────────┤                       │
     │                  │                      │                       │
```

**Code: `internal/gestion/get_artists.go`**
```go
func GetAllArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch artistes: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}
```

### 3. Requête Recherche: GET /api/search?q=metallica

```
Frontend JS           main.go            gestion.SearchArtists    api.FetchArtists()
     │                  │                      │                       │
     ├─ fetch() ───────>│                      │                       │
     │ ?q=metallica     ├─ route to /api/search ──>                   │
     │                  │                      ├─ fetch tous artistes ─>
     │                  │                      │                       │
     │                  │                      │  ┌──────────────────┐ │
     │                  │                      │  │ Filtre par nom   │ │
     │                  │                      │  │ strings.Contains │ │
     │                  │                      │  └──────────────────┘ │
     │                  │                      │                       │
     │<─ JSON results ──┼──────────────────────┤                       │
```

**Code: `internal/gestion/search.go`**
```go
func SearchArtists(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "Paramètre de recherche manquant", http.StatusBadRequest)
		return
	}
	
	query = strings.ToLower(query)
	
	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	
	results := make([]models.Artist, 0)
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, artist)
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
```

### 4. Requête Filtre: GET /api/filter?creationDateMin=2000&membersMin=4

```
Frontend JS           main.go          gestion.FilterArtists    api.FetchArtists()
     │                  │                      │                       │
     ├─ fetch() ───────>│                      │                       │
     │ ?creationDate... ├─ route to /api/filter ──>                   │
     │ &membersMin...   │                      ├─ extract params      │
     │                  │                      ├─ fetch artistes ─────>
     │                  │                      │                       │
     │                  │                      │  ┌──────────────────┐ │
     │                  │                      │  │ Appliquer filtres│ │
     │                  │                      │  │ - création date  │ │
     │                  │                      │  │ - membres count  │ │
     │                  │                      │  │ - location       │ │
     │                  │                      │  └──────────────────┘ │
     │                  │                      │                       │
     │<─ JSON filtered ──┼──────────────────────┤                       │
```

**Code: `internal/gestion/filter.go`**
```go
func FilterArtists(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	creationMin := params.Get("creationDateMin")
	creationMax := params.Get("creationDateMax")
	membersMin := params.Get("membersMin")
	membersMax := params.Get("membersMax")
	location := strings.ToLower(strings.TrimSpace(params.Get("location")))
	
	artists, err := api.FetchArtists()
	if err != nil {
		log.Printf("Erreur fetch: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	
	// Conversion paramètres (optimisation - une seule fois)
	var minYear, maxYear, minMembers, maxMembers int
	var hasMinYear, hasMaxYear, hasMinMembers, hasMaxMembers bool
	
	if creationMin != "" {
		if val, err := strconv.Atoi(creationMin); err == nil {
			minYear = val
			hasMinYear = true
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
		
		memberCount := len(artist.Members)
		if hasMinMembers && memberCount < minMembers {
			continue
		}
		if hasMaxMembers && memberCount > maxMembers {
			continue
		}
		
		results = append(results, artist)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
```

### 5. Requête Page Artiste: GET /artiste?id=1

```
Browser               main.go           gestion.ArtistePage    api.FetchArtistByID
   │                    │                     │                        │
   ├─ GET /artiste ────>│                     │                        │
   │  ?id=1             ├─ route ────────────>│                        │
   │                    │                     ├─ extract ID           │
   │                    │                     ├─ fetch artist ───────>
   │                    │                     │                        │
   │                    │                     │  ┌─────────────────┐   │
   │                    │                     │  │ HTTP GET         │   │
   │                    │                     │  │ /artists/1       │   │
   │                    │                     │  └─────────────────┘   │
   │                    │                     │                        │
   │                    │                     ├─ parse template       │
   │                    │                     ├─ execute (data)       │
   │                    │                     │                        │
   │<─ HTML artiste ───┼─────────────────────┤                        │
```

**Code: `internal/gestion/artists.go`**
```go
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
	
	tmpl, err := template.ParseFiles("static/artiste.html")
	if err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
	
	if err := tmpl.Execute(w, artist); err != nil {
		log.Printf("Erreur execution: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}
```

## Structure des Fichiers

### `main.go`
- **Responsabilité**: Configuration du serveur et des routes
- **Fonctions**:
  - `main()`: Lance le serveur sur :8080
  - `printBanner()`: Affiche le message de démarrage
  - `printServerInfo()`: Liste les endpoints disponibles

```go
func main() {
	printBanner()
	
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	
	http.HandleFunc("/", gestion.Home)
	http.HandleFunc("/artiste", gestion.ArtistePage)
	http.HandleFunc("/api/search", gestion.SearchArtists)
	http.HandleFunc("/api/filter", gestion.FilterArtists)
	http.HandleFunc("/api/artist/", gestion.ArtistDetailsAPI)
	http.HandleFunc("/api/artists", gestion.GetAllArtists)
	
	printServerInfo()
	
	port := ":8080"
	log.Fatal(http.ListenAndServe(port, nil))
}
```

### `internal/api/api.go`
- **Responsabilité**: Communication avec l'API externe Groupie Tracker
- **Configuration**:
  - `BaseURL = "https://groupietrackers.herokuapp.com/api"`
  - `Timeout = 10 * time.Second`
  - `httpClient` partagé pour toutes les requêtes

- **Fonctions**:
  - `FetchArtists()` → GET /artists → []models.Artist
  - `FetchArtistByID(id)` → GET /artists/{id} → *models.Artist
  - `FetchRelationByID(id)` → GET /relation/{id} → *models.Relation
  - `FetchRelations()` → GET /relation → *models.RelationData
  - `FetchLocations()` → GET /locations → *models.DateLocation

### `internal/models/models.go`
- **Responsabilité**: Structures de données

```go
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
```

### `internal/gestion/*.go` - Handlers HTTP

| Fichier | Route | Fonction |
|---------|-------|----------|
| `home.go` | `GET /` | Servir static/accueil.html |
| `artists.go` | `GET /artiste?id=X` | Servir template artiste.html |
| `get_artists.go` | `GET /api/artists` | Retourner JSON de tous les artistes |
| `search.go` | `GET /api/search?q=XXX` | Rechercher artistes par nom |
| `filter.go` | `GET /api/filter?...` | Filtrer par critères multiples |
| `api_details.go` | `GET /api/artist/{id}/concerts` | Retourner concerts d'un artiste |

## Points d'Optimisation

### 1. Cache Frontend
```javascript
let cachedArtists = null;

async function loadArtists() {
	const response = await fetch('/api/artists');
	const artists = await response.json();
	cachedArtists = artists;
	displayArtists(artists);
}
```

### 2. Conversion Paramètres Une Seule Fois
```go
// Conversion avant la boucle, pas à l'intérieur
minYear, _ := strconv.Atoi(params.Get("creationDateMin"))
for _, artist := range artists {
	if artist.CreationDate < minYear { continue }
}
```

### 3. Débounce Recherche (300ms)
```javascript
let searchTimeout;
searchInput.addEventListener('input', function() {
	clearTimeout(searchTimeout);
	searchTimeout = setTimeout(() => {
		fetch(`/api/search?q=${this.value}`);
	}, 300);
});
```

## Endpoints Disponibles

| Méthode | Route | Paramètres | Retour |
|---------|-------|-----------|--------|
| GET | `/` | - | HTML |
| GET | `/artiste` | `id` (int) | HTML |
| GET | `/api/artists` | - | JSON []Artist |
| GET | `/api/search` | `q` (string) | JSON []Artist |
| GET | `/api/filter` | `creationDateMin`, `creationDateMax`, `membersMin`, `membersMax`, `location` | JSON []Artist |
| GET | `/api/artist/{id}/concerts` | `id` (dans URL) | JSON []ConcertData |

## Installation et Lancement

```bash
go mod download
go build -o groupie-tracker.exe
./groupie-tracker.exe
```

Puis ouvrir http://localhost:8080
