package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

//NullTime special type for scan sql rows with Null data for time type variables
type NullTime struct {
	Time  time.Time
	Valid bool
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

	lenStr := len([]rune(asString))

	if lenStr > 4 {
		asString, err = strconv.Unquote(string(data))
		if err != nil {
			nt.Valid = false
			return fmt.Errorf("NullTime unmarshal error: invalid input %s\n%v", asString, err)
		}
	}

	if asString == "null" || asString == "nil" || asString == "" {
		nt.Valid = false
	} else {
		nt.Time, err = time.Parse("2006-01-02T15:04:05", asString)
		if err != nil {
			nt.Valid = false
			return fmt.Errorf("NullTime unmarshal error: invalid input %s\n%v", asString, err)
		}
		nt.Valid = true
	}

	return nil
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	var asString string

	if nt.Valid {
		asString = nt.Time.Format("2006-01-02T15:04:05")
		asString = strconv.Quote(asString)
	} else {
		asString = "null"
	}

	return []byte(asString), nil
}
