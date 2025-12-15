package cms

import (
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/settings"
	"github.com/dtekltd/common/system"
	"gorm.io/datatypes"
)

var setting *CmsSetting

func Settings() *CmsParams {
	if setting == nil {
		setting = defaultSettings()
		if err := database.DB.Take(setting, "module=?", MODULE).Error; err != nil {
			system.Logger.Error("error load cms-settings", err.Error())
			if err := database.DB.Save(setting).Error; err != nil {
				system.Logger.Error("error save cms-settings", err.Error())
				system.Logger.Error(err.Error())
			}
		}
	}
	params := setting.Params.Data()
	return &params
}

func SaveSettings() error {
	// make sure setting instance exists
	_ = Settings()
	return database.DB.Save(setting).Error
}

func UpdateSettings(params *CmsParams) error {
	_ = Settings()
	jsonParams := datatypes.NewJSONType(*params)
	setting.Params = &jsonParams
	return database.DB.Save(setting).Error
}

func defaultSettings() *CmsSetting {
	setting := &CmsSetting{
		Setting: &settings.Setting{
			Module: MODULE,
			Name:   "CMS Settings",
		},
	}
	params := CmsParams{
		HomePageID:   0,
		BlogPageSize: 12,
	}
	jsonParams := datatypes.NewJSONType(params)
	setting.Params = &jsonParams
	return setting
}
