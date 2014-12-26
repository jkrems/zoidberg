package main

import "fmt"

func main() {
	_, tokens := Tokenize("num.berg", "42")
	fmt.Println("Tokens:", tokens)
}
