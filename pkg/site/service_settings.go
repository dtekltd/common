package site

import (
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/settings"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
)

var setting *SiteSetting

func Settings() *SiteParams {
	if setting == nil {
		setting = defaultSettings()
		if err := database.DB.Take(setting, "module=?", MODULE).Error; err != nil {
			system.Logger.Error("error load site-settings", err.Error())
			if err := database.DB.Save(setting).Error; err != nil {
				system.Logger.Error("error save site-setting", err.Error())
			}
		}
	}
	return setting.Params.Data()
}

func SaveSettings() error {
	// make sure setting instance exists
	_ = Settings()
	return database.DB.Save(setting).Error
}

func UpdateSettings(params *SiteParams) error {
	_ = Settings()
	jsonParams := datatypes.NewJSONType(params)
	setting.Params = &jsonParams
	return database.DB.Save(setting).Error
}

func defaultSettings() *SiteSetting {
	setting := &SiteSetting{
		Setting: &settings.Setting{
			Module: MODULE,
			Name:   "Site Settings",
		},
	}
	params := &SiteParams{
		Name:      "Go! CMS",
		Logo:      "",
		Slogan:    "Simpler. Faster. Cheaper.",
		Copyright: "Copyright 2025 - Go! CMS",
		Contact: &ContactInfo{
			Name:    "Peter Phan",
			Number:  "+84 90 817 2887",
			Email:   "peter.phan07@gmail.com",
			Address: "348/58 Hoang Van Thu, Ward 4, Tan Binh Dist., HCM City",
		},
		CustomParams: &types.Params{
			"key": "value",
		},
	}
	jsonParams := datatypes.NewJSONType(params)
	setting.Params = &jsonParams
	return setting
}
