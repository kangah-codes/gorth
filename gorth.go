package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
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

	// STACK MANIPULATION
	DROP_OP
	SWAP_OP
	DUP_OP
	OVER_OP
	ROT_OP

	// OPERATORS
	PRINT_OP
	DUMP_OP
)

var operatorMap = map[string]Token{
	// MATH OPS
	"+":   ADD_OP,
	"-":   SUB_OP,
	"*":   MUL_OP,
	"/":   DIV_OP,
	"^":   POW_OP,
	"%":   MOD_OP,
	"inc": INC_OP,
	"dec": DEC_OP,
	"mod": MOD_OP,
	"pow": POW_OP,

	// STACK MANIPULATION
	"drop": DROP_OP,
	"swap": SWAP_OP,
	"dup":  DUP_OP,
	"over": OVER_OP,
	"rot":  ROT_OP,

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

	// MATH OPS
	ADD_OP: "ADD_OP",

	// PRINT OPS
	PRINT_OP: "PRINT_OP",
	DUMP_OP:  "DUMP_OP",

	// STACK MANIPULATION
	DROP_OP: "DROP_OP",
	SWAP_OP: "SWAP_OP",
	DUP_OP:  "DUP_OP",
	OVER_OP: "OVER_OP",
	ROT_OP:  "ROT_OP",
}

type Token int
type StackElement struct {
	Type     Token
	Value    string
	Position Position
}

type Node struct {
	Value string
	Left  *Node
	Right *Node
}

func IsOperator(s string) bool {
	_, ok := operatorMap[s]
	return ok
}

func IsDecimal(c rune) bool {
	return c == '.'
}

func PrintUsage() {
	fmt.Println("Usage: gorth <filename> [options]")
	fmt.Println("  filename: the name of the .gorth file to execute")
	fmt.Println("  options:")
	fmt.Println("    -d: optional enable debug mode")
	fmt.Println("    -s: optional enable strict mode")
}

type Gorth struct {
	StrictMode     bool
	DebugMode      bool
	ExecutionStack []StackElement
}

func NewGorth(s bool, d bool) *Gorth {
	return &Gorth{
		StrictMode:     s,
		DebugMode:      d,
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

	// var sysret syscall.Errno

	// bytes := []byte(val.Value + "\n")
	// _, _, sysret = syscall.Syscall(syscall.SYS_WRITE, 1, uintptr(unsafe.Pointer(&bytes[0])), uintptr(len(val.Value)))
	// if sysret < 0 {
	// 	return fmt.Errorf("error: %v", syscall.Errno(-sysret))
	// }
	//

	fmt.Println(val.Value)

	return nil
}

func (g *Gorth) Dump() error {
	err := g.Print()
	if err != nil {
		return err
	}

	err = g.Drop()
	if err != nil {
		return err
	}

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

func (g *Gorth) Divide() error {
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

		result := result2 / result1
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
		result := result2 / result1
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
		result := result2 / result1
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else {
		return fmt.Errorf("error: cannot divide %s and %s", tokenMap[val1.Type], tokenMap[val2.Type])
	}

	return nil
}

func (g *Gorth) Pow() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	if val1.Type == INT && val2.Type == INT {
		result1, err := strconv.ParseFloat(val1.Value, 64)
		if err != nil {
			return err
		}
		result2, err := strconv.ParseFloat(val2.Value, 64)
		if err != nil {
			return err
		}

		result := math.Pow(result2, result1)
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
		result := math.Pow(result2, result1)
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
		result := math.Pow(result2, result1)
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position})
	} else {
		return fmt.Errorf("error: cannot exponentiate %s and %s", tokenMap[val1.Type], tokenMap[val2.Type])
	}

	return nil
}

func (g *Gorth) Mod() error {
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

		result := result2 % result1
		g.Push(StackElement{Type: INT, Value: fmt.Sprintf("%v", result), Position: val2.Position})
	} else {
		return fmt.Errorf("error: cannot perform module on %s and %s", tokenMap[val1.Type], tokenMap[val2.Type])
	}

	return nil
}

