package dbase

import (
	"database/sql"

	"github.com/slevchyk/erp_mobile_local_srv/models"
)

func ScanFirebaseToken(rows *sql.Rows, ft *models.FirebaseTokens) error {
	return rows.Scan(
		&ft.ID,
		&ft.UserID,
		&ft.Token)
}

func ScanChannel(rows *sql.Rows, c *models.Channel) error {
	return rows.Scan(
		&c.ID, &c.UserID,
		&c.UpdateID,
		&c.Type,
		&c.Title,
		&c.News,
		&c.Date)
}

func ScanTiming(rows *sql.Rows, t *models.Timing) error {
	return rows.Scan(
		&t.ID,
		&t.MobID,
		&t.AccID,
		&t.UserID,
		&t.Date,
		&t.Status,
		&t.IsTurnstile,
		&t.StartedAt,
		&t.EndedAt,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt)
}

func ScanProfile(rows *sql.Rows, p *models.Profile) error {
	return rows.Scan(
		&p.ID,
		&p.Blocked,
		&p.UserID,
		&p.Pin,
		&p.InfoCard,
		&p.LastName,
		&p.FirstName,
		&p.MiddleName,
		&p.ITN,
		&p.Phone,
		&p.Birthday,
		&p.Email,
		&p.Gender,
		&p.Address,
		&p.PassportType,
		&p.PassportSeries,
		&p.PassportNumber,
		&p.PassportIssued,
		&p.PassportDate,
		&p.PassportExpiry,
		&p.CivilStatus,
		&p.Children,
		&p.JobPosition,
		&p.Education,
		&p.Specialty,
		&p.AdditionalEducation,
		&p.LastWorkPlace,
		&p.Skills,
		&p.Languages,
		&p.Disability,
		&p.Pensioner,
		&p.PhotoName,
		&p.PhotoData)
}

func ScanHelpDesk(rows *sql.Rows, hd *models.HelpDesk) error {
	return rows.Scan(
		&hd.ID,
		&hd.UserID,
		&hd.Date,
		&hd.Title,
		&hd.Body,
		&hd.Status,
		&hd.Answer,
		&hd.AnsweredBy,
		&hd.AnsweredAt,
		&hd.IsModifiedByMob,
		&hd.IsModifiedByAcc,
	)
}

func ScanPayDesk(rows *sql.Rows, pd *models.PayDesk) error {
	return rows.Scan(
		&pd.ID,
		&pd.PayDeskType,
		&pd.UserID,
		&pd.CurrencyAccID,
		&pd.CostItemAccID,
		&pd.IncomeItemAccID,
		&pd.FromPayOfficeAccID,
		&pd.ToPayOfficeAccID,
		&pd.Amount,
		&pd.Payment,
		&pd.DocumentNumber,
		&pd.DocumentDate,
		&pd.FilePaths,
		&pd.FilesQuantity,
		&pd.IsChecked,
		&pd.CreatedAt,
		&pd.UpdatedAt,
		&pd.IsDeleted,
		&pd.IsModifiedByMob,
		&pd.IsModifiedByAcc,
	)
}

func ScanCostItem(rows *sql.Rows, ci *models.CostItem) error {
	return rows.Scan(
		&ci.ID,
		&ci.AccID,
		&ci.Name,
		&ci.CreatedAt,
		&ci.UpdatedAt,
		&ci.IsDeleted,
	)
}

func ScanIncomeItem(rows *sql.Rows, ii *models.IncomeItem) error {
	return rows.Scan(
		&ii.ID,
		&ii.AccID,
		&ii.Name,
		&ii.CreatedAt,
		&ii.UpdatedAt,
		&ii.IsDeleted,
	)
}

func ScanPayOffice(rows *sql.Rows, po *models.PayOffice) error {
	return rows.Scan(
		&po.ID,
		&po.AccID,
		&po.CurrencyAccID,
		&po.Name,
		&po.CreatedAt,
		&po.UpdatedAt,
		&po.IsDeleted,
	)
}

func ScanPayOfficeBalance(rows *sql.Rows, pob *models.PayOfficeBalance) error {
	return rows.Scan(
		&pob.AccID,
		&pob.Balance,
		&pob.UpdatedAt,
	)
}

func ScanCurrency(rows *sql.Rows, c *models.Currency) error {
	return rows.Scan(
		&c.ID,
		&c.AccID,
		&c.Code,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.IsDeleted,
	)
}

func ScanUserGrants(rows *sql.Rows, ug *models.UserGrants) error {
	return rows.Scan(
		&ug.UserID,
		&ug.ObjectType,
		&ug.ObjectAccID,
		&ug.IsVisible,
		&ug.IsAvailable,
	)
}
