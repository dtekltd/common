package types

import (
	"strings"

	"github.com/dtekltd/common/utils"
)

type Params map[string]any

// type Params interface {
// 	Get(key string) any
// 	Set(key string, val any)
// }

func (p Params) IsEmpty() bool {
	return len(p) == 0
}

func (p Params) Copy() (Params, error) {
	dest := Params{}
	if err := utils.Copy(&p, &dest); err != nil {
		return nil, err
	}
	return dest, nil
}

func (p Params) Set(key string, val any) {
	params := p
	keys := strings.Split(key, ".")
	l := len(keys) - 1
	for i, k := range keys {
		if existParams, ok := params[k]; !ok {
			// key not exists
			if i == l {
				params[k] = val
			} else {
				newParams := Params{}
				params[k] = newParams
				params = newParams
			}
		} else {
			// key exists
			if i == l {
				params[k] = val
			} else {
				params = Params(existParams.(map[string]any))
			}
		}
	}
}

func (p Params) Get(key string, defaultVals ...any) any {
	var val any
	var ok bool
	params := p
	keys := strings.Split(key, ".")
	l := len(keys) - 1
	for i, k := range keys {
		if val, ok = params[k]; ok {
			if i == l {
				return val
			} else {
				params = Params(val.(map[string]any))
			}
		}
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	return nil
}

func (p Params) GetParams(key string, defaultVals ...Params) Params {
	if val := p.Get(key); val != nil {
		if params, ok := val.(map[string]any); ok {
			return Params(params)
		}
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	// empty Params
	return Params{}
}

func (p Params) GetSliceParams(key string) []Params {
	var slicesParams []Params
	if val := p.Get(key); val != nil {
		if data, ok := val.([]any); ok {
			slicesParams = make([]Params, 0, len(data))
			for _, val := range data {
				if params, ok := val.(map[string]any); ok {
					slicesParams = append(slicesParams, Params(params))
				}
			}
		}
	}
	return slicesParams
}

func (p Params) GetString(key string, defaultVals ...string) string {
	if val := p.Get(key); val != nil {
		return val.(string)
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	return ""
}

func (p Params) GetBool(key string, defaultVals ...bool) bool {
	if val := p.Get(key); val != nil {
		return val.(bool)
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	return false
}

func (p Params) GetFloat64(key string, defaultVals ...float64) float64 {
	if val := p.Get(key); val != nil {
		return utils.ToFloat64(val)
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	return 0
}

func (p Params) GetInt(key string, defaultVals ...int) int {
	if val := p.Get(key); val != nil {
		return utils.ToInt(val)
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	return 0
}

func (p Params) GetInt64(key string, defaultVals ...int64) int64 {
	if val := p.Get(key); val != nil {
		return utils.ToInt64(val)
	}
	if len(defaultVals) > 0 {
		return defaultVals[0]
	}
	return 0
}

func (p Params) GetUint(key string) uint {
	if val := p.Get(key); val != nil {
		return uint(utils.ToInt(val))
	}
	return 0
}

func (p Params) GetUint64(key string) uint64 {
	if val := p.Get(key); val != nil {
		return uint64(utils.ToInt64(val))
	}
	return 0
}
