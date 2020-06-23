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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id SERIAL PRIMARY KEY,
			blocked BOOLEAN DEFAULT false,
			user_id TEXT DEFAULT '',
			pin TEXT DEFAULT '',
			info_card INTEGER DEFAULT 0,
			last_name TEXT DEFAULT '',
			first_name TEXT DEFAULT '',
			middle_name TEXT DEFAULT '',
			itn TEXT DEFAULT '',
			phone TEXT DEFAULT '',
			birthday TIMESTAMP,
			email TEXT DEFAULT '',
			gender TEXT DEFAULT '',
			address TEXT DEFAULT '',
			passport_type TEXT DEFAULT '',
			passport_series TEXT DEFAULT '',
			passport_number TEXT DEFAULT '',
			passport_issued TEXT DEFAULT '',
			passport_date TIMESTAMP,
			passport_expiry TIMESTAMP,
			civil_status TEXT DEFAULT '',			
			children TEXT DEFAULT '',
			job_position TEXT DEFAULT '',
			education INTEGER DEFAULT 0,
			specialty TEXT DEFAULT '',
			additional_education TEXT DEFAULT '',
			last_work_place TEXT DEFAULT '',
			skills TEXT DEFAULT '',
			languages TEXT DEFAULT '',
			disability BOOLEAN DEFAULT false,
			pensioner BOOLEAN DEFAULT false,
			photo_name TEXT DEFAULT '',
			photo_data TEXT DEFAULT '');
			`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS help_desk (
			id SERIAL PRIMARY KEY,
			user_id TEXT,
			date TIMESTAMP,
			title TEXT,
			body TEXT,
			status TEXT,
			answer TEXT,
			answered_by TEXT,
			answered_at TIMESTAMP,
			is_modified_by_mob BOOLEAN,
			is_modified_by_acc BOOLEAN);
			`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pay_desk (
			id SERIAL PRIMARY KEY,
			user_id TEXT,
			amount FLOAT,
			payment TEXT,
			document_number TEXT,
			document_date TIMESTAMP,
			file_paths TEXT DEFAULT '', 
			files_quantity INTEGER DEFAULT 0,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,			
			is_deleted BOOLEAN DEFAULT false,
			is_modified_by_mob BOOLEAN,
			is_modified_by_acc BOOLEAN);
			`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cost_items (
			id SERIAL PRIMARY KEY,
			acc_id TEXT,
			name TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,			
			is_deleted BOOLEAN DEFAULT false);
			`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS income_items (
			id SERIAL PRIMARY KEY,
			acc_id TEXT,
			name TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,			
			is_deleted BOOLEAN DEFAULT false);
			`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pay_offices (
			id SERIAL PRIMARY KEY,
			acc_id TEXT,
			name TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,			
			is_deleted BOOLEAN DEFAULT false);
			`)
	if err != nil {
		log.Fatal(err)
	}

}
