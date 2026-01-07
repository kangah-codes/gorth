//go:build legacy_parser_tests
// +build legacy_parser_tests

package parser

import (
	"fmt"
	"gorth/lexer"
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCurr lexer.TokenType
		wantPeek lexer.TokenType
	}{
		{
			name:     "empty input",
			input:    "",
			wantCurr: lexer.EOF,
			wantPeek: lexer.EOF,
		},
		{
			name:     "single token",
			input:    "42",
			wantCurr: lexer.INT,
			wantPeek: lexer.EOF,
		},
		{
			name:     "two tokens",
			input:    "42 +",
			wantCurr: lexer.INT,
			wantPeek: lexer.OP_PLUS,
		},
		{
			name:     "multiple tokens",
			input:    "42 + 10",
			wantCurr: lexer.INT,
			wantPeek: lexer.OP_PLUS,
		},
		{
			name:     "string literal",
			input:    `"hello"`,
			wantCurr: lexer.STRING,
			wantPeek: lexer.EOF,
		},
		{
			name:     "boolean literal",
			input:    "true",
			wantCurr: lexer.BOOL,
			wantPeek: lexer.EOF,
		},
		{
			name:     "identifier",
			input:    "myVar",
			wantCurr: lexer.IDENT,
			wantPeek: lexer.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			if p.lexer != l {
				t.Errorf("NewParser() lexer = %v, want %v", p.lexer, l)
			}

			if p.currentToken != tt.wantCurr {
				t.Errorf("NewParser() currentToken = %v, want %v", p.currentToken, tt.wantCurr)
			}

			if p.peekToken != tt.wantPeek {
				t.Errorf("NewParser() peekToken = %v, want %v", p.peekToken, tt.wantPeek)
			}

			// Verify that parser is properly initialized
			if p.loopDepth != 0 {
				t.Errorf("NewParser() loopDepth = %v, want 0", p.loopDepth)
			}
		})
	}
}
func TestCurrentTokenIs(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		checkToken lexer.TokenType
		want       bool
	}{
		{
			name:       "current token matches",
			input:      "42",
			checkToken: lexer.INT,
			want:       true,
		},
		{
			name:       "current token does not match",
			input:      "42",
			checkToken: lexer.STRING,
			want:       false,
		},
		{
			name:       "EOF token matches",
			input:      "",
			checkToken: lexer.EOF,
			want:       true,
		},
		{
			name:       "EOF token does not match",
			input:      "",
			checkToken: lexer.INT,
			want:       false,
		},
		{
			name:       "string literal matches",
			input:      `"hello"`,
			checkToken: lexer.STRING,
			want:       true,
		},
		{
			name:       "boolean literal matches",
			input:      "true",
			checkToken: lexer.BOOL,
			want:       true,
		},
		{
			name:       "identifier matches",
			input:      "myVar",
			checkToken: lexer.IDENT,
			want:       true,
		},
		{
			name:       "operator matches",
			input:      "+ 10",
			checkToken: lexer.OP_PLUS,
			want:       true,
		},
		{
			name:       "operator does not match",
			input:      "+ 10",
			checkToken: lexer.OP_SUBTRACT,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			got := p.currentTokenIs(tt.checkToken)
			if got != tt.want {
				t.Errorf("currentTokenIs(%v) = %v, want %v", tt.checkToken, got, tt.want)
			}
		})
	}
}
func TestPeekTokenIs(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		checkToken lexer.TokenType
		want       bool
	}{
		{
			name:       "peek token matches",
			input:      "42 +",
			checkToken: lexer.OP_PLUS,
			want:       true,
		},
		{
			name:       "peek token does not match",
			input:      "42 +",
			checkToken: lexer.STRING,
			want:       false,
		},
		{
			name:       "EOF peek token matches",
			input:      "42",
			checkToken: lexer.EOF,
			want:       true,
		},
		{
			name:       "EOF peek token does not match",
			input:      "42",
			checkToken: lexer.INT,
			want:       false,
		},
		{
			name:       "string literal peek matches",
			input:      `42 "hello"`,
			checkToken: lexer.STRING,
			want:       true,
		},
		{
			name:       "boolean literal peek matches",
			input:      "42 true",
			checkToken: lexer.BOOL,
			want:       true,
		},
		{
			name:       "identifier peek matches",
			input:      "42 myVar",
			checkToken: lexer.IDENT,
			want:       true,
		},
		{
			name:       "operator peek matches",
			input:      "42 -",
			checkToken: lexer.OP_SUBTRACT,
			want:       true,
		},
		{
			name:       "operator peek does not match",
			input:      "42 -",
			checkToken: lexer.OP_PLUS,
			want:       false,
		},
		{
			name:       "multiple tokens - check second",
			input:      "42 + 10",
			checkToken: lexer.OP_PLUS,
			want:       true,
		},
		{
			name:       "empty input peek",
			input:      "",
			checkToken: lexer.EOF,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			got := p.peekTokenIs(tt.checkToken)
			if got != tt.want {
				t.Errorf("peekTokenIs(%v) = %v, want %v", tt.checkToken, got, tt.want)
			}
		})
	}
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantCurrents   []lexer.TokenType
		wantPeeks      []lexer.TokenType
		wantLiterals   []string
		wantError      bool
		nextTokenCalls int
	}{
		{
			name:           "single token progression",
			input:          "42",
			wantCurrents:   []lexer.TokenType{lexer.EOF},
			wantPeeks:      []lexer.TokenType{lexer.EOF},
			wantLiterals:   []string{""},
			wantError:      false,
			nextTokenCalls: 1,
		},
		{
			name:           "two token progression",
			input:          "42 +",
			wantCurrents:   []lexer.TokenType{lexer.OP_PLUS, lexer.EOF},
			wantPeeks:      []lexer.TokenType{lexer.EOF, lexer.EOF},
			wantLiterals:   []string{"+", ""},
			wantError:      false,
			nextTokenCalls: 2,
		},
		{
			name:           "multiple token progression",
			input:          "42 + 10",
			wantCurrents:   []lexer.TokenType{lexer.OP_PLUS, lexer.INT, lexer.EOF},
			wantPeeks:      []lexer.TokenType{lexer.INT, lexer.EOF, lexer.EOF},
			wantLiterals:   []string{"+", "10", ""},
			wantError:      false,
			nextTokenCalls: 3,
		},
		{
			name:           "string literal progression",
			input:          `"hello" world`,
			wantCurrents:   []lexer.TokenType{lexer.IDENT, lexer.EOF},
			wantPeeks:      []lexer.TokenType{lexer.EOF, lexer.EOF},
			wantLiterals:   []string{"world", ""},
			wantError:      false,
			nextTokenCalls: 2,
		},
		{
			name:           "boolean and identifier",
			input:          "true myVar",
			wantCurrents:   []lexer.TokenType{lexer.IDENT, lexer.EOF},
			wantPeeks:      []lexer.TokenType{lexer.EOF, lexer.EOF},
			wantLiterals:   []string{"myVar", ""},
			wantError:      false,
			nextTokenCalls: 2,
		},
		{
			name:           "operators sequence",
			input:          "+ - *",
			wantCurrents:   []lexer.TokenType{lexer.OP_SUBTRACT, lexer.OP_MULTIPLY, lexer.EOF},
			wantPeeks:      []lexer.TokenType{lexer.OP_MULTIPLY, lexer.EOF, lexer.EOF},
			wantLiterals:   []string{"-", "*", ""},
			wantError:      false,
			nextTokenCalls: 3,
		},
		{
			name:           "empty input progression",
			input:          "",
			wantCurrents:   []lexer.TokenType{lexer.EOF},
			wantPeeks:      []lexer.TokenType{lexer.EOF},
			wantLiterals:   []string{""},
			wantError:      false,
			nextTokenCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			for i := 0; i < tt.nextTokenCalls; i++ {
				err := p.nextToken()

				if tt.wantError && err == nil {
					t.Errorf("nextToken() call %d expected error, but got none", i+1)
					continue
				}
				if !tt.wantError && err != nil {
					t.Errorf("nextToken() call %d unexpected error: %v", i+1, err)
					continue
				}

				if i < len(tt.wantCurrents) {
					if p.currentToken != tt.wantCurrents[i] {
						t.Errorf("nextToken() call %d currentToken = %v, want %v", i+1, p.currentToken, tt.wantCurrents[i])
					}
				}

				if i < len(tt.wantPeeks) {
					if p.peekToken != tt.wantPeeks[i] {
						t.Errorf("nextToken() call %d peekToken = %v, want %v", i+1, p.peekToken, tt.wantPeeks[i])
					}
				}

				if i < len(tt.wantLiterals) {
					if p.currentLiteral != tt.wantLiterals[i] {
						t.Errorf("nextToken() call %d currentLiteral = %q, want %q", i+1, p.currentLiteral, tt.wantLiterals[i])
					}
				}
			}
		})
	}
}

func TestNextTokenPositionTracking(t *testing.T) {
	input := "42\n+ 10"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Initial state: current=42, peek=+
	initialPos := p.currentPosition

	// Call nextToken to advance: current=+, peek=10
	err := p.nextToken()
	if err != nil {
		t.Fatalf("nextToken() unexpected error: %v", err)
	}

	// Verify position was updated correctly
	if p.currentPosition.Line == initialPos.Line && p.currentPosition.Column == initialPos.Column {
		t.Error("nextToken() did not update currentPosition")
	}

	// Verify current token is now the operator
	if p.currentToken != lexer.OP_PLUS {
		t.Errorf("nextToken() currentToken = %v, want %v", p.currentToken, lexer.OP_PLUS)
	}

	// Verify peek token is now the number
	if p.peekToken != lexer.INT {
		t.Errorf("nextToken() peekToken = %v, want %v", p.peekToken, lexer.INT)
	}
}
func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantStmtCount int
		wantError     bool
		errorContains string
	}{
		{
			name:          "empty input",
			input:         "",
			wantStmtCount: 0,
			wantError:     false,
		},
		{
			name:          "single integer literal",
			input:         "42",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "multiple literals",
			input:         `42 "hello" true`,
			wantStmtCount: 3,
			wantError:     false,
		},
		{
			name:          "binary operations",
			input:         "42 10 +",
			wantStmtCount: 3,
			wantError:     false,
		},
		{
			name:          "variable declaration",
			input:         "VAR x",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "constant definition",
			input:         "CONST PI 3.14",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "assignment",
			input:         "42 := VAR x",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "array literal",
			input:         "[1, 2, 3]",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "stack operations",
			input:         "42 10 DUP SWAP",
			wantStmtCount: 4,
			wantError:     false,
		},
		{
			name:          "io operations",
			input:         "42 DUMP",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "while loop",
			input:         "WHILE 42 END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "do-while loop",
			input:         "DO 42 WHILE true END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "do-if statement",
			input:         "DO 42 IF true END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "do-if-else statement",
			input:         "DO 42 ELSE 10 IF true END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "break statement in while loop",
			input:         "WHILE BREAK END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "continue statement in while loop",
			input:         "WHILE CONTINUE END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "illegal token error",
			input:         "42 @",
			wantStmtCount: 0,
			wantError:     true,
			errorContains: "illegal token",
		},
		{
			name:          "unexpected token error",
			input:         "42 {",
			wantStmtCount: 0,
			wantError:     true,
			errorContains: "unexpected token",
		},
		{
			name:          "break outside loop error",
			input:         "BREAK",
			wantStmtCount: 0,
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue outside loop error",
			input:         "CONTINUE",
			wantStmtCount: 0,
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "complex program",
			input:         `VAR x 42 := x CONST PI 3.14 [1, 2, 3] DUMP`,
			wantStmtCount: 6,
			wantError:     false,
		},
		{
			name:          "nested loops with break",
			input:         "WHILE WHILE BREAK END END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "mixed data types",
			input:         `42 3.14 "hello" true false NULL`,
			wantStmtCount: 6,
			wantError:     false,
		},
		{
			name:          "all unary operators",
			input:         "42 ! ++ --",
			wantStmtCount: 4,
			wantError:     false,
		},
		{
			name:          "all binary operators",
			input:         "1 2 + 3 4 - 5 6 * 7 8 / 9 10 % 11 12 ** 1 2 == 3 4 != 5 6 > 7 8 < 9 10 >= 11 12 <= true false && true false ||",
			wantStmtCount: 30,
			wantError:     false,
		},
		{
			name:          "all stack operations",
			input:         "1 2 DROP SWAP DUP OVER ROT DEL CLEAR 5 PICK",
			wantStmtCount: 10,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()

			if tt.wantError {
				if err == nil {
					t.Errorf("Parse() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Parse() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if program == nil {
				t.Error("Parse() returned nil program")
				return
			}

			if len(program.Statements) != tt.wantStmtCount {
				t.Errorf("Parse() statement count = %d, want %d", len(program.Statements), tt.wantStmtCount)
			}
		})
	}
}

