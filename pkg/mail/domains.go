package mail

import (
	"fmt"
	"html/template"
	"regexp"

	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
)

func GetMessage(alias string) *Message {
	m := &Message{}
	if err := db.Take(m, "alias=?", alias).Error; err != nil {
		return nil
	}
	return m
}

func (m *Message) Save() error {
	return db.Save(m).Error
}

func (m *Message) Update() error {
	// return db.Model(m).Updates(map[string]any{
	// 	"processed_at": m.ProcessedAt,
	// 	"frequency":    m.Frequency,
	// }).Error
	return db.Model(m).Select("processed_at", "frequency").Updates(m).Error
}

func (m *Message) Delete() error {
	return db.Delete(m, m.ID).Error
}

func (m *Message) SetOptions(opts types.Params) {
	if len(opts) > 0 {
		jsonOpts := datatypes.NewJSONType(opts)
		m.Options = &jsonOpts
	}
}

func (m *Message) ParseBodyTmplt() error {
	opts := types.Params{}
	if m.Options != nil {
		opts = m.Options.Data()
	}
	if m.bodyTmplt == nil && regexp.MustCompile(`{{\s*\.\w+\s*}}`).MatchString(m.Body) {
		if tmplt, err := template.New("mail.body").Parse(m.Body); err != nil {
			opts["__error"] = fmt.Sprintf("Parse body template error: %s", err.Error())
			// jsonOpts := datatypes.NewJSONType(opts)
			// m.Options = &jsonOpts
			return err
		} else {
			m.bodyTmplt = tmplt
			opts["__tmpl"] = true
		}
	} else {
		opts["__tmpl"] = false
	}
	m.SetOptions(opts)
	return nil
}

func (m *Message) GetBodyTmplt() (*template.Template, error) {
	opts := m.Options.Data()
	if val := opts.Get("__tmpl"); val == nil {
		if err := m.ParseBodyTmplt(); err != nil {
			return nil, err
		}
	} else if val.(bool) && m.bodyTmplt == nil {
		// only parse when has __tmpl
		if err := m.ParseBodyTmplt(); err != nil {
			return nil, err
		}
	}
	return m.bodyTmplt, nil
}

func (m *Instance) Save() error {
	// return db.Session(&gorm.Session{
	// 	FullSaveAssociations: false,
	// }).Save(m).Error
	msg := m.Message
	if msg != nil {
		m.Message = nil
	}
	err := db.Save(m).Error
	if msg != nil {
		m.Message = msg
	}
	return err
}

func (m *Instance) Delete() error {
	return db.Delete(m, m.ID).Error
}

func (m *Instance) SetOptions(opts types.Params) {
	if len(opts) > 0 {
		jsonOpts := datatypes.NewJSONType(opts)
		m.Options = &jsonOpts
	}
}
