package database

import (
	"database/sql/driver"
	"encoding/json"
)

// AnyType handles any type with proper null handling
//
// Usage example:
//
//	type MyModel struct {
//	   gorm.Model
//	   Data AnyType `gorm:"type:jsonb"` // Use jsonb for PostgreSQL or json for MySQL
//	}
type AnyType struct {
	Val   any
	Valid bool
}

// Scan implements sql.Scanner
func (a *AnyType) Scan(value any) error {
	if value == nil {
		a.Val, a.Valid = nil, false
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		a.Val = value
		a.Valid = true
		return nil
	}

	// Try to unmarshal as JSON first
	var jsonValue any
	if err := json.Unmarshal(bytes, &jsonValue); err == nil {
		a.Val = jsonValue
	} else {
		a.Val = string(bytes)
	}

	a.Valid = true
	return nil
}

// Value implements driver.Valuer
func (a AnyType) Value() (driver.Value, error) {
	if !a.Valid {
		return nil, nil
	}

	switch v := a.Val.(type) {
	case string, int, float64, bool:
		return v, nil
	default:
		return json.Marshal(v)
	}
}
