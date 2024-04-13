package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

const (
	EOF = iota
	IDENTIFIER
	ILLEGAL

	// TYPES
	INT
	STRING
	BOOL
	FLOAT

	// MATH OPS
	ADD_OP
	SUB_OP
	MUL_OP
	DIV_OP
	POW_OP
	MOD_OP
	INC_OP
	DEC_OP

	// OPERATORS
	PRINT_OP
	DUMP_OP
)

var identifierMap = map[string]Token{
	// MATH OPS
	"+": ADD_OP,
	"-": SUB_OP,
	"*": MUL_OP,
	"/": DIV_OP,
	"^": POW_OP,
	"%": MOD_OP,

	// PRINT OPS
	"print": PRINT_OP,
	"dump":  DUMP_OP,
}

var tokenMap = map[Token]string{
	EOF:        "EOF",
	IDENTIFIER: "IDENTIFIER",
	ILLEGAL:    "ILLEGAL",
	INT:        "INT",
	STRING:     "STRING",
	FLOAT:      "FLOAT",
	BOOL:       "BOOL",
	ADD_OP:     "ADD_OP",
	PRINT_OP:   "PRINT_OP",
	DUMP_OP:    "DUMP_OP",
}

type Token int
type StackElement struct {
	Type     Token
	Value    string
	Position Position
}

func isDecimal(c rune) bool {
	return c == '.'
}

type Gorth struct {
	StrictMode     bool
	DebugMode      bool
	ExecutionStack []StackElement
}

func NewGorth() *Gorth {
	return &Gorth{
		StrictMode:     false,
		DebugMode:      false,
		ExecutionStack: make([]StackElement, 0),
	}
}

func (g *Gorth) Pop() (StackElement, error) {
	if len(g.ExecutionStack) == 0 {
		return StackElement{}, fmt.Errorf("error: cannot pop from an empty stack")
	}

	element := g.ExecutionStack[len(g.ExecutionStack)-1]
	g.ExecutionStack = g.ExecutionStack[:len(g.ExecutionStack)-1]

	return element, nil
}

func (g *Gorth) Push(e StackElement) {
	g.ExecutionStack = append(g.ExecutionStack, e)
}

func (g *Gorth) Peek() (StackElement, error) {
	if len(g.ExecutionStack) == 0 {
		return StackElement{}, fmt.Errorf("error: cannot peek an empty stack")
	}

	return g.ExecutionStack[len(g.ExecutionStack)-1], nil
}

func (g *Gorth) Print() error {
	// print the top of the stack
	val, err := g.Peek()

	if err != nil {
		return err
	}

	fmt.Print(val.Value)

	return nil
}

func (g *Gorth) Dump() error {
	// this function will print the topmost element of the stack and drop it
	val, err := g.Pop()

	if err != nil {
		return err
	}

	fmt.Println(val.Value)

	return nil
}

