package mailer

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/types"
	"github.com/dtekltd/common/utils"
	"gopkg.in/gomail.v2"
)

var (
	dialer        *gomail.Dialer
	Supporter     string
	DefaultSender string
)

func init() {
	Supporter = system.Env("MAILER_SUPPORTER", "Tuan Phan <peter.phan07@gmail.com>")
	DefaultSender = system.Env("MAILER_SENDER", "No-Reply <no.reply@gmail.com>")
	port, _ := strconv.Atoi(system.Env("SMTP_PORT"))
	dialer = gomail.NewDialer(system.Env("SMTP_HOST"), port, system.Env("SMTP_USER"), system.Env("SMTP_PASS"))
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
}

func Send(to string, subject string, message string, isText ...bool) error {
	opts := &types.Params{
		"to": to,
	}
	return SendEx(opts, subject, message, isText...)
}

func SendEx(opts *types.Params, subject string, message string, isText ...bool) error {
	m := gomail.NewMessage()
	m.SetHeader("subject", subject)

	from := DefaultSender
	if val := opts.GetString("from"); val != "" {
		from = val
	}
	m.SetHeader("From", from)

	if val, _ := getOptVal("to", opts); val != nil {
		m.SetHeader("To", (*val)...)
	}
	if val, _ := getOptVal("cc", opts); val != nil {
		m.SetHeader("Cc", (*val)...)
	}
	if val, _ := getOptVal("bcc", opts); val != nil {
		m.SetHeader("Bcc", (*val)...)
	}

	if len(isText) > 0 && isText[0] {
		m.SetBody("text/html", message)
		m.SetBody("text/html", strings.ReplaceAll(message, "\n", "<br>"))
	} else {
		m.SetBody("text/html", utils.StripHtmlTags(message))
		m.SetBody("text/html", message)
	}

	// send the email
	return dialer.DialAndSend(m)
}

func NotifySupport(subject string, message any) error {
	htmlMesasge := "<pre>\n" + utils.ToString(message) + "</pre>"
	return Send(Supporter, subject, htmlMesasge)
}

func getOptVal(key string, opts *types.Params) (*[]string, error) {
	if val := opts.Get(key); val != nil {
		var value []string
		switch v := val.(type) {
		case string:
			if str := val.(string); str != "" {
				pattern := `[,;\s]+`
				re := regexp.MustCompile(pattern)
				value = re.Split(val.(string), -1)
			}
		case []string:
			value = val.([]string)
		default:
			return nil, fmt.Errorf("unknown type: %T", v)
		}
		return &value, nil
	}
	return nil, nil
}
