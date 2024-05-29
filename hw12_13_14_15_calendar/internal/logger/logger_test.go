package logger

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	cases := []struct {
		name          string
		level         string
		expectedLevel int
		regex         []string
	}{
		{
			name:          "debug level",
			level:         "DEBUG",
			expectedLevel: DEBUG,
			regex: []string{
				"\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} DEBUG Test debug",
				"\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} INFO Test info",
				"\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} ERROR Test error",
			},
		},
		{
			name:          "info level",
			level:         "INFO",
			expectedLevel: INFO,
			regex: []string{
				"\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} INFO Test info",
				"\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} ERROR Test error",
			},
		},
		{
			name:          "error level",
			level:         "ERROR",
			expectedLevel: ERROR,
			regex: []string{
				"\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} ERROR Test error",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			logger, err := New(tc.level, &output)
			require.Nil(t, err)
			require.Equal(t, tc.expectedLevel, logger.severityLevel)

			logger.Debug("Test debug")
			logger.Info("Test info")
			logger.Error("Test error")

			logMessages := output.String()
			for i := range tc.regex {
				match, _ := regexp.MatchString(tc.regex[i], logMessages)
				require.True(t, match, "Log message not match pattern. Actual message: "+logMessages)
			}
		})
	}
}

func TestLoggerFail(t *testing.T) {
	_, err := New("UNKNOWN_LEVEL", os.Stderr)
	require.Equal(t, ErrUnknownLoggerLevel, err)
}
