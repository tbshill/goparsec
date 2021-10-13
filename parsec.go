package goparsec

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	ErrNoInput = errors.New("Ran out of input")
)

type TextParser func(string) (string, string, error)

// Drop drops the token from a parser and replaces it with an empty string.
func Drop(p TextParser) TextParser {
	return func(in string) (tok string, rem string, err error) {
		_, rem, err = p(in)
		return
	}
}

func checkInputSize(p TextParser) TextParser {
	return func(in string) (string, string, error) {
		if len(in) == 0 {
			return "", "", ErrNoInput
		}
		return p(in)
	}
}

// ExpectByte expects a specific byte to be the next character in the sequence
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

// ExpectRune expects a rune to be the next character in the sequence
// This allows unicode support
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

// ExpectString expects a predefined string to be the next sequence of
// characters in the input
func ExpectString(s string) TextParser {
	return checkInputSize(func(in string) (string, string, error) {

		slen := len(s)

		if slen > len(in) {
			return "", in, expectStringError(s, in)
		}
		if s != in[:slen] {
			return "", in, expectStringError(s, in[:slen])
		}
		return in[:slen], in[slen:], nil
	})
}

// ExpectCaseInsensitiveString expects a predefined string to be the next
// sequence of characters in teh input (ignoring case)
func ExpectCaseInsensitiveString(s string) TextParser {
	return checkInputSize(func(in string) (string, string, error) {
		slen := len(s)

		if slen > len(in) {
			return "", in, expectStringError(s, in)
		}
		if strings.ToUpper(s) != strings.ToUpper(in[:slen]) {
			return "", in, expectStringError(s, in[:slen])
		}
		return in[:slen], in[slen:], nil
	})
}

func expectStringError(expect, got string) error {
	return fmt.Errorf("Expect '%s', Got '%s'", expect, got)
}

// ExpectRuneFrom expects any rune in the string to be the next character in the sequence
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

// ExpectAnyRune will take off the next rune from the front of the input
var ExpectAnyRune = checkInputSize(func(in string) (string, string, error) {
	_, size := utf8.DecodeRuneInString(in)
	return in[:size], in[size:], nil
})

// And joins a sequence of parsers together such that all parsers must
// succeed in sequence. The token represents the cumulated output from all
// of the parsers.
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

// Or joins a sequence of parsers together such that only one of the parsers
// must succeed. The parsers are processed in order. The token
// is the token from the first parser to return without an error
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

// Optional will attempt to parse with a given parser. If it fails
// no error is thrown and the remaining represents the original input
func Optional(parser TextParser) TextParser {
	return func(in string) (string, string, error) {
		tok, rem, _ := parser(in)
		return tok, rem, nil
	}
}

// Repeat will repeat a parser until it fails. The token is the
// cumulative token from all successive parses.
func Repeat(parser TextParser) TextParser {
	return func(in string) (string, string, error) {
		tok, rem, err := parser(in)
		if err != nil {
			return "", in, err
		}

		for {
			t, r, err := parser(rem)
			if err != nil {
				return tok, rem, nil
			}
			tok += t
			rem = r
		}
	}
}

// ExpectEOI Expects the End Of Input
func ExpectEOI(in string) (string, string, error) {
	if in != "" {
		return "", in, expectEOIError()
	}
	return "", "", nil
}

// ExpectUntil will consume any rune until it the parser sucessfully parses the input
// The token will be all text upto (but not including) token from the parser
func ExpectUntil(parser TextParser) TextParser {
	return func(in string) (string, string, error) {
		var tmpTok, tok, rem string
		var err, err1 error

		rem = in
		for {
			_, _, err = parser(rem)
			if err != nil {
				// ExpectAnyRun forces consumption of a character so that the
				// parser can be applied on the next character
				tmpTok, rem, err1 = ExpectAnyRune(rem)
				if err1 == ErrNoInput {
					return "", in, ErrNoInput
				}
				tok += tmpTok
			} else {
				return tok, rem, nil
			}
		}
	}
}

// ExpectThrough will consume any rune until the parser successfully parses the input
// The token is all of the upto and including the token parsed by the parser.
func ExpectThrough(parser TextParser) TextParser {
	return func(in string) (string, string, error) {
		var (
			tmpTok, tok, rem string
			err, err1        error
		)
		rem = in
		for {
			tmpTok, rem, err = parser(rem)
			if err != nil {
				// ExpectAnyRune forces consumption of a character so that the
				// parser can be applied on the next character
				tmpTok, rem, err1 = ExpectAnyRune(rem)
				if err1 == ErrNoInput {
					return "", in, ErrNoInput
				}
				tok += tmpTok

			} else {
				tok += tmpTok
				return tok, rem, err
			}
		}
	}
}

func expectEOIError() error {
	return fmt.Errorf("Expected the end of input")
}

var (
	// ExpectDigit expects a digit 0-9
	ExpectDigit = ExpectRuneFrom("1234567890")

	// ExpectLetter expects a character from a-zA-Z
	ExpectLetter = ExpectRuneFrom("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// ExpectWhiteSpace expects a space, tab, carriage return, or newline
	ExpectWhiteSpace = ExpectRuneFrom(" \t\r\n")

	// ExpectUnixNewLine expects an \n
	ExpectUnixNewLine = ExpectRune('\n')

	// ExpectWindowsNewLine expects a \r\n
	ExpectWindowsNewLine = ExpectString("\r\n")

	// ExpectNewLine first checks for Unix, and then Windows newlines
	ExpectNewLine = Or(ExpectUnixNewLine, ExpectWindowsNewLine)
)
