package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fail()
	}

	cases := []struct {
		name         string
		cmdPath      string
		expectedCode int
	}{
		{
			name:    "Existed command",
			cmdPath: path + "/testdata/echo.sh",
		},
		{
			name:         "Existed command exit code 2",
			cmdPath:      path + "/testdata/exit_code_2.sh",
			expectedCode: 2,
		},
		{
			name:         "Command not found",
			cmdPath:      path + "/testdata/notfound.sh",
			expectedCode: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			returnCode := RunCmd([]string{tc.cmdPath}, Environment{})
			require.Equal(t, tc.expectedCode, returnCode)
		})
	}
}
