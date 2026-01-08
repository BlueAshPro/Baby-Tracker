# Baby Track

**Projet web en Go — Visualisation et filtrage d’artistes musicaux**

---

## Présentation du projet

**Baby Track** est une application web développée en **Go (Golang)** qui consomme l’API publique **Groupie Tracker** afin d’afficher des informations sur des artistes et groupes de musique :

- nom  
- membres  
- dates de création  
- concerts et localisations  

Le projet repose sur une **architecture serveur simple** 

- du routage HTTP  
- de la consommation d’API  
- de la séparation des responsabilités  
- du rendu HTML côté serveur  

---

## Architecture Générale

```
┌───────────────────────────────────────────────────────────────┐
│                     Frontend (HTML / CSS)                     │
│              accueil.html • artiste.html                      │
└───────────────────────────┬───────────────────────────────────┘
                            │ HTTP Requêtes
                            │
┌───────────────────────────▼───────────────────────────────────┐
│                     main.go (Router HTTP)                     │
│        - Déclaration des routes                               │
│        - Lancement du serveur                                 │
└───────────────────────────┬───────────────────────────────────┘
                            │
            ┌───────────────┼────────────────┐
            │               │                │
┌───────────▼───────────┐ ┌─▼────────────┐ ┌─▼────────────────┐
│ gestion / handlers    │ │ gestion /    │ │ internal / api   │
│                       │ │ recherche &  │ │                  │
│ - Pages HTML          │ │ filter       │ │ - Appels HTTP    │
│ - Endpoints API       │ │              │ │ - API externe    │
└───────────┬───────────┘ └─┬────────────┘ └─┬────────────────┘
            │               │                │
            └───────────────┴────────────────┘
                            │
┌───────────────────────────▼───────────────────────────────────┐
│              API – Groupie Tracker                            │
│        - Récupération des artistes                            │
│        - Relations concerts / lieux                           │
└───────────────────────────────────────────────────────────────┘
```

---

## Principe de Fonctionnement

### 1. Page d’accueil `/`

- Sert une page HTML statique  
- Chargement initial de tous les artistes via l’API interne  
- Affichage sous forme de cartes  

### 2. API interne `/api/artists`

- Point d’entrée JSON  
- Sert de passerelle entre le frontend et l’API externe  
- Centralise les appels réseau  

### 3. Recherche `/api/search`

- Recherche textuelle côté serveur  
- Insensible à la casse  
- Optimisée pour limiter les appels réseau  

### 4. Filtres `/api/filter`

Filtres combinables :
- date de création  
- nombre de membres  
- localisation  

### 5. Page artiste `/artiste?id=X`

- Page dynamique générée côté serveur  
- Données spécifiques à un artiste  
- Intégration des concerts et localisations  

---

## Structure du Projet

```
groupie-tracker/
│
├── main.go # Point d’entrée du serveur
│
├── internal/
│ ├── api/ # Communication API externe
│ ├── gestion/ # Handlers HTTP
│ └── models/ # Structures de données
│
├── static/
│ ├── accueil.html
│ └── artiste.html
│
├── css/
│ └── style.css
│
└── go.mod / go.sum
```

---

## API Externe Utilisée

**Groupie Tracker API (publique)**  
Elle fournit toutes les données nécessaires au projet.

Documentation officielle :  
👉 https://groupietrackers.herokuapp.com/api

Endpoints principaux :
- `/artists`  
- `/artists/{id}`  
- `/relation/{id}`  
- `/locations`  

---

## Documentation et Références Officielles

### Go / Backend

- Documentation officielle Go  
  https://go.dev/doc/

- Package `net/http`  
  https://pkg.go.dev/net/http

- Templates HTML en Go  
  https://pkg.go.dev/html/template

- Encodage JSON  
  https://pkg.go.dev/encoding/json

### HTTP & Web

- HTTP Status Codes  
  https://developer.mozilla.org/fr/docs/Web/HTTP/Status

- Méthodes HTTP  
  https://developer.mozilla.org/fr/docs/Web/HTTP/Methods

### Frontend

- Fetch API  
  https://developer.mozilla.org/fr/docs/Web/API/Fetch_API

- Manipulation du DOM  
  https://developer.mozilla.org/fr/docs/Web/API/Document_Object_Model
