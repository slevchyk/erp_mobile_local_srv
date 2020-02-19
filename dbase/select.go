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

func SelectTimingById(db *sql.DB, id int) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			t.id,
			t.mob_id,
			t.acc_id,
			t.user_id,
			t.date,
			t.status,
			t.is_turnstile,
			t.started_at,
			t.ended_at,
			t.created_at,
			t.updated_at,
			t.deleted_at					
		FROM 
			timing t
		WHERE
		t.id=$1`, id)
}

func SelectTimingByMobIdUserIdDate(db *sql.DB, id int, userID, date string) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			t.id,
			t.mob_id,
			t.acc_id,
			t.user_id,
			t.date,
			t.status,
			t.is_turnstile,
			t.started_at,
			t.ended_at,
			t.created_at,
			t.updated_at,
			t.deleted_at					
		FROM 
			timing t
		WHERE
		t.mob_id=$1
		and t.user_id=$2
		and t.date=$3`, id, userID, date)
}

func SelectTimingByExtIdUerIdDate(db *sql.DB, id int, userID, date string) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			t.id,
			t.mob_id,
			t.acc_id,
			t.user_id,
			t.date,
			t.status,
			t.is_turnstile,
			t.started_at,
			t.ended_at,
			t.created_at,
			t.updated_at,
			t.deleted_at					
		FROM 
			timing t
		WHERE
		t.ext_id=$1
		and t.user_id=$2
		and t.date=$3`, id, userID, date)
}
