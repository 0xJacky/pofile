package pofile

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// TranslateError represents a specific translation error
type TranslateError struct {
	Type    string
	Message string
	Err     error
}

func (e *TranslateError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Common error types
var (
	// ErrKeyNotFound is returned when a translation key is not found in the dictionary
	ErrKeyNotFound = errors.New("translation key not found")
	// ErrJSONMarshal is returned when JSON marshaling fails
	ErrJSONMarshal = errors.New("failed to marshal JSON data")
	// ErrJSONUnmarshal is returned when JSON unmarshaling fails
	ErrJSONUnmarshal = errors.New("failed to unmarshal JSON data")
	// ErrInvalidArgument is returned when an invalid argument type is provided
	ErrInvalidArgument = errors.New("invalid argument type")
)

// toCount converts various numeric types to int
func toCount(v any) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case int64:
		return int(val), nil
	default:
		return 0, &TranslateError{
			Type:    "invalid_count",
			Message: fmt.Sprintf("expected number as count but got %T", v),
			Err:     ErrInvalidArgument,
		}
	}
}

// getPluralForm returns the appropriate form based on count and available forms
func getPluralForm(forms []string, count int) string {
	if count == 1 && len(forms) > 0 {
		return forms[0] // singular
	}
	if len(forms) > 1 {
		return forms[1] // plural
	}
	if len(forms) > 0 {
		return forms[0] // fallback to singular if plural not available
	}
	return "" // should never happen if forms is validated before
}

// Translate returns the translated string for the given key and optional arguments.
// It may return an error if the key is not found or if there are issues with JSON processing.
func (d Dict) Translate(key string, args ...any) (string, error) {
	msgStr, ok := d[key]
	if !ok {
		return key, &TranslateError{
			Type:    "key_not_found",
			Message: fmt.Sprintf("key '%s' not found in dictionary", key),
			Err:     ErrKeyNotFound,
		}
	}

	// No arguments case - return the message string directly
	if len(args) == 0 {
		return toString(msgStr), nil
	}

	// Process placeholder replacement with single map argument
	if len(args) == 1 {
		// Handle string array case
		if strs, ok := args[0].([]string); ok {
			if len(strs) == 0 {
				return toString(msgStr), &TranslateError{
					Type:    "empty_string_array",
					Message: "string array is empty",
					Err:     ErrInvalidArgument,
				}
			}
			return strs[0], nil
		}

		// Replace placeholders with any type of data
		return replacePlaceholders(toString(msgStr), args[0])
	}

	// Handle plurality cases
	if len(args) == 2 {
		return handlePlurality(msgStr, args)
	}

	// Unhandled case
	return toString(msgStr), &TranslateError{
		Type:    "invalid_argument",
		Message: fmt.Sprintf("unexpected argument pattern with %d arguments", len(args)),
		Err:     ErrInvalidArgument,
	}
}

// replacePlaceholders replaces %{key} style placeholders in text with values from data
func replacePlaceholders(text string, data any) (string, error) {
	// Convert data to a map using json marshaling/unmarshaling
	var jsonMap map[string]any

	// Handle common types directly
	switch v := data.(type) {
	case map[string]any:
		jsonMap = v
	case map[string]string:
		// Convert map[string]string to map[string]any
		jsonMap = make(map[string]any, len(v))
		for k, val := range v {
			jsonMap[k] = val
		}
	default:
		// For structs and other complex types, use JSON marshaling/unmarshaling
		jsonData, err := json.Marshal(data)
		if err != nil {
			return text, &TranslateError{
				Type:    "json_marshal",
				Message: "failed to marshal JSON data",
				Err:     err,
			}
		}

		if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
			return text, &TranslateError{
				Type:    "json_unmarshal",
				Message: "failed to unmarshal JSON data",
				Err:     err,
			}
		}
	}

	re := regexp.MustCompile(`%\{([^}]+)\}`)
	result := re.ReplaceAllStringFunc(text, func(match string) string {
		key := match[2 : len(match)-1] // Extract key from %{key}
		if val, ok := jsonMap[key]; ok {
			return toString(val)
		}
		return match
	})
	return result, nil
}

// handlePlurality processes various plurality cases with 2 or 3 arguments
func handlePlurality(msgStr any, args []any) (string, error) {
	// Check if we're dealing with a string array in the dictionary
	strs, isStrArray := msgStr.([]string)

	// Case: nil/string, count
	if len(args) == 2 {
		if !isStrArray {
			return toString(msgStr), &TranslateError{
				Type:    "invalid_message_format",
				Message: "message is not a string array for plurality handling",
				Err:     ErrInvalidArgument,
			}
		}

		if len(strs) == 0 {
			return toString(msgStr), &TranslateError{
				Type:    "empty_string_array",
				Message: "message string array is empty",
				Err:     ErrInvalidArgument,
			}
		}

		count, err := toCount(args[1])
		if err != nil {
			return toString(msgStr), err
		}

		return getPluralForm(strs, count), nil
	}

	// This should never happen given our caller checks
	return toString(msgStr), &TranslateError{
		Type:    "invalid_argument",
		Message: "unexpected plurality pattern",
		Err:     ErrInvalidArgument,
	}
}

// Helper function to convert any to string
func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int64, float64, bool:
		return strings.TrimSpace(strings.Trim(fmt.Sprintf("%v", val), "\""))
	default:
		jsonData, _ := json.Marshal(val)
		return string(jsonData)
	}
}
