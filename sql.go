package inf

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
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

func (z *Dec) Float64() float64 {
	f, _ := strconv.ParseFloat(z.String(), 64)
	return f
}

func (z *Dec) SetFromFloat(f float64) *Dec {
	switch {
	case math.IsInf(f, 0):
		panic("cannot create a decimal from an infinte float")
	case math.IsNaN(f):
		panic("cannot create a decimal from an NaN float")
	}

	s := strconv.FormatFloat(f, 'e', -1, 64)

	// Determine the decimal's exponent.
	var e10 int64
	e := strings.IndexByte(s, 'e')
	for i := e + 2; i < len(s); i++ {
		e10 = e10*10 + int64(s[i]-'0')
	}
	switch s[e+1] {
	case '-':
		e10 = -e10
	case '+':
	default:
		panic(fmt.Sprintf("malformed float: %v -> %s", f, s))
	}
	e10++

	// Determine the decimal's mantissa.
	var mant int64
	i := 0
	neg := false
	if s[0] == '-' {
		i++
		neg = true
	}
	for ; i < e; i++ {
		if s[i] == '.' {
			continue
		}
		mant = mant*10 + int64(s[i]-'0')
		e10--
	}
	if neg {
		mant = -mant
	}

	return z.SetUnscaled(mant).SetScale(Scale(-e10))
}
