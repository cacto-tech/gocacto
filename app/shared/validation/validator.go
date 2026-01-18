package validation

import (
	"fmt"
	"reflect"
	"strings"

	"cacto-cms/app/shared/errors"
)

// Validator provides validation functionality
type Validator struct {
	errors []FieldError
}

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// New creates a new validator
func New() *Validator {
	return &Validator{
		errors: make([]FieldError, 0),
	}
}

// Validate validates a struct using tags
func (v *Validator) Validate(s interface{}) error {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validation: expected struct, got %s", val.Kind())
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}

		// Get validation tags
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		// Parse validation rules
		rules := strings.Split(tag, ",")
		fieldName := getFieldName(field)

		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if err := v.validateRule(fieldName, fieldVal, rule); err != nil {
				v.errors = append(v.errors, FieldError{
					Field:   fieldName,
					Message: err.Error(),
				})
			}
		}
	}

	if len(v.errors) > 0 {
		return v.Err()
	}

	return nil
}

// validateRule validates a single rule
func (v *Validator) validateRule(fieldName string, val reflect.Value, rule string) error {
	parts := strings.SplitN(rule, "=", 2)
	ruleName := parts[0]
	ruleValue := ""

	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	case "required":
		return v.validateRequired(fieldName, val)
	case "min":
		return v.validateMin(fieldName, val, ruleValue)
	case "max":
		return v.validateMax(fieldName, val, ruleValue)
	case "email":
		return v.validateEmail(fieldName, val)
	case "url":
		return v.validateURL(fieldName, val)
	default:
		return nil
	}
}

// validateRequired checks if field is required
func (v *Validator) validateRequired(fieldName string, val reflect.Value) error {
	if isEmpty(val) {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// validateMin checks minimum length/value
func (v *Validator) validateMin(fieldName string, val reflect.Value, minStr string) error {
	if isEmpty(val) {
		return nil // Skip if empty (use required for that)
	}

	var min int
	fmt.Sscanf(minStr, "%d", &min)

	switch val.Kind() {
	case reflect.String:
		if len(val.String()) < min {
			return fmt.Errorf("%s must be at least %d characters", fieldName, min)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() < int64(min) {
			return fmt.Errorf("%s must be at least %d", fieldName, min)
		}
	case reflect.Slice, reflect.Array:
		if val.Len() < min {
			return fmt.Errorf("%s must have at least %d items", fieldName, min)
		}
	}

	return nil
}

// validateMax checks maximum length/value
func (v *Validator) validateMax(fieldName string, val reflect.Value, maxStr string) error {
	if isEmpty(val) {
		return nil
	}

	var max int
	fmt.Sscanf(maxStr, "%d", &max)

	switch val.Kind() {
	case reflect.String:
		if len(val.String()) > max {
			return fmt.Errorf("%s must be at most %d characters", fieldName, max)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() > int64(max) {
			return fmt.Errorf("%s must be at most %d", fieldName, max)
		}
	case reflect.Slice, reflect.Array:
		if val.Len() > max {
			return fmt.Errorf("%s must have at most %d items", fieldName, max)
		}
	}

	return nil
}

// validateEmail validates email format
func (v *Validator) validateEmail(fieldName string, val reflect.Value) error {
	if isEmpty(val) {
		return nil
	}

	email := val.String()
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("%s must be a valid email address", fieldName)
	}

	return nil
}

// validateURL validates URL format
func (v *Validator) validateURL(fieldName string, val reflect.Value) error {
	if isEmpty(val) {
		return nil
	}

	url := val.String()
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("%s must be a valid URL", fieldName)
	}

	return nil
}

// isEmpty checks if value is empty
func isEmpty(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.String:
		return val.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Slice, reflect.Array, reflect.Map:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	default:
		return false
	}
}

// getFieldName gets the field name from struct tag or field name
func getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		parts := strings.Split(jsonTag, ",")
		return parts[0]
	}
	return strings.ToLower(field.Name)
}

// Err returns validation error
func (v *Validator) Err() error {
	if len(v.errors) == 0 {
		return nil
	}

	messages := make([]string, len(v.errors))
	for i, err := range v.errors {
		messages[i] = err.Message
	}

	return errors.NewValidation(strings.Join(messages, "; "))
}

// Errors returns all field errors
func (v *Validator) Errors() []FieldError {
	return v.errors
}

// ValidateStruct is a convenience function
func ValidateStruct(s interface{}) error {
	v := New()
	return v.Validate(s)
}
