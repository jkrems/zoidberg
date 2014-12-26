package main

import "fmt"

func main() {
	tokens := Tokenize("num.berg", "val a = 42\nf(x) { return x + a }")
	fmt.Println("Tokens:")
	for tok := range tokens {
		fmt.Printf("- %q\n", tok)
	}
}