func TestParseStatementTypes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  string
		wantError bool
	}{
		{
			name:      "integer literal",
			input:     "42",
			wantType:  "*parser.IntLiteral",
			wantError: false,
		},
		{
			name:      "float literal",
			input:     "3.14",
			wantType:  "*parser.FloatLiteral",
			wantError: false,
		},
		{
			name:      "string literal",
			input:     `"hello"`,
			wantType:  "*parser.StringLiteral",
			wantError: false,
		},
		{
			name:      "boolean literal",
			input:     "true",
			wantType:  "*parser.BoolLiteral",
			wantError: false,
		},
		{
			name:      "identifier",
			input:     "myVar",
			wantType:  "*parser.Identifier",
			wantError: false,
		},
		{
			name:      "null literal",
			input:     "NULL",
			wantType:  "*parser.NullLiteral",
			wantError: false,
		},
		{
			name:      "binary expression",
			input:     "+",
			wantType:  "*parser.BinaryExpression",
			wantError: false,
		},
		{
			name:      "unary expression",
			input:     "!",
			wantType:  "*parser.UnaryExpression",
			wantError: false,
		},
		{
			name:      "stack operation",
			input:     "DUP",
			wantType:  "*parser.StackOp",
			wantError: false,
		},
		{
			name:      "io statement",
			input:     "DUMP",
			wantType:  "*parser.IOStmt",
			wantError: false,
		},
		{
			name:      "variable declaration",
			input:     "VAR x",
			wantType:  "*parser.VarDeclaration",
			wantError: false,
		},
		{
			name:      "constant declaration",
			input:     "CONST PI 3.14",
			wantType:  "*parser.ConstDeclaration",
			wantError: false,
		},
		{
			name:      "assignment",
			input:     ":= VAR x",
			wantType:  "*parser.Assignment",
			wantError: false,
		},
		{
			name:      "array literal",
			input:     "[1, 2, 3]",
			wantType:  "*parser.ArrayLiteral",
			wantError: false,
		},
		{
			name:      "while statement",
			input:     "WHILE 1 END",
			wantType:  "*parser.WhileStmt",
			wantError: false,
		},
		{
			name:      "if statement",
			input:     "DO 1 IF true END",
			wantType:  "*parser.IfStmt",
			wantError: false,
		},
		{
			name:      "break statement",
			input:     "WHILE BREAK END",
			wantType:  "*parser.BreakStmt",
			wantError: false,
		},
		{
			name:      "continue statement",
			input:     "WHILE CONTINUE END",
			wantType:  "*parser.ContinueStmt",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()

			if tt.wantError {
				if err == nil {
					t.Errorf("Parse() expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if program == nil || len(program.Statements) == 0 {
				t.Error("Parse() returned no statements")
				return
			}

			// For nested structures like WHILE, find the target statement type
			var targetStmt Node
			if tt.wantType == "*parser.BreakStmt" || tt.wantType == "*parser.ContinueStmt" {
				if whileStmt, ok := program.Statements[0].(*WhileStmt); ok && len(whileStmt.Body) > 0 {
					targetStmt = whileStmt.Body[0]
				}
			} else {
				targetStmt = program.Statements[0]
			}

			gotType := fmt.Sprintf("%T", targetStmt)
			if gotType != tt.wantType {
				t.Errorf("Parse() statement type = %s, want %s", gotType, tt.wantType)
			}
		})
	}
}
func TestParseStatement(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  string
		wantError bool
		errorMsg  string
	}{
		// Literal types
		{
			name:     "integer literal",
			input:    "42",
			wantType: "*parser.IntLiteral",
		},
		{
			name:     "float literal",
			input:    "3.14",
			wantType: "*parser.FloatLiteral",
		},
		{
			name:     "string literal",
			input:    `"hello world"`,
			wantType: "*parser.StringLiteral",
		},
		{
			name:     "boolean literal true",
			input:    "true",
			wantType: "*parser.BoolLiteral",
		},
		{
			name:     "boolean literal false",
			input:    "false",
			wantType: "*parser.BoolLiteral",
		},
		{
			name:     "null literal",
			input:    "NULL",
			wantType: "*parser.NullLiteral",
		},
		{
			name:     "identifier",
			input:    "myVariable",
			wantType: "*parser.Identifier",
		},

		// Variable and constant declarations
		{
			name:     "variable declaration",
			input:    "VAR x",
			wantType: "*parser.VarDeclaration",
		},
		{
			name:     "constant definition with int",
			input:    "CONST PI 3",
			wantType: "*parser.ConstDeclaration",
		},
		{
			name:     "constant definition with float",
			input:    "CONST PI 3.14",
			wantType: "*parser.ConstDeclaration",
		},
		{
			name:     "constant definition with string",
			input:    `CONST MSG "hello"`,
			wantType: "*parser.ConstDeclaration",
		},
		{
			name:     "constant definition with bool",
			input:    "CONST FLAG true",
			wantType: "*parser.ConstDeclaration",
		},
		{
			name:     "constant definition with null",
			input:    "CONST EMPTY NULL",
			wantType: "*parser.ConstDeclaration",
		},

		// Binary operators
		{
			name:     "addition operator",
			input:    "+",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "subtraction operator",
			input:    "-",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "multiplication operator",
			input:    "*",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "division operator",
			input:    "/",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "power operator",
			input:    "**",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "modulo operator",
			input:    "%",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "equality operator",
			input:    "==",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "not equal operator",
			input:    "!=",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "greater than operator",
			input:    ">",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "less than operator",
			input:    "<",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "greater than or equal operator",
			input:    ">=",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "less than or equal operator",
			input:    "<=",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "logical and operator",
			input:    "&&",
			wantType: "*parser.BinaryExpression",
		},
		{
			name:     "logical or operator",
			input:    "||",
			wantType: "*parser.BinaryExpression",
		},

		// Unary operators
		{
			name:     "not operator",
			input:    "!",
			wantType: "*parser.UnaryExpression",
		},
		{
			name:     "increment operator",
			input:    "++",
			wantType: "*parser.UnaryExpression",
		},
		{
			name:     "decrement operator",
			input:    "--",
			wantType: "*parser.UnaryExpression",
		},

		// Stack operations
		{
			name:     "drop operation",
			input:    "DROP",
			wantType: "*parser.StackOp",
		},
		{
			name:     "swap operation",
			input:    "SWAP",
			wantType: "*parser.StackOp",
		},
		{
			name:     "dup operation",
			input:    "DUP",
			wantType: "*parser.StackOp",
		},
		{
			name:     "over operation",
			input:    "OVER",
			wantType: "*parser.StackOp",
		},
		{
			name:     "rot operation",
			input:    "ROT",
			wantType: "*parser.StackOp",
		},
		{
			name:     "del operation",
			input:    "DEL",
			wantType: "*parser.StackOp",
		},
		{
			name:     "clear operation",
			input:    "CLEAR",
			wantType: "*parser.StackOp",
		},
		{
			name:     "pick operation",
			input:    "PICK",
			wantType: "*parser.StackOp",
		},

		// IO operations
		{
			name:     "dump operation",
			input:    "DUMP",
			wantType: "*parser.IOStmt",
		},
		{
			name:     "dumpln operation",
			input:    "DUMPLN",
			wantType: "*parser.IOStmt",
		},

		// Assignment
		{
			name:     "assignment to variable",
			input:    ":= VAR x",
			wantType: "*parser.Assignment",
		},
		{
			name:     "assignment to constant",
			input:    ":= CONST x",
			wantType: "*parser.Assignment",
		},
		{
			name:     "assignment to identifier",
			input:    ":= myVar",
			wantType: "*parser.Assignment",
		},

		// Array literal
		{
			name:     "empty array",
			input:    "[]",
			wantType: "*parser.ArrayLiteral",
		},
		{
			name:     "array with integers",
			input:    "[1, 2, 3]",
			wantType: "*parser.ArrayLiteral",
		},
		{
			name:     "array with mixed types",
			input:    `[1, "hello", true]`,
			wantType: "*parser.ArrayLiteral",
		},

		// Control flow
		{
			name:     "do statement",
			input:    "DO 42 IF true END",
			wantType: "*parser.IfStmt",
		},
		{
			name:     "while statement",
			input:    "WHILE 1 END",
			wantType: "*parser.WhileStmt",
		},
		{
			name:     "break statement",
			input:    "BREAK",
			wantType: "*parser.BreakStmt",
		},
		{
			name:     "continue statement",
			input:    "CONTINUE",
			wantType: "*parser.ContinueStmt",
		},

		// Error cases
		{
			name:      "illegal token",
			input:     "@",
			wantError: true,
			errorMsg:  "illegal token",
		},
		{
			name:      "unexpected token",
			input:     "{",
			wantError: true,
			errorMsg:  "unexpected token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseStatement()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseStatement() expected error, but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("parseStatement() error = %v, want error containing %q", err, tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("parseStatement() unexpected error: %v", err)
				return
			}

			if stmt == nil {
				t.Error("parseStatement() returned nil statement")
				return
			}

			gotType := fmt.Sprintf("%T", stmt)
			if gotType != tt.wantType {
				t.Errorf("parseStatement() statement type = %s, want %s", gotType, tt.wantType)
			}
		})
	}
}

func TestParseStatementPosition(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "integer with position",
			input: "42",
		},
		{
			name:  "string with position",
			input: `"hello"`,
		},
		{
			name:  "operator with position",
			input: "+",
		},
		{
			name:  "break with position",
			input: "BREAK",
		},
		{
			name:  "continue with position",
			input: "CONTINUE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseStatement()
			if err != nil {
				t.Errorf("parseStatement() unexpected error: %v", err)
				return
			}

			// Check if statement has position information
			switch s := stmt.(type) {
			case *IntLiteral:
				if s.Position.Line == 0 && s.Position.Column == 0 {
					t.Error("parseStatement() IntLiteral missing position information")
				}
			case *StringLiteral:
				if s.Position.Line == 0 && s.Position.Column == 0 {
					t.Error("parseStatement() StringLiteral missing position information")
				}
			case *BinaryExpression:
				if s.Position.Line == 0 && s.Position.Column == 0 {
					t.Error("parseStatement() BinaryExpression missing position information")
				}
			case *BreakStmt:
				if s.Pos.Line == 0 && s.Pos.Column == 0 {
					t.Error("parseStatement() BreakStmt missing position information")
				}
			case *ContinueStmt:
				if s.Pos.Line == 0 && s.Pos.Column == 0 {
					t.Error("parseStatement() ContinueStmt missing position information")
				}
			}
		})
	}
}

func TestParseStatementValues(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue string
	}{
		{
			name:      "integer value",
			input:     "123",
			wantValue: "123",
		},
		{
			name:      "float value",
			input:     "3.14159",
			wantValue: "3.14159",
		},
		{
			name:      "string value",
			input:     `"test string"`,
			wantValue: "test string",
		},
		{
			name:      "boolean true value",
			input:     "true",
			wantValue: "true",
		},
		{
			name:      "boolean false value",
			input:     "false",
			wantValue: "false",
		},
		{
			name:      "null value",
			input:     "NULL",
			wantValue: "NULL",
		},
		{
			name:      "identifier value",
			input:     "myIdentifier",
			wantValue: "myIdentifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseStatement()
			if err != nil {
				t.Errorf("parseStatement() unexpected error: %v", err)
				return
			}

			var gotValue string
			switch s := stmt.(type) {
			case *IntLiteral:
				gotValue = s.Value
			case *FloatLiteral:
				gotValue = s.Value
			case *StringLiteral:
				gotValue = s.Value
			case *BoolLiteral:
				gotValue = s.Value
			case *NullLiteral:
				gotValue = s.Value
			case *Identifier:
				gotValue = s.Value
			default:
				t.Errorf("parseStatement() unexpected statement type: %T", s)
				return
			}

			if gotValue != tt.wantValue {
				t.Errorf("parseStatement() value = %q, want %q", gotValue, tt.wantValue)
			}
		})
	}
}
func TestParseVarDeclaration(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantName      string
		wantError     bool
		errorContains string
	}{
		{
			name:      "valid variable declaration",
			input:     "VAR myVar",
			wantName:  "myVar",
			wantError: false,
		},
		{
			name:      "valid variable with underscore",
			input:     "VAR my_var",
			wantName:  "my_var",
			wantError: false,
		},
		{
			name:      "valid variable with numbers",
			input:     "VAR var123",
			wantName:  "var123",
			wantError: false,
		},
		{
			name:      "valid single letter variable",
			input:     "VAR x",
			wantName:  "x",
			wantError: false,
		},
		{
			name:      "valid long variable name",
			input:     "VAR veryLongVariableName",
			wantName:  "veryLongVariableName",
			wantError: false,
		},
		{
			name:          "missing identifier after VAR",
			input:         "VAR",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "integer after VAR",
			input:         "VAR 42",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "string after VAR",
			input:         "VAR \"hello\"",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "boolean after VAR",
			input:         "VAR true",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "operator after VAR",
			input:         "VAR +",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "keyword after VAR",
			input:         "VAR WHILE",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "null after VAR",
			input:         "VAR NULL",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "EOF after VAR",
			input:         "VAR",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			// Ensure we start at VAR token
			if !p.currentTokenIs(lexer.VAR) {
				t.Fatalf("Test setup error: expected VAR token, got %v", p.currentToken)
			}

			stmt, err := p.parseVarDeclaration()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseVarDeclaration() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("parseVarDeclaration() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseVarDeclaration() unexpected error: %v", err)
				return
			}

			if stmt == nil {
				t.Error("parseVarDeclaration() returned nil statement")
				return
			}

			varDecl, ok := stmt.(*VarDeclaration)
			if !ok {
				t.Errorf("parseVarDeclaration() returned %T, want *VarDeclaration", stmt)
				return
			}

			if varDecl.Name != tt.wantName {
				t.Errorf("parseVarDeclaration() Name = %q, want %q", varDecl.Name, tt.wantName)
			}

			// Verify position is set
			if varDecl.Pos.Line == 0 && varDecl.Pos.Column == 0 {
				t.Error("parseVarDeclaration() position not set")
			}
		})
	}
}

