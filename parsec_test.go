package parsec

import (
	"testing"
)

func TestCheckInputSize(t *testing.T) {
	_, _, err := checkInputSize(ExpectAnyRune)("")
	if err != ErrNoInput {
		t.Fatalf("checkInputSize did not throw an error when an empty string was passed")
	}
	_, _, err = checkInputSize(ExpectAnyRune)("Hello")
	if err == ErrNoInput {
		t.Fatalf("checkInputSize threw an error when a string was passed")
	}
}

func BenchmarkCheckInputSize_NoErr(b *testing.B) {
    for n := 0; n < b.N; n++ {
		_, _, err := checkInputSize(ExpectAnyRune)("ABC")
		_ = err
	}
}

func BenchmarkCheckInputSize_Err(b *testing.B) {
    for n := 0; n < b.N; n++ {
		_, _, err := checkInputSize(ExpectAnyRune)("")
		_ = err
	}
}

func TestExpectByte(t *testing.T) {
	tok, rem, err := ExpectByte(byte(97))("abc")
	if err != nil {
		t.Errorf("ExpectByte returned an error on a valid input")
	}
	if tok != "a" {
		t.Errorf("ExpectByte returned the wrong token")
	}
	if rem != "bc" {
		t.Errorf("ExpectByte returned the wrong remaining string")
	}
}

func BenchmarkExpectByte(b *testing.B) {
    for n := 0; n < b.N; n++ {
		tok, rem, err := ExpectByte(byte(97))("abc")
		_, _, _ = tok, rem, err
	}
}
