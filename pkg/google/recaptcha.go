package google

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dtekltd/common/pkg/site"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func VerifyRecaptcha(token string) (*RecaptchaResponse, error) {
	params := site.Settings().CustomParams
	secretKey := params.GetString("recaptcha.secret")
	res := &RecaptchaResponse{}
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		map[string][]string{
			"secret":   {secretKey},
			"response": {token},
		})
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return res, err
	}

	return res, nil
}
