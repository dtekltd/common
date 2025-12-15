package authApi

import (
	"fmt"
	"time"

	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/jwt"
	"github.com/dtekltd/common/pkg/auth"
	"github.com/dtekltd/common/pkg/google"
	"github.com/dtekltd/common/pkg/mail"
	"github.com/dtekltd/common/pkg/runtime"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/types"
	"github.com/dtekltd/common/utils"
	"github.com/gofiber/fiber/v2"
)

func account(ctx *fiber.Ctx) error {
	if acc, err := users.FindAccount(ctx.Locals("usID")); err != nil {
		return api.ErrorNotFoundResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, acc)
	}
}

func login(ctx *fiber.Ctx, keyManager *jwt.KeyManager) error {
	req := &auth.LoginReq{}
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, "Parse failed! "+err.Error())
	}
	if req.Email == "" || req.Password == "" {
		return api.ErrorBadRequestResp(ctx, "Both \"email\" and \"password\" are required.")
	}

	if acc, err := auth.Login(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	} else {
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
}

func register(ctx *fiber.Ctx, keyManager *jwt.KeyManager) error {
	req := &auth.RegisterReq{}
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if req.Email == "" || req.Password == "" {
		return api.ErrorBadRequestResp(ctx, "Both \"email\" and \"password\" are required.")
	}
	if acc, err := auth.Register(req); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, acc)
	}
}

func resetPassword(ctx *fiber.Ctx) error {
	type ResetPasswordReq struct {
		OTP            string `json:"otp"`
		Key            string `json:"key"`
		Email          string `json:"email"`
		Password       string `json:"password"`
		RecaptchaToken string `json:"recaptchaToken"`
	}

	req := &ResetPasswordReq{}
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if req.Email == "" {
		return api.ErrorBadRequestResp(ctx, "\"email\" is required.")
	}
	if req.OTP == "" && req.RecaptchaToken == "" {
		return api.ErrorBadRequestResp(ctx, "Either OTP or Recaptcha token is required.")
	}
	if req.OTP != "" && req.Password == "" {
		return api.ErrorBadRequestResp(ctx, "New \"password\" is required.")
	}

	if req.RecaptchaToken != "" {
		res, err := google.VerifyRecaptcha(req.RecaptchaToken)
		if err != nil {
			return api.ErrorBadRequestResp(ctx, err.Error())
		}
		if !res.Success {
			return api.ErrorBadRequestResp(ctx, "Recaptcha verify fail.")
		}
	}

	acc := &users.Account{}
	if err := database.DB.Take(acc, "email=?", req.Email).Error; err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	if req.OTP != "" {
		otp := runtime.GetString(req.Key)
		if otp != req.OTP {
			return api.ErrorBadRequestResp(ctx, "OTP is incorrect!")
		}

		hassPass := utils.HashPassword(req.Password)
		database.DB.Model(acc).Update("password", hassPass)

		// send mail in a new routine
		go func() {
			msg := mail.GetMessage("account-password-was-reset")
			if msg == nil {
				system.Logger.Error("Mail template `account-password-was-reset` is missing")
				return
			}

			params := types.Params{
				"to":          acc.Email,
				"Name":        acc.Name,
				"Email":       acc.Email,
				"__immediate": true,
			}

			if err := mail.SendMessage(msg, params); err != nil {
				system.Logger.Error("Send mail error:", err.Error())
			}
		}()

		return api.SuccessResp(ctx, true)
	} else {
		// reset account with temp password
		msg := mail.GetMessage("account-password-reset-temp")
		if msg == nil {
			system.Logger.Error("Mail template `account-password-reset-temp` is missing")
			return api.ErrorInternalServerErrorResp(ctx, "Mail template `account-password-reset-temp` is missing")
		}

		if tmpPass, err := utils.GeneratePassword(8); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		} else {
			// send mail in a new routine
			go func() {
				params := types.Params{
					"to":          acc.Email,
					"Name":        acc.Name,
					"Email":       acc.Email,
					"Password":    tmpPass,
					"__immediate": true,
				}

				if err := mail.SendMessage(msg, params); err != nil {
					system.Logger.Error("Send mail error:", err.Error())
				}
			}()

			hassPass := utils.HashPassword(tmpPass)
			database.DB.Model(acc).Update("password", hassPass)
			return api.SuccessResp(ctx, true)
		}
	}
}

func resetPasswordRequest(ctx *fiber.Ctx) error {
	type ResetPasswordReq struct {
		Type           string `json:"type"`
		Email          string `json:"email"`
		RecaptchaToken string `json:"recaptchaToken"`
	}

	req := &ResetPasswordReq{}
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if req.Email == "" {
		return api.ErrorBadRequestResp(ctx, "\"email\" is required.")
	}

	if req.RecaptchaToken != "" {
		res, err := google.VerifyRecaptcha(req.RecaptchaToken)
		if err != nil {
			return api.ErrorBadRequestResp(ctx, err.Error())
		}
		if !res.Success {
			return api.ErrorBadRequestResp(ctx, "Recaptcha verify fail.")
		}
	}

	acc := &users.Account{}
	if err := database.DB.Take(acc, "email=?", req.Email).Error; err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	if req.Type == "otp" {
		msg := mail.GetMessage("account-password-reset-request-otp")
		if msg == nil {
			system.Logger.Error("Mail template `account-password-reset-request-otp` is missing")
			return api.ErrorInternalServerErrorResp(ctx, "Mail template `account-password-reset-otp` is missing")
		}
		otp, err := utils.GenerateOTP(6)
		if err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}

		params := types.Params{
			"to":          acc.Email,
			"Name":        acc.Name,
			"Email":       acc.Email,
			"OTP":         otp,
			"__immediate": true,
		}

		// send mail in a new routine
		go func() {
			if err := mail.SendMessage(msg, params); err != nil {
				system.Logger.Error("Send mail error:", err.Error())
			}
		}()

		// 15 minutes expiration
		now := time.Now().Unix()
		key := fmt.Sprintf("opt-reset-password-%d-%d", acc.ID, now)
		expiredTime := 15 * 60
		if err := runtime.Set(key, otp, expiredTime); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}

		return api.SuccessResp(ctx, &fiber.Map{
			"key":    key,
			"otp":    otp,
			"expire": now + int64(expiredTime),
		})
	} else {
		// TODO: need to implement sending password reset link
		// account-password-reset-link

		// msg := mail.GetMessage("account-password-reset-link")
		// if msg == nil {
		// 	return api.ErrorInternalServerErrorResp(ctx, "Mail template `account-password-reset-link` is missing")
		// }

		// if tmpPass, err := utils.GeneratePassword(8); err != nil {
		// 	return api.ErrorInternalServerErrorResp(ctx, err.Error())
		// } else {
		// 	params := types.Params{
		// 		"to":          acc.Email,
		// 		"Name":        acc.Name,
		// 		"Email":       acc.Email,
		// 		"Password":    tmpPass,
		// 		"__immediate": true,
		// 	}

		// 	// send mail in a new routine
		// 	go func() {
		// 		if err := mail.SendMessage(msg, params); err != nil {
		// 			system.Logger.Error("Send mail error:", err.Error())
		// 		}
		// 	}()

		// 	hassPass := utils.HashPassword(tmpPass)
		// 	database.DB.Model(acc).Update("password", hassPass)
		// 	return api.SuccessResp(ctx, true)
		// }

		return api.ErrorNotFoundResp(ctx, "Feature is not supported yet!")
	}
}