func TestParseVarDeclarationPosition(t *testing.T) {
	input := "VAR test"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Capture the initial position
	initialPos := p.currentPosition

	stmt, err := p.parseVarDeclaration()
	if err != nil {
		t.Fatalf("parseVarDeclaration() unexpected error: %v", err)
	}

	varDecl := stmt.(*VarDeclaration)

	// Position should be set to the VAR token position
	if varDecl.Pos.Line != initialPos.Line || varDecl.Pos.Column != initialPos.Column {
		t.Errorf("parseVarDeclaration() position = %v, want %v", varDecl.Pos, initialPos)
	}
}

func TestParseVarDeclarationTokenAdvancement(t *testing.T) {
	input := "VAR myVar 42"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Initially at VAR
	if !p.currentTokenIs(lexer.VAR) {
		t.Fatalf("Test setup error: expected VAR token, got %v", p.currentToken)
	}

	_, err := p.parseVarDeclaration()
	if err != nil {
		t.Fatalf("parseVarDeclaration() unexpected error: %v", err)
	}

	// After parsing, current token should be at the identifier
	if !p.currentTokenIs(lexer.IDENT) {
		t.Errorf("parseVarDeclaration() current token = %v, want %v", p.currentToken, lexer.IDENT)
	}

	if p.currentLiteral != "myVar" {
		t.Errorf("parseVarDeclaration() current literal = %q, want %q", p.currentLiteral, "myVar")
	}

	// Peek token should be the next token (42)
	if !p.peekTokenIs(lexer.INT) {
		t.Errorf("parseVarDeclaration() peek token = %v, want %v", p.peekToken, lexer.INT)
	}
}
func TestParseConstDeclaration(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantName      string
		wantError     bool
		errorContains string
	}{
		{
			name:      "valid const declaration",
			input:     "CONST myConst",
			wantName:  "myConst",
			wantError: false,
		},
		{
			name:      "valid const with underscore",
			input:     "CONST my_const",
			wantName:  "my_const",
			wantError: false,
		},
		{
			name:      "valid const with numbers",
			input:     "CONST const123",
			wantName:  "const123",
			wantError: false,
		},
		{
			name:      "valid single letter const",
			input:     "CONST x",
			wantName:  "x",
			wantError: false,
		},
		{
			name:      "valid long const name",
			input:     "CONST veryLongConstantName",
			wantName:  "veryLongConstantName",
			wantError: false,
		},
		{
			name:          "missing identifier after CONST",
			input:         "CONST",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "integer after CONST",
			input:         "CONST 42",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "string after CONST",
			input:         "CONST \"hello\"",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "boolean after CONST",
			input:         "CONST true",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "operator after CONST",
			input:         "CONST +",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "keyword after CONST",
			input:         "CONST WHILE",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "null after CONST",
			input:         "CONST NULL",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "EOF after CONST",
			input:         "CONST",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			// Ensure we start at CONST token
			if !p.currentTokenIs(lexer.CONST) {
				t.Fatalf("Test setup error: expected CONST token, got %v", p.currentToken)
			}

			stmt, err := p.parseConstDeclaration()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseConstDeclaration() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("parseConstDeclaration() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseConstDeclaration() unexpected error: %v", err)
				return
			}

			if stmt == nil {
				t.Error("parseConstDeclaration() returned nil statement")
				return
			}

			constDecl, ok := stmt.(*ConstDeclaration)
			if !ok {
				t.Errorf("parseConstDeclaration() returned %T, want *ConstDeclaration", stmt)
				return
			}

			if constDecl.Name != tt.wantName {
				t.Errorf("parseConstDeclaration() Name = %q, want %q", constDecl.Name, tt.wantName)
			}

			// Verify Value is nil for inline target form
			if constDecl.Value != nil {
				t.Errorf("parseConstDeclaration() Value = %v, want nil (inline target form)", constDecl.Value)
			}

			// Verify position is set
			if constDecl.Pos.Line == 0 && constDecl.Pos.Column == 0 {
				t.Error("parseConstDeclaration() position not set")
			}
		})
	}
}

func TestParseConstDeclarationPosition(t *testing.T) {
	input := "CONST test"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Capture the initial position
	initialPos := p.currentPosition

	stmt, err := p.parseConstDeclaration()
	if err != nil {
		t.Fatalf("parseConstDeclaration() unexpected error: %v", err)
	}

	constDecl := stmt.(*ConstDeclaration)

	// Position should be set to the CONST token position
	if constDecl.Pos.Line != initialPos.Line || constDecl.Pos.Column != initialPos.Column {
		t.Errorf("parseConstDeclaration() position = %v, want %v", constDecl.Pos, initialPos)
	}
}

func TestParseConstDeclarationTokenAdvancement(t *testing.T) {
	input := "CONST myConst 42"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Initially at CONST
	if !p.currentTokenIs(lexer.CONST) {
		t.Fatalf("Test setup error: expected CONST token, got %v", p.currentToken)
	}

	_, err := p.parseConstDeclaration()
	if err != nil {
		t.Fatalf("parseConstDeclaration() unexpected error: %v", err)
	}

	// After parsing, current token should be at the identifier
	if !p.currentTokenIs(lexer.IDENT) {
		t.Errorf("parseConstDeclaration() current token = %v, want %v", p.currentToken, lexer.IDENT)
	}

	if p.currentLiteral != "myConst" {
		t.Errorf("parseConstDeclaration() current literal = %q, want %q", p.currentLiteral, "myConst")
	}

	// Peek token should be the next token (42)
	if !p.peekTokenIs(lexer.INT) {
		t.Errorf("parseConstDeclaration() peek token = %v, want %v", p.peekToken, lexer.INT)
	}
}

func TestParseConstDeclarationInlineTarget(t *testing.T) {
	input := "CONST PI"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	stmt, err := p.parseConstDeclaration()
	if err != nil {
		t.Fatalf("parseConstDeclaration() unexpected error: %v", err)
	}

	constDecl := stmt.(*ConstDeclaration)

	// Verify this is an inline target form (no value)
	if constDecl.Value != nil {
		t.Errorf("parseConstDeclaration() inline target should have nil Value, got %v", constDecl.Value)
	}

	if constDecl.Name != "PI" {
		t.Errorf("parseConstDeclaration() Name = %q, want %q", constDecl.Name, "PI")
	}
}
func TestParseAssignment(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantTargetType string
		wantTargetName string
		wantError      bool
		errorContains  string
	}{
		{
			name:           "assignment to variable declaration",
			input:          ":= VAR myVar",
			wantTargetType: "*parser.VarDeclaration",
			wantTargetName: "myVar",
			wantError:      false,
		},
		{
			name:           "assignment to const declaration",
			input:          ":= CONST myConst",
			wantTargetType: "*parser.ConstDeclaration",
			wantTargetName: "myConst",
			wantError:      false,
		},
		{
			name:           "assignment to identifier",
			input:          ":= myIdent",
			wantTargetType: "*parser.Identifier",
			wantTargetName: "myIdent",
			wantError:      false,
		},
		{
			name:           "assignment to variable with underscore",
			input:          ":= VAR my_var",
			wantTargetType: "*parser.VarDeclaration",
			wantTargetName: "my_var",
			wantError:      false,
		},
		{
			name:           "assignment to const with numbers",
			input:          ":= CONST const123",
			wantTargetType: "*parser.ConstDeclaration",
			wantTargetName: "const123",
			wantError:      false,
		},
		{
			name:           "assignment to single letter identifier",
			input:          ":= x",
			wantTargetType: "*parser.Identifier",
			wantTargetName: "x",
			wantError:      false,
		},
		{
			name:           "assignment to long identifier",
			input:          ":= veryLongIdentifierName",
			wantTargetType: "*parser.Identifier",
			wantTargetName: "veryLongIdentifierName",
			wantError:      false,
		},
		{
			name:          "missing target after :=",
			input:         ":=",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "integer after :=",
			input:         ":= 42",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "string after :=",
			input:         ":= \"hello\"",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "boolean after :=",
			input:         ":= true",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "operator after :=",
			input:         ":= +",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "null after :=",
			input:         ":= NULL",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "keyword after :=",
			input:         ":= WHILE",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "bracket after :=",
			input:         ":= [",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "invalid var declaration",
			input:         ":= VAR",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "invalid const declaration",
			input:         ":= CONST",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
		{
			name:          "var declaration with invalid identifier",
			input:         ":= VAR 123",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "const declaration with invalid identifier",
			input:         ":= CONST true",
			wantError:     true,
			errorContains: "expected identifier after CONST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			// Ensure we start at OP_ASSIGN token
			if !p.currentTokenIs(lexer.OP_ASSIGN) {
				t.Fatalf("Test setup error: expected OP_ASSIGN token, got %v", p.currentToken)
			}

			stmt, err := p.parseAssignment()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseAssignment() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("parseAssignment() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseAssignment() unexpected error: %v", err)
				return
			}

			if stmt == nil {
				t.Error("parseAssignment() returned nil statement")
				return
			}

			assignment, ok := stmt.(*Assignment)
			if !ok {
				t.Errorf("parseAssignment() returned %T, want *Assignment", stmt)
				return
			}

			// Check target type
			gotTargetType := fmt.Sprintf("%T", assignment.Target)
			if gotTargetType != tt.wantTargetType {
				t.Errorf("parseAssignment() target type = %s, want %s", gotTargetType, tt.wantTargetType)
			}

			// Check target name
			var gotTargetName string
			switch target := assignment.Target.(type) {
			case *VarDeclaration:
				gotTargetName = target.Name
			case *ConstDeclaration:
				gotTargetName = target.Name
			case *Identifier:
				gotTargetName = target.Value
			}

			if gotTargetName != tt.wantTargetName {
				t.Errorf("parseAssignment() target name = %q, want %q", gotTargetName, tt.wantTargetName)
			}

			// Verify position is set
			if assignment.Pos.Line == 0 && assignment.Pos.Column == 0 {
				t.Error("parseAssignment() position not set")
			}
		})
	}
}

func TestParseAssignmentPosition(t *testing.T) {
	input := ":= myVar"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Capture the initial position
	initialPos := p.currentPosition

	stmt, err := p.parseAssignment()
	if err != nil {
		t.Fatalf("parseAssignment() unexpected error: %v", err)
	}

	assignment := stmt.(*Assignment)

	// Position should be set to the OP_ASSIGN token position
	if assignment.Pos.Line != initialPos.Line || assignment.Pos.Column != initialPos.Column {
		t.Errorf("parseAssignment() position = %v, want %v", assignment.Pos, initialPos)
	}
}

