package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

// QueryStruct parses a specific group of query parameters into a struct
// Example: api-end-point?filter[name]=peter&filter[sex]=male
// param specifies the query parameter group (e.g., "filter")
// out should be a pointer to a struct
func QueryStruct(ctx *fiber.Ctx, out any, param string) error {
	// Verify that out is a pointer to a struct
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("out must be a pointer to a struct")
	}

	// Get all query parameters
	query := ctx.Queries()

	// Filter and convert only the parameters for the specified group
	data := make(map[string]any)
	if param == "" {
		for key, value := range query {
			keys := strings.Split(key, "[")
			current := data

			for i, k := range keys {
				k = strings.TrimSuffix(k, "]")
				if i == len(keys)-1 {
					current[k] = value
				} else {
					if _, exists := current[k]; !exists {
						current[k] = make(map[string]any)
					}
					current = current[k].(map[string]any)
				}
			}
		}
	} else if err := parseParam(param, query, data); err != nil {
		return err
	}

	// Use mapstructure to decode the map into the struct
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   out,
		TagName:  "query",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToBoolHookFunc(),
			stringToIntHookFunc(),
			stringToUint64HookFunc(),
			stringToFloat64HookFunc(),
			mapstructure.StringToTimeHookFunc("2006-01-02"),
			mapstructure.StringToSliceHookFunc(","),
		),
		WeaklyTypedInput: true, // Allows more flexible type conversion
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode query params: %w", err)
	}

	return nil
}

func parseParam(param string, query map[string]string, out map[string]any) error {
	prefix := param + "["

	for key, value := range query {
		if strings.HasPrefix(key, prefix) {
			// Extract the nested key (e.g., "name" from "filter[name]")
			nestedKey := strings.TrimSuffix(strings.TrimPrefix(key, prefix), "]")

			// Handle multiple levels of nesting (e.g., "filter[user][name]")
			keys := strings.Split(nestedKey, "][")
			current := out

			for i, k := range keys {
				if i == len(keys)-1 {
					current[k] = value
				} else {
					if _, exists := current[k]; !exists {
						current[k] = make(map[string]any)
					}
					current = current[k].(map[string]any)
				}
			}
		} else if key == param {
			// Handle non-nested case where param is the direct key
			// e.g., "filter=somevalue" would set the whole struct to "somevalue"
			// This might not be what you want, so you could return an error here
			return fmt.Errorf("flat parameter '%s' found, expected nested parameters like '%s[name]'", param, param)
		}
	}

	return nil
}

// stringToBoolHookFunc converts string to bool
func stringToBoolHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() == reflect.Bool {
			str := strings.ToLower(data.(string))
			switch str {
			case "true", "1", "yes", "on":
				return true, nil
			case "false", "0", "no", "off", "":
				return false, nil
			default:
				return nil, fmt.Errorf("invalid boolean value: %s", str)
			}
		}

		return data, nil
	}
}

// stringToFloat64HookFunc converts string to float64
func stringToFloat64HookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64 {
			str := data.(string)
			if str == "" {
				return 0.0, nil
			}
			return strconv.ParseFloat(str, 64)
		}

		return data, nil
	}
}

// stringToUint64HookFunc converts string to uint64
func stringToUint64HookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() == reflect.Uint || t.Kind() == reflect.Uint8 ||
			t.Kind() == reflect.Uint16 || t.Kind() == reflect.Uint32 ||
			t.Kind() == reflect.Uint64 {
			str := data.(string)
			if str == "" {
				return uint64(0), nil
			}
			return strconv.ParseUint(str, 10, 64)
		}

		return data, nil
	}
}

// stringToIntHookFunc converts string to int
func stringToIntHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() == reflect.Int || t.Kind() == reflect.Int8 ||
			t.Kind() == reflect.Int16 || t.Kind() == reflect.Int32 ||
			t.Kind() == reflect.Int64 {
			str := data.(string)
			if str == "" {
				return 0, nil
			}
			return strconv.Atoi(str)
		}

		return data, nil
	}
}
