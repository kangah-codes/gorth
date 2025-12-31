package main

import (
	"fmt"
	"os"
	"strings"

	runtime "gorth/goruntime"
	"gorth/lexer"
	"gorth/parser"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	filename := os.Args[1]
	debugMode := false

	// Check for -d flag
	for _, arg := range os.Args[2:] {
		if arg == "-d" {
			debugMode = true
		} else {
			fmt.Printf("Unknown option: %s\n", arg)
			printUsage()
			os.Exit(1)
		}
	}

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	input := string(content)

	// Lexing phase
	lex := lexer.NewLexer(strings.NewReader(input))

	if debugMode {
		fmt.Println("=== LEXER OUTPUT ===")
		fmt.Printf("%-15s %-15s %s\n", "Token", "Literal", "Line:Col")
		fmt.Println(strings.Repeat("-", 50))

		tokens := []lexer.Token{}
		for {
			tok, err := lex.NextToken()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Lexer error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%-15s %-15s %s\n", tok.Type, tok.Literal, tok.Pos)
			tokens = append(tokens, tok)

			if tok.Type == lexer.EOF {
				break
			}
		}
		fmt.Println()

		// Re-create lexer for parsing
		lex = lexer.NewLexer(strings.NewReader(input))
	}

	// Parsing phase
	parse := parser.NewParser(lex)
	program, err := parse.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	if debugMode {
		fmt.Println("=== AST ===")
		fmt.Println(parser.PrintAST(program, 0))
		fmt.Println()

		fmt.Println("=== STACK SIMULATION ===")
		fmt.Println(parser.SimulateStack(program))
		fmt.Println()

		fmt.Println("=== EXECUTION ===")
	}

	// Execute the program
	runtime := runtime.NewRuntime()
	err = runtime.Execute(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}

	if debugMode {
		fmt.Println()
		fmt.Println("=== FINAL STATE ===")
		runtime.PrintState()
	}
}

func printUsage() {
	fmt.Println("Usage: gorth <filename> [options]")
	fmt.Println("  filename: the name of the .gorth file to execute")
	fmt.Println("  options:")
	fmt.Println("    -d: optional enable debug mode")
}