func (g *Gorth) Increment() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	if val.Type != INT && val.Type != FLOAT {
		return fmt.Errorf("error: cannot increment %s", tokenMap[val.Type])
	}

	if val.Type == INT {
		result, err := strconv.Atoi(val.Value)
		if err != nil {
			return err
		}
		g.Push(StackElement{Type: INT, Value: fmt.Sprintf("%v", result+1), Position: val.Position})
	} else if val.Type == FLOAT {
		result, err := strconv.ParseFloat(val.Value, 64)
		if err != nil {
			return err
		}
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result+1.0), Position: val.Position})
	} else {
		return fmt.Errorf("error: cannot increment %s", tokenMap[val.Type])
	}

	return nil
}

func (g *Gorth) Decrement() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	if val.Type != INT && val.Type != FLOAT {
		return fmt.Errorf("error: cannot decrement %s", tokenMap[val.Type])
	}

	if val.Type == INT {
		result, err := strconv.Atoi(val.Value)
		if err != nil {
			return err
		}
		g.Push(StackElement{Type: INT, Value: fmt.Sprintf("%v", result-1), Position: val.Position})
	} else if val.Type == FLOAT {
		result, err := strconv.ParseFloat(val.Value, 64)
		if err != nil {
			return err
		}
		g.Push(StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result-1.0), Position: val.Position})
	} else {
		return fmt.Errorf("error: cannot decrement %s", tokenMap[val.Type])
	}

	return nil
}

func (g *Gorth) Drop() error {
	if len(g.ExecutionStack) == 0 {
		return fmt.Errorf("error: cannot drop from an empty stack")
	}

	_, err := g.Pop()
	if err != nil {
		return err
	}

	return nil
}

func (g *Gorth) Swap() error {
	if len(g.ExecutionStack) < 2 {
		return fmt.Errorf("error: cannot swap with less than 2 elements in the stack")
	}

	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	g.Push(val1)
	g.Push(val2)

	return nil
}

func (g *Gorth) Dup() error {
	if len(g.ExecutionStack) == 0 {
		return fmt.Errorf("error: cannot duplicate from an empty stack")
	}

	val, err := g.Peek()
	if err != nil {
		return err
	}

	g.Push(val)

	return nil
}

func (g *Gorth) Over() error {
	if len(g.ExecutionStack) < 2 {
		return fmt.Errorf("error: cannot perform over with less than 2 elements in the stack")
	}

	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	g.Push(val2)
	g.Push(val1)
	g.Push(val2)

	return nil
}

func (g *Gorth) Rot() error {
	if len(g.ExecutionStack) < 3 {
		return fmt.Errorf("error: cannot perform rot with less than 3 elements in the stack")
	}

	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	val3, err := g.Pop()
	if err != nil {
		return err
	}

	g.Push(val2)
	g.Push(val1)
	g.Push(val3)

	return nil
}

