package dbase

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/slevchyk/erp_mobile_local_srv/models"
)

func ConnectDB(cfg models.DBConfig) (*sql.DB, error) {

	dbConnection := fmt.Sprintf("postgres://%v:%v@localhost/%v?sslmode=disable", cfg.User, cfg.Password, cfg.Name)
	db, err := sql.Open("postgres", dbConnection)

	return db, err
}

func InitDB(db *sql.DB) {

	var err error

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			id INTEGER,
			user_id TEXT,
			update_id INTEGER ,
			type TEXT,
			title TEXT,			
			new TEXT,
			date_time TEXT);
			`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS firebase_tokens (
			id SERIAL PRIMARY KEY,
			user_id TEXT,
			token TEXT);
			`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS timing (
			id SERIAL PRIMARY KEY,
			ext_id TEXT,
			user_id TEXT,
			date TEXT,
			operation TEXT,
			started_at TEXT,
			ended_at TEXT,
			created_at TEXT,
			updated_at TEXT,
			deleted_at TEXT,
			is_turnstile BOOLEAN);
			`)
	if err != nil {
		log.Fatal(err)
	}

}
