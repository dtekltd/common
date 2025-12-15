package mail

import (
	"bytes"
	"fmt"
	"maps"
	"time"

	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/mailer"
	"github.com/dtekltd/common/pkg/site"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SendMessage send the message with optional params
//
//	types.Params{
//		"to": acc.Email,
//		"Name": acc.Name,
//		"Email": acc.Email,
//		"Password": tmpPass,
//		"__immediate": true,
//	}
func SendMessage(msg *Message, params types.Params) error {
	opts := types.Params{}
	if msg.Options != nil {
		maps.Copy(opts, msg.Options.Data())
	}
	if params == nil {
		params = types.Params{}
	}
	if customParams := opts.GetParams("params"); customParams != nil {
		maps.Copy(params, customParams)
	}

	if msg.Type == "system" {
		if len(params) == 0 {
			return fmt.Errorf("mail-message #%d system requires params", msg.ID)
		}
	} else {
		if msg.ProcessedAt > 0 && !params.GetBool("__force") {
			return fmt.Errorf("mail-message #%d has already processed at %s",
				msg.ID, time.Unix(int64(msg.ProcessedAt), 0).Format(time.UnixDate))
		}
	}

	defer func() {
		// mark as processed & save
		msg.ProcessedAt = uint64(time.Now().Unix())
		msg.Update()
	}()

	if target := opts.GetString("target"); target == "" {
		ins := Instance{
			Message:   msg,
			MessageID: msg.ID,
			Priority:  opts.GetInt("priority", 10),
		}
		msg.Frequency += 1
		if len(params) > 0 {
			ins.SetOptions(params)
			if params.GetBool("__immediate") {
				return SendInstance(&ins)
			}
		}
		ins.Save()
	} else {
		switch target {
		case "active-accounts":
			query := database.DB.Model(&users.Account{}).Where("state=?", 10)
			if count, err := sendToAccounts(msg, query); err != nil {
				return err
			} else {
				msg.Frequency += count
			}
		}
	}
	return nil
}

func SendInstance(ins *Instance) error {
	if ins.Message == nil {
		return fmt.Errorf("missing mail message")
	}

	defer ins.Save()

	// should use a copy of option
	opts := maps.Clone(ins.Message.Options.Data())
	body := ins.Message.Body

	if ins.Options != nil {
		opts2 := maps.Clone(ins.Options.Data())
		// merge opts2 to opts
		maps.Copy(opts, opts2)
		if tmplt, err := ins.Message.GetBodyTmplt(); err != nil {
			return err
		} else if tmplt != nil {
			var out bytes.Buffer
			opts2["Site"] = site.Settings()
			if err = tmplt.Execute(&out, opts2); err != nil {
				opts := ins.Options.Data()
				opts["__error"] = fmt.Sprintf("Execute body template error: %s", err.Error())
				return err
			}
			body = out.String()
		}
	}

	if system.IsPROD() {
		// only send real email in PROD only!
		if err := mailer.SendEx(&opts, ins.Message.Subject, body, false); err != nil {
			ins.FailedCount += 1
			return err
		}
	} else {
		system.Logger.Infof("[NONE-PROD] Email instance #%d was marked as sent", ins.ID)
	}
	ins.SentAt = uint64(time.Now().Unix())
	return nil
}

func sendToAccounts(msg *Message, query *gorm.DB) (int, error) {
	rows := []users.Account{}
	if err := query.Select("name", "email", "phone").Find(&rows).Error; err != nil {
		return 0, err
	}

	messages := []Instance{}
	priority := msg.Options.Data().GetInt("priority", 10)

	for _, row := range rows {
		opts := types.Params{
			"to":    row.Email,
			"name":  row.Name,
			"phone": row.Phone,
		}
		jsonOpts := datatypes.NewJSONType(opts)
		messages = append(messages, Instance{
			Message:   msg,
			MessageID: msg.ID,
			Options:   &jsonOpts,
			Priority:  priority,
		})
	}

	if err := db.CreateInBatches(&messages, 100).Error; err != nil {
		return 0, err
	}
	return len(rows), nil
}
