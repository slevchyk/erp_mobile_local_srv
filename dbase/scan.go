package dbase

import (
	"database/sql"
	"github.com/slevchyk/erp_mobile_local_srv/models"
)

func ScanFirebaseToken(rows *sql.Rows, ft *models.FirebaseTokens) error {
	return rows.Scan(&ft.ID, &ft.UserID, &ft.Token)
}
