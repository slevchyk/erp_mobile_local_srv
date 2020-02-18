package dbase

import (
	"database/sql"
)

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

func SelectChannelById(db *sql.DB, id int) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			c.id,
			c.user_id,
			c.update_id,
			c.type,
			c.title,
			c.news,
			c.date		
		FROM 
			channels c
		WHERE
			c.id=$1`, id)
}

func SelectChannelsByUserIdUpdateId(db *sql.DB, userID string, updateID int) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			c.id,
			c.user_id,
			c.update_id,
			c.type,
			c.title,
			c.news,
			c.date		
		FROM 
			channels c
		WHERE
		c.user_id=$1
			and c.update_id>=$2`, userID, updateID)
}
