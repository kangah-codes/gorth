package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Operation int

const (
	// Arithmetic operations
	ADD_OP Operation = iota
	SUB_OP
	MUL_OP
	DIV_OP
	MOD_OP
	EXP_OP
	INC_OP
	DEC_OP
	NEG_OP

	// Stack manipulation operations
	SWAP_OP
	DUP_OP
	DROP_OP
	DUMP_OP

	// Print operation
	PRINT_OP

	// Logical operations
	AND_OP
	OR_OP
	NOT_OP
	EQUAL_OP
	NOT_EQUAL_OP
	EQUAL_TYP_OP // equal types
)

var operatorMap = map[string]Operation{
	"+":     ADD_OP,
	"-":     SUB_OP,
	"*":     MUL_OP,
	"/":     DIV_OP,
	"%":     MOD_OP,
	"^":     EXP_OP,
	"++":    INC_OP,
	"--":    DEC_OP,
	"neg":   NEG_OP,
	"swap":  SWAP_OP,
	"dup":   DUP_OP,
	"drop":  DROP_OP,
	"dump":  DUMP_OP,
	"print": PRINT_OP,
	"&&":    AND_OP,
	"||":    OR_OP,
	"!":     NOT_OP,
	"==":    EQUAL_OP,
	"!=":    NOT_EQUAL_OP,
	"===":   EQUAL_TYP_OP,
}

type Type int

const (
	Int Type = iota
	String
	Bool
	Float
	Operator
)

type StackElement struct {
	Type  Type
	Value interface{} // Use interface{} to support both int and string values
}

type Gorth struct {
	ExecStack  []StackElement
	DebugMode  bool
	StrictMode bool
}

func NewGorth(debugMode, strictMode bool) *Gorth {
	return &Gorth{
		ExecStack:  []StackElement{},
		DebugMode:  debugMode,
		StrictMode: strictMode,
	}
}

// ReadGorthFile reads a .gorth file and returns the contents as a slice of strings.
func ReadGorthFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// ParseProgram parses the lines of a .gorth file and converts them into a program of StackElements.
func ParseProgram(lines []string) ([]StackElement, error) {
	var program []StackElement
	r := regexp.MustCompile(`"[^"]*"|\S+`)

	for _, line := range lines {
		parts := r.FindAllString(line, -1)
		for _, part := range parts {
			element, err := parseElement(strings.Trim(part, `"`))
			if err != nil {
				return nil, err
			}
			program = append(program, element)
		}
	}

	return program, nil
}

// parseElement parses a string and converts it into a StackElement.
func parseElement(s string) (StackElement, error) {
	var element StackElement
	if op, ok := operatorMap[s]; ok {
		element.Type = Operator
		element.Value = op
	} else if val, err := strconv.Atoi(s); err == nil {
		element.Type = Int
		element.Value = val
	} else if s == "true" || s == "false" {
		element.Type = Bool
		element.Value = s == "true"
	} else if val, err := strconv.ParseFloat(s, 64); err == nil {
		element.Type = Float
		element.Value = val
	} else if s[0] == '"' && s[len(s)-1] == '"' {
		element.Type = String
		element.Value = s
	} else {
		return StackElement{}, errors.New("ERROR: invalid token")
	}
	return element, nil
}

func (g *Gorth) Push(val StackElement) error {
	if len(g.ExecStack) >= 1000 {
		return errors.New("ERROR: stack overflow")
	}
	g.ExecStack = append(g.ExecStack, val)
	return nil
}

func (g *Gorth) Pop() (StackElement, error) {
	if len(g.ExecStack) < 1 {
		return StackElement{}, errors.New("ERROR: cannot pop from an empty stack")
	}
	val := g.ExecStack[len(g.ExecStack)-1]
	g.ExecStack = g.ExecStack[:len(g.ExecStack)-1]
	return val, nil
}

func (g *Gorth) Drop() error {
	if len(g.ExecStack) < 1 {
		return errors.New("ERROR: cannot drop from an empty stack")
	}
	g.Pop()
	return nil
}

func (g *Gorth) Dump() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	switch val.Type {
	case Int, String, Bool:
		fmt.Println(val.Value)
	default:
		return errors.New("ERROR: top element is not a printable type")
	}
	return nil
}

