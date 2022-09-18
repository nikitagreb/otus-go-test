package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	partSize int64 = 100
)

func Copy(fromPath string, toPath string, offset int64, limit int64) error {
	in, err := os.Open(fromPath)
	if err != nil {
		return printErr(err)
	}
	defer closeFile(in)

	fileInfo, err := in.Stat()
	if err != nil {
		return printErr(err)
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	out, err := os.Create(toPath)
	if err != nil {
		return printErr(err)
	}
	defer closeFile(out)

	if _, err = in.Seek(offset, io.SeekStart); err != nil {
		return printErr(err)
	}

	total := fileSize - offset
	if limit == 0 || limit > total {
		limit = total
	}

	bar := pb.StartNew(int(limit))

	for {
		size := partSize
		if size > limit {
			size = limit
		}

		n, err := io.CopyN(out, in, size)

		eof := errors.Is(err, io.EOF)

		if err != nil && !eof {
			return printErr(err)
		}

		bar.Add(int(n))
		limit -= n

		if eof || limit <= 0 {
			break
		}
	}

	bar.Finish()

	return nil
}

func printErr(err error) error {
	return fmt.Errorf("error while copy file: %w", err)
}

func closeFile(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println("error while close file: ", err)
	}
}
