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

func InsertTiming(db *sql.DB, t models.Timing) (int64, error) {
	var lastInsertId int64
	err := db.QueryRow(`
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		t.MobID, t.AccID, t.UserID, t.Date, t.Status, t.IsTurnstile, t.StartedAt, t.EndedAt, t.CreatedAt, t.UpdatedAt, t.DeletedAt).Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertProfile(db *sql.DB, p models.Profile) (int64, error) {
	var lastInsertId int64
	err := db.QueryRow(`
		INSERT INTO
			profiles (	
				blocked,
				user_id,
				pin,
				info_card,
				last_name,
				first_name,
				middle_name,
				itn,
				phone,
				birthday,
				email,
				gender,
				address,
				passport_type,
				passport_series,
				passport_number,
				passport_issued,
				passport_date,
				passport_expiry,
				civil_status,
				job_position,
				children,
				education,
				specialty,
				additional_education,
				last_work_place,
				skills,
				languages,
				disability,
				pensioner,
				photo,
				photo_data
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32) RETURNING id`,
		p.Blocked,
		p.UserID,
		p.Pin,
		p.InfoCard,
		p.LastName,
		p.FirstName,
		p.MiddleName,
		p.ITN,
		p.Photo,
		p.Birthday,
		p.Email,
		p.Gender,
		p.Address,
		p.PassportType,
		p.PassportSeries,
		p.PassportNumber,
		p.PassportIssued,
		p.PassportDate,
		p.PassportExpiry,
		p.CivilStatus,
		p.JobPosition,
		p.Children,
		p.Education,
		p.Specialty,
		p.AdditionalEducation,
		p.LastWorkPlace,
		p.Skills,
		p.Languages,
		p.Disability,
		p.Pensioner,
		p.Photo,
		p.PhotoData,
	).Scan(&lastInsertId)

	return lastInsertId, err
}
