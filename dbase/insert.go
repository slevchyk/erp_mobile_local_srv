package dbase

import (
	"database/sql"

	"github.com/slevchyk/erp_mobile_local_srv/models"
)

func InsertChannel(db *sql.DB, c models.Channel) (sql.Result, error) {

	stmt, _ := db.Prepare(`
		INSERT INTO
			channels (
				id,
				user_id,
				type,
				update_id,
				title,
				new,				 
				date_time
				)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`)

	return stmt.Exec(c.ID, c.UserID, c.Type, c.UpdateID, c.Title, c.News, c.Date)
}

func InsertToken(db *sql.DB, ft models.FirebaseTokens) (sql.Result, error) {

	stmt, _ := db.Prepare(`
		INSERT INTO
			firebase_tokens (				
				user_id,
				token
				)
		VALUES ($1, $2);`)

	return stmt.Exec(ft.UserID, ft.Token)
}
