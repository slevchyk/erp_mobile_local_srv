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
			news TEXT,
			date TEXT);
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
			mob_id INTEGER,
			acc_id TEXT,
			user_id TEXT,
			date TIMESTAMP,
			status TEXT,
			is_turnstile BOOLEAN,
			started_at TIMESTAMP ,
			ended_at TIMESTAMP,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			deleted_at TIMESTAMP);
			`)
	if err != nil {
		log.Fatal(err)
	}

}
