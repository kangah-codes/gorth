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
	MAX_STACK_SIZE = 999_999
)

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

	// Stack manipulation operations
	SWAP_OP
	DUP_OP
	DROP_OP
	DUMP_OP
	ROT_OP

	// Print operation
	PRINT_OP

	// Logical operations
	AND_OP
	OR_OP
	NOT_OP
	EQUAL_OP
	NOT_EQUAL_OP
	EQUAL_TYP_OP // equal types
	GT_THAN_OP
	LS_THAN_OP
	GT_THAN_EQ_OP
	LS_THAN_EQ_OP
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
	"swap":  SWAP_OP,
	"dup":   DUP_OP,
	"drop":  DROP_OP,
	"dump":  DUMP_OP,
	"print": PRINT_OP,
	"rot":   ROT_OP,
	"&&":    AND_OP,
	"||":    OR_OP,
	"!":     NOT_OP,
	"==":    EQUAL_OP,
	"!=":    NOT_EQUAL_OP,
	"===":   EQUAL_TYP_OP,
	">":     GT_THAN_OP,
	"<":     LS_THAN_OP,
	">=":    GT_THAN_EQ_OP,
	"<=":    LS_THAN_EQ_OP,
}

type Type int

const (
	Int Type = iota
	String
	Bool
	Float
	Operator
	Identifier
	SpecialSymbol
)

type StackElement struct {
	Type  Type
	Value interface{} // Use interface{} to support both int and string values
}

type Variable struct {
	Type  Type
	Value interface{}
	Name  string
}

type Gorth struct {
	ExecStack    []StackElement
	VarStack     map[string]Variable
	DebugMode    bool
	StrictMode   bool
	MaxStackSize int
}

func NewGorth(debugMode, strictMode bool) *Gorth {
	return &Gorth{
		ExecStack:    []StackElement{},
		DebugMode:    debugMode,
		StrictMode:   strictMode,
		MaxStackSize: MAX_STACK_SIZE,
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
		// if line starts with a comment, ignore it
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Tokenizer
func Tokenize(s string) ([]StackElement, error) {
	var tokens []StackElement

	// define regex patterns
	integerRegex := regexp.MustCompile(`^-?\d+$`)
	floatRegex := regexp.MustCompile(`^-?\d+\.\d+$`)
	stringRegex := regexp.MustCompile(`^".*"$`)
	boolRegex := regexp.MustCompile(`^(true|false)$`)
	operatorRegex := regexp.MustCompile(`^(\+|-|\*|/|%|\^|\+\+|--|neg|swap|dup|drop|dump|print|&&|\|\||!|==|!=|===|>|<|>=|<=)$`)

	// split the string into tokens
	r := regexp.MustCompile(`"[^"]*"|\S+`)
	parts := r.FindAllString(s, -1)

	// parse each token
	for _, part := range parts {
		if integerRegex.MatchString(part) {
			val, _ := strconv.Atoi(part)
			tokens = append(tokens, StackElement{Type: Int, Value: val})
		} else if floatRegex.MatchString(part) {
			val, _ := strconv.ParseFloat(part, 64)
			tokens = append(tokens, StackElement{Type: Float, Value: val})
		} else if stringRegex.MatchString(part) {
			value := strings.Trim(part, `"`)
			tokens = append(tokens, StackElement{Type: String, Value: value})
		} else if boolRegex.MatchString(part) {
			val := part == "true"
			tokens = append(tokens, StackElement{Type: Bool, Value: val})
		} else if operatorRegex.MatchString(part) {
			tokens = append(tokens, StackElement{Type: Operator, Value: operatorMap[part]})
		} else {
			return nil, errors.New("ERROR: invalid token")
		}
	}

	return tokens, nil
}

func (g *Gorth) Push(val StackElement) error {
	if len(g.ExecStack) >= g.MaxStackSize {
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

func (g *Gorth) Rot() error {
	if len(g.ExecStack) < 3 {
		return errors.New("ERROR: at least 3 elements need to be on stack to perform ROT_OP")
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

	err = g.Push(val2)
	if err != nil {
		return err
	}

	err = g.Push(val1)
	if err != nil {
		return err
	}

	err = g.Push(val3)
	if err != nil {
		return err
	}

	return nil
}

func (g *Gorth) Peek() (StackElement, error) {
	if len(g.ExecStack) < 1 {
		return StackElement{}, errors.New("ERROR: cannot PEEK_OP at an empty stack")
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
		return errors.New("ERROR: cannot perform MUL_OP on different types")
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
		return errors.New("ERROR: cannot perform DIV_OP on different types")
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

	switch {
	case val1.Type == Bool && val1.Value.(bool):
		g.Push(StackElement{Type: Bool, Value: false})
	case val1.Type == Bool && !val1.Value.(bool):
		g.Push(StackElement{Type: Bool, Value: true})
	case val1.Type == Int && val1.Value.(int) != 0:
		g.Push(StackElement{Type: Int, Value: val1.Value.(int) * -1})
	case val1.Type == Int && val1.Value.(int) == 0:
		g.Push(StackElement{Type: Int, Value: 0})
	case val1.Type == Float && val1.Value.(float64) != 0:
		g.Push(StackElement{Type: Float, Value: val1.Value.(float64) * -1})
	case val1.Type == Float && val1.Value.(float64) == 0:
		g.Push(StackElement{Type: Float, Value: 0.0})
	default:
		return errors.New("ERROR: cannot perform NOT_OP on non boolean types")
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

func (g *Gorth) GreaterThan() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	switch {
	case val1.Type == Int && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(int) > val1.Value.(int)}) // Comparing val2 to val1
	case val1.Type == Float && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) > val1.Value.(float64)}) // Comparing val2 to val1
	case val1.Type == Int && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) > float64(val1.Value.(int))}) // Comparing val2 to val1
	case val1.Type == Float && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: float64(val2.Value.(int)) > val1.Value.(float64)}) // Comparing val2 to val1
	default:
		return errors.New("ERROR: cannot perform GT_THAN_OP on different types")
	}
	return nil
}

