package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var b strings.Builder
	var prev rune
	var havePrev bool
	escaped := false

	for i := 0; i < len(s); {
		r, w := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && w == 1 {
			return "", errors.New("invalid UTF-8 encoding")
		}
		i += w

		if escaped {
			b.WriteRune(r)
			prev = r
			havePrev = true
			escaped = false
			continue
		}

		switch {
		case r == '\\':
			if i >= len(s) {
				return "", errors.New("dangling escape at end of string")
			}
			escaped = true

		case unicode.IsDigit(r):
			if !havePrev {
				return "", errors.New("digit cannot start or follow nothing")
			}
			count := int(r - '0')
			for k := 1; k < count; k++ {
				b.WriteRune(prev)
			}

		default:
			b.WriteRune(r)
			prev = r
			havePrev = true
		}
	}

	return b.String(), nil
}


func main() {
	tests := []string{
		"a4bc2d5e",
		"abcd",
		"45",
		"",
		"qwe\\4\\5",
		"qwe\\45",
	}

	for _, t := range tests {
		result, err := Unpack(t)
		if err != nil {
			fmt.Printf("input: %-8q -> error: %v\n", t, err)
		} else {
			fmt.Printf("input: %-8q -> output: %q\n", t, result)
		}
	}
}