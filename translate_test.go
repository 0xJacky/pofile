package pofile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateBasic(t *testing.T) {
	// Create a test dictionary
	dict := Dict{
		"hello":   "你好",
		"welcome": "欢迎，%{name}",
		"apple":   []string{"苹果", "苹果们"},
	}

	// Test translation without arguments
	t.Run("Simple translation without args", func(t *testing.T) {
		result, err := dict.Translate("hello")
		assert.NoError(t, err)
		assert.Equal(t, "你好", result, "Basic translation should match")
	})

	// Test JSON placeholder replacement
	t.Run("JSON replacement placeholder", func(t *testing.T) {
		result, err := dict.Translate("welcome", map[string]string{"name": "张三"})
		assert.NoError(t, err)
		assert.Equal(t, "欢迎，张三", result, "Welcome message with placeholder should be replaced")
	})

	// Test string array for plural forms
	t.Run("Plurality handling", func(t *testing.T) {
		// Test singular form
		result, err := dict.Translate("apple", nil, 1)
		assert.NoError(t, err)
		assert.Equal(t, "苹果", result, "Should get singular form")

		// Test plural form
		result, err = dict.Translate("apple", nil, 2)
		assert.NoError(t, err)
		assert.Equal(t, "苹果们", result, "Should get plural form")
	})
}

func TestTranslateComplex(t *testing.T) {
	// Create test dictionary with more complex patterns
	dict := Dict{
		"userInfo": "用户 %{name} 在 %{time} 共享了 %{count} 个文件",
		"complex":  "这有 %{num} 个选项：%{options}",
	}

	// Test more complex string formatting
	t.Run("Complex message format with multiple placeholders", func(t *testing.T) {
		result, err := dict.Translate("userInfo", map[string]interface{}{
			"name":  "李四",
			"time":  "2023-05-15",
			"count": 5,
		})
		assert.NoError(t, err)
		assert.Equal(t, "用户 李四 在 2023-05-15 共享了 5 个文件", result, "Complex message format should be properly replaced")
	})

	t.Run("Complex message format with struct", func(t *testing.T) {
		type User struct {
			Name  string `json:"name"`
			Time  string `json:"time"`
			Count int    `json:"count"`
		}
		user := User{
			Name:  "李四",
			Time:  "2023-05-15",
			Count: 5,
		}
		result, err := dict.Translate("userInfo", user)
		assert.NoError(t, err)
		assert.Equal(t, "用户 李四 在 2023-05-15 共享了 5 个文件", result, "Complex message format should be properly replaced")
	})

	// Test nested data structures
	t.Run("Nested data structures in JSON args", func(t *testing.T) {
		result, err := dict.Translate("complex", map[string]interface{}{
			"num":     3,
			"options": []string{"选项A", "选项B", "选项C"},
		})
		assert.NoError(t, err)
		assert.Contains(t, result, "这有 3 个选项：", "Should contain number of options")
		assert.Contains(t, result, "选项A", "Should contain option A")
		assert.Contains(t, result, "选项B", "Should contain option B")
		assert.Contains(t, result, "选项C", "Should contain option C")
	})
}

func TestTranslateErrors(t *testing.T) {
	dict := Dict{
		"key":    "值",
		"plural": []string{"单数", "复数"},
	}

	// Test key not found error
	t.Run("Key not found error", func(t *testing.T) {
		result, err := dict.Translate("nonexistent")
		assert.Error(t, err)
		assert.Equal(t, "nonexistent", result, "Should return key when not found")

		var tErr *TranslateError
		assert.ErrorAs(t, err, &tErr, "Error should be a TranslateError")
		assert.Equal(t, "key_not_found", tErr.Type, "Error type should be key_not_found")
	})

	// Test invalid argument type
	t.Run("Invalid argument type", func(t *testing.T) {
		_, err := dict.Translate("key", 123) // Not a map or []string
		assert.Error(t, err)

		var tErr *TranslateError
		assert.ErrorAs(t, err, &tErr, "Error should be a TranslateError")
		assert.Equal(t, "invalid_argument", tErr.Type, "Error type should be invalid_argument")
	})

	// Test empty string array
	t.Run("Empty string array", func(t *testing.T) {
		_, err := dict.Translate("key", []string{})
		assert.Error(t, err)

		var tErr *TranslateError
		assert.ErrorAs(t, err, &tErr, "Error should be a TranslateError")
		assert.Equal(t, "empty_string_array", tErr.Type, "Error type should be empty_string_array")
	})
}

func TestToString(t *testing.T) {
	dict := Dict{"test": 123}

	result, err := dict.Translate("test")
	assert.NoError(t, err)
	assert.Equal(t, "123", result, "Should convert number to string")

	dict = Dict{"bool": true}
	result, err = dict.Translate("bool")
	assert.NoError(t, err)
	assert.Equal(t, "true", result, "Should convert boolean to string")
}
