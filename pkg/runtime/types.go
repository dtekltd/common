package runtime

import (
	"gorm.io/gorm"
)

var db *gorm.DB

type Runtime struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Key         string `json:"key" gorm:"type:varchar(128);unique_index"`
	Data        string `json:"data" gorm:"type:text"`
	CreatedAt   int64  `json:"createdAt,omitempty"`
	UpdatedAt   int64  `json:"updatedAt,omitempty"`
	ExpiredTime int    `json:"expiredTime,omitempty"`
}

// TableName overrides the table name
func (Runtime) TableName() string {
	return "com_runtimes"
}

func Migrate(_db *gorm.DB) {
	db = _db
	db.AutoMigrate(&Runtime{})
}
