package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var baseDir = func() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(filename)
}()

func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) {
	files := []string{
		filepath.Join(baseDir, "templates", "index.html"),
		filepath.Join(baseDir, "templates", "conference-ouest.html"),
		filepath.Join(baseDir, "templates", "conference-est.html"),
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

func serveTeamPage(w http.ResponseWriter, r *http.Request, folder string) {
	slug := strings.TrimPrefix(r.URL.Path, "/"+folder+"/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(baseDir, "templates", folder, slug+".html")
	if _, err := os.Stat(filePath); err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func main() {
	fs := http.FileServer(http.Dir(filepath.Join(baseDir, "static")))
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

	http.HandleFunc("/conference-ouest/", func(w http.ResponseWriter, r *http.Request) {
		serveTeamPage(w, r, "conference-ouest")
	})

	http.HandleFunc("/conference-est/", func(w http.ResponseWriter, r *http.Request) {
		serveTeamPage(w, r, "conference-est")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	fmt.Printf("Serveur prêt sur http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Erreur lancement serveur : %s\n", err.Error())
	}
}
