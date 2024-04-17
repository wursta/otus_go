package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	cases := []struct {
		name     string
		dirPath  string
		expected Environment
	}{
		{
			name:    "Directory env",
			dirPath: "./testdata/env",
			expected: Environment{
				"BAR": EnvValue{
					Value: "bar",
				},
				"EMPTY": EnvValue{
					NeedRemove: true,
				},
				"FOO": EnvValue{
					Value: "   foo\nwith new line",
				},
				"HELLO": EnvValue{
					Value: "\"hello\"",
				},
				"UNSET": EnvValue{
					NeedRemove: true,
				},
			},
		},
		{
			name:    "Directory env2",
			dirPath: "./testdata/env2",
			expected: Environment{
				"TAB": EnvValue{
					Value: "ddd",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := ReadDir(tc.dirPath)
			require.Nil(t, err)
			require.Equal(t, tc.expected, env)
		})
	}
}

func TestErrReadDir(t *testing.T) {
	cases := []struct {
		name     string
		dirPath  string
		expected Environment
	}{
		{
			name:    "Directory without env files #1",
			dirPath: "./testdata",
		},
		{
			name:    "Directory without env files #2",
			dirPath: "./",
		},
		{
			name:    "Directory with unsupported files",
			dirPath: "./testdata/env3",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadDir(tc.dirPath)
			require.NotNil(t, err)
		})
	}
}