func (g *Gorth) Peek() (StackElement, error) {
	if len(g.ExecStack) < 1 {
		return StackElement{}, errors.New("ERROR: cannot peek at an empty stack")
	}
	return g.ExecStack[len(g.ExecStack)-1], nil
}

func (g *Gorth) Print() error {
	val, err := g.Peek()
	if err != nil {
		return err
	}

	switch val.Type {
	case Int, String, Bool:
		fmt.Println(val.Value)
	default:
		return errors.New("ERROR: top element is not a printable type")
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
	switch {
	// integer addition
	case val1.Type == Int && val2.Type == Int:
		sum := val1.Value.(int) + val2.Value.(int)
		g.Push(StackElement{Type: Int, Value: sum})
	// string concatenation
	case val1.Type == String && val2.Type == String:
		concat := val2.Value.(string) + val1.Value.(string)
		g.Push(StackElement{Type: String, Value: concat})
	// float addition
	case val1.Type == Float && val2.Type == Float:
		sum := val1.Value.(float64) + val2.Value.(float64)
		g.Push(StackElement{Type: Float, Value: sum})
	// mixed type addition
	case (val1.Type == Int && val2.Type == Float) || (val1.Type == Float && val2.Type == Int):
		sum := val2.Value.(float64) + float64(val1.Value.(int))
		g.Push(StackElement{Type: Float, Value: sum})
	default:
		return errors.New("ERROR: cannot perform ADD_OP on different types")
	}
	return nil
}

func (g *Gorth) Sub() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}
	val2, err := g.Pop()
	if err != nil {
		return err
	}
	switch {
	// integer subtraction
	case val1.Type == Int && val2.Type == Int:
		sub := val2.Value.(int) - val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: sub})
	// float subtraction
	case val1.Type == Float && val2.Type == Float:
		sub := val2.Value.(float64) - val1.Value.(float64)
		g.Push(StackElement{Type: Float, Value: sub})
	default:
		return errors.New("ERROR: cannot perform SUB_OP on different types")
	}
	return nil
}

func (g *Gorth) Mul() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}
	val2, err := g.Pop()
	if err != nil {
		return err
	}
	switch {
	// integer multiplication
	case val1.Type == Int && val2.Type == Int:
		mul := val1.Value.(int) * val2.Value.(int)
		g.Push(StackElement{Type: Int, Value: mul})
	// string multiplication
	case val1.Type == String && val2.Type == Int:
		str := val1.Value.(string)
		num := val2.Value.(int)
		var concat string
		for i := 0; i < num; i++ {
			concat += str
		}
		g.Push(StackElement{Type: String, Value: concat})
	// string multiplication
	case val1.Type == Int && val2.Type == String:
		str := val2.Value.(string)
		num := val1.Value.(int)
		var concat string
		for i := 0; i < num; i++ {
			concat += str
		}
		g.Push(StackElement{Type: String, Value: concat})
	// float multiplication
	case val1.Type == Float && val2.Type == Float:
		mul := val1.Value.(float64) * val2.Value.(float64)
		g.Push(StackElement{Type: Float, Value: mul})
	// mixed type multiplication
	case (val1.Type == Int && val2.Type == Float) || (val1.Type == Float && val2.Type == Int):
		mul := val2.Value.(float64) * float64(val1.Value.(int))
		g.Push(StackElement{Type: Float, Value: mul})
	default:
		return errors.New("ERROR: cannot perform MUL_OP on different types")
	}
	return nil
}

func (g *Gorth) Div() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}
	val2, err := g.Pop()
	if err != nil {
		return err
	}
	switch {
	// integer division
	case val1.Type == Int && val2.Type == Int:
		div := val2.Value.(int) / val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: div})
	// float division
	case val1.Type == Float && val2.Type == Float:
		div := val2.Value.(float64) / val1.Value.(float64)
		g.Push(StackElement{Type: Float, Value: div})
	// mixed type division
	case (val1.Type == Int && val2.Type == Float) || (val1.Type == Float && val2.Type == Int):
		div := val2.Value.(float64) / float64(val1.Value.(int))
		g.Push(StackElement{Type: Float, Value: div})
	default:
		return errors.New("ERROR: cannot perform DIV_OP on different types")
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
	switch {
	// integer modulo
	case val1.Type == Int && val2.Type == Int:
		mod := val2.Value.(int) % val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: mod})
	default:
		return errors.New("ERROR: cannot perform MOD_OP on different types")
	}
	return nil
}

