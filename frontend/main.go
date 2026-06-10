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

	fmt.Printf("Serveur prêt sur http://localhost%s\n", ":8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erreur lancement serveur : %s\n", err.Error())
	}
}