func TestParseAssignmentTargetDetails(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		checkValue func(*testing.T, *Assignment)
	}{
		{
			name:  "identifier target details",
			input: ":= testIdent",
			checkValue: func(t *testing.T, a *Assignment) {
				ident, ok := a.Target.(*Identifier)
				if !ok {
					t.Errorf("Expected *Identifier, got %T", a.Target)
					return
				}
				if ident.Value != "testIdent" {
					t.Errorf("Identifier.Value = %q, want %q", ident.Value, "testIdent")
				}
				if ident.Token != lexer.IDENT {
					t.Errorf("Identifier.Token = %v, want %v", ident.Token, lexer.IDENT)
				}
				if ident.Position.Line == 0 && ident.Position.Column == 0 {
					t.Error("Identifier position not set")
				}
			},
		},
		{
			name:  "var declaration target details",
			input: ":= VAR testVar",
			checkValue: func(t *testing.T, a *Assignment) {
				varDecl, ok := a.Target.(*VarDeclaration)
				if !ok {
					t.Errorf("Expected *VarDeclaration, got %T", a.Target)
					return
				}
				if varDecl.Name != "testVar" {
					t.Errorf("VarDeclaration.Name = %q, want %q", varDecl.Name, "testVar")
				}
				if varDecl.Pos.Line == 0 && varDecl.Pos.Column == 0 {
					t.Error("VarDeclaration position not set")
				}
			},
		},
		{
			name:  "const declaration target details",
			input: ":= CONST testConst",
			checkValue: func(t *testing.T, a *Assignment) {
				constDecl, ok := a.Target.(*ConstDeclaration)
				if !ok {
					t.Errorf("Expected *ConstDeclaration, got %T", a.Target)
					return
				}
				if constDecl.Name != "testConst" {
					t.Errorf("ConstDeclaration.Name = %q, want %q", constDecl.Name, "testConst")
				}
				if constDecl.Value != nil {
					t.Errorf("ConstDeclaration.Value = %v, want nil (inline target)", constDecl.Value)
				}
				if constDecl.Pos.Line == 0 && constDecl.Pos.Column == 0 {
					t.Error("ConstDeclaration position not set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseAssignment()
			if err != nil {
				t.Fatalf("parseAssignment() unexpected error: %v", err)
			}

			assignment := stmt.(*Assignment)
			tt.checkValue(t, assignment)
		})
	}
}

func TestParseAssignmentTokenAdvancement(t *testing.T) {
	input := ":= myVar 42"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Initially at OP_ASSIGN
	if !p.currentTokenIs(lexer.OP_ASSIGN) {
		t.Fatalf("Test setup error: expected OP_ASSIGN token, got %v", p.currentToken)
	}

	_, err := p.parseAssignment()
	if err != nil {
		t.Fatalf("parseAssignment() unexpected error: %v", err)
	}

	// After parsing, current token should be at the identifier
	if !p.currentTokenIs(lexer.IDENT) {
		t.Errorf("parseAssignment() current token = %v, want %v", p.currentToken, lexer.IDENT)
	}

	if p.currentLiteral != "myVar" {
		t.Errorf("parseAssignment() current literal = %q, want %q", p.currentLiteral, "myVar")
	}

	// Peek token should be the next token (42)
	if !p.peekTokenIs(lexer.INT) {
		t.Errorf("parseAssignment() peek token = %v, want %v", p.peekToken, lexer.INT)
	}
}

func TestParseAssignmentInContext(t *testing.T) {
	// Test assignment as part of a complete parse
	tests := []struct {
		name          string
		input         string
		wantStmtCount int
		wantError     bool
	}{
		{
			name:          "value assignment to var",
			input:         "42 := VAR x",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "value assignment to const",
			input:         "3.14 := CONST PI",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "value assignment to identifier",
			input:         "100 := existing",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "multiple assignments",
			input:         "1 := VAR x 2 := VAR y",
			wantStmtCount: 4,
			wantError:     false,
		},
		{
			name:          "complex expression with assignment",
			input:         `"hello" := VAR msg 42 10 + := result`,
			wantStmtCount: 6,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()

			if tt.wantError {
				if err == nil {
					t.Errorf("Parse() expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if len(program.Statements) != tt.wantStmtCount {
				t.Errorf("Parse() statement count = %d, want %d", len(program.Statements), tt.wantStmtCount)
			}

			// Verify at least one assignment exists
			hasAssignment := false
			for _, stmt := range program.Statements {
				if _, ok := stmt.(*Assignment); ok {
					hasAssignment = true
					break
				}
			}

			if !hasAssignment {
				t.Error("Parse() expected at least one Assignment statement")
			}
		})
	}
}
func TestParseArray(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantElements  int
		wantError     bool
		errorContains string
		checkElements func(*testing.T, *ArrayLiteral)
	}{
		{
			name:         "empty array",
			input:        "[]",
			wantElements: 0,
			wantError:    false,
		},
		{
			name:         "single integer element",
			input:        "[42]",
			wantElements: 1,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				if intLit, ok := arr.Elements[0].(*IntLiteral); ok {
					if intLit.Value != "42" {
						t.Errorf("First element value = %q, want %q", intLit.Value, "42")
					}
				} else {
					t.Errorf("First element type = %T, want *IntLiteral", arr.Elements[0])
				}
			},
		},
		{
			name:         "multiple integer elements",
			input:        "[1, 2, 3]",
			wantElements: 3,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				expected := []string{"1", "2", "3"}
				for i, exp := range expected {
					if intLit, ok := arr.Elements[i].(*IntLiteral); ok {
						if intLit.Value != exp {
							t.Errorf("Element %d value = %q, want %q", i, intLit.Value, exp)
						}
					} else {
						t.Errorf("Element %d type = %T, want *IntLiteral", i, arr.Elements[i])
					}
				}
			},
		},
		{
			name:         "mixed type elements",
			input:        `[42, "hello", true]`,
			wantElements: 3,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				// Check integer
				if intLit, ok := arr.Elements[0].(*IntLiteral); ok {
					if intLit.Value != "42" {
						t.Errorf("First element value = %q, want %q", intLit.Value, "42")
					}
				} else {
					t.Errorf("First element type = %T, want *IntLiteral", arr.Elements[0])
				}
				// Check string
				if strLit, ok := arr.Elements[1].(*StringLiteral); ok {
					if strLit.Value != "hello" {
						t.Errorf("Second element value = %q, want %q", strLit.Value, "hello")
					}
				} else {
					t.Errorf("Second element type = %T, want *StringLiteral", arr.Elements[1])
				}
				// Check boolean
				if boolLit, ok := arr.Elements[2].(*BoolLiteral); ok {
					if boolLit.Value != "true" {
						t.Errorf("Third element value = %q, want %q", boolLit.Value, "true")
					}
				} else {
					t.Errorf("Third element type = %T, want *BoolLiteral", arr.Elements[2])
				}
			},
		},
		{
			name:         "array with floats",
			input:        "[3.14, 2.71]",
			wantElements: 2,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				expected := []string{"3.14", "2.71"}
				for i, exp := range expected {
					if floatLit, ok := arr.Elements[i].(*FloatLiteral); ok {
						if floatLit.Value != exp {
							t.Errorf("Element %d value = %q, want %q", i, floatLit.Value, exp)
						}
					} else {
						t.Errorf("Element %d type = %T, want *FloatLiteral", i, arr.Elements[i])
					}
				}
			},
		},
		{
			name:         "array with strings",
			input:        `["hello", "world"]`,
			wantElements: 2,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				expected := []string{"hello", "world"}
				for i, exp := range expected {
					if strLit, ok := arr.Elements[i].(*StringLiteral); ok {
						if strLit.Value != exp {
							t.Errorf("Element %d value = %q, want %q", i, strLit.Value, exp)
						}
					} else {
						t.Errorf("Element %d type = %T, want *StringLiteral", i, arr.Elements[i])
					}
				}
			},
		},
		{
			name:         "array with booleans",
			input:        "[true, false]",
			wantElements: 2,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				expected := []string{"true", "false"}
				for i, exp := range expected {
					if boolLit, ok := arr.Elements[i].(*BoolLiteral); ok {
						if boolLit.Value != exp {
							t.Errorf("Element %d value = %q, want %q", i, boolLit.Value, exp)
						}
					} else {
						t.Errorf("Element %d type = %T, want *BoolLiteral", i, arr.Elements[i])
					}
				}
			},
		},
		{
			name:         "array with null",
			input:        "[NULL, 42]",
			wantElements: 2,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				if nullLit, ok := arr.Elements[0].(*NullLiteral); ok {
					if nullLit.Value != "NULL" {
						t.Errorf("First element value = %q, want %q", nullLit.Value, "NULL")
					}
				} else {
					t.Errorf("First element type = %T, want *NullLiteral", arr.Elements[0])
				}
			},
		},
		{
			name:         "array with identifiers",
			input:        "[myVar, anotherVar]",
			wantElements: 2,
			wantError:    false,
			checkElements: func(t *testing.T, arr *ArrayLiteral) {
				expected := []string{"myVar", "anotherVar"}
				for i, exp := range expected {
					if ident, ok := arr.Elements[i].(*Identifier); ok {
						if ident.Value != exp {
							t.Errorf("Element %d value = %q, want %q", i, ident.Value, exp)
						}
					} else {
						t.Errorf("Element %d type = %T, want *Identifier", i, arr.Elements[i])
					}
				}
			},
		},
		{
			name:         "single element without comma",
			input:        "[42]",
			wantElements: 1,
			wantError:    false,
		},
		{
			name:         "multiple elements without trailing comma",
			input:        "[1, 2, 3]",
			wantElements: 3,
			wantError:    false,
		},
		{
			name:         "array with extra spaces",
			input:        "[ 1 , 2 , 3 ]",
			wantElements: 3,
			wantError:    false,
		},
		{
			name:          "missing closing bracket",
			input:         "[1, 2, 3",
			wantError:     true,
			errorContains: "expected ']' to close array",
		},
		{
			name:          "missing comma between elements",
			input:         "[1 2]",
			wantError:     true,
			errorContains: "expected ',' or ']' in array",
		},
		{
			name:          "invalid element type",
			input:         "[42, @]",
			wantError:     true,
			errorContains: "illegal token",
		},
		{
			name:          "unexpected token after comma",
			input:         "[1, }]",
			wantError:     true,
			errorContains: "unexpected token",
		},
		{
			name:          "EOF while parsing",
			input:         "[1,",
			wantError:     true,
			errorContains: "expected ']' to close array",
		},
		{
			name:          "double comma",
			input:         "[1,, 2]",
			wantError:     true,
			errorContains: "unexpected token",
		},
		{
			name:          "comma at start",
			input:         "[, 1]",
			wantError:     true,
			errorContains: "unexpected token",
		},
		{
			name:          "only comma",
			input:         "[,]",
			wantError:     true,
			errorContains: "unexpected token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			// Ensure we start at LSQUARE_BRACKET token
			if !p.currentTokenIs(lexer.LSQUARE_BRACKET) {
				t.Fatalf("Test setup error: expected LSQUARE_BRACKET token, got %v", p.currentToken)
			}

			stmt, err := p.parseArray()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseArray() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("parseArray() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseArray() unexpected error: %v", err)
				return
			}

			if stmt == nil {
				t.Error("parseArray() returned nil statement")
				return
			}

			arrayLit, ok := stmt.(*ArrayLiteral)
			if !ok {
				t.Errorf("parseArray() returned %T, want *ArrayLiteral", stmt)
				return
			}

			if len(arrayLit.Elements) != tt.wantElements {
				t.Errorf("parseArray() element count = %d, want %d", len(arrayLit.Elements), tt.wantElements)
			}

			// Verify position is set
			if arrayLit.Pos.Line == 0 && arrayLit.Pos.Column == 0 {
				t.Error("parseArray() position not set")
			}

			// Run custom element checks if provided
			if tt.checkElements != nil {
				tt.checkElements(t, arrayLit)
			}
		})
	}
}

func TestParseArrayPosition(t *testing.T) {
	input := "[1, 2, 3]"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Capture the initial position
	initialPos := p.currentPosition

	stmt, err := p.parseArray()
	if err != nil {
		t.Fatalf("parseArray() unexpected error: %v", err)
	}

	arrayLit := stmt.(*ArrayLiteral)

	// Position should be set to the LSQUARE_BRACKET token position
	if arrayLit.Pos.Line != initialPos.Line || arrayLit.Pos.Column != initialPos.Column {
		t.Errorf("parseArray() position = %v, want %v", arrayLit.Pos, initialPos)
	}
}

func TestParseArrayTokenAdvancement(t *testing.T) {
	input := "[42] + 10"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Initially at LSQUARE_BRACKET
	if !p.currentTokenIs(lexer.LSQUARE_BRACKET) {
		t.Fatalf("Test setup error: expected LSQUARE_BRACKET token, got %v", p.currentToken)
	}

	_, err := p.parseArray()
	if err != nil {
		t.Fatalf("parseArray() unexpected error: %v", err)
	}

	// After parsing, current token should be at the closing bracket
	if !p.currentTokenIs(lexer.RSQUARE_BRACKET) {
		t.Errorf("parseArray() current token = %v, want %v", p.currentToken, lexer.RSQUARE_BRACKET)
	}

	// Peek token should be the next token (+)
	if !p.peekTokenIs(lexer.OP_PLUS) {
		t.Errorf("parseArray() peek token = %v, want %v", p.peekToken, lexer.OP_PLUS)
	}
}

