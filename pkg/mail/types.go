package mail

import (
	"html/template"

	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var db *gorm.DB

// type Options struct {
// 	From     string           `json:"from,omitempty"`
// 	To       []string         `json:"to,omitempty"`
// 	Cc       []string         `json:"cc,omitempty"`
// 	Bcc      []string         `json:"bcc,omitempty"`
// 	Target   string           `json:"target,omitempty"`
// 	Template string           `json:"template,omitempty"`
// }

// type Options settings.Setting

type Message struct {
	ID          uint64                            `json:"id,omitempty" gorm:"primaryKey"`
	Type        string                            `json:"type,omitempty" gorm:"type:varchar(16);index"`
	Alias       *string                           `json:"alias,omitempty" gorm:"type:varchar(128);unique"`
	Subject     string                            `json:"subject,omitempty" gorm:"type:varchar(128);"`
	Body        string                            `json:"body,omitempty" gorm:"type:text"`
	bodyTmplt   *template.Template                `json:"-" gorm:"-"`
	Priority    int                               `json:"priority"`
	Frequency   int                               `json:"frequency"`
	Options     *datatypes.JSONType[types.Params] `json:"options,omitempty" gorm:"type:text"`
	CreatedAt   uint64                            `json:"created_at,omitempty"`
	ProcessedAt uint64                            `json:"processed_at,omitempty"`
	// Instances     []MailInstance                       `json:"instances,omitempty"`
}

type Instance struct {
	ID          uint64                            `json:"id,omitempty" gorm:"primaryKey"`
	MessageID   uint64                            `json:"messageId"`
	Options     *datatypes.JSONType[types.Params] `json:"options,omitempty" gorm:"type:varchar(500)"`
	Priority    int                               `json:"priority"`
	FailedCount int                               `json:"failedCount"`
	CreatedAt   uint64                            `json:"created_at,omitempty"`
	SentAt      uint64                            `json:"sent_at,omitempty"`
	Message     *Message                          `json:"message,omitempty"`
}

// TableName overrides the table name
func (Message) TableName() string {
	return "mail_messages"
}

func (Instance) TableName() string {
	return "mail_instances"
}

func Migrate(_db *gorm.DB) {
	db = _db
	db.AutoMigrate(&Message{}, &Instance{})
}

func DB(conns ...*gorm.DB) *gorm.DB {
	if len(conns) > 0 {
		db = conns[0]
	}
	return db
}
