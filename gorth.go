package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

const (
	EOF = iota
	ILLEGAL

	// TYPES
	INT
	STRING
	BOOL
	FLOAT
	PTR

	// POINTER MANIPULATION
	DEREF_OP

	// VARIABLE IDENTIFIER
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
	DELETE_OP

	// OPERATORS
	PRINT_OP
	PRINTLN_OP
	DUMP_OP

	// ASSIGNMENT
	ASSIGN_OP
)

const (
	PRIMITIVE_TYPES = "int|string|bool|float|ptr"
)

// gotta do this cos this idiotic language doesn't support constant arrays
var INCREMENTABLE_DEREMENTABLE_TYPES = []Token{INT, FLOAT, PTR}

var identifierMap = map[string]Token{
	// MATH OPS
	"+":   ADD_OP,
	"-":   SUB_OP,
	"*":   MUL_OP,
	"/":   DIV_OP,
	"^":   POW_OP,
	"%":   MOD_OP,
	"inc": INC_OP,
	"dec": DEC_OP,

	// ASSIGNMENT
	"=": ASSIGN_OP,

	// STACK MANIPULATION
	"drop": DROP_OP,
	"swap": SWAP_OP,
	"dup":  DUP_OP,
	"over": OVER_OP,
	"rot":  ROT_OP,

	// VARIABLE MANIPULATION
	"del": DELETE_OP,

	// PRINT OPS
	"print":   PRINT_OP,
	"println": PRINTLN_OP,
	"dump":    DUMP_OP,

	// TYPES
	"str":   STRING,
	"int":   INT,
	"bool":  BOOL,
	"float": FLOAT,
	"ptr":   PTR,
}

