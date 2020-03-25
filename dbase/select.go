package dbase

import (
	"database/sql"
	"time"
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

func SelectTimingById(db *sql.DB, id int64) (*sql.Rows, error) {
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

func SelectTimingByMobIdUserIdDate(db *sql.DB, id int64, userID string, date time.Time) (*sql.Rows, error) {
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
			AND t.user_id=$2
			AND t.date=$3`, id, userID, date)
}

func SelectTimingByAccIdUerIdDate(db *sql.DB, accID, userID string, date time.Time) (*sql.Rows, error) {
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
			t.acc_id=$1
			and t.user_id=$2
			and t.date=$3`, accID, userID, date)
}

func SelectTimingByUserIdDate(db *sql.DB, userID string, date time.Time) (*sql.Rows, error) {
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
			t.user_id=$1
			and t.date=$2`, userID, date)
}

func SelectTimingByUpdatedAt(db *sql.DB, date time.Time) (*sql.Rows, error) {
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
			t.updated_at>$1
			OR t.updated_at IS NULL`, date)
}

func SelectProfileByPhonePin(db *sql.DB, phone, pin string) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			p.id,
			p.blocked,
			p.user_id,
			p.pin,
			p.info_card,
			p.last_name,
			p.first_name,
			p.middle_name,
			p.itn,
			p.phone,
			p.birthday,
			p.email,
			p.gender,
			p.address,
			p.passport_type,
			p.passport_series,
			p.passport_number,
			p.passport_issued,
			p.passport_date,
			p.passport_expiry,
			p.civil_status,		
			p.children,
			p.job_position,
			p.education,
			p.specialty,
			p.additional_education,
			p.last_work_place,
			p.skills,
			p.languages,
			p.disability,
			p.pensioner,
			p.photo_name,
			p.photo_data
		FROM
			profiles p
		WHERE
			p.phone=$1
			AND p.pin=$2
			AND p.user_id<>$3`, phone, pin, "")

}

func SelectProfileByUserID(db *sql.DB, userID string) (*sql.Rows, error) {

	return db.Query(`
		SELECT
			p.id,
			p.blocked,
			p.user_id,
			p.pin,
			p.info_card,
			p.last_name,
			p.first_name,
			p.middle_name,
			p.itn,
			p.phone,
			p.birthday,
			p.email,
			p.gender,
			p.address,
			p.passport_type,
			p.passport_series,
			p.passport_number,
			p.passport_issued,
			p.passport_date,
			p.passport_expiry,
			p.civil_status,		
			p.children,
			p.job_position,
			p.education,
			p.specialty,
			p.additional_education,
			p.last_work_place,
			p.skills,
			p.languages,
			p.disability,
			p.pensioner,
			p.photo_name,
			p.photo_data
		FROM
			profiles p
		WHERE
			p.user_id=$1`, userID)
}

func SelectHelpDeskByID(db *sql.DB, id int) (*sql.Rows, error) {
	return db.Query(`
		SELECT
			hd.id,
		    hd.user_id,
		    hd.date,
		    hd.title,
		    hd.body,
		    hd.status,
		    hd.answer,
		    hd.answered_by,
		    hd.answered_at,
		    hd.is_modified_by_mob,
		    hd.is_modified_by_acc		       
		FROM 
			help_desk hd
		WHERE
			hd.id=$1`, id)
}

func SelectHelpDeskModifiedByMob(db *sql.DB) (*sql.Rows, error) {
	return db.Query(`
		SELECT
			hd.id,
		    hd.user_id,
		    hd.date,
		    hd.title,
		    hd.body,
		    hd.status,
		    hd.answer,
		    hd.answered_by,
		    hd.answered_at,
		    hd.is_modified_by_mob,
		    hd.is_modified_by_acc		       
		FROM 
			help_desk hd
		WHERE
			hd.is_modified_by_mob=true`)
}

func SelectHelpDeskModifiedByAcc(db *sql.DB, userID string) (*sql.Rows, error) {
	return db.Query(`
		SELECT
			hd.id,
		    hd.user_id,
		    hd.date,
		    hd.title,
		    hd.body,
		    hd.status,
		    hd.answer,
		    hd.answered_by,
		    hd.answered_at,
		    hd.is_modified_by_mob,
		    hd.is_modified_by_acc		       
		FROM 
			help_desk hd
		WHERE
			hd.is_modified_by_acc=true
			AND pd.user_id = $1`, userID)
}

func SelectPayDeskByID(db *sql.DB, id int) (*sql.Rows, error) {
	return db.Query(`
		SELECT
			pd.id,
		    pd.user_id,
		    pd.amount,
		    pd.payment,
		    pd.document_number,
			pd.document_date,
			pd.file_paths,
			pd.files_quantity,
		    pd.created_at,
			pd.updated_at,	
			is_deleted,			    		    	    
		    pd.is_modified_by_mob,
		    pd.is_modified_by_acc		       
		FROM 
			pay_desk pd
		WHERE
			pd.id=$1`, id)
}

func SelectPayDeskModifiedByMob(db *sql.DB) (*sql.Rows, error) {
	return db.Query(`
		SELECT
			pd.id,
		    pd.user_id,
		    pd.amount,
		    pd.payment,
		    pd.document_number,
			pd.document_date,
			pd.file_paths,
			pd.files_quantity,
		    pd.created_at,
			pd.updated_at,
			is_deleted,			    		    		    
		    pd.is_modified_by_mob,
		    pd.is_modified_by_acc		       
		FROM 
			pay_desk pd
		WHERE
			pd.is_modified_by_mob=true`)
}

func SelectPayDeskModifiedByAcc(db *sql.DB, userID string) (*sql.Rows, error) {
	return db.Query(`
		SELECT
			pd.id,
		    pd.user_id,
		    pd.amount,
		    pd.payment,
		    pd.document_number,
			pd.document_date,
			pd.file_paths,
			pd.files_quantity,
		    pd.created_at,
			pd.updated_at,
			is_deleted,			    		    
		    pd.is_modified_by_mob,
		    pd.is_modified_by_acc		       
		FROM 
			pay_desk pd
		WHERE
			pd.is_modified_by_acc=true
			AND pd.user_id = $1`, userID)
}
