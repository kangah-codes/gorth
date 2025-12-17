package main

import (
	"fmt"
	"strings"

	"gorth/lexer"
)

func main() {
	// Test input - modify this to test different code
	input := `
		"Hello, World!" print drop

		# or this, if in strict mode

		"Hello, World!" dump

		|| $name $fuck123 12.89 pi 3.14 =
	`

	// Create lexer
	lex := lexer.NewLexer(strings.NewReader(input))

	// Print all tokens
	fmt.Println("Tokens:")
	fmt.Println("-------")
	for {
		tok := lex.NextToken()
		fmt.Printf("%-15s %-15s %s\n", tok.Type, tok.Literal, tok.Pos)

		if tok.Type == lexer.EOF {
			break
		}
	}
}
