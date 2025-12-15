package database

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// GenericData uses GORM's datatypes.JSON with custom handling
//
// Usage example:
//
//	type MyModel struct {
//	   gorm.Model
//	   Data GenericData `gorm:"type:jsonb"` // Use jsonb for PostgreSQL or json for MySQL
//	}
type GenericData datatypes.JSON

func (g *GenericData) Unmarshal(v any) error {
	return json.Unmarshal(*g, v)
}

func (g GenericData) Marshal(v any) (GenericData, error) {
	bytes, err := json.Marshal(v)
	return GenericData(bytes), err
}
