package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

type NullableUint64 struct {
	Uint64 uint64
	Valid  bool // Valid is true if Int64 is not NULL
}

// Scan implement Scan method to convert from database value
func (n *NullableUint64) Scan(value any) error {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true

	switch v := value.(type) {
	case uint64:
		n.Uint64 = v
		return nil
	case int64:
		if v < 0 {
			return errors.New("cannot convert negative int64 to uint64")
		}
		n.Uint64 = uint64(v)
		return nil
	case string:
		uv, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		n.Uint64 = uv
		return nil
	default:
		return errors.New("unsupported type for NullableUint64")
	}
}

// Value implement Value method to convert to database value
func (n NullableUint64) Value() (driver.Value, error) {
	if n.Uint64 == 0 {
		return nil, nil
	}
	return n.Uint64, nil
}

func (n NullableUint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Uint64)
}

func (n *NullableUint64) UnmarshalJSON(data []byte) error {
	var value uint64
	err := json.Unmarshal(data, &value)
	if err != nil {
		var null bool
		err = json.Unmarshal(data, &null)
		if err != nil {
			return err
		}
		if null {
			n.Uint64 = 0
		} else {
			return errors.New("invalid uint64 value")
		}
	} else {
		n.Uint64 = value
	}
	return nil
}
