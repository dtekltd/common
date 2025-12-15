package site

// TableName overrides the table name
func (SiteSetting) TableName() string {
	return "com_settings"
}

func (m *SiteParams) IsDEV() bool {
	return m.Env == "DEV"
}

func (m *SiteParams) IsPROD() bool {
	return m.Env == "PROD"
}
