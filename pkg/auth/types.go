package auth

import (
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Referrer   string `json:"referrer"`
	ReferrerID uint64 `json:"referrerId"`
	AvatarUrl  string `json:"avatarUrl,omitempty"`
}

type OAuth2Info struct {
	AccountID uint64                            `json:"accountId"`
	Service   string                            `json:"service" gorm:"primaryKey;type:varchar(16);index:service_id"`
	ServiceID string                            `json:"serviceId" gorm:"primaryKey;type:varchar(36);index:service_id"`
	Data      *datatypes.JSONType[types.Params] `json:"data" gorm:"type:text"`
	CreatedAt uint64                            `json:"createdAt"`
	UpdatedAt uint64                            `json:"updatedAt"`
}

func (OAuth2Info) TableName() string {
	return "user_oauth2_info"
}

func Migrate() {
	database.DB.AutoMigrate(&OAuth2Info{})
}