func TestParseNestedArrayElements(t *testing.T) {
	// Test that array elements are properly parsed as statements
	tests := []struct {
		name      string
		input     string
		checkFunc func(*testing.T, *ArrayLiteral)
	}{
		{
			name:  "array with operators",
			input: "[+, -, *]",
			checkFunc: func(t *testing.T, arr *ArrayLiteral) {
				expectedOps := []lexer.TokenType{lexer.OP_PLUS, lexer.OP_SUBTRACT, lexer.OP_MULTIPLY}
				if len(arr.Elements) != 3 {
					t.Fatalf("Expected 3 elements, got %d", len(arr.Elements))
				}
				for i, expectedOp := range expectedOps {
					if binExpr, ok := arr.Elements[i].(*BinaryExpression); ok {
						if binExpr.Operator != expectedOp {
							t.Errorf("Element %d operator = %v, want %v", i, binExpr.Operator, expectedOp)
						}
					} else {
						t.Errorf("Element %d type = %T, want *BinaryExpression", i, arr.Elements[i])
					}
				}
			},
		},
		{
			name:  "array with unary operators",
			input: "[!, ++, --]",
			checkFunc: func(t *testing.T, arr *ArrayLiteral) {
				expectedOps := []lexer.TokenType{lexer.OP_NOT, lexer.OP_INC, lexer.OP_DEC}
				if len(arr.Elements) != 3 {
					t.Fatalf("Expected 3 elements, got %d", len(arr.Elements))
				}
				for i, expectedOp := range expectedOps {
					if unExpr, ok := arr.Elements[i].(*UnaryExpression); ok {
						if unExpr.Operator != expectedOp {
							t.Errorf("Element %d operator = %v, want %v", i, unExpr.Operator, expectedOp)
						}
					} else {
						t.Errorf("Element %d type = %T, want *UnaryExpression", i, arr.Elements[i])
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseArray()
			if err != nil {
				t.Fatalf("parseArray() unexpected error: %v", err)
			}

			arrayLit := stmt.(*ArrayLiteral)
			tt.checkFunc(t, arrayLit)
		})
	}
}

func TestParseArrayInContext(t *testing.T) {
	// Test array parsing as part of a complete parse
	tests := []struct {
		name          string
		input         string
		wantStmtCount int
		wantError     bool
	}{
		{
			name:          "array as single statement",
			input:         "[1, 2, 3]",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "array with other statements",
			input:         "42 [1, 2] DUP",
			wantStmtCount: 3,
			wantError:     false,
		},
		{
			name:          "array assignment",
			input:         "[1, 2, 3] := VAR arr",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "multiple arrays",
			input:         "[1] [2] [3]",
			wantStmtCount: 3,
			wantError:     false,
		},
		{
			name:          "empty arrays",
			input:         "[] [] []",
			wantStmtCount: 3,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()

			if tt.wantError {
				if err == nil {
					t.Errorf("Parse() expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if len(program.Statements) != tt.wantStmtCount {
				t.Errorf("Parse() statement count = %d, want %d", len(program.Statements), tt.wantStmtCount)
			}

			// Verify at least one array exists
			hasArray := false
			for _, stmt := range program.Statements {
				if _, ok := stmt.(*ArrayLiteral); ok {
					hasArray = true
					break
				}
			}

			if !hasArray {
				t.Error("Parse() expected at least one ArrayLiteral statement")
			}
		})
	}
}

func TestParseArrayElementPositions(t *testing.T) {
	input := "[42, \"hello\"]"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	stmt, err := p.parseArray()
	if err != nil {
		t.Fatalf("parseArray() unexpected error: %v", err)
	}

	arrayLit := stmt.(*ArrayLiteral)

	// Check that elements have position information
	for i, elem := range arrayLit.Elements {
		switch e := elem.(type) {
		case *IntLiteral:
			if e.Position.Line == 0 && e.Position.Column == 0 {
				t.Errorf("Element %d IntLiteral missing position", i)
			}
		case *StringLiteral:
			if e.Position.Line == 0 && e.Position.Column == 0 {
				t.Errorf("Element %d StringLiteral missing position", i)
			}
		}
	}
}
func TestParseDoStmt(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantType      string
		wantError     bool
		errorContains string
		checkStmt     func(*testing.T, Node)
	}{
		// Do-While tests
		{
			name:     "simple do-while",
			input:    "DO 42 WHILE true END",
			wantType: "*parser.WhileStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Body) != 1 {
					t.Errorf("Body length = %d, want 1", len(whileStmt.Body))
				}
				if len(whileStmt.Condition) != 1 {
					t.Errorf("Condition length = %d, want 1", len(whileStmt.Condition))
				}
				if intLit, ok := whileStmt.Body[0].(*IntLiteral); ok {
					if intLit.Value != "42" {
						t.Errorf("Body value = %q, want %q", intLit.Value, "42")
					}
				} else {
					t.Errorf("Body type = %T, want *IntLiteral", whileStmt.Body[0])
				}
			},
		},
		{
			name:     "do-while with multiple body statements",
			input:    "DO 42 10 + WHILE true false && END",
			wantType: "*parser.WhileStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Body) != 3 {
					t.Errorf("Body length = %d, want 3", len(whileStmt.Body))
				}
				if len(whileStmt.Condition) != 3 {
					t.Errorf("Condition length = %d, want 3", len(whileStmt.Condition))
				}
			},
		},
		{
			name:     "do-while with empty body",
			input:    "DO WHILE true END",
			wantType: "*parser.WhileStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Body) != 0 {
					t.Errorf("Body length = %d, want 0", len(whileStmt.Body))
				}
				if len(whileStmt.Condition) != 1 {
					t.Errorf("Condition length = %d, want 1", len(whileStmt.Condition))
				}
			},
		},
		{
			name:     "do-while with empty condition",
			input:    "DO 42 WHILE END",
			wantType: "*parser.WhileStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Body) != 1 {
					t.Errorf("Body length = %d, want 1", len(whileStmt.Body))
				}
				if len(whileStmt.Condition) != 0 {
					t.Errorf("Condition length = %d, want 0", len(whileStmt.Condition))
				}
			},
		},

		// Do-If tests
		{
			name:     "simple do-if",
			input:    "DO 42 IF true END",
			wantType: "*parser.IfStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.ThenBranch) != 1 {
					t.Errorf("ThenBranch length = %d, want 1", len(ifStmt.ThenBranch))
				}
				if len(ifStmt.ElseBranch) != 0 {
					t.Errorf("ElseBranch length = %d, want 0", len(ifStmt.ElseBranch))
				}
				if len(ifStmt.Condition) != 1 {
					t.Errorf("Condition length = %d, want 1", len(ifStmt.Condition))
				}
			},
		},
		{
			name:     "do-if with else",
			input:    "DO 42 ELSE 10 IF true END",
			wantType: "*parser.IfStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.ThenBranch) != 1 {
					t.Errorf("ThenBranch length = %d, want 1", len(ifStmt.ThenBranch))
				}
				if len(ifStmt.ElseBranch) != 1 {
					t.Errorf("ElseBranch length = %d, want 1", len(ifStmt.ElseBranch))
				}
				if len(ifStmt.Condition) != 1 {
					t.Errorf("Condition length = %d, want 1", len(ifStmt.Condition))
				}
				if intLit, ok := ifStmt.ElseBranch[0].(*IntLiteral); ok {
					if intLit.Value != "10" {
						t.Errorf("ElseBranch value = %q, want %q", intLit.Value, "10")
					}
				}
			},
		},
		{
			name:     "do-if-else with multiple statements",
			input:    "DO 1 2 + ELSE 3 4 * IF x y == END",
			wantType: "*parser.IfStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.ThenBranch) != 3 {
					t.Errorf("ThenBranch length = %d, want 3", len(ifStmt.ThenBranch))
				}
				if len(ifStmt.ElseBranch) != 3 {
					t.Errorf("ElseBranch length = %d, want 3", len(ifStmt.ElseBranch))
				}
				if len(ifStmt.Condition) != 3 {
					t.Errorf("Condition length = %d, want 3", len(ifStmt.Condition))
				}
			},
		},
		{
			name:     "do-if with empty then branch",
			input:    "DO IF true END",
			wantType: "*parser.IfStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.ThenBranch) != 0 {
					t.Errorf("ThenBranch length = %d, want 0", len(ifStmt.ThenBranch))
				}
			},
		},
		{
			name:     "do-if with empty else branch",
			input:    "DO 42 ELSE IF true END",
			wantType: "*parser.IfStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.ThenBranch) != 1 {
					t.Errorf("ThenBranch length = %d, want 1", len(ifStmt.ThenBranch))
				}
				if len(ifStmt.ElseBranch) != 0 {
					t.Errorf("ElseBranch length = %d, want 0", len(ifStmt.ElseBranch))
				}
			},
		},
		{
			name:     "do-if with empty condition",
			input:    "DO 42 IF END",
			wantType: "*parser.IfStmt",
			checkStmt: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.Condition) != 0 {
					t.Errorf("Condition length = %d, want 0", len(ifStmt.Condition))
				}
			},
		},

		// Error cases
		{
			name:          "missing END in do-while",
			input:         "DO 42 WHILE true",
			wantError:     true,
			errorContains: "expected END",
		},
		{
			name:          "missing END in do-if",
			input:         "DO 42 IF true",
			wantError:     true,
			errorContains: "expected END",
		},
		{
			name:          "missing IF before END",
			input:         "DO 42 ELSE 10 END",
			wantError:     true,
			errorContains: "expected IF before END",
		},
		{
			name:          "unexpected token after ELSE",
			input:         "DO 42 ELSE 10 WHILE true END",
			wantError:     true,
			errorContains: "expected IF before END",
		},
		{
			name:          "invalid statement in body",
			input:         "DO @ WHILE true END",
			wantError:     true,
			errorContains: "illegal token",
		},
		{
			name:          "invalid statement in condition",
			input:         "DO 42 WHILE @ END",
			wantError:     true,
			errorContains: "illegal token",
		},
		{
			name:          "invalid statement in else branch",
			input:         "DO 42 ELSE @ IF true END",
			wantError:     true,
			errorContains: "illegal token",
		},
		{
			name:          "EOF after DO",
			input:         "DO",
			wantError:     true,
			errorContains: "expected IF before END",
		},
		{
			name:          "EOF in body",
			input:         "DO 42",
			wantError:     true,
			errorContains: "expected IF before END",
		},
		{
			name:          "EOF in condition",
			input:         "DO 42 WHILE",
			wantError:     true,
			errorContains: "expected END",
		},
		{
			name:          "EOF in else branch",
			input:         "DO 42 ELSE",
			wantError:     true,
			errorContains: "expected IF before END",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			// Ensure we start at DO token
			if !p.currentTokenIs(lexer.DO) {
				t.Fatalf("Test setup error: expected DO token, got %v", p.currentToken)
			}

			stmt, err := p.parseDoStmt()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseDoStmt() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("parseDoStmt() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseDoStmt() unexpected error: %v", err)
				return
			}

			if stmt == nil {
				t.Error("parseDoStmt() returned nil statement")
				return
			}

			gotType := fmt.Sprintf("%T", stmt)
			if gotType != tt.wantType {
				t.Errorf("parseDoStmt() statement type = %s, want %s", gotType, tt.wantType)
			}

			// Verify position is set
			switch s := stmt.(type) {
			case *WhileStmt:
				if s.Position.Line == 0 && s.Position.Column == 0 {
					t.Error("parseDoStmt() WhileStmt position not set")
				}
			case *IfStmt:
				if s.Position.Line == 0 && s.Position.Column == 0 {
					t.Error("parseDoStmt() IfStmt position not set")
				}
			}

			// Run custom checks if provided
			if tt.checkStmt != nil {
				tt.checkStmt(t, stmt)
			}
		})
	}
}

func TestParseDoStmtPosition(t *testing.T) {
	input := "DO 42 WHILE true END"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Capture the initial position
	initialPos := p.currentPosition

	stmt, err := p.parseDoStmt()
	if err != nil {
		t.Fatalf("parseDoStmt() unexpected error: %v", err)
	}

	whileStmt := stmt.(*WhileStmt)

	// Position should be set to the DO token position
	if whileStmt.Position.Line != initialPos.Line || whileStmt.Position.Column != initialPos.Column {
		t.Errorf("parseDoStmt() position = %v, want %v", whileStmt.Position, initialPos)
	}
}

func TestParseDoStmtTokenAdvancement(t *testing.T) {
	input := "DO 42 WHILE true END + 10"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Initially at DO
	if !p.currentTokenIs(lexer.DO) {
		t.Fatalf("Test setup error: expected DO token, got %v", p.currentToken)
	}

	_, err := p.parseDoStmt()
	if err != nil {
		t.Fatalf("parseDoStmt() unexpected error: %v", err)
	}

	// After parsing, current token should be at END
	if !p.currentTokenIs(lexer.END) {
		t.Errorf("parseDoStmt() current token = %v, want %v", p.currentToken, lexer.END)
	}

	// Peek token should be the next token (+)
	if !p.peekTokenIs(lexer.OP_PLUS) {
		t.Errorf("parseDoStmt() peek token = %v, want %v", p.peekToken, lexer.OP_PLUS)
	}
}

