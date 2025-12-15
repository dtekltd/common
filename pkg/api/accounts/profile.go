package accountsApi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/utils"
	"github.com/gofiber/fiber/v2"
)

type UpdateAccountReq struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Username  string `json:"username,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
	Password  string `json:"password,omitempty"`
}

func saveAccount(ctx *fiber.Ctx) error {
	pId, _ := ctx.ParamsInt("id")
	id := uint64(pId)
	uiID := ctx.Locals("uiID").(uint64)

	if id == 0 {
		id = uiID
	}

	req := &UpdateAccountReq{}
	acc, _ := users.FindAccount(id)
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	if req.Password != "" {
		req.Password = utils.HashPassword(req.Password)
	} else {
		req.Password = acc.Password
	}

	if err := database.DB.Model(&users.Account{}).
		Where("id=?", id).
		Updates(req).Error; err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	// reload account
	acc, _ = users.FindAccount(id)
	return api.SuccessResp(ctx, acc)
}

var (
	uploadUrl = "/uploads/avatars/"
	uploadDir = "./uploads/avatars"
)

func uploadAvatar(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	sID := ctx.Locals("usID")
	acc, _ := users.FindAccount(sID)
	oldAvatarUrl := acc.AvatarUrl

	// fileName := fmt.Sprintf("%s-%s", sID, strings.Replace(file.Filename, " ", "_", -1))
	fileName := fmt.Sprintf("%s-%s", sID, utils.RemoveSignChars(file.Filename))
	filePath := filepath.Join(uploadDir, fileName)

	if err := ctx.SaveFile(file, filePath); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	avatarUrl := uploadUrl + fileName
	database.DB.Model(acc).Update("avatar_url", avatarUrl)

	// delete old avatar
	if avatarUrl != oldAvatarUrl {
		if err := os.Remove("." + oldAvatarUrl); err != nil {
			system.Logger.Errorf("Fail to delete old avatar: %v", err)
		}
	}

	return api.SuccessResp(ctx, avatarUrl)
}
