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