func TestParseDoStmtComplexExpressions(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		checkFunc func(*testing.T, Node)
	}{
		{
			name:  "do-while with complex body and condition",
			input: "DO VAR x 42 := x [1, 2, 3] WHILE x 10 > x 5 + DUMP END",
			checkFunc: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Body) != 5 {
					t.Errorf("Body length = %d, want 5", len(whileStmt.Body))
				}
				if len(whileStmt.Condition) != 6 {
					t.Errorf("Condition length = %d, want 6", len(whileStmt.Condition))
				}

				// Check body contains variable declaration
				if _, ok := whileStmt.Body[0].(*VarDeclaration); !ok {
					t.Errorf("First body statement type = %T, want *VarDeclaration", whileStmt.Body[0])
				}

				// Check body contains array literal
				if _, ok := whileStmt.Body[4].(*ArrayLiteral); !ok {
					t.Errorf("Fifth body statement type = %T, want *ArrayLiteral", whileStmt.Body[4])
				}
			},
		},
		{
			name:  "do-if-else with complex branches",
			input: "DO VAR result 100 ELSE CONST msg \"error\" msg DUMP IF x 0 == END",
			checkFunc: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.ThenBranch) != 2 {
					t.Errorf("ThenBranch length = %d, want 2", len(ifStmt.ThenBranch))
				}
				if len(ifStmt.ElseBranch) != 3 {
					t.Errorf("ElseBranch length = %d, want 3", len(ifStmt.ElseBranch))
				}
				if len(ifStmt.Condition) != 3 {
					t.Errorf("Condition length = %d, want 3", len(ifStmt.Condition))
				}

				// Check then branch contains variable declaration
				if _, ok := ifStmt.ThenBranch[0].(*VarDeclaration); !ok {
					t.Errorf("First then statement type = %T, want *VarDeclaration", ifStmt.ThenBranch[0])
				}

				// Check else branch contains constant declaration
				if _, ok := ifStmt.ElseBranch[0].(*ConstDeclaration); !ok {
					t.Errorf("First else statement type = %T, want *ConstDeclaration", ifStmt.ElseBranch[0])
				}
			},
		},
		{
			name:  "do-while with stack operations",
			input: "DO DUP SWAP ROT WHILE DROP END",
			checkFunc: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Body) != 3 {
					t.Errorf("Body length = %d, want 3", len(whileStmt.Body))
				}

				// Check all body statements are stack operations
				for i, bodyStmt := range whileStmt.Body {
					if _, ok := bodyStmt.(*StackOp); !ok {
						t.Errorf("Body statement %d type = %T, want *StackOp", i, bodyStmt)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseDoStmt()
			if err != nil {
				t.Fatalf("parseDoStmt() unexpected error: %v", err)
			}

			tt.checkFunc(t, stmt)
		})
	}
}

func TestParseDoStmtInContext(t *testing.T) {
	// Test do statements as part of complete programs
	tests := []struct {
		name          string
		input         string
		wantStmtCount int
		wantError     bool
	}{
		{
			name:          "do-while as single statement",
			input:         "DO 42 WHILE true END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "do-if as single statement",
			input:         "DO 42 IF true END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "multiple do statements",
			input:         "DO 1 WHILE true END DO 2 IF false END",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "do statement with other statements",
			input:         "VAR x DO x WHILE true END 42 DUMP",
			wantStmtCount: 4,
			wantError:     false,
		},
		{
			name:          "assignment to do statement result",
			input:         "DO 42 IF true END := VAR result",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "do statement in array",
			input:         "[DO 1 IF true END, 42]",
			wantStmtCount: 1,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()

			if tt.wantError {
				if err == nil {
					t.Errorf("Parse() expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if len(program.Statements) != tt.wantStmtCount {
				t.Errorf("Parse() statement count = %d, want %d", len(program.Statements), tt.wantStmtCount)
			}
		})
	}
}

func TestParseDoStmtEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		checkFunc func(*testing.T, Node)
	}{
		{
			name:      "do-while with only whitespace in body",
			input:     "DO   WHILE true END",
			wantError: false,
			checkFunc: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Body) != 0 {
					t.Errorf("Body should be empty, got %d statements", len(whileStmt.Body))
				}
			},
		},
		{
			name:      "do-if with multiple else statements",
			input:     "DO 1 ELSE 2 3 4 IF true END",
			wantError: false,
			checkFunc: func(t *testing.T, stmt Node) {
				ifStmt := stmt.(*IfStmt)
				if len(ifStmt.ElseBranch) != 3 {
					t.Errorf("ElseBranch length = %d, want 3", len(ifStmt.ElseBranch))
				}
			},
		},
		{
			name:      "do-while with operators in condition",
			input:     "DO 42 WHILE + - * / END",
			wantError: false,
			checkFunc: func(t *testing.T, stmt Node) {
				whileStmt := stmt.(*WhileStmt)
				if len(whileStmt.Condition) != 4 {
					t.Errorf("Condition length = %d, want 4", len(whileStmt.Condition))
				}
				for i, condStmt := range whileStmt.Condition {
					if _, ok := condStmt.(*BinaryExpression); !ok {
						t.Errorf("Condition statement %d type = %T, want *BinaryExpression", i, condStmt)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseDoStmt()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseDoStmt() expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("parseDoStmt() unexpected error: %v", err)
				return
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, stmt)
			}
		})
	}
}
func TestParseWhileStmt(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantBodyCount int
		wantError     bool
		errorContains string
		checkStmt     func(*testing.T, *WhileStmt)
	}{
		{
			name:          "simple while with single statement",
			input:         "WHILE 42 END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 1 {
					t.Errorf("Body length = %d, want 1", len(stmt.Body))
				}
				if stmt.Condition != nil {
					t.Errorf("Condition = %v, want nil", stmt.Condition)
				}
				if intLit, ok := stmt.Body[0].(*IntLiteral); ok {
					if intLit.Value != "42" {
						t.Errorf("Body value = %q, want %q", intLit.Value, "42")
					}
				} else {
					t.Errorf("Body type = %T, want *IntLiteral", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with multiple statements",
			input:         "WHILE 42 10 + DUP SWAP END",
			wantBodyCount: 4,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 4 {
					t.Errorf("Body length = %d, want 4", len(stmt.Body))
				}
				// Check first statement is integer
				if _, ok := stmt.Body[0].(*IntLiteral); !ok {
					t.Errorf("First body statement type = %T, want *IntLiteral", stmt.Body[0])
				}
				// Check second statement is integer
				if _, ok := stmt.Body[1].(*IntLiteral); !ok {
					t.Errorf("Second body statement type = %T, want *IntLiteral", stmt.Body[1])
				}
				// Check third statement is binary operation
				if _, ok := stmt.Body[2].(*BinaryExpression); !ok {
					t.Errorf("Third body statement type = %T, want *BinaryExpression", stmt.Body[2])
				}
				// Check fourth statement is stack operation
				if _, ok := stmt.Body[3].(*StackOp); !ok {
					t.Errorf("Fourth body statement type = %T, want *StackOp", stmt.Body[3])
				}
			},
		},
		{
			name:          "empty while loop",
			input:         "WHILE END",
			wantBodyCount: 0,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 0 {
					t.Errorf("Body length = %d, want 0", len(stmt.Body))
				}
			},
		},
		{
			name:          "while with string literal",
			input:         `WHILE "hello" END`,
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if strLit, ok := stmt.Body[0].(*StringLiteral); ok {
					if strLit.Value != "hello" {
						t.Errorf("Body value = %q, want %q", strLit.Value, "hello")
					}
				} else {
					t.Errorf("Body type = %T, want *StringLiteral", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with boolean literal",
			input:         "WHILE true false END",
			wantBodyCount: 2,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if boolLit, ok := stmt.Body[0].(*BoolLiteral); ok {
					if boolLit.Value != "true" {
						t.Errorf("First body value = %q, want %q", boolLit.Value, "true")
					}
				}
				if boolLit, ok := stmt.Body[1].(*BoolLiteral); ok {
					if boolLit.Value != "false" {
						t.Errorf("Second body value = %q, want %q", boolLit.Value, "false")
					}
				}
			},
		},
		{
			name:          "while with float literal",
			input:         "WHILE 3.14 END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if floatLit, ok := stmt.Body[0].(*FloatLiteral); ok {
					if floatLit.Value != "3.14" {
						t.Errorf("Body value = %q, want %q", floatLit.Value, "3.14")
					}
				} else {
					t.Errorf("Body type = %T, want *FloatLiteral", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with identifier",
			input:         "WHILE myVar END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if ident, ok := stmt.Body[0].(*Identifier); ok {
					if ident.Value != "myVar" {
						t.Errorf("Body value = %q, want %q", ident.Value, "myVar")
					}
				} else {
					t.Errorf("Body type = %T, want *Identifier", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with null literal",
			input:         "WHILE NULL END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if nullLit, ok := stmt.Body[0].(*NullLiteral); ok {
					if nullLit.Value != "NULL" {
						t.Errorf("Body value = %q, want %q", nullLit.Value, "NULL")
					}
				} else {
					t.Errorf("Body type = %T, want *NullLiteral", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with variable declaration",
			input:         "WHILE VAR x END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if varDecl, ok := stmt.Body[0].(*VarDeclaration); ok {
					if varDecl.Name != "x" {
						t.Errorf("Variable name = %q, want %q", varDecl.Name, "x")
					}
				} else {
					t.Errorf("Body type = %T, want *VarDeclaration", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with constant declaration",
			input:         "WHILE CONST PI 3.14 END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if constDecl, ok := stmt.Body[0].(*ConstDeclaration); ok {
					if constDecl.Name != "PI" {
						t.Errorf("Constant name = %q, want %q", constDecl.Name, "PI")
					}
				} else {
					t.Errorf("Body type = %T, want *ConstDeclaration", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with assignment",
			input:         "WHILE := VAR x END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if assignment, ok := stmt.Body[0].(*Assignment); ok {
					if varDecl, ok := assignment.Target.(*VarDeclaration); ok {
						if varDecl.Name != "x" {
							t.Errorf("Assignment target name = %q, want %q", varDecl.Name, "x")
						}
					} else {
						t.Errorf("Assignment target type = %T, want *VarDeclaration", assignment.Target)
					}
				} else {
					t.Errorf("Body type = %T, want *Assignment", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with array",
			input:         "WHILE [1, 2, 3] END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if arrayLit, ok := stmt.Body[0].(*ArrayLiteral); ok {
					if len(arrayLit.Elements) != 3 {
						t.Errorf("Array element count = %d, want 3", len(arrayLit.Elements))
					}
				} else {
					t.Errorf("Body type = %T, want *ArrayLiteral", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with all binary operators",
			input:         "WHILE + - * / % ** == != > < >= <= && || END",
			wantBodyCount: 14,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 14 {
					t.Errorf("Body length = %d, want 14", len(stmt.Body))
				}
				for i, bodyStmt := range stmt.Body {
					if _, ok := bodyStmt.(*BinaryExpression); !ok {
						t.Errorf("Body statement %d type = %T, want *BinaryExpression", i, bodyStmt)
					}
				}
			},
		},
		{
			name:          "while with all unary operators",
			input:         "WHILE ! ++ -- END",
			wantBodyCount: 3,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 3 {
					t.Errorf("Body length = %d, want 3", len(stmt.Body))
				}
				for i, bodyStmt := range stmt.Body {
					if _, ok := bodyStmt.(*UnaryExpression); !ok {
						t.Errorf("Body statement %d type = %T, want *UnaryExpression", i, bodyStmt)
					}
				}
			},
		},
		{
			name:          "while with all stack operations",
			input:         "WHILE DROP SWAP DUP OVER ROT DEL CLEAR PICK END",
			wantBodyCount: 8,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 8 {
					t.Errorf("Body length = %d, want 8", len(stmt.Body))
				}
				for i, bodyStmt := range stmt.Body {
					if _, ok := bodyStmt.(*StackOp); !ok {
						t.Errorf("Body statement %d type = %T, want *StackOp", i, bodyStmt)
					}
				}
			},
		},
		{
			name:          "while with io operations",
			input:         "WHILE DUMP DUMPLN END",
			wantBodyCount: 2,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 2 {
					t.Errorf("Body length = %d, want 2", len(stmt.Body))
				}
				for i, bodyStmt := range stmt.Body {
					if _, ok := bodyStmt.(*IOStmt); !ok {
						t.Errorf("Body statement %d type = %T, want *IOStmt", i, bodyStmt)
					}
				}
			},
		},
		{
			name:          "while with break statement",
			input:         "WHILE BREAK END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if _, ok := stmt.Body[0].(*BreakStmt); !ok {
					t.Errorf("Body type = %T, want *BreakStmt", stmt.Body[0])
				}
			},
		},
		{
			name:          "while with continue statement",
			input:         "WHILE CONTINUE END",
			wantBodyCount: 1,
			wantError:     false,
			checkStmt: func(t *testing.T, stmt *WhileStmt) {
				if _, ok := stmt.Body[0].(*ContinueStmt); !ok {
					t.Errorf("Body type = %T, want *ContinueStmt", stmt.Body[0])
				}
			},
		},
		{
			name:          "missing END",
			input:         "WHILE 42",
			wantError:     true,
			errorContains: "expected END",
		},
		{
			name:          "missing END with multiple statements",
			input:         "WHILE 42 10 +",
			wantError:     true,
			errorContains: "expected END",
		},
		{
			name:          "invalid statement in body",
			input:         "WHILE @ END",
			wantError:     true,
			errorContains: "illegal token",
		},
		{
			name:          "unexpected token in body",
			input:         "WHILE { END",
			wantError:     true,
			errorContains: "unexpected token",
		},
		{
			name:          "EOF after WHILE",
			input:         "WHILE",
			wantError:     true,
			errorContains: "expected END",
		},
		{
			name:          "EOF in body",
			input:         "WHILE 42 10",
			wantError:     true,
			errorContains: "expected END",
		},
		{
			name:          "invalid variable declaration",
			input:         "WHILE VAR END",
			wantError:     true,
			errorContains: "expected identifier after VAR",
		},
		{
			name:          "invalid assignment",
			input:         "WHILE := END",
			wantError:     true,
			errorContains: "expected identifier, VAR, or CONST after ':='",
		},
		{
			name:          "invalid array",
			input:         "WHILE [1, END",
			wantError:     true,
			errorContains: "expected ']' to close array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			// Ensure we start at WHILE token
			if !p.currentTokenIs(lexer.WHILE) {
				t.Fatalf("Test setup error: expected WHILE token, got %v", p.currentToken)
			}

			stmt, err := p.parseWhileStmt()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseWhileStmt() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("parseWhileStmt() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseWhileStmt() unexpected error: %v", err)
				return
			}

			if stmt == nil {
				t.Error("parseWhileStmt() returned nil statement")
				return
			}

			whileStmt, ok := stmt.(*WhileStmt)
			if !ok {
				t.Errorf("parseWhileStmt() returned %T, want *WhileStmt", stmt)
				return
			}

			if len(whileStmt.Body) != tt.wantBodyCount {
				t.Errorf("parseWhileStmt() body count = %d, want %d", len(whileStmt.Body), tt.wantBodyCount)
			}

			// Verify Condition is always nil for regular while statements
			if whileStmt.Condition != nil {
				t.Error("parseWhileStmt() Condition should be nil for regular while statements")
			}

			// Verify position is set
			if whileStmt.Position.Line == 0 && whileStmt.Position.Column == 0 {
				t.Error("parseWhileStmt() position not set")
			}

			// Run custom checks if provided
			if tt.checkStmt != nil {
				tt.checkStmt(t, whileStmt)
			}
		})
	}
}

