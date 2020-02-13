package dbase

import (
	"database/sql"
	"fmt"
	"github.com/slevchyk/erp_mobile_local_srv/models"
	_ "github.com/lib/pq"
	"log"
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
			id SERIAL PRIMARY KEY,
			ext_id INTEGER,
			user_id TEXT,
			update_id INTEGER,
			title TEXT,			
			new TEXT,
			date_time TIMESTAMP);
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
}