package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/forum_nba?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("erreur connexion bdd: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("la bdd ne repond pas: %v", err)
	}
	return db
}
