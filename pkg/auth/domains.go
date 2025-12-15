package auth

import (
	"fmt"

	"github.com/dtekltd/common/auth"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/users"
)

func (m *OAuth2Info) Save() error {
	return database.DB.Save(m).Error
}

func (m *OAuth2Info) UpdateField(field string, value any) error {
	return database.DB.Model(m).Update(field, value).Error
}

func FindOAuth2Info(service string, id string) (*OAuth2Info, error) {
	m := &OAuth2Info{}
	err := database.DB.Take(m, "service=? AND service_id=?", service, id).Error
	return m, err
}

func FindAuthAccount(id any) (*auth.AuthTokenData, error) {
	act, err := users.FindAccount(id)
	if err != nil {
		return nil, err
	}

	if act.ID == 0 {
		return nil, fmt.Errorf("account #%v does not exist", id)
	}

	return act.GetTokenData(), nil
}
