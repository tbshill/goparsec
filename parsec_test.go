package goparsec

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

func TestExpectRune(t *testing.T) {
	tok, rem, err := ExpectByte('a')("abc")
	if err != nil {
		t.Errorf("ExpectRune returned an error on a valid input")
	}
	if tok != "a" {
		t.Errorf("ExpectRune returned the wrong token")
	}
	if rem != "bc" {
		t.Errorf("ExpectRune returned the wrong remaining string")
	}
}

func BenchmarkExpectRune(b *testing.B) {
	for n := 0; n < b.N; n++ {
	tok, rem, err := ExpectByte('a')("abc")
		_, _, _ = tok, rem, err
	}
}

func TestExpectString(t *testing.T) {
	tok, rem, err := ExpectString("Hello")("Hello World")
	if err != nil {
		t.Errorf("ExpectString returned an error on a valid input: %v", err)
	}
	if tok != "Hello" {
		t.Errorf("ExpectString returned the wrong token: %s", tok)
	}
	if rem != " World" {
		t.Errorf("ExpectString returned the wrong remaining string: %s", rem)
	}
}

func TestExpectCaseInsensitiveString(t *testing.T) {
	tok, rem, err := ExpectCaseInsensitiveString("HELLO")("Hello World")
	if err != nil {
		t.Errorf("ExpectString returned an error on a valid input: %v", err)
	}
	if tok != "Hello" {
		t.Errorf("ExpectString returned the wrong token: %s", tok)
	}
	if rem != " World" {
		t.Errorf("ExpectString returned the wrong remaining string: %s", rem)
	}
}

func BenchmarkExpectString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tok, rem, err := ExpectString("Hello")("Hello World")
		_, _, _ = tok, rem, err
	}
}

func TestExpectRuneFrom(t *testing.T) {
	ex := "abc"
	abc := ExpectRuneFrom(ex)

	for _, r := range ex {
		tok, rem, err := abc(string(r)+"Hello")
		if err != nil {
			t.Errorf("ExpectString returned an error on a valid input: %v", err)
		}
		if tok != string(r) {
			t.Errorf("ExpectString returned the wrong token: %s", tok)
		}
		if rem != "Hello" {
			t.Errorf("ExpectString returned the wrong remaining string: %s", rem)
		}
	}
}

func TestRepeat(t *testing.T) {
	in := "aaabbb"
	as := Repeat(ExpectRune('a'))
	bs := Repeat(ExpectRune('b'))

	tok, rem, err := as(in)
	if err != nil {
		t.Errorf("Repeat returned an error when it had valid input")
	}
	if tok != "aaa" {
		t.Errorf("Repeat did not return the correct token")
	}
	if rem != "bbb" {
		t.Errorf("Repeat did not return the correct remaining string")
	}


	tok, rem, err = bs(rem)
	if err != nil {
		t.Errorf("Repeat returned an error when it had valid input")
	}
	if tok != "bbb" {
		t.Errorf("Repeat did not return the correct token")
	}
	if rem != "" {
		t.Errorf("Repeat did not return the correct remaining string")
	}
}

func TestOptional(t *testing.T) {
	in := "a,b"
	p := And(ExpectRune('a'), Optional(ExpectRune(',')), ExpectRune('b'))
	tok, rem, err := p(in)

	if err != nil {
		t.Errorf("Optional returned an error when it had valid input")
	}
	if tok != "a,b" {
		t.Errorf("Optional did not return the correct token")
	}
	if rem != "" {
		t.Errorf("Optional did not return the correct remaining string")
	}

	tok, rem, err = p("ab")

	if err != nil {
		t.Errorf("Optional returned an error when it had valid input")
	}
	if tok != "ab" {
		t.Errorf("Optional did not return the correct token")
	}
	if rem != "" {
		t.Errorf("Optional did not return the correct remaining string")
	}
}

func TestExpectUntil(t *testing.T) {
	in  := "aaabaaabbaaa"
	p := ExpectUntil(ExpectString("bb"))

	tok, rem, err := p(in)
	if err != nil {
		t.Errorf("ExpectUntil returned an error when it had valid input")
	}
	if tok != "aaabaaa" {
		t.Errorf("ExpectUntil did not return the correct token: %s", tok)
	}
	if rem != "bbaaa" {
		t.Errorf("ExpectUntil did not return the correct remaining string: %s", rem)
	}
}
