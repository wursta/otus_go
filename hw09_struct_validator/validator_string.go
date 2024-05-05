package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

var stringsValidators = map[string]func(string, string) error{
	"len": stringLenValidator,

	"regexp": stringRegexValidator,

	"in": stringInValidator,
}

func stringLenValidator(value, params string) (err error) {
	length, err := strconv.Atoi(params)
	if err != nil {
		err = ErrValidatorInvalidParam
		return
	}

	if len(value) != length {
		err = ErrLengthNotValid
	}

	return
}

func stringRegexValidator(value, params string) (err error) {
	regex, err := regexp.Compile(params)
	if err != nil {
		err = ErrValidatorInvalidParam
		return
	}

	match := regex.MatchString(value)

	if !match {
		err = ErrNotMatchPattern
	}

	return
}

func stringInValidator(value, params string) (err error) {
	allowedStrs := strings.Split(params, ",")

	contains := sliceContains(allowedStrs, value)

	if !contains {
		err = ErrValueNotFoundInAllowedList
	}

	return
}