func (g *Gorth) Add() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	if val1.Type == INT && val2.Type == INT {
		// convert values to integers
		result1, err1 := strconv.Atoi(val1.Value)
		if err1 != nil {
			return err1
		}
		result2, err2 := strconv.Atoi(val2.Value)
		if err2 != nil {
			return err2
		}
		result := result1 + result2
		g.Push(StackElement{Type: INT, Value: strconv.Itoa(result), Position: val2.Position})
	} else if val1.Type == FLOAT && val2.Type == FLOAT {
		// convert values to floats
		result1, err1 := strconv.ParseFloat(val1.Value, 64)
		if err1 != nil {
			return err1
		}
		result2, err2 := strconv.ParseFloat(val2.Value, 64)
		if err2 != nil {
			return err2
		}

		result := result1 + result2
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else if val1.Type == STRING && val2.Type == STRING {
		result := val1.Value + val2.Value
		g.Push(StackElement{Type: STRING, Value: result, Position: val2.Position})
	} else if val1.Type == INT && val2.Type == FLOAT || val1.Type == FLOAT && val2.Type == INT {
		// convert values to floats
		result1, err1 := strconv.ParseFloat(val1.Value, 64)
		if err1 != nil {
			return err1
		}
		result2, err2 := strconv.ParseFloat(val2.Value, 64)
		if err2 != nil {
			return err2
		}
		result := result1 + result2
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else {
		return fmt.Errorf("error: cannot add %s and %s", tokenMap[val1.Type], tokenMap[val2.Type])
	}

	return nil
}

func (g *Gorth) Subtract() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	if val1.Type == INT && val2.Type == INT {
		result1, err := strconv.Atoi(val1.Value)
		if err != nil {
			return err
		}
		result2, err := strconv.Atoi(val2.Value)
		if err != nil {
			return err
		}

		result := result2 - result1
		g.Push(StackElement{Type: INT, Value: fmt.Sprintf("%v", result), Position: val2.Position})
	} else if val1.Type == FLOAT && val2.Type == FLOAT {
		result1, err := strconv.ParseFloat(val1.Value, 64)
		if err != nil {
			return err
		}
		result2, err := strconv.ParseFloat(val2.Value, 64)
		if err != nil {
			return err
		}
		result := result2 - result1
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else if val1.Type == INT && val2.Type == FLOAT || val1.Type == FLOAT && val2.Type == INT {
		result1, err := strconv.ParseFloat(val1.Value, 64)
		if err != nil {
			return err
		}
		result2, err := strconv.ParseFloat(val2.Value, 64)
		if err != nil {
			return err
		}
		result := result2 - result1
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else {
		return fmt.Errorf("error: cannot subtract %s and %s", tokenMap[val1.Type], tokenMap[val2.Type])
	}

	return nil
}

func (g *Gorth) Multiply() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	if val1.Type == INT && val2.Type == INT {
		result1, err := strconv.Atoi(val1.Value)
		if err != nil {
			return err
		}
		result2, err := strconv.Atoi(val2.Value)
		if err != nil {
			return err
		}
		result := result1 * result2
		g.Push(StackElement{Type: INT, Value: fmt.Sprintf("%v", result), Position: val2.Position})
	} else if val1.Type == FLOAT && val2.Type == FLOAT {
		result1, err := strconv.ParseFloat(val1.Value, 64)
		if err != nil {
			return err
		}
		result2, err := strconv.ParseFloat(val2.Value, 64)
		if err != nil {
			return err
		}
		result := result1 * result2
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else if val1.Type == INT && val2.Type == FLOAT || val1.Type == FLOAT && val2.Type == INT {
		result1, err := strconv.ParseFloat(val1.Value, 64)
		if err != nil {
			return err
		}
		result2, err := strconv.ParseFloat(val2.Value, 64)
		if err != nil {
			return err
		}
		result := result1 * result2
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else if val1.Type == STRING && val2.Type == INT {
		val, err := strconv.Atoi(val2.Value)
		if err != nil {
			return err
		}
		result := ""
		for i := 0; i < val; i++ {
			result += val1.Value
		}
		g.Push(StackElement{Type: STRING, Value: result, Position: val2.Position})
	} else if val1.Type == INT && val2.Type == STRING {
		val, err := strconv.Atoi(val1.Value)
		if err != nil {
			return err
		}
		result := ""
		for i := 0; i < val; i++ {
			result += val2.Value
		}
		g.Push(StackElement{Type: STRING, Value: result, Position: val2.Position})
	} else {
		return fmt.Errorf("error: cannot multiply %s and %s", tokenMap[val1.Type], tokenMap[val2.Type])
	}

	return nil
}

func (g *Gorth) ExecuteStack(p []StackElement) {
	for _, e := range p {
		switch e.Type {
		case PRINT_OP:
			err := g.Print()
			if err != nil {
				panic(err)
			}
		case DUMP_OP:
			err := g.Dump()
			if err != nil {
				panic(err)
			}
		case ADD_OP:
			err := g.Add()
			if err != nil {
				panic(err)
			}
		case SUB_OP:
			err := g.Subtract()
			if err != nil {
				panic(err)
			}
		case MUL_OP:
			err := g.Multiply()
			if err != nil {
				panic(err)
			}
		default:
			g.Push(e)
		}
	}

	if g.StrictMode {
		if len(g.ExecutionStack) > 1 {
			panic("error: execution stack must be empty at the end of the program")
		}
	}
	if g.DebugMode {
		fmt.Println(g.ExecutionStack)
	}
}

/**
 * LEXER START
 */
type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

// Returns a new lexer with the given reader
func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token, and returns the position of the token, type of the token, and the value of the token
func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			panic(err)
		}

		l.pos.column++

		switch {
		case r == '\n':
			l.resetPosition()
		case r == '\t':
			l.pos.column += 4
		case unicode.IsLetter(r):
			startPos := l.pos
			l.backup()
			token, lit := l.lexIdentifier()
			return startPos, token, lit
		case unicode.IsDigit(r):
			startPos := l.pos
			l.backup()
			token, lit := l.lexNumber()
			return startPos, token, lit
		case unicode.IsSpace(r):
			continue
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			switch r {
			case '+':
				return l.pos, ADD_OP, string(r)
			case '-':
				return l.pos, SUB_OP, string(r)
			case '*':
				return l.pos, MUL_OP, string(r)
			case '/':
				return l.pos, DIV_OP, string(r)
			case '^':
				return l.pos, POW_OP, string(r)
			case '%':
				return l.pos, MOD_OP, string(r)
			case '"':
				startPos := l.pos
				token, lit := l.lexString()
				return startPos, token, lit
			}
		default:
			fmt.Printf("unknown rune: %v\n", string(r))
			return l.pos, ILLEGAL, string(r)
		}
	}
}

