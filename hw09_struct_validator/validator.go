package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrorTypeNotAllowed           = errors.New("validated type is not allowed")
	ErrUnknownValidator           = errors.New("unknown validator")
	ErrValidatorInvalidParam      = errors.New("validator has invalid parameter")
	ErrLengthNotValid             = errors.New("value length is not valid")
	ErrNotMatchPattern            = errors.New("value not match pattern")
	ErrMinValue                   = errors.New("value is more than min")
	ErrMaxValue                   = errors.New("value is less than max")
	ErrValueNotFoundInAllowedList = errors.New("value not found in allowed list of values")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type ValidatorInfo struct {
	Name   string
	Params string
}

type AllowedFieldTypes interface {
	~string | ~int | ~int64 | ~int32
}

func (v ValidationErrors) Error() string {
	errors := make([]string, len(v))

	for i, k := range v {
		errors[i] = fmt.Sprintf("%s: %v", k.Field, k.Err)
	}

	return strings.Join(errors, "\n")
}

func (v ValidationErrors) Unwrap() []error {
	errList := make([]error, len(v))

	for i, k := range v {
		errList[i] = k.Err
	}
	return errList
}

func Validate(v interface{}) error {
	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.Struct:
		err := validateStruct(r)
		if err != nil {
			return err
		}

		return nil
	default:
		return ErrorTypeNotAllowed
	}
}

func validateStruct(r reflect.Value) error {
	validationErrors := ValidationErrors{}

	for fIndex := 0; fIndex < r.NumField(); fIndex++ {
		fieldType := r.Type().Field(fIndex)
		tag, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}

		validatorsInfo := getValidatorsInfoFromTag(tag)
		if len(validatorsInfo) == 0 {
			continue
		}

		fieldValue := r.Field(fIndex)

		var err error
		for i := 0; i < len(validatorsInfo); i++ {
			switch fieldType.Type.Kind() {
			case reflect.Slice:
				sliceFieldType := fieldType.Type.Elem().Kind()
				for j := 0; j < fieldValue.Len(); j++ {
					err = validateStructField(
						&validationErrors,
						sliceFieldType,
						fieldType.Name+" #"+strconv.Itoa(j),
						fieldValue.Index(j),
						validatorsInfo[i],
					)

					if err != nil {
						return err
					}
				}
			default:
				err = validateStructField(
					&validationErrors,
					fieldType.Type.Kind(),
					fieldType.Name,
					fieldValue,
					validatorsInfo[i],
				)

				if err != nil {
					return err
				}
			}
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func validateStructField(
	validationErrors *ValidationErrors,
	fieldType reflect.Kind,
	fieldName string,
	fieldValue reflect.Value,
	validatorInfo ValidatorInfo,
) error {
	switch fieldType {
	case reflect.String:
		return validateField(validationErrors, fieldName, fieldValue.String(), validatorInfo, stringsValidators)
	case reflect.Int:
		return validateField(validationErrors, fieldName, int(fieldValue.Int()), validatorInfo, intValidators)
	default:
		return nil
	}
}

func validateField[T AllowedFieldTypes](
	validationErrors *ValidationErrors,
	fieldName string,
	value T,
	validator ValidatorInfo,
	allowedValidators map[string]func(T, string) error,
) error {
	validatorFunc, ok := allowedValidators[validator.Name]
	if !ok {
		return ErrUnknownValidator
	}
	err := validatorFunc(value, validator.Params)

	if err == nil {
		return nil
	}

	if errors.Is(err, ErrValidatorInvalidParam) {
		return err
	}

	*validationErrors = append(*validationErrors, ValidationError{
		Field: fieldName,
		Err:   err,
	})

	return nil
}

func getValidatorsInfoFromTag(tag string) []ValidatorInfo {
	info := []ValidatorInfo{}

	validators := strings.Split(tag, "|")

	for i := 0; i < len(validators); i++ {
		tmp := strings.Split(validators[i], ":")

		info = append(info, ValidatorInfo{
			Name:   tmp[0],
			Params: tmp[1],
		})
	}

	return info
}

func sliceContains[S ~[]E, E comparable](s S, v E) bool {
	for i := range s {
		if v == s[i] {
			return true
		}
	}
	return false
}
