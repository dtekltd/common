package authApi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/jwt"
	"github.com/dtekltd/common/pkg/auth"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random" // Should be a random string for security
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  system.Env("GOOGLE_REDIRECT_URL"),
		ClientID:     system.Env("GOOGLE_CLIENT_ID"),
		ClientSecret: system.Env("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type GoogleToken struct {
	Token string `json:"token"`
}

func googleLogin(ctx *fiber.Ctx, keyManager *jwt.KeyManager) error {
	var token GoogleToken
	if err := ctx.BodyParser(&token); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	userInfo, err := getGoogleUserInfo(token.Token)
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	return loginWithGoogle(ctx, keyManager, userInfo)
}

func googleCallback(ctx *fiber.Ctx, keyManager *jwt.KeyManager) error {
	state := ctx.Query("state")
	if state != oauthStateString {
		return api.ErrorBadRequestResp(ctx, "Invalid state parameter")
	}

	code := ctx.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, "Failed to exchange token: "+err.Error())
	}

	userInfo, err := getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, "Failed to get user info: "+err.Error())
	}

	return loginWithGoogle(ctx, keyManager, userInfo)
}

func getGoogleUserInfo(token string) (types.Params, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	var userInfo types.Params
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed decoding user info: %s", err.Error())
	}

	return userInfo, nil
}

func loginWithGoogle(ctx *fiber.Ctx, keyManager *jwt.KeyManager, info types.Params) error {
	id := info.GetString("id")
	acc, err := auth.EnsureAccount("google", id, info)
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	if token, err := auth.GenerateToken(keyManager, &types.Params{"id": acc.SID()}); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		policy, _ := auth.GetUserPolicyEx(acc.SID())
		return api.SuccessResp(ctx, fiber.Map{
			"token":   token,
			"account": acc,
			"policy":  policy,
		})
	}
}
