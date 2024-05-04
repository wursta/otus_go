package hw09structvalidator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var intValidators = map[string]func(int, string) error{
	"min": intMinValidator,

	"max": intMaxValidator,

	"in": intInValidator,
}

func intMinValidator(value int, params string) (err error) {
	min, err := strconv.Atoi(params)
	if err != nil {
		err = ErrValidatorInvalidParam
		return
	}

	if value < min {
		err = ErrMinValue
	}

	return
}

func intMaxValidator(value int, params string) (err error) {
	max, err := strconv.Atoi(params)
	if err != nil {
		err = ErrValidatorInvalidParam
		return
	}

	if value > max {
		err = ErrMaxValue
	}

	return
}

func intInValidator(value int, params string) (err error) {
	allowedStrs := strings.Split(params, ",")

	contains := slices.Contains(allowedStrs, fmt.Sprint(value))

	if !contains {
		err = ErrValueNotFoundInAllowedList
	}

	return
}
