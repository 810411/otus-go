package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const Backslash = 92

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var res strings.Builder
	ra := []rune(s)
	for i := 0; i < len(ra); i++ {
		if ra[i] == Backslash {
			i++
			res.WriteRune(ra[i])
		} else if n, err := strconv.Atoi(string(ra[i])); err == nil {
			if i == 0 || (i > 0 && unicode.IsDigit(ra[i-1]) && (i > 1 && ra[i-2] != Backslash)) {
				return "", ErrInvalidString
			}
			res.WriteString(strings.Repeat(string(ra[i-1]), n-1))
		} else {
			res.WriteRune(ra[i])
		}
	}
	return res.String(), nil
}
