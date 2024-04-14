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
	"syscall"
	"unicode"
	"unsafe"
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
	VARIABLE

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

	// ASSIGNMENT
	ASSIGN_OP
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

	// ASSIGNMENT
	"=": ASSIGN_OP,

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
	VARIABLE:   "VARIABLE",

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

	// ASSIGNMENT
	ASSIGN_OP: "ASSIGN_OP",
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

type ArithmeticFunc func(float64, float64) (float64, error)

var Add ArithmeticFunc = func(a, b float64) (float64, error) {
	return a + b, nil
}

var Multiply ArithmeticFunc = func(a, b float64) (float64, error) {
	return a * b, nil
}

var Subtract ArithmeticFunc = func(a, b float64) (float64, error) {
	return b - a, nil
}

var Divide ArithmeticFunc = func(a, b float64) (float64, error) {
	return b / a, nil
}

var Pow ArithmeticFunc = func(a, b float64) (float64, error) {
	return math.Pow(b, a), nil
}

var Mod ArithmeticFunc = func(a, b float64) (float64, error) {
	return math.Mod(b, a), nil
}

var Decrement ArithmeticFunc = func(a, b float64) (float64, error) {
	return b - a, nil
}

func performIntArithmetic(val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	result1, err1 := strconv.Atoi(val1.Value)
	if err1 != nil {
		return StackElement{}, err1
	}

	result2, err2 := strconv.Atoi(val2.Value)
	if err2 != nil {
		return StackElement{}, err2
	}

	result, err := op(float64(result2), float64(result1))
	if err != nil {
		return StackElement{}, err
	}

	return StackElement{Type: INT, Value: strconv.Itoa(int(result)), Position: val2.Position}, nil
}

func performFloatArithmetic(val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	result1, err1 := strconv.ParseFloat(val1.Value, 64)
	if err1 != nil {
		return StackElement{}, err1
	}

	result2, err2 := strconv.ParseFloat(val2.Value, 64)
	if err2 != nil {
		return StackElement{}, err2
	}

	result, err := op(result2, result1)
	if err != nil {
		return StackElement{}, err
	}

	return StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position}, nil
}

func performStringArithmetic(val1, val2 StackElement) (StackElement, error) {
	result := val2.Value + val1.Value
	return StackElement{Type: STRING, Value: result, Position: val2.Position}, nil
}

func performMixedArithmetic(val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	result1, err1 := strconv.ParseFloat(val1.Value, 64)
	if err1 != nil {
		return StackElement{}, err1
	}

	result2, err2 := strconv.ParseFloat(val2.Value, 64)
	if err2 != nil {
		return StackElement{}, err2
	}

	result, err := op(result2, result1)
	if err != nil {
		return StackElement{}, err
	}

	return StackElement{Type: FLOAT, Value: fmt.Sprintf("%f", result), Position: val2.Position}, nil
}

