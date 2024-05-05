package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InvalidValidator struct {
		Test string `validate:"test:19"`
	}

	ValidatorWithInvalidParam struct {
		Test int `validate:"min:test"`
	}
)

func TestValidateSuccess(t *testing.T) {
	tests := []interface{}{
		User{
			ID:    "77e1d514-a2ec-495c-96eb-674fd9e2a291",
			Age:   25,
			Email: "test@mail.ru",
			Role:  "stuff",
			Phones: []string{
				"89111112233",
				"89111112244",
			},
		},
		App{
			Version: "1.2.4",
		},
		Response{
			Code: 200,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			validationErrors := Validate(tt)
			require.Nil(t, validationErrors)
			_ = tt
		})
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				Phones: make([]string, 2),
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrLengthNotValid,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrMinValue,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrNotMatchPattern,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValueNotFoundInAllowedList,
				},
				ValidationError{
					Field: "Phones #0",
					Err:   ErrLengthNotValid,
				},
				ValidationError{
					Field: "Phones #1",
					Err:   ErrLengthNotValid,
				},
			},
		},
		{
			in: User{
				ID:    "77e1d514-a2ec-495c-96eb-674fd9e2a291",
				Age:   25,
				Email: "test@mail.ru",
				Role:  "testRole",
				Phones: []string{
					"89111112233",
					"89111112244",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Role",
					Err:   ErrValueNotFoundInAllowedList,
				},
			},
		},
		{
			in: User{
				ID:    "77e1d514-a2ec-495c-96eb-674fd9e2a291",
				Age:   25,
				Email: "test",
				Role:  "stuff",
				Phones: []string{
					"89111112233",
					"89111112244",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrNotMatchPattern,
				},
			},
		},
		{
			in: App{},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrLengthNotValid,
				},
			},
		},
		{
			in: Response{},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValueNotFoundInAllowedList,
				},
			},
		},
		{
			in: Response{
				Code: 503,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValueNotFoundInAllowedList,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			validationErrors := Validate(tt.in)

			require.EqualError(t, validationErrors, tt.expectedErr.Error())
			require.Equal(t, true, errors.As(validationErrors, &ValidationErrors{}))
			_ = tt
		})
	}
}

func TestValidatorInvalidErros(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          InvalidValidator{},
			expectedErr: ErrUnknownValidator,
		},
		{
			in:          ValidatorWithInvalidParam{},
			expectedErr: ErrValidatorInvalidParam,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.EqualError(t, err, tt.expectedErr.Error())
			require.Equal(t, false, errors.As(err, &ValidationErrors{}))
			_ = tt
		})
	}
}
