package auth

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/system"
)

var Enforcer *casbin.Enforcer

func init() {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(database.DB, "user", "casbin_rules")
	if err != nil {
		system.Logger.Panic("casbin create adapter", "error", err.Error())
	}
	// Load the model and use the GORM adapter
	Enforcer, err = casbin.NewEnforcer("./casbin-model.conf", adapter)
	if err != nil {
		system.Logger.Panic("casbin create enforcer", "error", err.Error())
	}

	Enforcer.LoadPolicy()
}

func ReloadPolicy() error {
	return Enforcer.LoadPolicy()
}

func GetUserPolicy(usID string) (*[][]string, error) {
	all := [][]string{
		// default! user role
		{"g", usID, "user", "*"},
	}
	// get grouping policies
	if policies, err := Enforcer.GetFilteredGroupingPolicy(0, usID); err != nil {
		return nil, err
	} else {
		if len(policies) > 0 {
			for _, p := range policies {
				g := append([]string{"g"}, p...)
				all = append(all, g)
				// get policies of group
				if policies, err := Enforcer.GetFilteredPolicy(0, p[1]); err != nil {
					return nil, err
				} else {
					for _, p := range policies {
						r := append([]string{"p"}, p...)
						all = append(all, r)
					}
				}
			}
		} else {
			// get policies
			if policies, err := Enforcer.GetFilteredPolicy(0, usID); err != nil {
				return nil, err
			} else {
				// all = policies
				for _, p := range policies {
					r := append([]string{"p"}, p...)
					all = append(all, r)
				}
			}
		}
		if policies, err := Enforcer.GetFilteredPolicy(0, "user"); err == nil {
			// default! user policies
			for _, p := range policies {
				r := append([]string{"p"}, p...)
				all = append(all, r)
			}
		}
	}
	return &all, nil
}

func GetUserPolicyText(usID string) (string, error) {
	text := ""
	if all, err := GetUserPolicy(usID); err != nil {
		return "", err
	} else {
		for _, p := range *all {
			text += fmt.Sprintf("%s\n", strings.Join(p, ", "))
		}
		return text, nil
	}
}

// GetUserPolicyEx get only [dom, obj, act] without sub
// if admin, only return 1 row [*, *, *]
func GetUserPolicyEx(usID string) (*map[string]any, error) {
	roles := []string{}
	routes := [][]string{}

	if ok, _ := Enforcer.HasRoleForUser(usID, "admin"); ok {
		// admin, do not get more...
		roles = append(roles, "admin")
		routes = append(routes, []string{"*", "*"})
	} else {
		// get grouping policies
		if policies, err := Enforcer.GetFilteredGroupingPolicy(0, usID); err != nil {
			return nil, err
		} else {
			if len(policies) > 0 {
				for _, p := range policies {
					roles = append(roles, p[1])
					// get policies of group
					if policies, err := Enforcer.GetFilteredPolicy(0, p[1]); err != nil {
						return nil, err
					} else {
						for _, p := range policies {
							routes = append(routes, p[2:])
						}
					}
				}
			} else {
				// get policies
				if policies, err := Enforcer.GetFilteredPolicy(0, usID); err != nil {
					return nil, err
				} else {
					// routes = policies
					for _, p := range policies {
						routes = append(routes, p[2:])
					}
				}
			}
			if policies, err := Enforcer.GetFilteredPolicy(0, "user"); err == nil {
				// default! user policies
				for _, p := range policies {
					routes = append(routes, p[2:])
				}
			}
		}
	}
	return &map[string]any{
		"roles":  roles,
		"routes": routes,
	}, nil
}

func HasRole(usID string, roles ...string) bool {
	for _, role := range roles {
		if ok, _ := Enforcer.HasRoleForUser(usID, role); ok {
			return true
		}
	}
	return false
}
