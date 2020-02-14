package dbase

import "database/sql"

func SelectFirebaseTokenByUserId(db *sql.DB, userID string) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			ft.id,
			ft.user_id,
			ft.token		
		FROM 
			firebase_tokens ft
		WHERE
			user_id=$1`, userID)
}

func SelectFirebaseTokenByUserIdToken(db *sql.DB, userID, token string) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			ft.id		
		FROM 
			firebase_tokens ft
		WHERE
			user_id=$1
			AND token=$2`, userID, token)
}
