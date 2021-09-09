package goparsec

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

var (
	ErrNoInput = errors.New("Ran out of input")
)

type TextParser func(string) (string, string, error)

func checkInputSize(p TextParser) TextParser {
	return func(in string) (string, string, error) {
		if len(in) == 0 {
			return "", "", ErrNoInput
		}
		return p(in)
	}
}

func ExpectByte(b byte) TextParser {
	return checkInputSize(func(in string) (string, string, error) {
		if in[0] != b {
			return "", in, expectByteError(b, in[0])
		}
		return in[:1], in[1:], nil
	})
}

func expectByteError(expect, got byte) error {
	return fmt.Errorf("Expected '%b', Got '%b'", expect, got)
}

func ExpectRune(r rune) TextParser {
	return checkInputSize(func(in string) (string, string, error) {
		got, s := utf8.DecodeRuneInString(in)
		if r != got {
			return "", in, expectRuneError(r, got)
		}
		return in[:s], in[s:], nil
	})
}

func expectRuneError(expect, got rune) error {
	return fmt.Errorf("Expected '%c', Got '%c'", expect, got)
}

func ExpectString(s string) TextParser {
	return checkInputSize(func(in string) (string, string, error) {

		slen := len(s)

		if slen < len(in) {
			return "", in, expectStringError(s, in)
		}
		if s != in[:slen] {
			return "", in, expectStringError(s, in[:slen])
		}
		return in[:slen], in[slen:], nil
	})
}

func expectStringError(expect, got string) error {
	return fmt.Errorf("Expect '%s', Got '%s'", expect, got)
}

func ExpectRuneFrom(s string) TextParser {
	m := make(map[rune]struct{})
	for _, r := range s {
		m[r] = struct{}{}
	}

	return checkInputSize(func(in string) (string, string, error) {
		r, size := utf8.DecodeRuneInString(in)
		if _, ok := m[r]; !ok {
			return "", in, expectRuneFromError(s, r)
		}
		return in[:size], in[size:], nil
	})
}

func expectRuneFromError(expect string, got rune) error {
	return fmt.Errorf("Expected rune from string '%s'. Got '%c'", expect, got)
}

var ExpectAnyRune = checkInputSize(func(in string) (string, string, error) {
	_, size := utf8.DecodeRuneInString(in)
	return in[:size], in[size:], nil
})

func And(parsers ...TextParser) TextParser {
	return func(in string) (string, string, error) {
		tok, rem := "", in
		for _, parser := range parsers {
			if tmpTok, tmpRem, err := parser(rem); err != nil {
				return "", in, err
			} else {
				tok += tmpTok
				rem = tmpRem
			}
		}
		return tok, rem, nil
	}
}

func Or(parsers ...TextParser) TextParser {
	return func(in string) (string, string, error) {
		for _, parser := range parsers {
			if tok, rem, err := parser(in); err == nil {
				return tok, rem, err
			}
		}
		return "", in, fmt.Errorf("No match")
	}
}

func ExpectEOI(in string) (string, string, error) {
	if in != "" {
		return "", in, expectEOIError()
	}
	return "", "", nil
}

func expectEOIError() error {
	return fmt.Errorf("Expected the end of input")
}

var (
	ExpectDigit      = ExpectRuneFrom("1234567890")
	ExpectLetter     = ExpectRuneFrom("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ExpectWhiteSpace = ExpectRuneFrom(" \t\r\n")
	ExpectUnixNewLine = ExpectRune('\n')
	ExpectWindowsNewLine = ExpectString("\r\n")
	ExpectNewLine = Or(ExpectWindowsNewLine, ExpectUnixNewLine)
)