func TestParseWhileStmtPosition(t *testing.T) {
	input := "WHILE 42 END"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Capture the initial position
	initialPos := p.currentPosition

	stmt, err := p.parseWhileStmt()
	if err != nil {
		t.Fatalf("parseWhileStmt() unexpected error: %v", err)
	}

	whileStmt := stmt.(*WhileStmt)

	// Position should be set to the WHILE token position
	if whileStmt.Position.Line != initialPos.Line || whileStmt.Position.Column != initialPos.Column {
		t.Errorf("parseWhileStmt() position = %v, want %v", whileStmt.Position, initialPos)
	}
}

func TestParseWhileStmtTokenAdvancement(t *testing.T) {
	input := "WHILE 42 END + 10"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	// Initially at WHILE
	if !p.currentTokenIs(lexer.WHILE) {
		t.Fatalf("Test setup error: expected WHILE token, got %v", p.currentToken)
	}

	_, err := p.parseWhileStmt()
	if err != nil {
		t.Fatalf("parseWhileStmt() unexpected error: %v", err)
	}

	// After parsing, current token should be at END
	if !p.currentTokenIs(lexer.END) {
		t.Errorf("parseWhileStmt() current token = %v, want %v", p.currentToken, lexer.END)
	}

	// Peek token should be the next token (+)
	if !p.peekTokenIs(lexer.OP_PLUS) {
		t.Errorf("parseWhileStmt() peek token = %v, want %v", p.peekToken, lexer.OP_PLUS)
	}
}

func TestParseWhileStmtComplexExpressions(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		checkFunc func(*testing.T, *WhileStmt)
	}{
		{
			name:  "while with complex body",
			input: "WHILE VAR x 42 := x [1, 2, 3] x 10 > := result result DUMP END",
			checkFunc: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 8 {
					t.Errorf("Body length = %d, want 8", len(stmt.Body))
				}

				// Check variable declaration
				if _, ok := stmt.Body[0].(*VarDeclaration); !ok {
					t.Errorf("First body statement type = %T, want *VarDeclaration", stmt.Body[0])
				}

				// Check assignment
				if _, ok := stmt.Body[2].(*Assignment); !ok {
					t.Errorf("Third body statement type = %T, want *Assignment", stmt.Body[2])
				}

				// Check array literal
				if _, ok := stmt.Body[4].(*ArrayLiteral); !ok {
					t.Errorf("Fifth body statement type = %T, want *ArrayLiteral", stmt.Body[4])
				}
			},
		},
		{
			name:  "while with nested structures",
			input: "WHILE [VAR x, CONST y 10, := result] DUP SWAP END",
			checkFunc: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 3 {
					t.Errorf("Body length = %d, want 3", len(stmt.Body))
				}

				// Check array contains complex elements
				if arrayLit, ok := stmt.Body[0].(*ArrayLiteral); ok {
					if len(arrayLit.Elements) != 3 {
						t.Errorf("Array element count = %d, want 3", len(arrayLit.Elements))
					}
					if _, ok := arrayLit.Elements[0].(*VarDeclaration); !ok {
						t.Errorf("First array element type = %T, want *VarDeclaration", arrayLit.Elements[0])
					}
				}
			},
		},
		{
			name:  "while with operators and literals mix",
			input: "WHILE 1 2 + 3.14 \"test\" true false NULL myVar + * / END",
			checkFunc: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 12 {
					t.Errorf("Body length = %d, want 12", len(stmt.Body))
				}

				// Verify mix of different statement types
				expectedTypes := []string{
					"*parser.IntLiteral",
					"*parser.IntLiteral",
					"*parser.BinaryExpression",
					"*parser.FloatLiteral",
					"*parser.StringLiteral",
					"*parser.BoolLiteral",
					"*parser.BoolLiteral",
					"*parser.NullLiteral",
					"*parser.Identifier",
					"*parser.BinaryExpression",
					"*parser.BinaryExpression",
					"*parser.BinaryExpression",
				}

				for i, expectedType := range expectedTypes {
					if i < len(stmt.Body) {
						gotType := fmt.Sprintf("%T", stmt.Body[i])
						if gotType != expectedType {
							t.Errorf("Body statement %d type = %s, want %s", i, gotType, expectedType)
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseWhileStmt()
			if err != nil {
				t.Fatalf("parseWhileStmt() unexpected error: %v", err)
			}

			whileStmt := stmt.(*WhileStmt)
			tt.checkFunc(t, whileStmt)
		})
	}
}

