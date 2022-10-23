package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrNotValidatable         = errors.New("that entity cannot be validated")
	ErrInvalidValidationRule  = errors.New("invalid validation rule")
	ErrInvalidValidationValue = errors.New("invalid value")
)

func (ve ValidationErrors) Error() string {
	var msg strings.Builder

	for _, v := range ve {
		if msg.Len() != 0 {
			msg.WriteString(", ")
		}
		msg.WriteString(v.Error())
	}

	return msg.String()
}

func (ve *ValidationErrors) Add(field string, err error) {
	*ve = append(*ve, ValidationError{
		Field: field,
		Err:   err,
	})
}

func (v ValidationError) Error() string {
	if v.Err == nil {
		return ""
	}
	return v.Err.Error()
}

func Validate(v interface{}) error {
	rVal := reflect.ValueOf(v)

	if rVal.Type().Kind() != reflect.Struct {
		return ErrNotValidatable
	}

	validateErrors := make(ValidationErrors, 0)

	for i := 0; i < rVal.Type().NumField(); i++ {
		tField := rVal.Type().Field(i)

		// if private field
		if !tField.IsExported() {
			continue
		}

		strRules := tField.Tag.Get("validate")
		if strRules == "" {
			continue
		}

		field := rVal.Field(i)
		rules := strings.Split(strRules, "|")

		if err := runValidation(rules, field, tField, &validateErrors); err != nil {
			return err
		}
	}

	if len(validateErrors) > 0 {
		return validateErrors
	}

	return nil
}

func runValidation(rules []string, field reflect.Value, tField reflect.StructField, vErrors *ValidationErrors) error {
	for _, r := range rules {
		if isIterable(tField.Type.Kind()) {
			if err := runValidateMultiRule(r, field, tField, vErrors); err != nil {
				return err
			}
			continue
		}
		if err := validateOneRule(r, field, tField, vErrors); err != nil {
			return err
		}
	}
	return nil
}

func runValidateMultiRule(rule string, field reflect.Value, tField reflect.StructField, vErrs *ValidationErrors) error {
	slice := reflect.ValueOf(field.Interface())
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i)
		err := runRuleValidation(rule, item.Kind(), field)
		if err == nil {
			continue
		}
		if errors.Is(err, ErrInvalidValidationValue) {
			vErrs.Add(tField.Name, err)
		}
		return err
	}
	return nil
}

func validateOneRule(rule string, field reflect.Value, tField reflect.StructField, vErrors *ValidationErrors) error {
	if err := runRuleValidation(rule, tField.Type.Kind(), field); err != nil {
		if errors.Is(err, ErrInvalidValidationValue) {
			vErrors.Add(tField.Name, err)
			return nil
		}
		return err
	}
	return nil
}

func runRuleValidation(rule string, kind reflect.Kind, field reflect.Value) error {
	ruleSplit := strings.Split(rule, ":")
	if len(ruleSplit) != 2 {
		return ErrInvalidValidationRule
	}

	err := validate(ruleSplit[0], ruleSplit[1], kind, field)
	if err != nil {
		return err
	}

	return nil
}

func validate(rName string, rValStr string, kind reflect.Kind, field reflect.Value) error {
	if kind == reflect.String {
		switch rName {
		case "len":
			return validateStrLen(rValStr, field.String())
		case "in":
			return validateStrIn(rValStr, field.String())
		case "regexp":
			return validateStrRegexp(rValStr, field.String())
		}
	}

	if kind == reflect.Int {
		switch rName {
		case "min":
			return validateIntMin(rValStr, field.Int())
		case "max":
			return validateIntMax(rValStr, field.Int())
		case "in":
			return validateIntIn(rValStr, field.Int())
		}
	}

	return nil
}

func isIterable(kind reflect.Kind) bool {
	return kind == reflect.Slice || kind == reflect.Array
}

func validateStrLen(ruleValStr string, fieldValStr string) error {
	ruleVal, err := strconv.Atoi(ruleValStr)
	if err != nil {
		return ErrInvalidValidationRule
	}

	if ruleVal == len(fieldValStr) {
		return nil
	}

	return wrapValidationErr(fmt.Sprintf(
		"value length (%d) is not match required length (%s)",
		len(fieldValStr),
		ruleValStr,
	))
}

func validateStrIn(ruleValStr string, fieldValStr string) error {
	in := strings.Split(ruleValStr, ",")
	for _, n := range in {
		if n == fieldValStr {
			return nil
		}
	}

	return wrapValidationErr(fmt.Sprintf(
		"value (%s) is not in (%s)",
		fieldValStr,
		ruleValStr,
	))
}

func validateStrRegexp(ruleValStr string, fieldValStr string) error {
	match, err := regexp.MatchString(ruleValStr, fieldValStr)
	if match && err == nil {
		return nil
	}

	return wrapValidationErr(fmt.Sprintf(
		"value (%s) is not matched regexp (%s)",
		fieldValStr,
		ruleValStr,
	))
}

func validateIntMin(ruleValStr string, fieldValStr int64) error {
	min, err := strconv.Atoi(ruleValStr)
	if err != nil {
		return ErrInvalidValidationRule
	}
	if int64(min) <= fieldValStr {
		return nil
	}

	return wrapValidationErr(fmt.Sprintf(
		"value (%d) is less than (%d)",
		fieldValStr,
		min,
	))
}

func validateIntMax(ruleValStr string, fieldValStr int64) error {
	min, err := strconv.Atoi(ruleValStr)
	if err != nil {
		return ErrInvalidValidationRule
	}
	if int64(min) >= fieldValStr {
		return nil
	}

	return wrapValidationErr(fmt.Sprintf(
		"value (%d) is more than (%d)",
		fieldValStr,
		min,
	))
}

func validateIntIn(ruleValStr string, fieldValStr int64) error {
	in := strings.Split(ruleValStr, ",")
	sInt := make([]int64, 0, len(in))
	for _, s := range in {
		i, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidValidationRule
		}
		sInt = append(sInt, int64(i))
	}
	for _, i := range sInt {
		if i == fieldValStr {
			return nil
		}
	}

	return wrapValidationErr(fmt.Sprintf(
		"value (%d) is not in (%s)",
		fieldValStr,
		ruleValStr,
	))
}

func wrapValidationErr(str string) error {
	return fmt.Errorf("%w: %s", ErrInvalidValidationValue, str)
}
