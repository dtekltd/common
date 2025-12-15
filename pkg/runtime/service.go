package runtime

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dtekltd/common/json"
	"gorm.io/gorm"
)

// Cleanup deletes all records where UpdatedAt + ExpiredTime < current time
func Cleanup() (int64, error) {
	currentTime := time.Now().Unix()
	result := db.Where("updated_at + expired_time < ?", currentTime).Delete(&Runtime{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func Get(key string, out any) error {
	return get(key, out)
}

func GetString(key string, defaults ...string) string {
	defaultVal := ""
	if len(defaults) > 0 {
		defaultVal = defaults[0]
	}
	var out any
	if err := get(key, &out); err != nil {
		return defaultVal
	}
	if str, ok := out.(string); ok {
		return str
	}
	return defaultVal
}

func GetInt(key string, defaults ...int) int {
	defaultVal := 0
	if len(defaults) > 0 {
		defaultVal = defaults[0]
	}
	var out any
	if err := get(key, &out); err != nil {
		return defaultVal
	}
	switch v := out.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func GetFloat64(key string, defaults ...float64) float64 {
	defaulVal := 0.0
	if len(defaults) > 0 {
		defaulVal = defaults[0]
	}
	var out any
	if err := get(key, &out); err != nil {
		return defaulVal
	}
	switch v := out.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaulVal
}

func GetBool(key string, defauls ...bool) bool {
	defaultVal := false
	if len(defauls) > 0 {
		defaultVal = defauls[0]
	}
	var out any
	if err := get(key, &out); err != nil {
		return defaultVal
	}
	if b, ok := out.(bool); ok {
		return b
	}
	return defaultVal
}

func Set(key string, val any, exp int) error {
	runtime, err := record(key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if runtime == nil {
		runtime = &Runtime{
			Key: key,
		}
	}
	str, err := json.ToJSON(val)
	if err != nil {
		return err
	}

	runtime.Data = str
	runtime.ExpiredTime = exp
	if err := db.Save(runtime).Error; err != nil {
		return err
	}
	return nil
}

func get(key string, out any) error {
	runtime, err := record(key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if runtime != nil {
		if runtime.ExpiredTime == 0 || (time.Now().Unix()-runtime.UpdatedAt) < int64(runtime.ExpiredTime) {
			if err := json.FromJSON(runtime.Data, out); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("key %s expired", key)
		}
	}
	return nil
}

func record(key string) (*Runtime, error) {
	var runtime Runtime
	if key == "" {
		return nil, fmt.Errorf("runtime key cannot be empty")
	}
	if err := db.Where("`key`=?", key).First(&runtime).Error; err != nil {
		return nil, err
	}
	return &runtime, nil
}
