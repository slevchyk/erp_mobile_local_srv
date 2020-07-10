package models

import "fmt"

//ConvertibleBoolean special type for scan sql rows with Null data for bool type variables
type ConvertibleBoolean bool

//UnmarshalJSON override typical method
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

// Scan implements the Scanner interface.
func (bit *ConvertibleBoolean) Scan(value interface{}) error {
	if value == nil {
		*bit = false
	} else if value.(bool) {
		*bit = true
	} else {
		*bit = false
	}

	return nil
}
