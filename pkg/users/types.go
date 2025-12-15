package users

import (
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
)

// ReferrerID   *uint64           `json:"referrerId" gorm:"null"`
type Account struct {
	ID         uint64                   `json:"id,omitempty" gorm:"primaryKey"`
	ReferrerID *database.NullableUint64 `json:"referrerId,omitempty"`
	PublicID   string                   `json:"publicId,omitempty" gorm:"type:varchar(8);uniqueIndex"`
	Name       string                   `json:"name,omitempty" gorm:"type:varchar(64);"`
	Email      string                   `json:"email,omitempty" gorm:"type:varchar(64);uniqueIndex"`
	Phone      string                   `json:"phone,omitempty" gorm:"type:varchar(16)"`
	Username   string                   `json:"username,omitempty" gorm:"type:varchar(64);uniqueIndex"`
	AvatarUrl  string                   `json:"avatarUrl,omitempty" gorm:"type:varchar(256)"`
	Password   string                   `json:"-" gorm:"type:varchar(128)"`
	State      int                      `json:"state,omitempty"`
	IsAdmin    bool                     `json:"isAdmin,omitempty"`
	CreatedAt  uint64                   `json:"created_at,omitempty"`
	UpdatedAt  uint64                   `json:"updated_at,omitempty"`
	DeletedAt  uint64                   `json:"-"`
	Referrer   *Account                 `json:"referrer,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Referrers  *AccountReferrers        `json:"referrers,omitempty" gorm:"foreignKey:ID"`
}

type AccountReferrers struct {
	ID  uint64 `json:"id" gorm:"primaryKey;autoIncrement:false"`
	R1  uint64 `json:"r1"`
	R2  uint64 `json:"r2"`
	R3  uint64 `json:"r3"`
	R4  uint64 `json:"r4"`
	R5  uint64 `json:"r5"`
	R6  uint64 `json:"r6"`
	R7  uint64 `json:"r7"`
	R8  uint64 `json:"r8"`
	R9  uint64 `json:"r9"`
	R10 uint64 `json:"r10"`
}

type Message struct {
	ID        uint64                            `json:"id,omitempty" gorm:"primaryKey"`
	AccountID uint64                            `json:"accountId"`
	Title     string                            `json:"title" gorm:"type:varchar(128)"`
	Body      string                            `json:"body" gorm:"type:varchar(512)"`
	Action    string                            `json:"action" gorm:"type:varchar(32)"`
	Data      *datatypes.JSONType[types.Params] `json:"params,omitempty" gorm:"type:text"`
	CreatedAt uint64                            `json:"created_at"`
	UpdatedAt uint64                            `json:"updated_at"`
	Account   *Account                          `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName overrides the table name
func (Account) TableName() string {
	return "user_accounts"
}

func (AccountReferrers) TableName() string {
	return "user_account_referrers"
}

func (Message) TableName() string {
	return "user_messages"
}

func Migrate() {
	database.DB.AutoMigrate(&Account{}, &AccountReferrers{}, &Message{})
}
