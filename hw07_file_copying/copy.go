package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetNegative        = errors.New("offset cannot be negative")
	ErrLimitNegative         = errors.New("limit cannot be negative")
)

func Copy(fromPath, toPath string, offset, limit int64, showProgress bool) error {
	if offset < 0 {
		return ErrOffsetNegative
	}
	if limit < 0 {
		return ErrLimitNegative
	}

	reader, size, err := getReader(fromPath, offset)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := getWriter(toPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	totalWrite := size - offset
	if limit > 0 && limit < totalWrite {
		totalWrite = limit
	}

	if showProgress {
		bar := pb.Full.Start64(totalWrite)
		reader = bar.NewProxyReader(reader)
		defer bar.Finish()
	}

	_, err = io.CopyN(writer, reader, totalWrite)
	if err != nil {
		return err
	}

	return nil
}

func getReader(path string, offset int64) (io.ReadCloser, int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		if err != nil {
			file.Close()
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	if fileInfo.Size() == 0 {
		err = ErrUnsupportedFile
		return nil, 0, err
	}

	if offset > 0 {
		if fileInfo.Size() < offset {
			err = ErrOffsetExceedsFileSize
			return nil, 0, err
		}

		_, err = file.Seek(offset, 0)
		if err != nil {
			err = ErrOffsetExceedsFileSize
			return nil, 0, err
		}
	}

	return file, fileInfo.Size(), nil
}

func getWriter(path string) (io.WriteCloser, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}
