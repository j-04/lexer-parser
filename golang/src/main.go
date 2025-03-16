package main

import (
	"os"

	"github.com/j-04/lexer-parser/golang/src/lexer"
)

func main() {
	bytes, _ := os.ReadFile("./examples/01.lang")
	source := string(bytes)
	tokens := lexer.Tokenize(string(source))
	for _, token := range tokens {
		token.Debug()
	}	
}