package dbase

import (
	"database/sql"

	"github.com/slevchyk/erp_mobile_local_srv/models"
)

func UpdateChannel(db *sql.DB, c models.Channel) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				channels
			SET
				user_id = $1,
				update_id = $2,
				type = $3,
				title = $4,
				news = $5,
				date = $6				
			WHERE
				id=$7
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(c.UserID, c.UpdateID, c.Type, c.Title, c.News, c.Date, c.ID)

	return res, err
}

func UpdateFirebaseTokens(db *sql.DB, ft models.FirebaseTokens) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				firebase_tokens
			SET
				token=$1
			WHERE
				user_id=$2
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(ft.Token, ft.UserID)

	return res, err
}
func UpdateTiming(db *sql.DB, t models.Timing) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				timing
			SET
				mob_id = $1,
				acc_id = $2,
				user_id = $3,
				date = $4,
				status = $5,
				is_turnstile = $6,
				started_at = $7,
				ended_at = $8,
				created_at = $9,
				updated_at = $10,
				deleted_at = $11			
			WHERE
				id=$12
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(t.MobID, t.AccID, t.UserID, t.Date, t.Status, t.IsTurnstile, t.StartedAt, t.EndedAt, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.ID)

	return res, err
}
