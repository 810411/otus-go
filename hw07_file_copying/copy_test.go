package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	inTest := "in_file"
	dataTest := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	testCases := []struct {
		desc     string
		fromFile string
		toFile   string
		limit    int64
		offset   int64
		want     []byte
		err      error
	}{
		{
			desc:     "Test full copy file",
			fromFile: inTest,
			toFile:   "out_file",
			limit:    0,
			offset:   0,
			want:     dataTest,
			err:      nil,
		},
		{
			desc:     "Test offset",
			fromFile: inTest,
			toFile:   "out_file_offset",
			limit:    0,
			offset:   5,
			want:     []byte{5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			err:      nil,
		},
		{
			desc:     "Test offset limit",
			fromFile: inTest,
			toFile:   "out_file_offset_limit",
			limit:    10,
			offset:   5,
			want:     []byte{5, 6, 7, 8, 9, 0, 1, 2, 3, 4},
			err:      nil,
		},
		{
			desc:     "Test err",
			fromFile: inTest,
			toFile:   "out_file_err",
			limit:    10,
			offset:   50,
			want:     nil,
			err:      errors.New("any error"),
		},
		{
			desc:     "Test file not exists",
			fromFile: "file_not_exists",
			toFile:   "out_file_err",
			limit:    0,
			offset:   0,
			want:     nil,
			err:      errors.New("any error"),
		},
		{
			desc:     "Test copy non regular",
			fromFile: "/dev/urandom",
			toFile:   "out_file_err",
			limit:    0,
			offset:   0,
			want:     nil,
			err:      errors.New("any error"),
		},
	}

	err := fileCreate(inTest, dataTest)
	if err != nil {
		t.Error(fmt.Errorf("file creation error %v", err))
	}
	defer func() {
		err = os.Remove(inTest)
		if err != nil {
			t.Error(fmt.Errorf("file deleting error %v", err))
		}
	}()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			err := Copy(tC.fromFile, tC.toFile, tC.offset, tC.limit)
			if err != nil && tC.err == nil || err == nil && tC.err != nil {
				t.Errorf("%s, got=%v, expected=%v", tC.desc, err, tC.err)
			}

			if err == nil {
				got, err := ioutil.ReadFile(tC.toFile)
				if err != nil {
					t.Errorf("%s, error ReadFile %q %v", tC.desc, tC.toFile, err)
				}
				if !bytes.Equal(tC.want, got) {
					t.Errorf("%s, got=%v, expected=%v", tC.desc, got, tC.want)
				}
			}

			if _, err := os.Stat(tC.toFile); err == nil {
				err := os.Remove(tC.toFile)
				if err != nil {
					t.Error(fmt.Errorf("file deleting error %q %v", tC.toFile, err))
				}
			}
		})
	}
}

func fileCreate(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
