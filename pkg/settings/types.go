package settings

import (
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
)

type Setting struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	Module    string `json:"module" gorm:"type:varchar(64);unique_index"`
	Name      string `json:"name" gorm:"type:varchar(128)"`
	CreatedAt uint64 `json:"createdAt"`
	UpdatedAt uint64 `json:"updatedAt"`
}

type GeneralSetting struct {
	*Setting
	Params *datatypes.JSONType[types.Params] `json:"params" gorm:"type:text"`
}

// TableName overrides the table name
func (GeneralSetting) TableName() string {
	return "com_settings"
}

func Migrate() {
	database.DB.AutoMigrate(&GeneralSetting{})
}

type SettingParams interface {
	Get(key string) any
	Set(key string, val any)
}