func performVariableArithmetic(g *Gorth, val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	if _, ok := (*g.VariableMap)[val1.Value]; !ok {
		return StackElement{}, fmt.Errorf("error: variable %s is not defined", val1.Value)
	}

	if _, ok := (*g.VariableMap)[val2.Value]; !ok {
		return StackElement{}, fmt.Errorf("error: variable %s is not defined", val2.Value)
	}

	val1Value := (*g.VariableMap)[val1.Value].Value
	val2Value := (*g.VariableMap)[val2.Value].Value

	var result float64
	var err error

	switch {
	case val1Value.Type == INT && val2Value.Type == INT:
		result1, err1 := strconv.Atoi(val1Value.Value)
		if err1 != nil {
			return StackElement{}, err1
		}

		result2, err2 := strconv.Atoi(val2Value.Value)
		if err2 != nil {
			return StackElement{}, err2
		}

		result, err = op(float64(result1), float64(result2))
	case val1Value.Type == FLOAT && val2Value.Type == FLOAT:
		result1, err1 := strconv.ParseFloat(val1Value.Value, 64)
		if err1 != nil {
			return StackElement{}, err1
		}

		result2, err2 := strconv.ParseFloat(val2Value.Value, 64)
		if err2 != nil {
			return StackElement{}, err2
		}

		result, err = op(result1, result2)
	case (val1Value.Type == INT && val2Value.Type == FLOAT) || (val1Value.Type == FLOAT && val2Value.Type == INT):
		result1, err1 := strconv.ParseFloat(val1Value.Value, 64)
		if err1 != nil {
			return StackElement{}, err1
		}

		result2, err2 := strconv.ParseFloat(val2Value.Value, 64)
		if err2 != nil {
			return StackElement{}, err2
		}

		result, err = op(result1, result2)
	case val1Value.Type == STRING && val2Value.Type == STRING:
		result := val1Value.Value + val2Value.Value

		(*g.VariableMap)[val2.Value] = Variable{Name: val2.Value, Value: StackElement{Type: STRING, Value: result, Position: val2.Position}, Const: false}
		return StackElement{Type: VARIABLE, Value: val2.Value, Position: val2.Position}, nil
	default:
		return StackElement{}, fmt.Errorf("error: cannot perform %s on %s and %s", fmt.Sprintf("%v", op), tokenMap[val1Value.Type], tokenMap[val2Value.Type])
	}

	if err != nil {
		return StackElement{}, err
	}

	var resultType Token
	if val1Value.Type == INT && val2Value.Type == INT {
		resultType = INT
	} else {
		resultType = FLOAT
	}

	(*g.VariableMap)[val2.Value] = Variable{Name: val2.Value, Value: StackElement{Type: resultType, Value: fmt.Sprintf("%v", result), Position: val2.Position}, Const: false}
	return StackElement{Type: VARIABLE, Value: val2.Value, Position: val2.Position}, nil
}