func (g *Gorth) LessThan() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	switch {
	case val1.Type == Int && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(int) < val1.Value.(int)})
	case val1.Type == Float && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) < val1.Value.(float64)})
	case val1.Type == Int && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) < float64(val1.Value.(int))})
	case val1.Type == Float && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: float64(val2.Value.(int)) < val1.Value.(float64)})
	default:
		return errors.New("ERROR: cannot perform LS_THAN_OP on different types")
	}
	return nil
}

func (g *Gorth) GreaterThanEqual() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	switch {
	case val1.Type == Int && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(int) >= val1.Value.(int)})
	case val1.Type == Float && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) >= val1.Value.(float64)})
	case val1.Type == Int && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) >= float64(val1.Value.(int))})
	case val1.Type == Float && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: float64(val2.Value.(int)) >= val1.Value.(float64)})
	default:
		return errors.New("ERROR: cannot perform GT_THAN_EQ_OP on different types")
	}
	return nil
}

func (g *Gorth) LessThanEqual() error {
	val1, err := g.Pop()
	if err != nil {
		return err
	}

	val2, err := g.Pop()
	if err != nil {
		return err
	}

	switch {
	case val1.Type == Int && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(int) <= val1.Value.(int)})
	case val1.Type == Float && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) <= val1.Value.(float64)})
	case val1.Type == Int && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val2.Value.(float64) <= float64(val1.Value.(int))})
	case val1.Type == Float && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: float64(val2.Value.(int)) <= val1.Value.(float64)})
	default:
		return errors.New("ERROR: cannot perform LS_THAN_EQ_OP on different types")
	}
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
			case GT_THAN_OP:
				err := g.GreaterThan()
				if err != nil {
					return err
				}
			case LS_THAN_OP:
				err := g.LessThan()
				if err != nil {
					return err
				}
			case GT_THAN_EQ_OP:
				err := g.GreaterThanEqual()
				if err != nil {
					return err
				}
			case LS_THAN_EQ_OP:
				err := g.LessThanEqual()
				if err != nil {
					return err
				}
			case ROT_OP:
				err := g.Rot()
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
	program, err := Tokenize(strings.Join(lines, " "))

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
