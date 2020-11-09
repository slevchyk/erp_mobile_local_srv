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
				photo_name = $31,
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
		p.PhotoName,
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
				is_modified_by_mob = $9,
				is_modified_by_acc = $10				
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
		hd.IsModifiedByAcc,
		hd.ID)
	return res, err
}

func UpdatePayDesk(db *sql.DB, pd models.PayDesk) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				pay_desk
			SET
				pay_desk_type = $1,
				user_id = $2, 
				currency_acc_id = $3,
				cost_item_acc_id = $4,
				income_item_acc_id = $5,	
				from_pay_office_acc_id = $6,
				to_pay_office_acc_id = $7,
				amount = $8,
				payment = $9,
				document_number = $10,
				document_date = $11,
				file_paths = $12,
				files_quantity = $13,
				is_checked = $14,
				is_read_only = $15,
				created_at = $16,
				updated_at = $17,
				is_deleted =$18,			    				
				is_modified_by_mob = $19,
				is_modified_by_acc = $20			
			WHERE
				id=$21
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		pd.PayDeskType,
		pd.UserID,
		pd.CurrencyAccID,
		pd.CostItemAccID,
		pd.IncomeItemAccID,
		pd.FromPayOfficeAccID,
		pd.ToPayOfficeAccID,
		pd.Amount,
		pd.Payment,
		pd.DocumentNumber,
		pd.DocumentDate,
		pd.FilePaths,
		pd.FilesQuantity,
		pd.IsChecked,
		pd.IsReadOnly,
		pd.CreatedAt,
		pd.UpdatedAt,
		pd.IsDeleted,
		pd.IsModifiedByMob,
		pd.IsModifiedByAcc,
		pd.ID)
	return res, err
}

func UpdateCostItem(db *sql.DB, ci models.CostItem) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				cost_items
			SET
				acc_id = $1,
				name = $2,				
				created_at = $3,
				updated_at = $4,
				is_deleted =$5				
			WHERE
				id=$6
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		ci.AccID,
		ci.Name,
		ci.CreatedAt,
		ci.UpdatedAt,
		ci.IsDeleted,
		ci.ID)
	return res, err
}

func UpdateIncomeItem(db *sql.DB, ii models.IncomeItem) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				income_items
			SET
				acc_id = $1,
				name = $2,				
				created_at = $3,
				updated_at = $4,
				is_deleted =$5				
			WHERE
				id=$6
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		ii.AccID,
		ii.Name,
		ii.CreatedAt,
		ii.UpdatedAt,
		ii.IsDeleted,
		ii.ID)
	return res, err
}

func UpdatePayOffice(db *sql.DB, po models.PayOffice) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				pay_offices
			SET
				acc_id = $1,
				currency_acc_id = $2,
				name = $3,				
				created_at = $4,
				updated_at = $5,
				is_deleted =$6				
			WHERE
				id=$7
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		po.AccID,
		po.CurrencyAccID,
		po.Name,
		po.CreatedAt,
		po.UpdatedAt,
		po.IsDeleted,
		po.ID)
	return res, err
}

func UpdatePayOfficeBalance(db *sql.DB, pob models.PayOfficeBalance) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				pay_offices_balance
			SET
				balance = $1,
				updated_at = $2				
			WHERE
				acc_id=$3
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		pob.Balance,
		pob.UpdatedAt,
		pob.AccID)
	return res, err
}

func UpdateCurrency(db *sql.DB, c models.Currency) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				currency
			SET
				acc_id = $1,
				code = $2,
				name = $3,				
				created_at = $4,
				updated_at = $5,
				is_deleted =$6				
			WHERE
				id=$7
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		c.AccID,
		c.Code,
		c.Name,
		c.CreatedAt,
		c.UpdatedAt,
		c.IsDeleted,
		c.ID)
	return res, err
}

func UpdateUserGrants(db *sql.DB, ug models.UserGrants) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				user_grants
			SET				
				is_visible = $1,				
				is_available = $2,
				is_receiver = $3
			WHERE
				user_id=$4
				AND odject_type = $5
				AND odject_acc_id = $6
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		ug.IsVisible,
		ug.IsAvailable,
		ug.IsReceiver,
		ug.UserID,
		ug.ObjectType,
		ug.ObjectAccID)
	return res, err
}

func UpdatePayDeskImages(db *sql.DB, pdi models.PayDeskImage) (sql.Result, error) {

	var err error

	stmt, err := db.Prepare(`
			UPDATE
				pay_desk_images
			SET				
				is_deleted = $1 
			WHERE
				pid =$2
				AND 
				image_name =$3
			`)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(
		pdi.IsDeleted,
		pdi.PID,
		pdi.ImageName)
	return res, err
}