var tokenMap = map[Token]string{
	EOF:      "EOF",
	ILLEGAL:  "ILLEGAL",
	INT:      "INT",
	STRING:   "STRING",
	FLOAT:    "FLOAT",
	BOOL:     "BOOL",
	VARIABLE: "VARIABLE",
	PTR:      "PTR",

	// PTR OPS
	DEREF_OP: "DEREF_OP",

	// MATH OPS
	ADD_OP: "ADD_OP",

	// PRINT OPS
	PRINT_OP:   "PRINT_OP",
	PRINTLN_OP: "PRINTLN_OP",
	DUMP_OP:    "DUMP_OP",

	// STACK MANIPULATION
	DROP_OP: "DROP_OP",
	SWAP_OP: "SWAP_OP",
	DUP_OP:  "DUP_OP",
	OVER_OP: "OVER_OP",
	ROT_OP:  "ROT_OP",

	// ASSIGNMENT
	ASSIGN_OP: "ASSIGN_OP",

	// VARIABLE MANIPULATION
	DELETE_OP: "DELETE_OP",
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

func PtrToInt(p string) int {
	// remove * from the string
	result, err := strconv.Atoi(strings.Trim(p, "*"))
	if err != nil {
		panic(err)
	}

	return result
}

func PtreDerefToValue(p string) string {
	return strings.Trim(p, "&")
}

func ContainsValue[T any](arr []T, target T) bool {
	for _, value := range arr {
		if reflect.DeepEqual(value, target) {
			return true
		}
	}

	return false
}

func PerformIntArithmetic(val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
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

func PerformFloatArithmetic(val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
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

func PerformStringArithmetic(val1, val2 StackElement) (StackElement, error) {
	result := val2.Value + val1.Value
	return StackElement{Type: STRING, Value: result, Position: val2.Position}, nil
}

func PerformMixedArithmetic(val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
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

func PerformVariableArithmetic(g *Gorth, val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	if _, ok := (g.SemanticAnalyser.SymbolTable.Variables)[val1.Value]; !ok {
		return StackElement{}, fmt.Errorf("error: variable %s is not defined", val1.Value)
	}

	if _, ok := (g.SemanticAnalyser.SymbolTable.Variables)[val2.Value]; !ok {
		return StackElement{}, fmt.Errorf("error: variable %s is not defined", val2.Value)
	}

	val1Value := (g.SemanticAnalyser.SymbolTable.Variables)[val1.Value].Value
	val2Value := (g.SemanticAnalyser.SymbolTable.Variables)[val2.Value].Value

	var result float64
	var err error

	fmt.Println("Incrementing pointer ", tokenMap[val1.Type], " by ", tokenMap[val2.Type])

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
	// increasing the value of a ptr
	case val2Value.Type == PTR && val1Value.Type == INT:
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

		(g.SemanticAnalyser.SymbolTable.Variables)[val2.Value] = Variable{Name: val2.Value, Value: StackElement{Type: STRING, Value: result, Position: val2.Position}, Const: false}
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

	(g.SemanticAnalyser.SymbolTable.Variables)[val2.Value] = Variable{Name: val2.Value, Value: StackElement{Type: resultType, Value: fmt.Sprintf("%v", result), Position: val2.Position}, Const: false, Type: resultType}
	return StackElement{Type: VARIABLE, Value: val2.Value, Position: val2.Position}, nil
}

func PerformVariableAndValueArithmetic(g *Gorth, val1, val2 StackElement, op ArithmeticFunc) (StackElement, error) {
	if _, ok := (g.SemanticAnalyser.SymbolTable.Variables)[val1.Value]; !ok {
		return StackElement{}, fmt.Errorf("error: variable %s is not defined", val1.Value)
	}

	val1Value := (g.SemanticAnalyser.SymbolTable.Variables)[val1.Value].Value

	var result float64
	var err error

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
	// incrementing a variable pointer by an int
	// pointers are stored as ints and can only be incremented by ints
	// case val1Value.Type == PTR && val2.Type == INT:
	// 	result1 := PtrToInt(val1Value.Value)

	// 	result2, err2 := strconv.Atoi(val2.Value)
	// 	if err2 != nil {
	// 		return StackElement{}, err2
	// 	}

	// 	result, err = op(float64(result1), float64(result2))
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

		(g.SemanticAnalyser.SymbolTable.Variables)[val1.Value] = Variable{Name: val1.Value, Value: StackElement{Type: STRING, Value: result, Position: val1Value.Position}, Const: false, Type: STRING}
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

	(g.SemanticAnalyser.SymbolTable.Variables)[val1.Value] = Variable{Name: val1.Value, Value: StackElement{Type: resultType, Value: fmt.Sprintf("%v", result), Position: val1Value.Position}, Const: false, Type: resultType}
	return StackElement{Type: VARIABLE, Value: val1.Value, Position: val1Value.Position}, nil
}

func IsOperator(s string) bool {
	_, ok := identifierMap[s]
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
	Type  Token
}

type Gorth struct {
	StrictMode       bool
	DebugMode        bool
	ExecutionStack   []StackElement
	Parser           *Parser
	Lexer            *Lexer
	SemanticAnalyser *SemanticAnalyser
}

func NewGorth(s bool, d bool, p *Parser, l *Lexer, a *SemanticAnalyser) *Gorth {
	return &Gorth{
		StrictMode:       s,
		DebugMode:        d,
		ExecutionStack:   make([]StackElement, 0),
		Parser:           p,
		Lexer:            l,
		SemanticAnalyser: a,
	}
}

func (g *Gorth) ExecStackRepr() {
	// pretty print the stack
	fmt.Println("Stack:")
	for i := len(g.ExecutionStack) - 1; i >= 0; i-- {
		fmt.Printf("\t\tType: %v, Value: %v, Line: %v, Col: %v\n", tokenMap[g.ExecutionStack[i].Type], g.ExecutionStack[i].Value, g.ExecutionStack[i].Position.line, g.ExecutionStack[i].Position.column)
	}
	fmt.Printf("\t\t%v\n", g.ExecutionStack)
}

func (g *Gorth) Pop() (StackElement, error) {
	if len(g.ExecutionStack) == 0 {
		return StackElement{}, fmt.Errorf("error: cannot pop from an empty stack")
	}

	element := g.ExecutionStack[len(g.ExecutionStack)-1]
	g.ExecutionStack = g.ExecutionStack[:len(g.ExecutionStack)-1]

	return element, nil
}

func (g *Gorth) Dereference() error {
	// adds the value of the dereferenced ptr to the stack
	val, err := g.Pop()
	if err != nil {
		return err
	}

	var derefVal StackElement
	derefPtr := PtreDerefToValue(val.Value)
	derefType := g.SemanticAnalyser.InferType(derefPtr)

	// we're dereferencing by using a raw int value
	if derefType == INT {
		if PtrToInt(derefPtr) > len(g.ExecutionStack)-1 {
			return fmt.Errorf("error: execution stack index %v out of range with length %v", derefPtr, len(g.ExecutionStack))
		}

		ref := g.ExecutionStack[PtrToInt(derefPtr)]
		derefVal = ref
	} else {
		// we're dereferencing using a variable containing a ptr
		// first check if it exists
		if _, ok := g.SemanticAnalyser.SymbolTable.Variables[derefPtr]; !ok {
			return fmt.Errorf("error: variable %v is not defined", derefPtr)
		}

		derefVal = g.ExecutionStack[PtrToInt(g.SemanticAnalyser.SymbolTable.Variables[derefPtr].Value.Value)]
	}

	g.Push(derefVal)
	return nil

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
	if err != nil {
		return err
	}

	if val.Type == VARIABLE {
		// check if the variable is defined
		if _, ok := (g.SemanticAnalyser.SymbolTable.Variables)[PtreDerefToValue(val.Value)]; !ok {
			return fmt.Errorf("error: variable %s is not defined", val.Value)
		}

		// if we're dealing with a ptr deref get the actual value
		if strings.Contains(val.Value, "&") {
			val = g.ExecutionStack[PtrToInt(g.SemanticAnalyser.SymbolTable.Variables[PtreDerefToValue(val.Value)].Value.Value)]
		} else {
			val = (g.SemanticAnalyser.SymbolTable.Variables)[PtreDerefToValue(val.Value)].Value
		}
	}

	// we're dealing with a raw number deref
	if val.Type == INT && strings.Contains(val.Value, "&") {
		val = g.ExecutionStack[PtrToInt(PtreDerefToValue(val.Value))]
	}

	var sysret syscall.Errno

	bytes := []byte(val.Value)
	_, _, sysret = syscall.Syscall(syscall.SYS_WRITE, 1, uintptr(unsafe.Pointer(&bytes[0])), uintptr(len(bytes)))
	if sysret != 0 {
		return fmt.Errorf("error: %v", syscall.Errno(-sysret))
	}

	return nil
}

func (g *Gorth) Println() error {
	// print the top of the stack
	val, err := g.Peek()
	if err != nil {
		return err
	}

	if val.Type == VARIABLE {
		// check if the variable is defined
		if _, ok := (g.SemanticAnalyser.SymbolTable.Variables)[PtreDerefToValue(val.Value)]; !ok {
			return fmt.Errorf("error: variable %s is not defined", val.Value)
		}

		// if we're dealing with a ptr deref get the actual value
		if strings.Contains(val.Value, "&") {
			val = g.ExecutionStack[PtrToInt(g.SemanticAnalyser.SymbolTable.Variables[PtreDerefToValue(val.Value)].Value.Value)]
		} else {
			val = (g.SemanticAnalyser.SymbolTable.Variables)[PtreDerefToValue(val.Value)].Value
		}
	}

	// we're dealing with a raw number deref
	if val.Type == INT && strings.Contains(val.Value, "&") {
		val = g.ExecutionStack[PtrToInt(PtreDerefToValue(val.Value))]
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
		return PerformIntArithmetic(val1, val2, op)
	case val1.Type == FLOAT && val2.Type == FLOAT:
		return PerformFloatArithmetic(val1, val2, op)
	case val1.Type == STRING && val2.Type == STRING:
		return PerformStringArithmetic(val1, val2)
	case (val1.Type == INT && val2.Type == FLOAT) || (val1.Type == FLOAT && val2.Type == INT):
		return PerformMixedArithmetic(val1, val2, op)
	case val1.Type == PTR && val2.Type == INT || val2.Type == PTR && val1.Type == INT || val1.Type == PTR && val2.Type == PTR:
		return PerformIntArithmetic(val1, val2, op)
	case val1.Type == VARIABLE && val2.Type == VARIABLE:
		return PerformVariableArithmetic(g, val1, val2, op)
	case val1.Type == VARIABLE:
		return PerformVariableAndValueArithmetic(g, val1, val2, op)
	case val2.Type == VARIABLE:
		return PerformVariableAndValueArithmetic(g, val2, val1, op)
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

	// if value is a variable, point to it's value instead
	if val.Type == VARIABLE {
		// check if it exists
		if _, ok := g.SemanticAnalyser.SymbolTable.Variables[val.Value]; !ok {
			return fmt.Errorf("error: variable %v is not defined", val.Value)
		}

		// check if the variable value is actually an incrementable type
		if !ContainsValue(INCREMENTABLE_DEREMENTABLE_TYPES, g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Type) {
			return fmt.Errorf("error: variable %v is not an incrementable type. expected INT, FLOAT, PTR got %v instead", val.Value, tokenMap[g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Type])
		}

		// we turn it back into a string to appease the go gods... since literals are stored as strings
		// also i anticipated this would fail for regular variable int increments, but guess it still workds cause ptrtoint returns the string as a number
		result, err := g.PerformArithmetic(StackElement{Type: g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Type, Value: fmt.Sprintf("%v", PtrToInt(g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Value)), Position: g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Position}, StackElement{Type: INT, Value: "1", Position: val.Position}, Add)
		if err != nil {
			return err
		}

		g.SemanticAnalyser.SymbolTable.Variables[val.Value] = Variable{Value: result, Type: g.SemanticAnalyser.InferType(result.Value), Const: false, Name: val.Value}
		g.Push(result)
		return nil
	}

	// check if the variable value is actually an incrementable type
	if !ContainsValue(INCREMENTABLE_DEREMENTABLE_TYPES, val.Type) {
		return fmt.Errorf("error: %v is not an incrementable type. expected INT, FLOAT, PTR got %v instead", val.Value, tokenMap[val.Type])
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

	// if value is a variable, point to it's value instead
	if val.Type == VARIABLE {
		// check if it exists
		if _, ok := g.SemanticAnalyser.SymbolTable.Variables[val.Value]; !ok {
			return fmt.Errorf("error: variable %v is not defined", val.Value)
		}

		// check if the variable value is actually an incrementable type
		if !ContainsValue(INCREMENTABLE_DEREMENTABLE_TYPES, g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Type) {
			return fmt.Errorf("error: variable %v is not an decrementable type. expected INT, FLOAT, PTR got %v instead", val.Value, tokenMap[g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Type])
		}

		// we turn it back into a string to appease the go gods... since literals are stored as strings
		// also i anticipated this would fail for regular variable int increments, but guess it still workds cause ptrtoint returns the string as a number
		result, err := g.PerformArithmetic(StackElement{Type: INT, Value: "1", Position: val.Position}, StackElement{Type: g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Type, Value: fmt.Sprintf("%v", PtrToInt(g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Value)), Position: g.SemanticAnalyser.SymbolTable.Variables[val.Value].Value.Position}, Add)
		if err != nil {
			return err
		}

		g.SemanticAnalyser.SymbolTable.Variables[val.Value] = Variable{Value: result, Type: g.SemanticAnalyser.InferType(result.Value), Const: false, Name: val.Value}
		g.Push(result)
		return nil
	}

	// check if the variable value is actually an incrementable type
	if !ContainsValue(INCREMENTABLE_DEREMENTABLE_TYPES, val.Type) {
		return fmt.Errorf("error: %v is not an decrementable type. expected INT, FLOAT, PTR got %v instead", val.Value, tokenMap[val.Type])
	}

	result, err := g.PerformArithmetic(StackElement{Type: INT, Value: "1", Position: val.Position}, val, Subtract)
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

func (g *Gorth) Delete() error {
	if len(g.ExecutionStack) < 1 {
		return errors.New("error: cannot delete with less than 1 element in the stack")
	}

	val, err := g.Pop()
	if err != nil {
		return err
	}

	if val.Type == VARIABLE {
		if _, ok := g.SemanticAnalyser.SymbolTable.Variables[val.Value]; !ok {
			return fmt.Errorf("error: variable %s is not defined", val.Value)
		}

		delete(g.SemanticAnalyser.SymbolTable.Variables, val.Value)
	} else {
		return errors.New("error: cannot delete non-variable element, consider using drop instead?")
	}

	return nil
}

func (g *Gorth) AssignVar() error {
	if len(g.ExecutionStack) < 2 {
		return fmt.Errorf("error: cannot assign variable with less than 2 elements in the stack")
	}

	// potential variable value or type
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	// means we have a type declaration
	if _, ok := identifierMap[val1.Value]; ok && strings.Contains(PRIMITIVE_TYPES, val1.Value) {
		// val3 becomes the actual variable value
		val3, err := g.Pop()
		if err != nil {
			return err
		}

		// val4 becomes the variable name
		val4, err := g.Pop()
		if err != nil {
			return err
		}

		// check if the variable name is using a reserved keyword
		if _, ok := identifierMap[val4.Value]; ok {
			return fmt.Errorf("error: %v is a reserved keyword", val4.Value)
		}

		// check if the variable is already defined
		if _, ok := (g.SemanticAnalyser.SymbolTable.Variables)[strings.ToLower(val3.Value)]; ok && (g.SemanticAnalyser.SymbolTable.Variables)[val3.Value].Const {
			return fmt.Errorf("error: variable %s is already defined", val3.Value)
		}

		// always make strings lowercase in the variable map
		// in this language, variables will be case insensitive
		(g.SemanticAnalyser.SymbolTable.Variables)[strings.ToLower(val4.Value)] = Variable{Name: val4.Value, Value: val3, Const: false, Type: identifierMap[val1.Value]}

		// TODO: DO NOT PUSH DECLARED VARIABLES ONTO THE STACK, THEY WILL BE CONSUMED SO YOU HAVE TO ADD THEM WHEN YOU WANT TO USE THEM
		// push the variable value to the stack
		// g.Push(StackElement{Type: VARIABLE, Value: val3.Value, Position: val3.Position})

		return nil
	}

	// variable name
	val2, err := g.Pop()
	if err != nil {
		return err
	}

	// check if the variable name is using a reserved keyword
	if _, ok := identifierMap[val2.Value]; ok {
		return fmt.Errorf("error: %v is a reserved keyword", val2.Value)
	}

	// check if the variable is already defined
	if _, ok := (g.SemanticAnalyser.SymbolTable.Variables)[strings.ToLower(val2.Value)]; ok && (g.SemanticAnalyser.SymbolTable.Variables)[val2.Value].Const {
		return fmt.Errorf("error: variable %s is already defined", val2.Value)
	}

	// always make strings lowercase in the variable map
	// in this language, variables will be case insensitive
	(*&g.SemanticAnalyser.SymbolTable.Variables)[strings.ToLower(val2.Value)] = Variable{Name: val2.Value, Value: val1, Const: false, Type: g.SemanticAnalyser.InferType(val1.Value)}

	// TODO: DO NOT PUSH DECLARED VARIABLES ONTO THE STACK, THEY WILL BE CONSUMED SO YOU HAVE TO ADD THEM WHEN YOU WANT TO USE THEM
	// push the variable value to the stack
	// g.Push(StackElement{Type: VARIABLE, Value: val1.Value, Position: val2.Position})

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
		case PRINTLN_OP:
			err := g.Println()
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
		case DELETE_OP:
			err := g.Delete()
			if err != nil {
				panic(err)
			}

		// ASSIGNMENT
		case ASSIGN_OP:
			err := g.AssignVar()
			if err != nil {
				panic(err)
			}

		// DEREF_OP
		case DEREF_OP:
			g.Push(e)
			err := g.Dereference()
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
		fmt.Println("Program Stack at the end of execution:")
		fmt.Print("\t")
		g.ExecStackRepr()
		fmt.Println("Symbol Table at the end of execution:")
		fmt.Print("\t")
		g.SemanticAnalyser.VariablesRepr()
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
		pos:    Position{line: 1, column: 1},
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
			l.JumpToNextLine()
		case r == '\t':
			// a tab is 4 spaces
			l.pos.column += 4
		case unicode.IsLetter(r):
			startPos := l.pos
			l.Backup()
			token, lit := l.LexIdentifier()
			return startPos, token, lit
		case unicode.IsDigit(r):
			startPos := l.pos
			l.Backup()
			token, lit := l.LexNumber()
			return startPos, token, lit
		case unicode.IsSpace(r):
			continue
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			switch r {
			case '#':
				// means the entire line is a comment
				// consume the rest of the line
				for {
					r, _, err = l.reader.ReadRune()
					if err != nil || r == '\n' {
						break
					}
				}

				l.JumpToNextLine()
			case '+':
				return l.pos, ADD_OP, string(r)
			case '-':
				return l.pos, SUB_OP, string(r)
			case '*':
				nextR, err := l.PeekNextChar()
				if err != nil {
					panic(err)
				}

				if unicode.IsDigit(nextR) {
					l.Backup()
					token, digit := l.LexNumber()

					if token == FLOAT {
						panic("err: pointers cannot be floats")
					}

					return l.pos, PTR, fmt.Sprintf("*%s", digit)
				}

				return l.pos, MUL_OP, string(r)
			case '&':
				// we are dereferencing a pointer
				nextR, err := l.PeekNextChar()
				if err != nil {
					panic(err)
				}

				if unicode.IsDigit(nextR) {
					// means it's not a variable pointer
					l.Backup()
					token, digit := l.LexNumber()

					if token != INT {
						panic("err: cannot dereference by a float")
					}

					return l.pos, DEREF_OP, fmt.Sprintf("&%s", digit)
				} else {
					// means it's a variable pointer cos it's a string
					// we need to get the variable name
					// and check if it's a pointer and then dereference by its value
					l.Backup()
					_, lit := l.LexIdentifier()

					return l.pos, DEREF_OP, fmt.Sprintf("&%s", lit)
				}

			case '/':
				return l.pos, DIV_OP, string(r)
			case '^':
				return l.pos, POW_OP, string(r)
			case '%':
				return l.pos, MOD_OP, string(r)
			case '"':
				startPos := l.pos
				token, lit := l.LexString()
				return startPos, token, lit
			case '`':
				startPos := l.pos
				token, lit := l.LexMultiLineString()
				return startPos, token, lit
			case '=':
				return l.pos, ASSIGN_OP, string(r)
			}
		default:
			return l.pos, ILLEGAL, string(r)
		}
	}
}

// backup moves the reader back one rune
func (l *Lexer) Backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) PeekNextChar() (rune, error) {
	r, _, err := l.reader.ReadRune()

	if err != nil {
		if err == io.EOF {
			return 0, nil
		}

		return 0, err
	}

	return r, nil
}

func (l *Lexer) LexIdentifier() (Token, string) {
	var lit string

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if IsOperator(lit) {
					return identifierMap[lit], lit
				} else {
					// TODO: Research which scenarios can cause this
					return VARIABLE, lit
				}
			}

			panic(err)
		}

		l.pos.column++
		if string(r) == "\n" {
			// when we reach a newLine we return the current literal since we can't continue reading the identifier
			if IsOperator(lit) {
				return identifierMap[lit], lit
			}

			// anything other than an operator is assumed to be a variable
			return VARIABLE, lit
		} else if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.Backup()

			// fmt.Printf("Identifier: %v \n", lit)

			// check if lit does not exist in identifierMap
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

					// fmt.Printf("Next Rune: %v \n Is Space: %v \n", string(nextR), unicode.IsSpace(nextR))

					if !unicode.IsSpace(nextR) {
						break
					}
					counts++
				}

				// fmt.Printf("This is how many spaces were skipped: %v", counts)

				// Backup the reader
				for i := 0; i < counts; i++ {
					l.Backup()
				}

				return VARIABLE, lit
			}

			return identifierMap[lit], lit
		}
	}
}

