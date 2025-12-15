package users

import (
	"fmt"

	"github.com/dtekltd/common/auth"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/site"
)

// Account

func (m *Account) SID() string {
	return fmt.Sprintf("%d", m.ID)
}

func (m *Account) GetTokenData() *auth.AuthTokenData {
	return &auth.AuthTokenData{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		PublicID:  m.PublicID,
		AvatarUrl: m.AvatarUrl,
		IsAdmin:   m.IsAdmin,
	}
}

func (m *Account) Save() error {
	return database.DB.Save(m).Error
}

func (m *Account) UpdateField(field string, value any) error {
	return database.DB.Model(m).Update(field, value).Error
}

func (m *Account) GetReferrer() *Account {
	if m.Referrer == nil && m.ReferrerID != nil && m.ReferrerID.Uint64 > 0 {
		acc := &Account{}
		if err := database.DB.Take(acc, m.ReferrerID.Uint64).Error; err == nil {
			m.Referrer = acc
		}
	}
	return m.Referrer
}

func (m *Account) GetReferrers() *AccountReferrers {
	if m.Referrers == nil {
		r := &AccountReferrers{}
		if err := database.DB.Take(r, m.ID).Error; err == nil {
			m.Referrers = r
		}
	}
	return m.Referrers
}

func FindAccount(cond ...any) (*Account, error) {
	m := &Account{}
	if err := database.DB.Take(m, cond...).Error; err != nil {
		return nil, err
	}
	return m, nil
}

// AccountReferrers

func (m *AccountReferrers) GetReferrerSlices() []uint64 {
	return []uint64{m.R1, m.R2, m.R3, m.R4, m.R5, m.R6}
}

// utils
func FindReferrerIDByPublicID(referrer string) uint64 {
	var id uint64
	if referrer != "" {
		var ids []uint64
		if err := database.DB.Model(&Account{}).Select("id").
			Where("(username=?) OR (public_id=?)", referrer, referrer).
			Find(&ids).Error; err == nil {
			if len(ids) > 0 {
				return ids[0]
			}
		}
	}
	if id == 0 {
		if val := site.Settings().CustomParams.GetUint64("user.defaultReferrerId"); val != 0 {
			id = val
		}
	}
	return id
}

func FindReferrerIDByID(id uint64) uint64 {
	var ids []database.NullableUint64
	if err := database.DB.Model(&Account{}).Select("referrer_id").
		Where("id=?", id).Find(&ids).Error; err == nil {
		if len(ids) > 0 {
			return ids[0].Uint64
		}
	}
	return 0
}