func (g *Gorth) Exp() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}
	val2, err := g.Pop()
	if err != nil {
		return err
	}
	switch {
	// integer exponentiation
	case val1.Type == Int && val2.Type == Int:
		exp := int(math.Pow(float64(val2.Value.(int)), float64(val1.Value.(int))))
		g.Push(StackElement{Type: Int, Value: exp})
	// float exponentiation
	case val1.Type == Float && val2.Type == Float:
		exp := math.Pow(val2.Value.(float64), val1.Value.(float64))
		g.Push(StackElement{Type: Float, Value: exp})
	// mixed type exponentiation
	case (val1.Type == Int && val2.Type == Float) || (val1.Type == Float && val2.Type == Int):
		exp := math.Pow(float64(val2.Value.(int)), val1.Value.(float64))
		g.Push(StackElement{Type: Float, Value: exp})
	default:
		return errors.New("ERROR: cannot perform EXP_OP on different types")
	}
	return nil
}

func (g *Gorth) Inc() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}
	switch {
	case val.Type == Int:
		g.Push(StackElement{Type: Int, Value: val.Value.(int) + 1})
	default:
		return errors.New("ERROR: cannot perform INC_OP on different types")
	}
	return nil
}

func (g *Gorth) Dec() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}
	switch {
	case val.Type == Int:
		g.Push(StackElement{Type: Int, Value: val.Value.(int) - 1})
	default:
		return errors.New("ERROR: cannot perform DEC_OP on different types")
	}
	return nil
}

func (g *Gorth) Neg() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}
	switch {
	case val.Type == Int:
		g.Push(StackElement{Type: Int, Value: -val.Value.(int)})
	case val.Type == Float:
		g.Push(StackElement{Type: Float, Value: -val.Value.(float64)})
	default:
		return errors.New("ERROR: cannot perform NEG_OP on different types")
	}
	return nil
}

func (g *Gorth) Swap() error {
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
	val, err := g.Peek()
	if err != nil {
		return err
	}
	g.Push(val)
	return nil
}

func (g *Gorth) And() error {
	// checks if the top two elements are both true
	// only works if both elements are boolean
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	// check types
	if val1.Type != Bool || val2.Type != Bool {
		return errors.New("ERROR: cannot perform AND_OP on non boolean types")
	}

	// check values
	if val1.Value.(bool) && val2.Value.(bool) {
		g.Push(StackElement{Type: Bool, Value: true})
	} else {
		g.Push(StackElement{Type: Bool, Value: false})
	}

	return nil
}

func (g *Gorth) Or() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	// check types
	if val1.Type != Bool || val2.Type != Bool {
		return errors.New("ERROR: cannot perform OR_OP on non boolean types")
	}

	if val1.Value.(bool) || val2.Value.(bool) {
		g.Push(StackElement{Type: Bool, Value: true})
	} else {
		g.Push(StackElement{Type: Bool, Value: false})
	}

	return nil
}

func (g *Gorth) Not() error {
	// flips the top of the stack
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	// check types
	if val1.Type != Bool {
		return errors.New("ERROR: cannot perform NOT_OP on non boolean types")
	}

	if val1.Value.(bool) {
		g.Push(StackElement{Type: Bool, Value: false})
	} else {
		g.Push(StackElement{Type: Bool, Value: true})
	}
	return nil
}

func (g *Gorth) Equal() error {
	// checks if the top elements are equal
	// equality checking is independent of type
	// maybe bad language design lol
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	g.Push(StackElement{Type: Bool, Value: val1.Value == val2.Value})
	return nil
}

func (g *Gorth) NotEqual() {
	// checks if the top elements are not equal
	// equality checking is independent of type
	// maybe bad language design lol
	g.Equal()
	g.Not()
}

func (g *Gorth) EqualType() error {
	// checks if the top elements are equal
	// equality checking is dependent on type
	// maybe bad language design lol
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	g.Push(StackElement{Type: Bool, Value: val1.Type == val2.Type})
	return nil
}

func (g *Gorth) PrintStack() {
	fmt.Printf("Program stack: %v\n", g.ExecStack)
}

