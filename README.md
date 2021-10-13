# goparsec

This is a simple combinatoric parsing library written in Go. If you have never heard of combinatoric parsing, checkout my blog post: https://medium.com/braden-shill-cs-review/combinatoric-parsing-with-golang-1f5b4af10324

```go
package main

import (
	"fmt"
	. "github.com/tbshill/goparsec"
)

func main() {
	input := "LDA R0 R1"

	var (
		ExpectLDA             = ExpectString("LDA")
		ExpectPC              = ExpectString("PC")
		ExpectNumericRegister = And(ExpectRune('R'), ExpectDigit)
		ExpectRegister        = Or(ExpectPC, ExpectNumericRegister)
	)

	parser := And(
		ExpectLDA,
		ExpectWhiteSpace,
		ExpectRegister,
		ExpectWhiteSpace,
		ExpectRegister)

	token, rem, err := parser(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed Text:", token)
	fmt.Println("Remaining Text:", rem)
}
```


