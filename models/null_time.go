package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

//NullTime special type for scan sql rows with Null data for time type variables
type NullTime struct {
	Time  time.Time `json:"time"`
	Valid bool      `json:"valid"` // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var err error
	asString := string(data)

	if asString == "null" {
		nt.Valid = false
	} else if asString == "0" || asString == "false" {
		nt.Time, err = time.Parse(time.RFC3339, asString)
		if err != nil {
			nt.Valid = false
			return err
		}
		nt.Valid = true
	} else {
		return fmt.Errorf("NullTime unmarshal error: invalid input %s", asString)
	}
	return nil
}
