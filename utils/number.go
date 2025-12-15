package utils

import "strconv"

func IntWithDefault(val int, defaultVal int) int {
	if val != 0 {
		return val
	}
	return defaultVal
}

func Int64WithDefault(val int64, defaultVal int64) int64 {
	if val != 0 {
		return val
	}
	return defaultVal
}

func FloatWithDefault(val float64, defaultVal float64) float64 {
	if val != 0 {
		return val
	}
	return defaultVal
}

func ToInt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

func ToInt64(value any) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func ToFloat64(value any) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float64:
		return v
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0
}
