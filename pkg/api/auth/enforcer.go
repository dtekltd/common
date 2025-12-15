package authApi

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type enforcerHandler struct {
	Enforcer *casbin.Enforcer
}

func NewEnforcerHandler() *enforcerHandler {
	return &enforcerHandler{
		Enforcer: auth.Enforcer,
	}
}

func (h *enforcerHandler) reloadPolicy(ctx *fiber.Ctx) error {
	if err := h.Enforcer.LoadPolicy(); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}
	return api.SuccessResp(ctx, true)
}

// enforce ["Peter", "settings", "site", "view"]
func (h *enforcerHandler) enforce(ctx *fiber.Ctx) error {
	p := []interface{}{}
	if err := ctx.BodyParser(&p); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		if ok, err := h.Enforcer.Enforce(p...); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		} else {
			return api.SuccessResp(ctx, ok)
		}
	}
}

// getPolicy get policies or grouping policies
func (h *enforcerHandler) getPolicy(ctx *fiber.Ctx) error {
	t := ctx.Query("type")
	index := ctx.QueryInt("index")
	values := strings.Split(ctx.Query("values"), ",")
	var policies [][]string
	var err error
	if t == "p" {
		policies, err = h.Enforcer.GetFilteredPolicy(index, values...)
	} else {
		policies, err = h.Enforcer.GetFilteredGroupingPolicy(index, values...)
	}
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, policies)
	}
}

// addPolicy add policies or grouping policies
func (h *enforcerHandler) addPolicy(ctx *fiber.Ctx) error {
	type Req struct {
		Type     string     `json:"type"`
		Policies [][]string `json:"policies"`
	}
	rq := &Req{}
	if err := ctx.BodyParser(&rq); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		var ok bool
		var err error
		if rq.Type == "g" {
			ok, err = h.Enforcer.AddGroupingPolicies(rq.Policies)
		} else {
			ok, err = h.Enforcer.AddPolicies(rq.Policies)
		}
		if err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		} else {
			return api.SuccessResp(ctx, ok)
		}
	}
}

// deletePolicy delete policies or grouping policies
func (h *enforcerHandler) deletePolicy(ctx *fiber.Ctx) error {
	type Req struct {
		Type   string   `json:"type"`
		Index  int      `json:"index"`
		Fields []string `json:"fields"`
	}
	rq := &Req{}
	if err := ctx.BodyParser(&rq); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		var ok bool
		var err error
		if rq.Type == "g" {
			ok, err = h.Enforcer.RemoveFilteredGroupingPolicy(rq.Index, rq.Fields...)
		} else {
			ok, err = h.Enforcer.RemoveFilteredPolicy(rq.Index, rq.Fields...)
		}
		if err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		} else {
			return api.SuccessResp(ctx, ok)
		}
	}
}

func (h *enforcerHandler) getAllPolicy(ctx *fiber.Ctx) error {
	if policy, err := h.Enforcer.GetPolicy(); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, policy)
	}
}

func (h *enforcerHandler) deleteUserPolicy(ctx *fiber.Ctx) error {
	if user := ctx.Query("user"); user != "" {
		if ok, err := h.Enforcer.DeleteUser(user); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		} else {
			return api.SuccessResp(ctx, ok)
		}
	}
	return api.SuccessResp(ctx, false)
}

func (h *enforcerHandler) getUserPolicies(ctx *fiber.Ctx) error {
	sid := ctx.Query("sub", ctx.Locals("usID").(string))
	asText := ctx.Query("text") != ""
	if asText {
		if text, err := auth.GetUserPolicyText(sid); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		} else {
			return api.SuccessResp(ctx, text)
		}
	} else {
		if all, err := auth.GetUserPolicyEx(sid); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		} else {
			return api.SuccessResp(ctx, all)
		}
	}
}

func (h *enforcerHandler) getUserPolicyEx(ctx *fiber.Ctx) error {
	usID := ctx.Locals("usID").(string)
	policies, _ := auth.GetUserPolicyEx(usID)
	return api.SuccessResp(ctx, policies)
}