func TestParseWhileStmtInContext(t *testing.T) {
	// Test while statements as part of complete programs
	tests := []struct {
		name          string
		input         string
		wantStmtCount int
		wantError     bool
	}{
		{
			name:          "while as single statement",
			input:         "WHILE 42 END",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "multiple while statements",
			input:         "WHILE 1 END WHILE 2 END",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "while with other statements",
			input:         "VAR x WHILE x END 42 DUMP",
			wantStmtCount: 4,
			wantError:     false,
		},
		{
			name:          "assignment to while result",
			input:         "WHILE 42 END := VAR result",
			wantStmtCount: 2,
			wantError:     false,
		},
		{
			name:          "while in array",
			input:         "[WHILE 1 END, 42]",
			wantStmtCount: 1,
			wantError:     false,
		},
		{
			name:          "nested statements around while",
			input:         "CONST x 10 WHILE VAR y y := x END [1, 2]",
			wantStmtCount: 3,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()

			if tt.wantError {
				if err == nil {
					t.Errorf("Parse() expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if len(program.Statements) != tt.wantStmtCount {
				t.Errorf("Parse() statement count = %d, want %d", len(program.Statements), tt.wantStmtCount)
			}

			// Verify at least one while statement exists
			hasWhile := false
			var findWhile func(Node) bool
			findWhile = func(n Node) bool {
				switch node := n.(type) {
				case *WhileStmt:
					return true
				case *ArrayLiteral:
					for _, elem := range node.Elements {
						if findWhile(elem) {
							return true
						}
					}
				}
				return false
			}

			for _, stmt := range program.Statements {
				if findWhile(stmt) {
					hasWhile = true
					break
				}
			}

			if !hasWhile {
				t.Error("Parse() expected at least one WhileStmt")
			}
		})
	}
}

func TestParseWhileStmtEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		checkFunc func(*testing.T, *WhileStmt)
	}{
		{
			name:      "while with only whitespace in body",
			input:     "WHILE   END",
			wantError: false,
			checkFunc: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 0 {
					t.Errorf("Body should be empty, got %d statements", len(stmt.Body))
				}
			},
		},
		{
			name:      "while with control flow statements",
			input:     "WHILE BREAK CONTINUE END",
			wantError: false,
			checkFunc: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 2 {
					t.Errorf("Body length = %d, want 2", len(stmt.Body))
				}
				if _, ok := stmt.Body[0].(*BreakStmt); !ok {
					t.Errorf("First body statement type = %T, want *BreakStmt", stmt.Body[0])
				}
				if _, ok := stmt.Body[1].(*ContinueStmt); !ok {
					t.Errorf("Second body statement type = %T, want *ContinueStmt", stmt.Body[1])
				}
			},
		},
		{
			name:      "while with all literal types",
			input:     "WHILE 42 3.14 \"string\" true false NULL identifier END",
			wantError: false,
			checkFunc: func(t *testing.T, stmt *WhileStmt) {
				if len(stmt.Body) != 7 {
					t.Errorf("Body length = %d, want 7", len(stmt.Body))
				}
				expectedTypes := []string{
					"*parser.IntLiteral",
					"*parser.FloatLiteral",
					"*parser.StringLiteral",
					"*parser.BoolLiteral",
					"*parser.BoolLiteral",
					"*parser.NullLiteral",
					"*parser.Identifier",
				}
				for i, expectedType := range expectedTypes {
					gotType := fmt.Sprintf("%T", stmt.Body[i])
					if gotType != expectedType {
						t.Errorf("Body statement %d type = %s, want %s", i, gotType, expectedType)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			stmt, err := p.parseWhileStmt()

			if tt.wantError {
				if err == nil {
					t.Errorf("parseWhileStmt() expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("parseWhileStmt() unexpected error: %v", err)
				return
			}

			whileStmt := stmt.(*WhileStmt)
			if tt.checkFunc != nil {
				tt.checkFunc(t, whileStmt)
			}
		})
	}
}

func TestParseWhileStmtBodyElementPositions(t *testing.T) {
	input := "WHILE 42 \"hello\" myVar END"
	l := lexer.NewLexer(strings.NewReader(input))
	p := NewParser(l)

	stmt, err := p.parseWhileStmt()
	if err != nil {
		t.Fatalf("parseWhileStmt() unexpected error: %v", err)
	}

	whileStmt := stmt.(*WhileStmt)

	// Check that body elements have position information
	for i, elem := range whileStmt.Body {
		switch e := elem.(type) {
		case *IntLiteral:
			if e.Position.Line == 0 && e.Position.Column == 0 {
				t.Errorf("Body element %d IntLiteral missing position", i)
			}
		case *StringLiteral:
			if e.Position.Line == 0 && e.Position.Column == 0 {
				t.Errorf("Body element %d StringLiteral missing position", i)
			}
		case *Identifier:
			if e.Position.Line == 0 && e.Position.Column == 0 {
				t.Errorf("Body element %d Identifier missing position", i)
			}
		}
	}
}
func TestValidateLoopControl(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantError     bool
		errorContains string
	}{
		// Valid cases
		{
			name:      "empty program",
			input:     "",
			wantError: false,
		},
		{
			name:      "program without break/continue",
			input:     "42 10 + VAR x := x",
			wantError: false,
		},
		{
			name:      "break inside while loop",
			input:     "WHILE BREAK END",
			wantError: false,
		},
		{
			name:      "continue inside while loop",
			input:     "WHILE CONTINUE END",
			wantError: false,
		},
		{
			name:      "break and continue inside while loop",
			input:     "WHILE BREAK CONTINUE END",
			wantError: false,
		},
		{
			name:      "break inside do-while loop",
			input:     "DO BREAK WHILE true END",
			wantError: false,
		},
		{
			name:      "continue inside do-while loop",
			input:     "DO CONTINUE WHILE true END",
			wantError: false,
		},
		{
			name:      "break inside nested while loops",
			input:     "WHILE WHILE BREAK END END",
			wantError: false,
		},
		{
			name:      "continue inside nested while loops",
			input:     "WHILE WHILE CONTINUE END END",
			wantError: false,
		},
		{
			name:      "multiple breaks in same loop",
			input:     "WHILE BREAK 42 BREAK END",
			wantError: false,
		},
		{
			name:      "multiple continues in same loop",
			input:     "WHILE CONTINUE 42 CONTINUE END",
			wantError: false,
		},
		{
			name:      "break/continue in different loops",
			input:     "WHILE BREAK END WHILE CONTINUE END",
			wantError: false,
		},
		{
			name:      "break/continue with other statements",
			input:     "WHILE VAR x 42 := x BREAK x DUMP CONTINUE END",
			wantError: false,
		},
		{
			name:      "break/continue in do-if statement inside while",
			input:     "WHILE DO BREAK ELSE CONTINUE IF true END END",
			wantError: false,
		},
		{
			name:      "deeply nested while loops with break/continue",
			input:     "WHILE WHILE WHILE BREAK CONTINUE END END END",
			wantError: false,
		},
		{
			name:      "break in while condition",
			input:     "DO 42 WHILE BREAK END",
			wantError: false,
		},
		{
			name:      "continue in while condition",
			input:     "DO 42 WHILE CONTINUE END",
			wantError: false,
		},
		{
			name:      "if statement without break/continue",
			input:     "DO 42 IF true END",
			wantError: false,
		},
		{
			name:      "if statement with break/continue outside while",
			input:     "WHILE DO 42 BREAK ELSE CONTINUE IF true END END",
			wantError: false,
		},

		// Invalid cases
		{
			name:          "break outside loop",
			input:         "BREAK",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue outside loop",
			input:         "CONTINUE",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "break in regular statements",
			input:         "42 BREAK 10",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue in regular statements",
			input:         "42 CONTINUE 10",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "break after while loop",
			input:         "WHILE 42 END BREAK",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue after while loop",
			input:         "WHILE 42 END CONTINUE",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "break in if statement outside while",
			input:         "DO BREAK IF true END",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue in if statement outside while",
			input:         "DO CONTINUE IF true END",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "break in else branch outside while",
			input:         "DO 42 ELSE BREAK IF true END",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue in else branch outside while",
			input:         "DO 42 ELSE CONTINUE IF true END",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "break in if condition outside while",
			input:         "DO 42 IF BREAK END",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue in if condition outside while",
			input:         "DO 42 IF CONTINUE END",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "multiple breaks outside loops",
			input:         "BREAK 42 BREAK",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "multiple continues outside loops",
			input:         "CONTINUE 42 CONTINUE",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
		{
			name:          "mixed break/continue outside loops",
			input:         "BREAK CONTINUE",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "break with variable declaration outside loop",
			input:         "VAR x BREAK := x",
			wantError:     true,
			errorContains: "BREAK can only be used inside a WHILE loop",
		},
		{
			name:          "continue with assignment outside loop",
			input:         "42 CONTINUE := VAR x",
			wantError:     true,
			errorContains: "CONTINUE can only be used inside a WHILE loop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			err = validateLoopControl(program)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateLoopControl() expected error, but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateLoopControl() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("validateLoopControl() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateLoopControlComplexCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		checkFunc func(*testing.T, error)
	}{
		{
			name:      "break in deeply nested while loops",
			input:     "WHILE WHILE WHILE WHILE BREAK END END END END",
			wantError: false,
		},
		{
			name:      "continue in deeply nested while loops",
			input:     "WHILE WHILE WHILE WHILE CONTINUE END END END END",
			wantError: false,
		},
		{
			name:      "mixed break/continue in nested loops",
			input:     "WHILE BREAK WHILE CONTINUE END BREAK END",
			wantError: false,
		},
		{
			name:      "break/continue in complex do-while structure",
			input:     "DO VAR x WHILE y BREAK CONTINUE END := x WHILE x BREAK END",
			wantError: false,
		},
		{
			name:      "break/continue in mixed control structures",
			input:     "WHILE DO BREAK ELSE CONTINUE IF true END DO BREAK WHILE false END END",
			wantError: false,
		},
		{
			name:      "break in array inside while loop",
			input:     "WHILE [1, 2, 3] END",
			wantError: false,
		},
		{
			name:      "break in assignment inside while loop",
			input:     "WHILE 42 := VAR x BREAK END",
			wantError: false,
		},
		{
			name:      "break with position information",
			input:     "BREAK",
			wantError: true,
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				// Check that error contains position information
				if !strings.Contains(err.Error(), "line") || !strings.Contains(err.Error(), "col") {
					t.Errorf("Error should contain position information: %v", err)
				}
			},
		},
		{
			name:      "continue with position information",
			input:     "CONTINUE",
			wantError: true,
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				// Check that error contains position information
				if !strings.Contains(err.Error(), "line") || !strings.Contains(err.Error(), "col") {
					t.Errorf("Error should contain position information: %v", err)
				}
			},
		},
		{
			name:      "while loop with all statement types",
			input:     "WHILE 42 \"hello\" true VAR x CONST y 10 := x [1, 2] + DUMP BREAK CONTINUE END",
			wantError: false,
		},
		{
			name:      "do-if with complex nesting",
			input:     "WHILE DO VAR x ELSE CONST y 10 IF x 0 == END BREAK END",
			wantError: false,
		},
		{
			name:      "multiple while loops with break/continue",
			input:     "WHILE BREAK END WHILE CONTINUE END WHILE BREAK CONTINUE END",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			program, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			err = validateLoopControl(program)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateLoopControl() expected error, but got none")
					return
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, err)
				}
				return
			}

			if err != nil {
				t.Errorf("validateLoopControl() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateLoopControlEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		program   *Program
		wantError bool
		checkFunc func(*testing.T, error)
	}{
		{
			name: "empty program",
			program: &Program{
				Statements: []Node{},
			},
			wantError: false,
		},
		{
			name: "program with only literals",
			program: &Program{
				Statements: []Node{
					&IntLiteral{Value: "42"},
					&StringLiteral{Value: "hello"},
					&BoolLiteral{Value: "true"},
				},
			},
			wantError: false,
		},
		{
			name: "while loop with break",
			program: &Program{
				Statements: []Node{
					&WhileStmt{
						Body: []Node{
							&BreakStmt{Pos: lexer.Position{Line: 1, Column: 1}},
						},
						Condition: nil,
						Position:  lexer.Position{Line: 1, Column: 1},
					},
				},
			},
			wantError: false,
		},
		{
			name: "while loop with continue",
			program: &Program{
				Statements: []Node{
					&WhileStmt{
						Body: []Node{
							&ContinueStmt{Pos: lexer.Position{Line: 1, Column: 1}},
						},
						Condition: nil,
						Position:  lexer.Position{Line: 1, Column: 1},
					},
				},
			},
			wantError: false,
		},
		{
			name: "break outside loop",
			program: &Program{
				Statements: []Node{
					&BreakStmt{Pos: lexer.Position{Line: 2, Column: 5}},
				},
			},
			wantError: true,
			checkFunc: func(t *testing.T, err error) {
				expected := "BREAK can only be used inside a WHILE loop at line 2, col 5"
				if err.Error() != expected {
					t.Errorf("Error message = %q, want %q", err.Error(), expected)
				}
			},
		},
		{
			name: "continue outside loop",
			program: &Program{
				Statements: []Node{
					&ContinueStmt{Pos: lexer.Position{Line: 3, Column: 10}},
				},
			},
			wantError: true,
			checkFunc: func(t *testing.T, err error) {
				expected := "CONTINUE can only be used inside a WHILE loop at line 3, col 10"
				if err.Error() != expected {
					t.Errorf("Error message = %q, want %q", err.Error(), expected)
				}
			},
		},
		{
			name: "break in if statement outside while",
			program: &Program{
				Statements: []Node{
					&IfStmt{
						ThenBranch: []Node{
							&BreakStmt{Pos: lexer.Position{Line: 4, Column: 15}},
						},
						ElseBranch: nil,
						Condition:  nil,
						Position:   lexer.Position{Line: 4, Column: 1},
					},
				},
			},
			wantError: true,
			checkFunc: func(t *testing.T, err error) {
				expected := "BREAK can only be used inside a WHILE loop at line 4, col 15"
				if err.Error() != expected {
					t.Errorf("Error message = %q, want %q", err.Error(), expected)
				}
			},
		},
		{
			name: "continue in if else branch outside while",
			program: &Program{
				Statements: []Node{
					&IfStmt{
						ThenBranch: []Node{
							&IntLiteral{Value: "42"},
						},
						ElseBranch: []Node{
							&ContinueStmt{Pos: lexer.Position{Line: 5, Column: 20}},
						},
						Condition: nil,
						Position:  lexer.Position{Line: 5, Column: 1},
					},
				},
			},
			wantError: true,
			checkFunc: func(t *testing.T, err error) {
				expected := "CONTINUE can only be used inside a WHILE loop at line 5, col 20"
				if err.Error() != expected {
					t.Errorf("Error message = %q, want %q", err.Error(), expected)
				}
			},
		},
		{
			name: "break in if condition outside while",
			program: &Program{
				Statements: []Node{
					&IfStmt{
						ThenBranch: []Node{
							&IntLiteral{Value: "42"},
						},
						ElseBranch: nil,
						Condition: []Node{
							&BreakStmt{Pos: lexer.Position{Line: 6, Column: 25}},
						},
						Position: lexer.Position{Line: 6, Column: 1},
					},
				},
			},
			wantError: true,
			checkFunc: func(t *testing.T, err error) {
				expected := "BREAK can only be used inside a WHILE loop at line 6, col 25"
				if err.Error() != expected {
					t.Errorf("Error message = %q, want %q", err.Error(), expected)
				}
			},
		},
		{
			name: "nested while with correct depth",
			program: &Program{
				Statements: []Node{
					&WhileStmt{
						Body: []Node{
							&WhileStmt{
								Body: []Node{
									&BreakStmt{Pos: lexer.Position{Line: 7, Column: 1}},
									&ContinueStmt{Pos: lexer.Position{Line: 8, Column: 1}},
								},
								Condition: []Node{
									&BreakStmt{Pos: lexer.Position{Line: 9, Column: 1}},
								},
								Position: lexer.Position{Line: 7, Column: 1},
							},
						},
						Condition: nil,
						Position:  lexer.Position{Line: 7, Column: 1},
					},
				},
			},
			wantError: false,
		},
		{
			name: "if statement inside while with break/continue",
			program: &Program{
				Statements: []Node{
					&WhileStmt{
						Body: []Node{
							&IfStmt{
								ThenBranch: []Node{
									&BreakStmt{Pos: lexer.Position{Line: 10, Column: 1}},
								},
								ElseBranch: []Node{
									&ContinueStmt{Pos: lexer.Position{Line: 11, Column: 1}},
								},
								Condition: []Node{
									&BreakStmt{Pos: lexer.Position{Line: 12, Column: 1}},
								},
								Position: lexer.Position{Line: 10, Column: 1},
							},
						},
						Condition: nil,
						Position:  lexer.Position{Line: 10, Column: 1},
					},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLoopControl(tt.program)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateLoopControl() expected error, but got none")
					return
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, err)
				}
				return
			}

			if err != nil {
				t.Errorf("validateLoopControl() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateLoopControlDepthTracking(t *testing.T) {
	tests := []struct {
		name          string
		createProgram func() *Program
		wantError     bool
		description   string
	}{
		{
			name: "depth 0 - top level break",
			createProgram: func() *Program {
				return &Program{
					Statements: []Node{
						&BreakStmt{Pos: lexer.Position{Line: 1, Column: 1}},
					},
				}
			},
			wantError:   true,
			description: "break at depth 0 should error",
		},
		{
			name: "depth 1 - break inside while",
			createProgram: func() *Program {
				return &Program{
					Statements: []Node{
						&WhileStmt{
							Body: []Node{
								&BreakStmt{Pos: lexer.Position{Line: 1, Column: 1}},
							},
							Position: lexer.Position{Line: 1, Column: 1},
						},
					},
				}
			},
			wantError:   false,
			description: "break at depth 1 should be valid",
		},
		{
			name: "depth 2 - break inside nested while",
			createProgram: func() *Program {
				return &Program{
					Statements: []Node{
						&WhileStmt{
							Body: []Node{
								&WhileStmt{
									Body: []Node{
										&BreakStmt{Pos: lexer.Position{Line: 1, Column: 1}},
									},
									Position: lexer.Position{Line: 1, Column: 1},
								},
							},
							Position: lexer.Position{Line: 1, Column: 1},
						},
					},
				}
			},
			wantError:   false,
			description: "break at depth 2 should be valid",
		},
		{
			name: "depth reset - break outside nested while",
			createProgram: func() *Program {
				return &Program{
					Statements: []Node{
						&WhileStmt{
							Body: []Node{
								&IntLiteral{Value: "42"},
							},
							Position: lexer.Position{Line: 1, Column: 1},
						},
						&BreakStmt{Pos: lexer.Position{Line: 2, Column: 1}},
					},
				}
			},
			wantError:   true,
			description: "break after while should error (depth reset to 0)",
		},
		{
			name: "if statement preserves depth",
			createProgram: func() *Program {
				return &Program{
					Statements: []Node{
						&WhileStmt{
							Body: []Node{
								&IfStmt{
									ThenBranch: []Node{
										&BreakStmt{Pos: lexer.Position{Line: 1, Column: 1}},
									},
									Position: lexer.Position{Line: 1, Column: 1},
								},
							},
							Position: lexer.Position{Line: 1, Column: 1},
						},
					},
				}
			},
			wantError:   false,
			description: "break inside if within while should be valid (if preserves depth)",
		},
		{
			name: "complex nesting",
			createProgram: func() *Program {
				return &Program{
					Statements: []Node{
						&WhileStmt{
							Body: []Node{
								&IfStmt{
									ThenBranch: []Node{
										&WhileStmt{
											Body: []Node{
												&BreakStmt{Pos: lexer.Position{Line: 1, Column: 1}},
											},
											Position: lexer.Position{Line: 1, Column: 1},
										},
									},
									ElseBranch: []Node{
										&ContinueStmt{Pos: lexer.Position{Line: 2, Column: 1}},
									},
									Position: lexer.Position{Line: 1, Column: 1},
								},
							},
							Position: lexer.Position{Line: 1, Column: 1},
						},
					},
				}
			},
			wantError:   false,
			description: "complex nesting with proper depth tracking should be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := tt.createProgram()
			err := validateLoopControl(program)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateLoopControl() expected error for %s, but got none", tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("validateLoopControl() unexpected error for %s: %v", tt.description, err)
				}
			}
		})
	}
}