func (g *Gorth) ExecuteProgram(program []StackElement) error {
	for _, op := range program {
		if op.Type == Operator {
			switch op.Value {
			case ADD_OP:
				err := g.Add()
				if err != nil {
					return err
				}
			case SUB_OP:
				err := g.Sub()
				if err != nil {
					return err
				}
			case MUL_OP:
				err := g.Mul()
				if err != nil {
					return err
				}
			case DIV_OP:
				err := g.Div()
				if err != nil {
					return err
				}
			case MOD_OP:
				err := g.Mod()
				if err != nil {
					return err
				}
			case EXP_OP:
				err := g.Exp()
				if err != nil {
					return err
				}
			case INC_OP:
				err := g.Inc()
				if err != nil {
					return err
				}
			case DEC_OP:
				err := g.Dec()
				if err != nil {
					return err
				}
			case NEG_OP:
				err := g.Neg()
				if err != nil {
					return err
				}
			case SWAP_OP:
				err := g.Swap()
				if err != nil {
					return err
				}
			case DUP_OP:
				err := g.Dup()
				if err != nil {
					return err
				}
			case DROP_OP:
				err := g.Drop()
				if err != nil {
					return err
				}
			case DUMP_OP:
				err := g.Dump()
				if err != nil {
					return err
				}
			case PRINT_OP:
				err := g.Print()
				if err != nil {
					return err
				}
			case AND_OP:
				err := g.And()
				if err != nil {
					return err
				}
			case OR_OP:
				err := g.Or()
				if err != nil {
					return err
				}
			case NOT_OP:
				err := g.Not()
				if err != nil {
					return err
				}
			case EQUAL_OP:
				err := g.Equal()
				if err != nil {
					return err
				}
			case NOT_EQUAL_OP:
				g.NotEqual()
			case EQUAL_TYP_OP:
				err := g.EqualType()
				if err != nil {
					return err
				}
			}
		} else {
			err := g.Push(op)
			if err != nil {
				return err
			}
		}
	}

	if g.StrictMode {
		if len(g.ExecStack) > 0 {
			return fmt.Errorf("ERROR: unconsumed elements remain on the stack\n\t%v", g.ExecStack)
		}
	}

	if g.DebugMode {
		fmt.Printf("Program stack at end of execution\n\t%v\n", g.ExecStack)
	}

	return nil
}

func PrintUsage() {
	fmt.Println("Usage: gorth <filename> [options]")
	fmt.Println("  filename: the name of the .gorth file to execute")
	fmt.Println("  options:")
	fmt.Println("    -d: optional enable debug mode")
	fmt.Println("    -s: optional enable strict mode")
}

func main() {
	// get system arguments
	args := os.Args[1:]

	// check if there are no arguments
	if len(args) == 0 {
		PrintUsage()
		return
	}

	// check if there are too many arguments
	if len(args) > 3 {
		panic("Too many arguments provided")
	}

	// check if the first argument is a .gorth file
	if !strings.HasSuffix(args[0], ".gorth") {
		panic(fmt.Sprintf("File %s is not a .gorth file", args[0]))
	}

	// check if the file exists
	_, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("File %s does not exist", args[0]))
	}

	// read the file
	lines, err := ReadGorthFile(args[0])
	if err != nil {
		panic(err)
	}

	// get the other arguments even if there are not in the correct order
	debugMode := false
	strictMode := false

	for _, arg := range args[1:] {
		switch arg {
		case "-d":
			debugMode = true
		case "-s":
			strictMode = true
		default:
			panic(fmt.Sprintf("Invalid option: %s", arg))
		}
	}

	// parse the program
	program, err := ParseProgram(lines)

	if err != nil {
		panic(err)
	}

	// create a new gorth instance
	g := NewGorth(debugMode, strictMode)

	start := time.Now()

	if g.DebugMode {
		fmt.Printf("Program stack at start of execution\n\t%v\n", g.ExecStack)
	}

	// execute the program
	err = g.ExecuteProgram(program)

	end := time.Now()

	if err != nil {
		fmt.Println("Program simulation failed")
		fmt.Println(err)
	} else {
		fmt.Printf("Program simulation completed in %v seconds\n", end.Sub(start).Seconds())
	}
}
