package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		files := []string{
			filepath.Join("templates", "index.html"),
			filepath.Join("templates", "accueil.html"),
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			http.Error(w, "Erreur de chargement du template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Title": "Bienvenue sur mon Frontend",
		}

		err = tmpl.ExecuteTemplate(w, "index", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Printf("Serveur prêt sur http://localhost%s\n", ":8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erreur lancement serveur : %s\n", err.Error())
	}
}
