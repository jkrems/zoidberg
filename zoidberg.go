package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	source := string(bytes)
	tokens := Tokenize(filename, source)
	fmt.Println("Tokens:")
	for tok := range tokens {
		fmt.Printf("- %q\n", tok)
	}
}
