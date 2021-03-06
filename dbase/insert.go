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
		t.DeletedAt).
		Scan(&lastInsertId)

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
				children,
				job_position,
				education,
				specialty,
				additional_education,
				last_work_place,
				skills,
				languages,
				disability,
				pensioner,
				photo_name,
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
		p.PhotoData).
		Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertHelpDesk(db *sql.DB, hd models.HelpDesk) (int64, error) {
	var lastInsertId int64

	err := db.QueryRow(`
		INSERT INTO
			help_desk (
				user_id,
				date,
			    title,
			    body,
			    status,
			    answer,
			    answered_by,
			    answered_at,
			    is_modified_by_mob,
			    is_modified_by_acc
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		hd.UserID,
		hd.Date,
		hd.Title,
		hd.Body,
		hd.Status,
		hd.Answer,
		hd.AnsweredBy,
		hd.AnsweredAt,
		hd.IsModifiedByMob,
		hd.IsModifiedByAcc).
		Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertPayDesk(db *sql.DB, pd models.PayDesk) (int64, error) {
	var lastInsertId int64

	err := db.QueryRow(`
		INSERT INTO
			pay_desk (
				pay_desk_type,
				user_id,
				currency_acc_id,
				cost_item_acc_id,
				income_item_acc_id,	
				from_pay_office_acc_id,
				to_pay_office_acc_id,
				amount,
			    payment,
			    document_number,
				document_date,
				file_paths,
				files_quantity,
				is_checked,
			    is_read_only,
			    created_at,
				updated_at,
				is_deleted,			    
			    is_modified_by_mob,
			    is_modified_by_acc
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20) RETURNING id`,
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
		pd.IsModifiedByAcc).
		Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertCostItem(db *sql.DB, ci models.CostItem) (int64, error) {

	var lastInsertId int64

	err := db.QueryRow(`
		INSERT INTO
			cost_items (
				acc_id,
				name,			    
			    created_at,
				updated_at,
				is_deleted
			)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		ci.AccID,
		ci.Name,
		ci.CreatedAt,
		ci.UpdatedAt,
		ci.IsDeleted).
		Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertIncomeItem(db *sql.DB, ii models.IncomeItem) (int64, error) {
	var lastInsertId int64

	err := db.QueryRow(`
		INSERT INTO
			income_items (
				acc_id,
				name,			    
			    created_at,
				updated_at,
				is_deleted
			)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		ii.AccID,
		ii.Name,
		ii.CreatedAt,
		ii.UpdatedAt,
		ii.IsDeleted).
		Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertPayOffice(db *sql.DB, po models.PayOffice) (int64, error) {
	var lastInsertId int64

	err := db.QueryRow(`
		INSERT INTO
			pay_offices (
				acc_id,
				name,	
				currency_acc_id,		    
			    created_at,
				updated_at,
				is_deleted
			)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		po.AccID,
		po.Name,
		po.CurrencyAccID,
		po.CreatedAt,
		po.UpdatedAt,
		po.IsDeleted).
		Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertPayOfficeBalance(db *sql.DB, pob models.PayOfficeBalance) (string, error) {
	var lastInsertAccId string

	err := db.QueryRow(`
		INSERT INTO
			pay_offices_balance (
				acc_id,
				balance,		    			    
				updated_at
			)
		VALUES ($1, $2, $3) RETURNING acc_id`,
		pob.AccID,
		pob.Balance,
		pob.UpdatedAt).
		Scan(&lastInsertAccId)

	return lastInsertAccId, err
}

func InsertCurrency(db *sql.DB, c models.Currency) (int64, error) {
	var lastInsertId int64

	err := db.QueryRow(`
		INSERT INTO
			currency (
				acc_id,
				code,
				name,			    
			    created_at,
				updated_at,
				is_deleted
			)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		c.AccID,
		c.Code,
		c.Name,
		c.CreatedAt,
		c.UpdatedAt,
		c.IsDeleted).
		Scan(&lastInsertId)

	return lastInsertId, err
}

func InsertUserGrants(db *sql.DB, ug models.UserGrants) (string, error) {
	var lastInsertUserId string

	err := db.QueryRow(`
		INSERT INTO
			user_grants (
				user_id,
				odject_type,
				odject_acc_id,
				is_visible,			    
				is_available,
				is_receiver
			)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`,
		ug.UserID,
		ug.ObjectType,
		ug.ObjectAccID,
		ug.IsVisible,
		ug.IsAvailable,
		ug.IsReceiver).
		Scan(&lastInsertUserId)

	return lastInsertUserId, err
}

func InsertPayDeskImage(db *sql.DB, pdi models.PayDeskImage) error {
	var lastInsertId int64

	err := db.QueryRow(`
		INSERT INTO
			pay_desk_images (
				pid,
				image_name,
				file,
				sha256,
				is_deleted
			)
		VALUES ($1, $2, $3, $4, $5) RETURNING pid`,
		pdi.PID,
		pdi.ImageName,
		pdi.File,
		pdi.Sha256,
		pdi.IsDeleted).
		Scan(&lastInsertId)

	return err
}

func InsertLogInfo(db *sql.DB, log models.LogInfo) error {
	var lastInsertUserId string

	err := db.QueryRow(`
		INSERT INTO
			logs_info_files (
				user_id,
				file_name,
				file,
				date
			)
		VALUES ($1, $2, $3, $4) RETURNING user_id`,
		log.UserID,
		log.FileName,
		log.File,
		log.Date).
		Scan(&lastInsertUserId)

	return err
}