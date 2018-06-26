package inf

import (
	"database/sql/driver"
	"errors"
)

var errInvalidType = errors.New("invalid type for inf.Dec")
var errParseFailure = errors.New("parse failure for inf.Dec")

func (z *Dec) Scan(value interface{}) error {
	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return errInvalidType
	}
	_, ok := z.SetString(s)
	if !ok {
		return errParseFailure
	}
	return nil
}

func (z *Dec) Value() (driver.Value, error) {
	return z.String(), nil
}
