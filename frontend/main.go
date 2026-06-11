package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) {
	files := []string{
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "conference-ouest.html"),
		filepath.Join("templates", "conference-est.html"),
		filepath.Join("templates", "team.html"),
		filepath.Join("templates", "inscription.html"),
		filepath.Join("templates", "connexion.html"),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erreur de chargement du template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	data := map[string]interface{}{
		"Title": "Bienvenue sur mon Frontend",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderTemplate(w, "index", data)
	})

	http.HandleFunc("/conference-ouest", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "conference-ouest", data)
	})

	http.HandleFunc("/conference-est", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "conference-est", data)
	})

	http.HandleFunc("/conference-est/celtics", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "team", data)
	})

	http.HandleFunc("/inscription", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "inscription", data)
	})

	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "connexion", data)
	})

	port := ":8081"
	fmt.Printf("Serveur prêt sur http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Erreur lancement serveur : %s\n", err.Error())
	}
}
