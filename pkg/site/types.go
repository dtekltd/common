package site

import (
	"github.com/dtekltd/common/pkg/settings"
	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
)

const MODULE = "site"

type NavItem struct {
	URL    string `json:"url,omitempty"`
	Name   string `json:"name,omitempty"`
	Icon   string `json:"icon,omitempty"`
	Target string `json:"target,omitempty"`
}

type (
	Menu        []NavItem
	Breadcrumbs []NavItem
)

type ContactInfo struct {
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Number  string `json:"number,omitempty"`
	Address string `json:"address,omitempty"`
}

type SiteParams struct {
	Env          string        `json:"env,omitempty"`
	Name         string        `json:"name,omitempty"`
	Logo         string        `json:"logo,omitempty"`
	Slogan       string        `json:"slogan,omitempty"`
	Website      string        `json:"website,omitempty"`
	Copyright    string        `json:"copyright,omitempty"`
	Contact      *ContactInfo  `json:"contact,omitempty"`
	MainMenu     *Menu         `json:"mainMenu,omitempty"`
	FooterMenu   *Menu         `json:"footerMenu,omitempty"`
	CustomParams *types.Params `json:"customParams,omitempty"`
}

type SiteSetting struct {
	*settings.Setting
	Params *datatypes.JSONType[*SiteParams] `json:"params" gorm:"type:text"`
}

type Meta struct {
	Title        string        `json:"title,omitempty"`
	Description  string        `json:"description,omitempty"`
	Canonical    string        `json:"canonical,omitempty"`
	FeatureImage string        `json:"featureImage,omitempty"`
	IntroText    string        `json:"introText,omitempty"`
	CustomBody   string        `json:"customBody,omitempty"`
	CustomHeader string        `json:"customHeader,omitempty"`
	CustomFooter string        `json:"customFooter,omitempty"`
	ProductInfo  *types.Params `json:"productInfo,omitempty"`
	CustomParams *types.Params `json:"customParams,omitempty"`
}

type Page struct {
	Title       string       `json:"title,omitempty"`
	Path        string       `json:"path,omitempty"`
	Content     string       `json:"content,omitempty"`
	Message     string       `json:"message,omitempty"`
	StatusCode  int          `json:"statusCode,omitempty"`
	Breadcrumbs *Breadcrumbs `json:"breadcrumbs,omitempty"`
	Meta        *Meta        `json:"meta,omitempty"`
}

type ContactReq struct {
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
}
