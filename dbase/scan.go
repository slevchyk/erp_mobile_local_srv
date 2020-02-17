package dbase

import (
	"database/sql"
	"github.com/slevchyk/erp_mobile_local_srv/models"
)

func ScanFirebaseToken(rows *sql.Rows, ft *models.FirebaseTokens) error {
	return rows.Scan(&ft.ID, &ft.UserID, &ft.Token)
}

func ScanChannel(rows *sql.Rows, c *models.Channel) error {
	return rows.Scan(&c.ID, &c.UserID, &c.UpdateID, &c.Type, &c.Title, &c.News, &c.DateTime)
}