func (g *Gorth) ExecuteStack(p []StackElement) {
	for _, e := range p {
		switch e.Type {
		// MISC OPERATIONS
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

		// ARITHMETIC OPERATIONS
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
		case DIV_OP:
			err := g.Divide()
			if err != nil {
				panic(err)
			}
		case POW_OP:
			err := g.Pow()
			if err != nil {
				panic(err)
			}
		case MOD_OP:
			err := g.Mod()
			if err != nil {
				panic(err)
			}
		case INC_OP:
			err := g.Increment()
			if err != nil {
				panic(err)
			}
		case DEC_OP:
			err := g.Decrement()
			if err != nil {
				panic(err)
			}

		// STACK OPERATIONS
		case DROP_OP:
			err := g.Drop()
			if err != nil {
				panic(err)
			}
		case SWAP_OP:
			err := g.Swap()
			if err != nil {
				panic(err)
			}
		case DUP_OP:
			err := g.Dup()
			if err != nil {
				panic(err)
			}
		case OVER_OP:
			err := g.Over()
			if err != nil {
				panic(err)
			}
		case ROT_OP:
			err := g.Rot()
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
		fmt.Println("Program Stack at the end of execution: \n\t", g.ExecutionStack)
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
				// check if lit does not exist in operatorMap
				if _, ok := operatorMap[lit]; !ok {
					panic(fmt.Errorf("unknown identifier: %s at line %d column %d", lit, l.pos.line, l.pos.column))
				}

				return operatorMap[lit], lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()

			// check if lit does not exist in operatorMap
			if _, ok := operatorMap[lit]; !ok {
				panic(fmt.Errorf("unknown identifier: %s at line %d column %d", lit, l.pos.line, l.pos.column))
			}

			return operatorMap[lit], lit
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
		} else if IsDecimal(r) {
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

func (p *Parser) BuildAST(s []StackElement) (*Node, error) {
	var stack []*Node
	unaryOps := "print|drop|dup|inc|dec"
	binaryOps := "+|-|*|/|mod|pow|swap|over"
	ternaryOps := "rot"

	// assert that all ops are included
	unOps := strings.Split(unaryOps, "|")
	for _, op := range unOps {
		_, ok := operatorMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator: %s, did you perhaps forget to add it to the operatorMap?", op)
		}
	}

	binOps := strings.Split(binaryOps, "|")
	for _, op := range binOps {
		_, ok := operatorMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator: %s, did you perhaps forget to add it to the operatorMap?", op)
		}
	}

	terOps := strings.Split(ternaryOps, "|")
	for _, op := range terOps {
		_, ok := operatorMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator: %s, did you perhaps forget to add it to the operatorMap?", op)
		}
	}

	for _, element := range s {
		if IsOperator(element.Value) {
			switch {
			case strings.Contains(unaryOps, element.Value):
				if len(stack) < 1 {
					return nil, fmt.Errorf("syntax error: %s requires at least 1 element on the stack\nline: %v col: %v", element.Value, element.Position.line, element.Position.column)
				}

				operand := stack[len(stack)-1]
				stack = stack[:len(stack)-1] // Remove the operand from the stack
				node := &Node{Value: element.Value, Left: operand}
				stack = append(stack, node)
			case strings.Contains(binaryOps, element.Value):
				if len(stack) < 2 {
					return nil, fmt.Errorf("syntax error: %s requires at least 2 elements on the stack\nline: %v col: %v", element.Value, element.Position.line, element.Position.column)
				}

				right := stack[len(stack)-1]
				left := stack[len(stack)-2]
				node := &Node{Value: element.Value, Left: left, Right: right}
				stack = stack[:len(stack)-2] // Remove two operands from the stack
				stack = append(stack, node)
			case strings.Contains(ternaryOps, element.Value):
				return nil, fmt.Errorf("ternary operators are not supported yet")
			}
		} else {
			stack = append(stack, &Node{Value: element.Value})
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}

func (p *Parser) PrintAST(root *Node, indent string) {
	if root == nil {
		return
	}

	fmt.Printf("%s%s\n", indent, root.Value)
	p.PrintAST(root.Left, indent+"  ")
	p.PrintAST(root.Right, indent+"  ")
}

/**
 * PARSER END
 */

func main() {
	// get the other arguments even if there are not in the correct order
	debugMode := flag.Bool("d", false, "enable debug mode")
	strictMode := flag.Bool("s", false, "enable strict mode")
	filePath := flag.String("f", "", "path to the file to execute")
	flag.Parse()

	// check if the first argument is a .gorth file
	if !strings.HasSuffix(*filePath, ".gorth") {
		panic(fmt.Sprintf("File %s is not a .gorth file", *filePath))
	}

	// check if the file exists
	_, err := os.Stat(*filePath)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("File %s does not exist", *filePath))
	}

	file, err := os.Open(*filePath)

	if err != nil {
		panic(err)
	}

	args := os.Args[1:]

	// check if there are no arguments
	if len(args) == 0 {
		PrintUsage()
		return
	}
	lexer := NewLexer(file)
	parser := NewParser()
	gorth := NewGorth(*strictMode, *debugMode)
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

	// build the AST
	root, err := parser.BuildAST(program)
	if err != nil {
		panic(fmt.Errorf("error building AST: %v", err))
	}

	// print the AST
	if gorth.DebugMode {
		parser.PrintAST(root, "")
	}

	gorth.ExecuteStack(program)
}
