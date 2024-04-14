package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	cases := []struct {
		name        string
		fromPath    string
		outPath     string
		comparePath string
		offset      int64
		limit       int64
	}{
		{
			name:        "offset 0 limit 0",
			fromPath:    "./testdata/input.txt",
			outPath:     "./out0_0.txt",
			comparePath: "./testdata/out_offset0_limit0.txt",
		},
		{
			name:        "offset 0 limit 10",
			fromPath:    "./testdata/input.txt",
			outPath:     "./out0_10.txt",
			comparePath: "./testdata/out_offset0_limit10.txt",
			limit:       10,
		},
		{
			name:        "offset 0 limit 1000",
			fromPath:    "./testdata/input.txt",
			outPath:     "./out0_1000.txt",
			comparePath: "./testdata/out_offset0_limit1000.txt",
			limit:       1000,
		},
		{
			name:        "offset 0 limit 10000",
			fromPath:    "./testdata/input.txt",
			outPath:     "./out0_10000.txt",
			comparePath: "./testdata/out_offset0_limit10000.txt",
			limit:       10000,
		},
		{
			name:        "offset 100 limit 1000",
			fromPath:    "./testdata/input.txt",
			outPath:     "./out100_1000.txt",
			comparePath: "./testdata/out_offset100_limit1000.txt",
			offset:      100,
			limit:       1000,
		},
		{
			name:        "offset 600 limit 1000",
			fromPath:    "./testdata/input.txt",
			outPath:     "./out6000_1000.txt",
			comparePath: "./testdata/out_offset6000_limit1000.txt",
			offset:      6000,
			limit:       1000,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outFile, err := os.Create(tc.outPath)
			if err != nil {
				t.Fatal(err)
			}
			outFile.Close()
			defer os.Remove(tc.outPath)

			err = Copy(tc.fromPath, tc.outPath, tc.offset, tc.limit, false)
			if err != nil {
				t.Fatal(err)
			}

			outData, err := os.ReadFile(tc.outPath)
			if err != nil {
				t.Fatal(err)
			}
			compareData, err := os.ReadFile(tc.comparePath)
			if err != nil {
				t.Fatal(err)
			}

			require.True(t, bytes.Equal(compareData, outData))
		})
	}
}

func TestFail(t *testing.T) {
	cases := []struct {
		name     string
		fromPath string
		outPath  string
		offset   int64
		limit    int64
		err      error
	}{
		{
			name:     "ErrLimitNegative",
			fromPath: "./testdata/input.txt",
			outPath:  "./out.txt",
			limit:    -1,
			err:      ErrLimitNegative,
		},
		{
			name:     "ErrOffsetNegative",
			fromPath: "./testdata/input.txt",
			outPath:  "./out.txt",
			offset:   -1,
			err:      ErrOffsetNegative,
		},
		{
			name:     "ErrOffsetExceedsFileSize",
			fromPath: "./testdata/input.txt",
			outPath:  "./out.txt",
			offset:   600000,
			err:      ErrOffsetExceedsFileSize,
		},
		{
			name:     "ErrUnsupportedFile",
			fromPath: "/dev/null",
			outPath:  "./out.txt",
			err:      ErrUnsupportedFile,
		},
		{
			name:     "ErrUnsupportedFile",
			fromPath: "/dev/urandom",
			outPath:  "./out.txt",
			err:      ErrUnsupportedFile,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outFile, err := os.Create(tc.outPath)
			if err != nil {
				t.Fatal(err)
			}
			outFile.Close()
			defer os.Remove(tc.outPath)

			err = Copy(tc.fromPath, tc.outPath, tc.offset, tc.limit, false)
			require.NotNil(t, err)
			require.ErrorIs(t, tc.err, err)
		})
	}
}
