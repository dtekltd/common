package utils

func MapVal(m map[string]any, key string, defaultVal ...any) any {
	val, ok := m[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return nil
	}
	return val
}

func MapIntVal(m map[string]any, key string, defaultVal ...int) int {
	return MapVal(m, key, defaultVal).(int)
}

func MapStringVal(m map[string]any, key string, defaultVal ...string) string {
	return MapVal(m, key, defaultVal).(string)
}
