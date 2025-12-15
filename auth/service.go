package auth

import (
	"errors"

	"github.com/dtekltd/common/http"
	"github.com/dtekltd/common/system"
)

func RefreshToken(token string) (string, error) {
	resp := &AuthRefreshResponse{}
	err := http.Post(system.Env("AUTH_API")+"/auth/refresh", nil, resp, map[string]string{
		"Authorization": "Bearer " + token,
	})
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", errors.New(resp.Error.Message)
	}

	return resp.Data, nil
}