func (l *Lexer) LexNumber() (Token, string) {
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

		if (unicode.IsSymbol(r) || unicode.IsPunct(r)) && r != '=' && len(lit) > 0 {
			panic(fmt.Errorf("error: invalid token %v at line %v col %v", string(r), l.pos.line, l.pos.column))
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
			l.Backup()
			break
		}
	}

	return tokenType, lit
}

func (l *Lexer) LexMultiLineString() (Token, string) {
	var lit string

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return STRING, lit
			}

			panic(err)
		}

		if r == '`' {
			break
		}

		lit += string(r)
	}

	return STRING, lit
}

// LexString scans the input for a string, and returns the string as a string
func (l *Lexer) LexString() (Token, string) {
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

		if r == '\n' {
			panic(fmt.Errorf("unexpected newline in string at line %d column %d", l.pos.line, l.pos.column))
		}

		lit += string(r)
	}

	return STRING, lit
}

// Resets the position of the lexer to the first col of the next line
func (l *Lexer) JumpToNextLine() {
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
	unaryOps := "print|drop|dup|dump|inc|dec|del"
	binaryOps := "+|-|*|/|mod|pow|swap|over|="
	ternaryOps := "rot"

	// assert that all ops are included
	unOps := strings.Split(unaryOps, "|")
	for _, op := range unOps {
		_, ok := identifierMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator: %s, did you perhaps forget to add it to the identifierMap?", op)
		}
	}

	binOps := strings.Split(binaryOps, "|")
	for _, op := range binOps {
		_, ok := identifierMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator: %s, did you perhaps forget to add it to the identifierMap?", op)
		}
	}

	terOps := strings.Split(ternaryOps, "|")
	for _, op := range terOps {
		_, ok := identifierMap[op]
		if !ok {
			return nil, fmt.Errorf("unknown operator: %s, did you perhaps forget to add it to the identifierMap?", op)
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

// START SEMANTIC ANALYSER
type SymbolTable struct {
	Variables map[string]Variable
}
type SemanticAnalyser struct {
	AST         *Node
	SymbolTable *SymbolTable
}

func NewSemanticAnalyser() *SemanticAnalyser {
	return &SemanticAnalyser{
		SymbolTable: &SymbolTable{
			Variables: make(map[string]Variable),
		},
	}
}

func (s *SemanticAnalyser) VariablesRepr() {
	// pretty print the variables
	fmt.Println("Variables:\t")
	for k, v := range s.SymbolTable.Variables {
		fmt.Printf("\t\tVariable: %s, Value: %v, Const: %v, Type: %v, Line: %v, Col: %v\n", k, v.Value, v.Const, tokenMap[v.Type], v.Value.Position.line, v.Value.Position.column)
	}
}

func (s *SemanticAnalyser) InferType(lit string) Token {
	pointerRegex := regexp.MustCompile(`\*\d+`)

	if _, err := strconv.Atoi(lit); err == nil {
		return INT
	}

	if _, err := strconv.ParseFloat(lit, 64); err == nil {
		return FLOAT
	}

	if lit == "true" || lit == "false" {
		return BOOL
	}

	if pointerRegex.MatchString(lit) {
		return PTR
	}

	return STRING
}

// END SEMANTIC ANALYSER

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

	parser := NewParser()
	lexer := NewLexer(file)
	semanticAnalyser := NewSemanticAnalyser()

	gorth := NewGorth(*strictMode, *debugMode, parser, lexer, semanticAnalyser)
	program := make([]StackElement, 0)

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

		semanticAnalyser.AST = root

		fmt.Println("Program AST: ")
		gorth.Parser.PrintAST(root, "")
	}

	fmt.Println("Program output:")
	gorth.ExecuteStack(program)
}
