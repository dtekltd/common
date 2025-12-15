package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// FlexibleData represents any JSON-serializable data
//
// Usage example:
//
//	type MyModel struct {
//	   gorm.Model
//	   Data FlexibleData `gorm:"type:jsonb"` // Use jsonb for PostgreSQL or json for MySQL
//	}
type FlexibleData struct {
	Val any `json:"value"`
}

// Scan implements the sql.Scanner any
func (f *FlexibleData) Scan(value any) error {
	if value == nil {
		f.Val = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &f.Val)
}

// Value implements the driver.Valuer any
func (f FlexibleData) Value() (driver.Value, error) {
	if f.Val == nil {
		return nil, nil
	}
	return json.Marshal(f.Val)
}
