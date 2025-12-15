package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/auth"
	"github.com/dtekltd/common/jwt"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validate both Auth-token and API-key
// optional prefix ex. "/api"
func AuthMiddleware(keyManager *jwt.KeyManager, prefixes ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		domain := "route"
		route := string(ctx.Context().Path())
		uid, sid := getAccountID(ctx, keyManager)

		// bypass, special account 1-admin
		if uid != 1 {
			if len(prefixes) > 0 {
				if route2, ok := strings.CutPrefix(route, prefixes[0]); ok {
					route = route2
				}
			}

			// check privilege
			if ok, err := Enforcer.Enforce(sid, domain, route, ctx.Method()); err != nil {
				return api.ErrorInternalServerErrorResp(ctx, err.Error())
			} else {
				system.Logger.Info("enforce:", ok, sid, domain, route, ctx.Method())
				if !ok {
					if uid == 0 {
						return api.ErrorUnauthorizedResp(ctx, "Unauthorized")
					} else {
						// check for default 'user' role!
						if ok, err := Enforcer.Enforce("user", domain, route, ctx.Method()); err != nil || !ok {
							return api.ErrorUnauthorizedResp(ctx, "Unauthorized")
						}
					}
				}
			}
		}

		ctx.Locals("uiID", uid)
		ctx.Locals("usID", sid)

		return ctx.Next()
	}
}

func HasRoleMiddleware(roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		usID := ctx.Locals("usID").(string)
		if HasRole(usID, roles...) {
			return ctx.Next()
		}
		return api.ErrorUnauthorizedResp(ctx, "Unauthorized")
	}
}

func IsAdminMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// check casbin role
		usID := ctx.Locals("usID").(string)
		if ok, _ := Enforcer.HasRoleForUser(usID, "admin"); ok {
			system.Logger.Info("admin:", ok, usID)
			return ctx.Next()
		}
		if acc, _ := GetAccount(ctx); acc != nil && acc.IsAdmin {
			system.Logger.Info("admin:", "is_admin", usID)
			return ctx.Next()
		}
		return api.ErrorUnauthorizedResp(ctx, fmt.Sprintf("Unauthorized - account #%d not found!", ctx.Locals("uiID").(uint64)))
	}
}

func getAccountID(ctx *fiber.Ctx, keyManager *jwt.KeyManager) (uint64, string) {
	if token := auth.ExtractToken(ctx, "header:Authorization,query:auth_token"); token != "" {
		data, err := ParseToken(keyManager, token)
		system.Logger.Info("token:", token[0:16], data)
		if err == nil {
			sid := (*data)["id"].(string)
			return uint64(utils.StringToInt(sid)), sid
		}
	}
	return 0, "0"
}

func IsGuest(ctx *fiber.Ctx) bool {
	uiID := ctx.Locals("uiID")
	return uiID == nil || uiID.(uint64) == 0
}

func GetAccount(ctx *fiber.Ctx) (*users.Account, error) {
	if acc := ctx.Locals("account"); acc != nil {
		return acc.(*users.Account), nil
	}
	if sid := ctx.Locals("usID"); sid != nil {
		if acc, err := users.FindAccount(sid); err != nil {
			return nil, err
		} else {
			ctx.Locals("account", acc)
			return acc, nil
		}
	}
	return nil, errors.New("auth token not found")
}
