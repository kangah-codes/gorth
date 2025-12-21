package main

import (
	"fmt"
	"strings"

	"gorth/lexer"
	"gorth/parser"
)

func main() {
	// Test input - modify this to test different code
	input := `
		3 4 + $sum =
		[1,2,3,4] $numbers = 
	`

	// Create lexer
	lex := lexer.NewLexer(strings.NewReader(input))

	// Print all tokens
	fmt.Println("Tokens:")
	fmt.Println("-------")
	for {
		tok, err := lex.NextToken()
		if err != nil {
			fmt.Println(fmt.Errorf("NextToken() threw error %v", err))
		}
		fmt.Printf("%-15s %-15s %s\n", tok.Type, tok.Literal, tok.Pos)

		if tok.Type == lexer.EOF {
			break
		}
	}

	newLexer := lexer.NewLexer(strings.NewReader(input))
	parse := parser.NewParser(newLexer)
	program, err := parse.Parse()
	if err != nil {
		fmt.Println("Error during parsing:", err)
		return
	}

	fmt.Println("Parsed Program AST:")
	fmt.Println("-------------------")
	fmt.Println(program.String())
	fmt.Println()
}