// backup moves the reader back one rune
func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) lexIdentifier() (Token, string) {
	var lit string

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// check if lit does not exist in identifierMap
				if _, ok := identifierMap[lit]; !ok {
					panic(fmt.Errorf("unknown identifier: %s at line %d column %d", lit, l.pos.line, l.pos.column))
				}

				return identifierMap[lit], lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()

			// check if lit does not exist in identifierMap
			if _, ok := identifierMap[lit]; !ok {
				panic(fmt.Errorf("unknown identifier: %s at line %d column %d", lit, l.pos.line, l.pos.column))
			}

			return identifierMap[lit], lit
		}
	}
}

func (l *Lexer) peekChar() rune {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0
		}

		panic(err)
	}

	l.backup()

	return r
}

func (l *Lexer) lexNumber() (Token, string) {
	var lit string
	var tokenType Token = INT
	position := l.pos

	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return tokenType, lit
			}

			panic(err)
		}

		if unicode.IsDigit(r) {
			lit += string(r)
		} else if isDecimal(r) {
			if tokenType == INT {
				tokenType = FLOAT
				lit += string(r)
			} else {
				panic(fmt.Errorf("unexpected decimal point at line %d column %d", position.line, position.column))
			}
		} else {
			l.backup()
			break
		}
	}

	return tokenType, lit
}

// lexString scans the input for a string, and returns the string as a string
func (l *Lexer) lexString() (Token, string) {
	var lit string

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return STRING, lit
			}

			panic(err)
		}

		if r == '"' {
			break
		}

		lit += string(r)
	}

	return STRING, lit
}

// Resets the position of the lexer to the beginning of the line
func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

/**
 * LEXER END
 */

/**
 * PARSER START
 */
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(pos Position, tok Token, lit string) (StackElement, error) {
	switch tok {
	default:
		return StackElement{Type: tok, Value: lit, Position: pos}, nil
	case ILLEGAL:
		return StackElement{}, fmt.Errorf("unknown token: %d at line %v col %v", tok, pos.line, pos.column)
	}
}

/**
 * PARSER END
 */

func main() {
	file, err := os.Open("./examples/hello_world.gorth")

	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	parser := NewParser()
	gorth := NewGorth()
	var program []StackElement = make([]StackElement, 0)

	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		element, err := parser.Parse(pos, tok, lit)

		if err != nil {
			panic(fmt.Errorf("error parsing token: %v", err))
		}

		// add the parsed token to the program stack
		program = append(program, element)

		// fmt.Printf("Token: %v, Literal: %s, Line: %d, Column: %d\n", tokenMap[tok], lit, pos.line, pos.column)
	}

	gorth.ExecuteStack(program)
}
