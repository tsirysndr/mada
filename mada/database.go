package mada

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

func OpenDatabaseConnection() (*sql.DB, error) {
	if os.Getenv("MADA_POSTGRES_URL") != "" {
		return sql.Open("postgres", os.Getenv("MADA_POSTGRES_URL"))
	}

	return OpenSQLiteConnection()
}

func OpenSQLiteConnection() (*sql.DB, error) {
	sql.Register("sqlite3_with_spatialite",
		&sqlite3.SQLiteDriver{
			Extensions: []string{"mod_spatialite"},
		})

	return sql.Open("sqlite3_with_spatialite", filepath.Join(CreateConfigDir(), "spatialmada.db"))
}

func OpenPostgresConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("MADA_POSTGRES_URL"))

	if err != nil {
		log.Fatal(err)
	}

	db.Exec("CREATE EXTENSION postgis;")

	return db, nil
}