func performVariableAndValueArithmetic(g *Gorth, val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	if _, ok := (*g.VariableMap)[val1.Value]; !ok {
		return StackElement{}, fmt.Errorf("error: variable %s is not defined", val1.Value)
	}

	val1Value := (*g.VariableMap)[val1.Value].Value

	var result float64
	var err error

	fmt.Println("VAL 1: ", val1Value, " Type: ", tokenMap[val1Value.Type])
	fmt.Println("VAL 2: ", val2, " Type: ", tokenMap[val2.Type])

	switch {
	case val1Value.Type == INT && val2.Type == INT:
		result1, err1 := strconv.Atoi(val1Value.Value)
		if err1 != nil {
			return StackElement{}, err1
		}

		result2, err2 := strconv.Atoi(val2.Value)
		if err2 != nil {
			return StackElement{}, err2
		}

		result, err = op(float64(result2), float64(result1))
	case val1Value.Type == FLOAT && val2.Type == FLOAT:
		result1, err1 := strconv.ParseFloat(val1Value.Value, 64)
		if err1 != nil {
			return StackElement{}, err1
		}

		result2, err2 := strconv.ParseFloat(val2.Value, 64)
		if err2 != nil {
			return StackElement{}, err2
		}

		result, err = op(result2, result1)
	case (val1Value.Type == INT && val2.Type == FLOAT) || (val1Value.Type == FLOAT && val2.Type == INT):
		result1, err1 := strconv.ParseFloat(val1Value.Value, 64)
		if err1 != nil {
			return StackElement{}, err1
		}

		result2, err2 := strconv.ParseFloat(val2.Value, 64)
		if err2 != nil {
			return StackElement{}, err2
		}

		// did this because if the op is subtraction or division, the order matters
		result, err = op(result2, result1)
	case val1Value.Type == STRING && val2.Type == STRING:
		result := val2.Value + val1Value.Value

		(*g.VariableMap)[val1.Value] = Variable{Name: val1.Value, Value: StackElement{Type: STRING, Value: result, Position: val1Value.Position}, Const: false}
		return StackElement{Type: VARIABLE, Value: val1.Value, Position: val1Value.Position}, nil
	default:
		return StackElement{}, fmt.Errorf("error: cannot perform %s on %s and %s", fmt.Sprintf("%v", op), tokenMap[val1Value.Type], tokenMap[val2.Type])
	}

	if err != nil {
		return StackElement{}, err
	}

	var resultType Token
	if val1Value.Type == INT && val2.Type == INT {
		resultType = INT
	} else {
		resultType = FLOAT
	}

	(*g.VariableMap)[val1.Value] = Variable{Name: val1.Value, Value: StackElement{Type: resultType, Value: fmt.Sprintf("%v", result), Position: val1Value.Position}, Const: false}
	return StackElement{Type: VARIABLE, Value: val1.Value, Position: val1Value.Position}, nil
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

type Variable struct {
	Name  string
	Value StackElement
	Const bool
}

type Gorth struct {
	StrictMode     bool
	DebugMode      bool
	ExecutionStack []StackElement
	Parser         *Parser
	Lexer          *Lexer
	VariableMap    *map[string]Variable
}

func NewGorth(s bool, d bool, p *Parser, l *Lexer) *Gorth {
	return &Gorth{
		StrictMode:     s,
		DebugMode:      d,
		ExecutionStack: make([]StackElement, 0),
		Parser:         p,
		Lexer:          l,
		VariableMap:    &map[string]Variable{},
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

func (g *Gorth) PopValues() (StackElement, StackElement, error) {
	if len(g.ExecutionStack) < 2 {
		return StackElement{}, StackElement{}, fmt.Errorf("error: cannot pop values from an empty stack")
	}

	val1, err := g.Pop()
	if err != nil {
		return StackElement{}, StackElement{}, err
	}

	val2, err := g.Pop()
	if err != nil {
		return StackElement{}, StackElement{}, err
	}

	return val1, val2, nil
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

	if val.Type == VARIABLE {
		// check if the variable is defined
		if _, ok := (*g.VariableMap)[val.Value]; !ok {
			return fmt.Errorf("error: variable %s is not defined", val.Value)
		}

		val = (*g.VariableMap)[val.Value].Value
	}

	if err != nil {
		return err
	}

	var sysret syscall.Errno

	bytes := append([]byte(val.Value), '\n')
	_, _, sysret = syscall.Syscall(syscall.SYS_WRITE, 1, uintptr(unsafe.Pointer(&bytes[0])), uintptr(len(bytes)))
	if sysret != 0 {
		return fmt.Errorf("error: %v", syscall.Errno(-sysret))
	}

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

func (g *Gorth) PerformArithmetic(val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	switch {
	case val1.Type == INT && val2.Type == INT:
		return performIntArithmetic(val1, val2, op)
	case val1.Type == FLOAT && val2.Type == FLOAT:
		return performFloatArithmetic(val1, val2, op)
	case val1.Type == STRING && val2.Type == STRING:
		return performStringArithmetic(val1, val2)
	case (val1.Type == INT && val2.Type == FLOAT) || (val1.Type == FLOAT && val2.Type == INT):
		return performMixedArithmetic(val1, val2, op)
	case val1.Type == VARIABLE && val2.Type == VARIABLE:
		return performVariableArithmetic(g, val1, val2, op)
	case val1.Type == VARIABLE:
		return performVariableAndValueArithmetic(g, val1, val2, op)
	case val2.Type == VARIABLE:
		return performVariableAndValueArithmetic(g, val2, val1, op)
	default:
		return StackElement{}, fmt.Errorf("error: cannot perform op on %s and %s", tokenMap[val1.Type], tokenMap[val2.Type])
	}
}

func (g *Gorth) Add() error {
	val1, val2, err := g.PopValues()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(val1, val2, Add)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Subtract() error {
	val1, val2, err := g.PopValues()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(val1, val2, Subtract)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Multiply() error {
	val1, val2, err := g.PopValues()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(val1, val2, Multiply)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Divide() error {
	val1, val2, err := g.PopValues()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(val1, val2, Divide)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Pow() error {
	val1, val2, err := g.PopValues()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(val1, val2, Pow)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Mod() error {
	val1, val2, err := g.PopValues()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(val1, val2, Mod)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Increment() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(val, StackElement{Type: INT, Value: "1", Position: val.Position}, Add)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Decrement() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	result, err := g.PerformArithmetic(StackElement{Type: INT, Value: "1", Position: val.Position}, val, Decrement)
	if err != nil {
		return err
	}

	g.Push(result)
	return nil
}

func (g *Gorth) Drop() error {
	if len(g.ExecutionStack) == 0 {
		return fmt.Errorf("error: cannot drop from an empty stack")
	}

	val, err := g.Pop()
	if err != nil {
		return err
	}

	if val.Type == VARIABLE {
		if _, ok := (*g.VariableMap)[val.Value]; !ok {
			return fmt.Errorf("error: variable %s is not defined", val.Value)
		}

		delete(*g.VariableMap, val.Value)
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

func (g *Gorth) AssignVar() error {
	if len(g.ExecutionStack) < 2 {
		return fmt.Errorf("error: cannot assign variable with less than 2 elements in the stack")
	}

	// variable name
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	// variable value
	val2, err := g.Pop()
	if err != nil {
		return err
	}

	// check if the variable is already defined
	if _, ok := (*g.VariableMap)[val1.Value]; ok && (*g.VariableMap)[val1.Value].Const {
		return fmt.Errorf("error: variable %s is already defined", val1.Value)
	}

	(*g.VariableMap)[val1.Value] = Variable{Name: val1.Value, Value: val2, Const: false}

	// push the variable value to the stack
	g.Push(StackElement{Type: VARIABLE, Value: val1.Value, Position: val2.Position})

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

		// ASSIGNMENT
		case ASSIGN_OP:
			err := g.AssignVar()
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
		fmt.Println("Variable Map at the end of execution: \n\t", *g.VariableMap)
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
			case '=':
				return l.pos, ASSIGN_OP, string(r)
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
				if IsOperator(lit) {
					return operatorMap[lit], lit
				} else {
					return IDENTIFIER, lit
				}
			}

			panic(err)
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()

			// check if lit does not exist in operatorMap
			if !IsOperator(lit) {
				var nextR rune
				var err error
				var counts int

				// Skip all whitespace runes until we get to the next rune
				for {
					nextR, _, err = l.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							break
						}
						panic(err)
					}

					if !unicode.IsSpace(nextR) {
						break
					}
					counts++
				}

				// Backup the reader
				for i := 0; i < counts; i++ {
					l.backup()
				}

				return VARIABLE, lit
			}

			return operatorMap[lit], lit
		}
	}
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
	unaryOps := "print|drop|dup|dump|inc|dec"
	binaryOps := "+|-|*|/|mod|pow|swap|over|="
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

	// Combine unary operations if any
	for len(stack) > 1 {
		// Pop the top two nodes from the stack
		right := stack[len(stack)-1]
		left := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		// Create a new node for the binary operation
		node := &Node{Value: "", Left: left, Right: right}

		// Push the binary operation node onto the stack
		stack = append(stack, node)
	}

	if len(stack) != 1 {
		for _, node := range stack {
			// get the pointer to the node
			if node != nil {
				fmt.Printf("Node: %v\n", node.Value)
			}
		}
		return nil, fmt.Errorf("invalid expression: %v", s)
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

	gorth := NewGorth(*strictMode, *debugMode, NewParser(), NewLexer(file))
	var program []StackElement = make([]StackElement, 0)

	for {
		pos, tok, lit := gorth.Lexer.Lex()
		if tok == EOF {
			break
		}

		element, err := gorth.Parser.Parse(pos, tok, lit)

		if err != nil {
			panic(fmt.Errorf("error parsing token: %v", err))
		}

		// add the parsed token to the program stack
		program = append(program, element)

		// fmt.Printf("Token: %v, Literal: %s, Line: %d, Column: %d\n", tokenMap[tok], lit, pos.line, pos.column)
	}

	// print the AST
	if gorth.DebugMode {
		root, err := gorth.Parser.BuildAST(program)
		if err != nil {
			panic(fmt.Errorf("error building AST: %v", err))
		}

		fmt.Println("Program AST: ")
		gorth.Parser.PrintAST(root, "")
	}

	gorth.ExecuteStack(program)
}
