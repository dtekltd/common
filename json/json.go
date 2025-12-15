package json

import (
	ejson "encoding/json"
	"fmt"
	"log"
)

// ToJSON converts any value to JSON string
func ToJSON(data any) (string, error) {
	if data == nil {
		return "null", nil
	}

	jsonBytes, err := ejson.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// ToJSONIndent converts any value to formatted JSON string
func ToJSONIndent(data any, prefix, indent string) (string, error) {
	if data == nil {
		return "null", nil
	}

	jsonBytes, err := ejson.MarshalIndent(data, prefix, indent)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// ToJSONBytes converts any value to JSON bytes
func ToJSONBytes(data any) ([]byte, error) {
	if data == nil {
		return []byte("null"), nil
	}

	jsonBytes, err := ejson.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return jsonBytes, nil
}

// FromJSON converts JSON string back to any type
func FromJSON(jsonStr string, out any) error {
	if jsonStr != "null" && jsonStr != "" {
		err := ejson.Unmarshal([]byte(jsonStr), &out)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}
	return nil
}

// FromJSONBytes converts JSON bytes back to any type
func FromJSONBytes(jsonBytes []byte, out any) error {
	if len(jsonBytes) != 0 && string(jsonBytes) != "null" {
		err := ejson.Unmarshal(jsonBytes, &out)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}
	return nil
}

// FromJSONToType converts JSON to a specific type
func FromJSONToType[T any](jsonStr string) (T, error) {
	var result T
	if jsonStr == "null" || jsonStr == "" {
		return result, nil
	}

	err := ejson.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return result, nil
}

// FromJSONBytesToType converts JSON bytes to a specific type
func FromJSONBytesToType[T any](jsonBytes []byte) (T, error) {
	var result T
	if len(jsonBytes) == 0 || string(jsonBytes) == "null" {
		return result, nil
	}

	err := ejson.Unmarshal(jsonBytes, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return result, nil
}

// SafeToJSON converts any to JSON with error logging instead of returning error
func SafeToJSON(data any) string {
	result, err := ToJSON(data)
	if err != nil {
		log.Printf("SafeToJSON error: %v", err)
		return "null"
	}
	return result
}

// SafeFromJSON converts JSON to any with error logging
func SafeFromJSON(jsonStr string) any {
	var result any
	if err := FromJSON(jsonStr, result); err != nil {
		log.Printf("SafeFromJSON error: %v", err)
		return nil
	}
	return result
}

// IsValidJSON checks if a string is valid JSON
func IsValidJSON(jsonStr string) bool {
	var val any
	return ejson.Unmarshal([]byte(jsonStr), &val) == nil
}

// PrettyPrint prints formatted JSON to stdout
func PrettyPrint(data any) error {
	jsonStr, err := ToJSONIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(jsonStr)
	return nil
}

// ConvertType converts any type to another type using JSON marshaling/unmarshaling
func ConvertType(source, destination any) error {
	jsonBytes, err := ejson.Marshal(source)
	if err != nil {
		return err
	}
	return ejson.Unmarshal(jsonBytes, destination)
}
