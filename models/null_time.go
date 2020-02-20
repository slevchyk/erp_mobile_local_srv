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
	asString := string(data)

	if asString == "null"{
		*nt.Time = nil
		*nt.Valid = false
	} else if asString == "0" || asString == "false" {
		*nt.Time = nil
		*nt.Valid = true
	} else {
		return fmt.Errorf("NullTime unmarshal error: invalid input %s", asString)
	}
	return nil
}

func (bit *ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		*bit = true
	} else if asString == "0" || asString == "false" {
		*bit = false
	} else {
		return fmt.Errorf("Boolean unmarshal error: invalid input %s", asString)
	}
	return nil
}
