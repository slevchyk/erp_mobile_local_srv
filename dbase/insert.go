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
				news,				 
				date
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

func InsertTiming(db *sql.DB, t models.Timing) (sql.Result, error) {

	stmt, _ := db.Prepare(`
		INSERT INTO
			timing (				
				mob_id,
				acc_id,
				user_id,
				date,
				status,
				is_turnstile,
				started_at,
				ended_at,
				created_at,
				updated_at,
				deleted_at
				)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)

	return stmt.Exec(t.MobID, t.AccID, t.UserID, t.Date, t.Status, t.IsTurnstile, t.StartedAt, t.EndedAt, t.CreatedAt, t.UpdatedAt, t.DeletedAt)
}
