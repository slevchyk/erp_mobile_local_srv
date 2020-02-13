package dbase

import (
	"database/sql"
	"github.com/slevchyk/erp_mobile_local_srv/models"
)

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
