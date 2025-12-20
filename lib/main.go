package main

import (
	"fmt"
	"strings"

	"gorth/lexer"
)

func main() {
	// Test input - modify this to test different code
	input := `
		$variable
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
