package main

import (
	"errors"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const MB = 1024 * 1024

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if !fromStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fromSize := fromStat.Size()
	if fromSize <= offset {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = fromSize - offset
	}

	in, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = in.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	barMax := limit
	bar := pb.Full.Start64(barMax)
	barReader := bar.NewProxyReader(in)

	buf := make([]byte, int(math.Min(float64(MB), float64(limit))))
	for {
		n, err := barReader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if _, err := out.Write(buf[:n]); err != nil {
			return err
		}
		limit -= int64(n)
		if limit <= 0 {
			break
		}
	}

	bar.SetCurrent(barMax)
	bar.Finish()

	return nil
}
