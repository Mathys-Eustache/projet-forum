package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL pour communiquer avec la base
	"github.com/joho/godotenv"
)

// InitDB initialise la connexion à la base de données MySQL
func InitDB() *sql.DB {
	// On charge les variables d'environnement (.env) silencieusement
	_ = godotenv.Load()

	// Configuration de la connexion (Data Source Name)
	dsn := "root:@tcp(127.0.0.1:3306)/forum_nba?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la BDD : %v", err)
	}

	// On vérifie que la base de données répond vraiment
	if err = db.Ping(); err != nil {
		log.Fatalf("La base de données ne répond pas : %v", err)
	}

	return db
}
