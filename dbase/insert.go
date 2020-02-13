package dbase

import (
	"database/sql"
	"github.com/slevchyk/erp_mobile_local_srv/models"
)

func InsertChannel(db *sql.DB, c models.Channel) (sql.Result, error)  {

	stmt, _ := db.Prepare(`
		INSERT INTO
			users (
				ext_id,
				user_id,
				update_id,
				title,
				new,				 
				date_time
				)
		VALUES ($1, $2, $3, $4, $5, $6);`)

	return stmt.Exec(c.ExtID, c.UserID, c.UpdateID, c.Title, c.New, c.DateTime)
}

func InsertToken(db *sql.DB, ft models.FirebaseTokens) (sql.Result, error)  {

	stmt, _ := db.Prepare(`
		INSERT INTO
			firebase_tokens (				
				user_id,
				token
				)
		VALUES ($1, $2);`)

	return stmt.Exec(ft.UserID, ft.Token)
}