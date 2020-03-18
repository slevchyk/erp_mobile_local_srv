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
	res, err := stmt.Exec(
		t.MobID,
		t.AccID,
		t.UserID,
		t.Date,
		t.Status,
		t.IsTurnstile,
		t.StartedAt,
		t.EndedAt,
		t.CreatedAt,
		t.UpdatedAt,
		t.DeletedAt,
		t.ID)

	return res, err
}

func UpdateProfile(db *sql.DB, p models.Profile) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				profiles
			SET
				blocked = $1,
				user_id = $2,
				pin = $3,
				info_card = $4,
				last_name = $5,
				first_name = $6,
				middle_name = $7,
				itn = $8,
				phone = $9,
				birthday = $10,
				email = $11,
				gender = $12,
				address = $13,
				passport_type = $14,
				passport_series = $15,
				passport_number = $16,
				passport_issued = $17,
				passport_date = $18,
				passport_expiry = $19,
				civil_status = $20,				
				children = $21,
				job_position = $22,
				education = $23,
				specialty = $24,
				additional_education = $25,
				last_work_place = $26,
				skills = $27,
				languages = $28,
				disability = $29,
				pensioner = $30,
				photo = $31,
				photo_data = $32			
			WHERE
				id=$33
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		p.Blocked,
		p.UserID,
		p.Pin,
		p.InfoCard,
		p.LastName,
		p.FirstName,
		p.MiddleName,
		p.ITN,
		p.Phone,
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
		p.Children,
		p.JobPosition,
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
		p.ID)

	return res, err
}

func UpdateHelpDesk(db *sql.DB, hd models.HelpDesk) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				help_desk
			SET
				user_id = $1,
				date = $2,
				title = $3,
				body = $4,
				status = $5,
				answer = $6,
				answered_by = $7,
				answered_at = $8,
				is_modified_mob = $9,
				is_modified_acc = $10				
			WHERE
				id=$11
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		hd.UserID,
		hd.Date,
		hd.Title,
		hd.Body,
		hd.Status,
		hd.Answer,
		hd.AnsweredBy,
		hd.AnsweredAt,
		hd.IsModifiedByMob,
		hd.IsModifiedAcc,
		hd.ID)
	return res, err
}
