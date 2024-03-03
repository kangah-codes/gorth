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

	// assignment operation
	VAR_ASSIGN_OP
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
	"=":     VAR_ASSIGN_OP,
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
	KeyWord
)

var typeMap = map[Type]string{
	Int:           "int",
	String:        "string",
	Bool:          "bool",
	Operator:      "operator",
	Identifier:    "identifier",
	SpecialSymbol: "special symbol",
	KeyWord:       "keyword",
}

type StackElement struct {
	Type  Type
	Value interface{} // Use interface{} to support both int and string values
}

func (s *StackElement) Repr() string {
	return fmt.Sprintf("Type: %v\nValue: %v", typeMap[s.Type], s.Value)
}

type Variable struct {
	Type  Type
	Value interface{}
	Name  string
	Const bool
}

type Gorth struct {
	ExecStack    []StackElement
	VariableMap  map[string]Variable
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

const (
	StateNormal = iota
	StateVarDeclaration
)

type Tokeniser struct {
	HandleToken func(s string) ([]StackElement, map[string]Variable, error)
}

type TokeniserStateMachine struct {
	CurrentState int
	States       map[int]Tokeniser
}

func (t *TokeniserStateMachine) SetState(state int) {
	t.CurrentState = state
}

// Tokenizer
func Tokenize(s string) ([]StackElement, map[string]Variable, error) {
	var tokens []StackElement
	var lastAddedVariable Variable
	variables := make(map[string]Variable)

	// Define regex patterns
	integerRegex := regexp.MustCompile(`^-?\d+$`)
	floatRegex := regexp.MustCompile(`^-?\d+\.\d+$`)
	stringRegex := regexp.MustCompile(`^".*"$`)
	boolRegex := regexp.MustCompile(`^(true|false)$`)
	operatorRegex := regexp.MustCompile(`^(\+|-|\*|/|%|\^|\+\+|--|neg|swap|dup|drop|dump|print|rot|&&|\|\||!|==|!=|===|>|<|>=|<=|=)$`)
	varNameRegex := regexp.MustCompile(`^\/[a-zA-Z_][a-zA-Z0-9_]*$`)
	// TODO: rename this
	keyWordRegex := regexp.MustCompile(`^(def|const|=)$`)
	// using variables : _varName
	varUsageRegex := regexp.MustCompile(`^_[a-zA-Z_][a-zA-Z0-9_]*$`)

	// Split the string into tokens
	r := regexp.MustCompile(`"[^"]*"|\S+`)
	parts := r.FindAllString(s, -1)

	// Current state
	state := StateNormal
	stateMachine := TokeniserStateMachine{}
	stateMachine.SetState(state)

	// Define state machine
	stateMachine.States = map[int]Tokeniser{
		StateNormal: {
			HandleToken: func(s string) ([]StackElement, map[string]Variable, error) {
				switch {
				case integerRegex.MatchString(s):
					val, _ := strconv.Atoi(s)
					tokens = append(tokens, StackElement{Type: Int, Value: val})
				case floatRegex.MatchString(s):
					val, _ := strconv.ParseFloat(s, 64)
					tokens = append(tokens, StackElement{Type: Float, Value: val})
				case stringRegex.MatchString(s):
					value := strings.Trim(s, `"`)
					tokens = append(tokens, StackElement{Type: String, Value: value})
				case boolRegex.MatchString(s):
					val := s == "true"
					tokens = append(tokens, StackElement{Type: Bool, Value: val})
				case operatorRegex.MatchString(s):
					tokens = append(tokens, StackElement{Type: Operator, Value: operatorMap[s]})
				case keyWordRegex.MatchString(s):
					// Reset back to normal state since we've encountered the def keyword which means we're done declaring variables
					if strings.TrimSpace(s) == "const" {
						// variable is a constant
						lastAddedVariable.Const = true
						variables[lastAddedVariable.Name] = lastAddedVariable
					}
				case operatorRegex.MatchString(s):
					// means an operator comes after a variable, most likely we are reassiging a variable
					if len(variables) < 1 {
						return nil, nil, errors.New("ERROR: no variable to assign to")
					}

					if strings.TrimSpace(s) == "=" && lastAddedVariable.Const {
						return nil, nil, errors.New("ERROR: cannot reassign a constant")
					}
				// had to add new syntax to check if a variable was being used
				// because the same syntax caused a bug where the last known variable was used even if
				// the variable name was not the same
				/* Example:
				 * /myName "Joshua" def
				 * /notMyName print
				 * This code should throw an error since we don't have notMyName,
				 * But since our state machine recognises that syntax for variable declaration
				 * We just jump into adding the print operation to the exec stack,
				 * Which then pops the topmost value which would be the last known variable, ie. myName
				 */
				case varUsageRegex.MatchString(s):
					// check if the variable exists
					// if it does, add it's value to the tokens
					variable, exists := variables[s[1:]]

					if !exists {
						return nil, nil, fmt.Errorf("variable %s has not been declared", s[1:])
					}

					// tokens = append(tokens, StackElement{Type: variable.Type, Value: variable.Value})
					tokens = append(tokens, StackElement{Type: Identifier, Value: variable.Name})
				default:
					return nil, nil, fmt.Errorf("invalid token: %s", s)
				}

				return tokens, nil, nil
			},
		},
		StateVarDeclaration: {
			HandleToken: func(part string) ([]StackElement, map[string]Variable, error) {
				// Assuming the value immediately follows the variable name
				// Add the variable and its value to the map
				// check if the variable map is not empty
				// if it is not empty, get the last token and add the value to the variable map
				// if it is empty, return an error
				if len(variables) > 0 {
					if operatorRegex.MatchString(part) {
						// idk why this would happen
						tokens = append(tokens, StackElement{Type: Operator, Value: operatorMap[part]})
					} else {
						switch {
						case integerRegex.MatchString(part):
							val, _ := strconv.Atoi(part)
							lastAddedVariable.Value = val
							lastAddedVariable.Type = Int
							variables[lastAddedVariable.Name] = lastAddedVariable
							tokens = append(tokens, StackElement{Type: Identifier, Value: lastAddedVariable.Name})
						case floatRegex.MatchString(part):
							val, _ := strconv.ParseFloat(part, 64)
							lastAddedVariable.Value = val
							lastAddedVariable.Type = Float
							variables[lastAddedVariable.Name] = lastAddedVariable
							tokens = append(tokens, StackElement{Type: Identifier, Value: lastAddedVariable.Name})
						case stringRegex.MatchString(part):
							value := strings.Trim(part, `"`)
							lastAddedVariable.Value = value
							lastAddedVariable.Type = String
							variables[lastAddedVariable.Name] = lastAddedVariable
							tokens = append(tokens, StackElement{Type: Identifier, Value: lastAddedVariable.Name})
						case boolRegex.MatchString(part):
							val := part == "true"
							lastAddedVariable.Value = val
							lastAddedVariable.Type = Bool
							variables[lastAddedVariable.Name] = lastAddedVariable
							tokens = append(tokens, StackElement{Type: Identifier, Value: lastAddedVariable.Name})
						default:
							return nil, nil, fmt.Errorf("invalid type: %s", part)
						}
					}
				}

				// Reset the state to normal
				stateMachine.SetState(StateNormal)

				return nil, variables, nil
			},
		},
	}

	// Parse each token
	for _, part := range parts {
		// Check if the variable name already exists
		if _, exists := variables[part]; exists {
			return nil, nil, fmt.Errorf("variable %s is already declared", part)
		}

		// set the machine state based on the current token
		if varNameRegex.MatchString(part) {
			// check if variable already exists in the map
			_, exists := variables[part[1:]]

			if exists {
				// just jump because we've already declared the variable
				// and we're probably just using it
				continue
			} else {
				stateMachine.SetState(StateVarDeclaration)
				varName := part[1:] // Remove the leading '/'
				variables[varName] = Variable{Name: varName, Type: Identifier}
				lastAddedVariable = variables[varName]
				continue
			}

		}

		_, _, err := stateMachine.States[stateMachine.CurrentState].HandleToken(part)

		if err != nil {
			return nil, nil, err
		}
	}

	return tokens, variables, nil
}

func (g *Gorth) GPrint(val interface{}) {
	if g.DebugMode {
		fmt.Println(val)
	}
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
	case Identifier:
		fmt.Println(g.VariableMap[val.Value.(string)].Value)
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
	case Int, String, Bool, Float:
		fmt.Println(val.Value)
	case Identifier:
		// we use value since we set the value of variables on the element stack to the name of the variable
		_, exists := g.VariableMap[val.Value.(string)]

		if !exists {
			return fmt.Errorf("ERROR: variable %v has not been declared", val.Value.(string))
		}

		switch g.VariableMap[val.Value.(string)].Type {
		case Int, String, Bool, Float:
			fmt.Println(g.VariableMap[val.Value.(string)].Value)
		}
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
		// for string concatenation, we reverse the order of the strings
		// since the first string to be popped is the second string and vice versa
		concat := val1.Value.(string) + val2.Value.(string)
		g.Push(StackElement{Type: String, Value: concat})
	// float addition
	case val1.Type == Float && val2.Type == Float:
		sum := val1.Value.(float64) + val2.Value.(float64)
		g.Push(StackElement{Type: Float, Value: sum})
	// mixed type addition
	case val1.Type == Int && val2.Type == Float:
		sum := val2.Value.(float64) + float64(val1.Value.(int))
		g.Push(StackElement{Type: Float, Value: sum})
	case val1.Type == Float && val2.Type == Int:
		sum := val1.Value.(float64) + float64(val2.Value.(int))
		g.Push(StackElement{Type: Float, Value: sum})
	// both variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Int:
			sum := g.VariableMap[val1.Value.(string)].Value.(int) + g.VariableMap[val2.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: sum})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Float:
			sum := g.VariableMap[val1.Value.(string)].Value.(float64) + g.VariableMap[val2.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: sum})
		case g.VariableMap[val1.Value.(string)].Type == String && g.VariableMap[val2.Value.(string)].Type == String:
			concat := g.VariableMap[val1.Value.(string)].Value.(string) + g.VariableMap[val2.Value.(string)].Value.(string)
			g.Push(StackElement{Type: String, Value: concat})
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Float:
			sum := float64(g.VariableMap[val1.Value.(string)].Value.(int)) + g.VariableMap[val2.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: sum})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Int:
			sum := g.VariableMap[val1.Value.(string)].Value.(float64) + float64(g.VariableMap[val2.Value.(string)].Value.(int))
			g.Push(StackElement{Type: Float, Value: sum})
		default:
			return errors.New("ERROR: cannot perform ADD_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Int:
			sum := g.VariableMap[val1.Value.(string)].Value.(int) + val2.Value.(int)
			g.Push(StackElement{Type: Int, Value: sum})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Float:
			sum := g.VariableMap[val1.Value.(string)].Value.(float64) + val2.Value.(float64)
			g.Push(StackElement{Type: Float, Value: sum})
		case g.VariableMap[val1.Value.(string)].Type == String && val2.Type == String:
			concat := g.VariableMap[val1.Value.(string)].Value.(string) + val2.Value.(string)
			g.Push(StackElement{Type: String, Value: concat})
		// one is an int and the other is a float
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Float:
			sum := float64(g.VariableMap[val1.Value.(string)].Value.(int)) + val2.Value.(float64)
			g.Push(StackElement{Type: Float, Value: sum})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Int:
			sum := g.VariableMap[val1.Value.(string)].Value.(float64) + float64(val2.Value.(int))
			g.Push(StackElement{Type: Float, Value: sum})
		default:
			return errors.New("ERROR: cannot perform ADD_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Int:
			sum := g.VariableMap[val2.Value.(string)].Value.(int) + val1.Value.(int)
			g.Push(StackElement{Type: Int, Value: sum})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Float:
			sum := g.VariableMap[val2.Value.(string)].Value.(float64) + val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: sum})
		case g.VariableMap[val2.Value.(string)].Type == String && val1.Type == String:
			concat := g.VariableMap[val2.Value.(string)].Value.(string) + val1.Value.(string)
			g.Push(StackElement{Type: String, Value: concat})
		// one is an int and the other is a float
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Float:
			sum := float64(g.VariableMap[val2.Value.(string)].Value.(int)) + val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: sum})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Int:
			sum := g.VariableMap[val2.Value.(string)].Value.(float64) + float64(val1.Value.(int))
			g.Push(StackElement{Type: Float, Value: sum})
		default:
			return errors.New("ERROR: cannot perform ADD_OP on different types")
		}
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
	// mixed number subtraction
	case val1.Type == Int && val2.Type == Float:
		sub := val2.Value.(float64) - float64(val1.Value.(int))
		g.Push(StackElement{Type: Float, Value: sub})
	case val1.Type == Float && val2.Type == Int:
		sub := float64(val2.Value.(int)) - val1.Value.(float64)
		g.Push(StackElement{Type: Float, Value: sub})
	// variable subtraction
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Int:
			sub := g.VariableMap[val2.Value.(string)].Value.(int) - g.VariableMap[val1.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: sub})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Float:
			sub := g.VariableMap[val2.Value.(string)].Value.(float64) - g.VariableMap[val1.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: sub})
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Float:
			sub := g.VariableMap[val2.Value.(string)].Value.(float64) - float64(g.VariableMap[val1.Value.(string)].Value.(int))
			g.Push(StackElement{Type: Float, Value: sub})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Int:
			sub := float64(g.VariableMap[val2.Value.(string)].Value.(int)) - g.VariableMap[val1.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: sub})
		default:
			return errors.New("ERROR: cannot perform SUB_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Int:
			sub := val2.Value.(int) - g.VariableMap[val1.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: sub})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Float:
			sub := val2.Value.(float64) - g.VariableMap[val1.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: sub})
		// one is an int and the other is a float
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Float:
			sub := val2.Value.(float64) - float64(g.VariableMap[val1.Value.(string)].Value.(int))
			g.Push(StackElement{Type: Float, Value: sub})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Int:
			sub := float64(val2.Value.(int)) - g.VariableMap[val1.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: sub})
		default:
			return errors.New("ERROR: cannot perform SUB_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Int:
			sub := g.VariableMap[val2.Value.(string)].Value.(int) - val1.Value.(int)
			g.Push(StackElement{Type: Int, Value: sub})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Float:
			sub := g.VariableMap[val2.Value.(string)].Value.(float64) - val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: sub})
		// one is an int and the other is a float
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Float:
			sub := float64(g.VariableMap[val2.Value.(string)].Value.(int)) - val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: sub})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Int:
			sub := g.VariableMap[val2.Value.(string)].Value.(float64) - float64(val1.Value.(int))
			g.Push(StackElement{Type: Float, Value: sub})
		default:
			return errors.New("ERROR: cannot perform SUB_OP on different types")
		}
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
	// one is float and the other is an int
	case val1.Type == Int && val2.Type == Float:
		mul := float64(val1.Value.(int)) * val2.Value.(float64)
		g.Push(StackElement{Type: Float, Value: mul})
	case val1.Type == Float && val2.Type == Int:
		mul := val1.Value.(float64) * float64(val2.Value.(int))
		g.Push(StackElement{Type: Float, Value: mul})
	// variable multiplication
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Int:
			mul := g.VariableMap[val1.Value.(string)].Value.(int) * g.VariableMap[val2.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: mul})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Float:
			mul := g.VariableMap[val1.Value.(string)].Value.(float64) * g.VariableMap[val2.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: mul})
		case g.VariableMap[val1.Value.(string)].Type == String && g.VariableMap[val2.Value.(string)].Type == Int:
			str := g.VariableMap[val1.Value.(string)].Value.(string)
			num := g.VariableMap[val2.Value.(string)].Value.(int)
			var concat string
			for i := 0; i < num; i++ {
				concat += str
			}
			g.Push(StackElement{Type: String, Value: concat})
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == String:
			str := g.VariableMap[val2.Value.(string)].Value.(string)
			num := g.VariableMap[val1.Value.(string)].Value.(int)
			var concat string
			for i := 0; i < num; i++ {
				concat += str
			}
			g.Push(StackElement{Type: String, Value: concat})
		// one is a float and one is an int
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Float:
			mul := float64(g.VariableMap[val1.Value.(string)].Value.(int)) * g.VariableMap[val2.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: mul})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Int:
			mul := g.VariableMap[val1.Value.(string)].Value.(float64) * float64(g.VariableMap[val2.Value.(string)].Value.(int))
			g.Push(StackElement{Type: Float, Value: mul})
		default:
			return errors.New("ERROR: cannot perform MUL_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Int:
			mul := g.VariableMap[val1.Value.(string)].Value.(int) * val2.Value.(int)
			g.Push(StackElement{Type: Int, Value: mul})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Float:
			mul := g.VariableMap[val1.Value.(string)].Value.(float64) * val2.Value.(float64)
			g.Push(StackElement{Type: Float, Value: mul})
		case g.VariableMap[val1.Value.(string)].Type == String && val2.Type == Int:
			str := g.VariableMap[val1.Value.(string)].Value.(string)
			num := val2.Value.(int)
			var concat string
			for i := 0; i < num; i++ {
				concat += str
			}
			g.Push(StackElement{Type: String, Value: concat})
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == String:
			str := val2.Value.(string)
			num := g.VariableMap[val1.Value.(string)].Value.(int)
			var concat string
			for i := 0; i < num; i++ {
				concat += str
			}
			g.Push(StackElement{Type: String, Value: concat})
		// one is a float and one is an int
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Float:
			mul := float64(g.VariableMap[val1.Value.(string)].Value.(int)) * val2.Value.(float64)
			g.Push(StackElement{Type: Float, Value: mul})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Int:
			mul := g.VariableMap[val1.Value.(string)].Value.(float64) * float64(val2.Value.(int))
			g.Push(StackElement{Type: Float, Value: mul})
		default:
			return errors.New("ERROR: cannot perform MUL_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Int:
			mul := g.VariableMap[val2.Value.(string)].Value.(int) * val1.Value.(int)
			g.Push(StackElement{Type: Int, Value: mul})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Float:
			mul := g.VariableMap[val2.Value.(string)].Value.(float64) * val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: mul})
		case g.VariableMap[val2.Value.(string)].Type == String && val1.Type == Int:
			str := g.VariableMap[val2.Value.(string)].Value.(string)
			num := val1.Value.(int)
			var concat string
			for i := 0; i < num; i++ {
				concat += str
			}
			g.Push(StackElement{Type: String, Value: concat})
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == String:
			str := val1.Value.(string)
			num := g.VariableMap[val2.Value.(string)].Value.(int)
			var concat string
			for i := 0; i < num; i++ {
				concat += str
			}
			g.Push(StackElement{Type: String, Value: concat})
		// one is a float and one is an int
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Float:
			mul := float64(g.VariableMap[val2.Value.(string)].Value.(int)) * val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: mul})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Int:
			mul := g.VariableMap[val2.Value.(string)].Value.(float64) * float64(val1.Value.(int))
			g.Push(StackElement{Type: Float, Value: mul})
		default:
			return errors.New("ERROR: cannot perform MUL_OP on different types")
		}
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
	// one is float and the other is an int
	case val1.Type == Int && val2.Type == Float:
		div := val2.Value.(float64) / float64(val1.Value.(int))
		g.Push(StackElement{Type: Float, Value: div})
	case val1.Type == Float && val2.Type == Int:
		div := float64(val2.Value.(int)) / val1.Value.(float64)
		g.Push(StackElement{Type: Float, Value: div})
	// variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Int:
			div := g.VariableMap[val2.Value.(string)].Value.(int) / g.VariableMap[val1.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: div})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Float:
			div := g.VariableMap[val2.Value.(string)].Value.(float64) / g.VariableMap[val1.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: div})
		// one is an int and the other is a float
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Float:
			div := g.VariableMap[val2.Value.(string)].Value.(float64) / float64(g.VariableMap[val1.Value.(string)].Value.(int))
			g.Push(StackElement{Type: Float, Value: div})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Int:
			div := float64(g.VariableMap[val2.Value.(string)].Value.(int)) / g.VariableMap[val1.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: div})
		default:
			return errors.New("ERROR: cannot perform DIV_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Int:
			div := val2.Value.(int) / g.VariableMap[val1.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: div})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Float:
			div := val2.Value.(float64) / g.VariableMap[val1.Value.(string)].Value.(float64)
			g.Push(StackElement{Type: Float, Value: div})
		// one is an int and the other is a float
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Float:
			div := val2.Value.(float64) / float64(g.VariableMap[val1.Value.(string)].Value.(int))
			g.Push(StackElement{Type: Float, Value: div})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Int:
			div := g.VariableMap[val1.Value.(string)].Value.(float64) / float64(val2.Value.(int))
			g.Push(StackElement{Type: Float, Value: div})
		default:
			return errors.New("ERROR: cannot perform DIV_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Int:
			div := g.VariableMap[val2.Value.(string)].Value.(int) / val1.Value.(int)
			g.Push(StackElement{Type: Int, Value: div})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Float:
			div := g.VariableMap[val2.Value.(string)].Value.(float64) / val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: div})
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Float:
			div := float64(g.VariableMap[val2.Value.(string)].Value.(int)) / val1.Value.(float64)
			g.Push(StackElement{Type: Float, Value: div})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Int:
			div := g.VariableMap[val2.Value.(string)].Value.(float64) / float64(val1.Value.(int))
			g.Push(StackElement{Type: Float, Value: div})
		default:
			return errors.New("ERROR: cannot perform DIV_OP on different types")
		}
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
		if val1.Value.(int) == 0 {
			return errors.New("ERROR: cannot divide by zero")
		}
		mod := val2.Value.(int) % val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: mod})
	// variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Int:
			if g.VariableMap[val1.Value.(string)].Value.(int) == 0 {
				return errors.New("ERROR: cannot divide by zero")
			}
			mod := g.VariableMap[val2.Value.(string)].Value.(int) % g.VariableMap[val1.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: mod})
		default:
			return errors.New("ERROR: cannot perform MOD_OP on different types")
		}
	// one is a variable and the other is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Int:
			if val1.Value.(int) == 0 {
				return errors.New("ERROR: cannot divide by zero")
			}
			mod := val2.Value.(int) % g.VariableMap[val1.Value.(string)].Value.(int)
			g.Push(StackElement{Type: Int, Value: mod})
		default:
			return errors.New("ERROR: cannot perform MOD_OP on different types")
		}
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
	case val1.Type == Int && val2.Type == Float:
		exp := math.Pow(val2.Value.(float64), float64(val1.Value.(int)))
		g.Push(StackElement{Type: Float, Value: exp})
	case val1.Type == Float && val2.Type == Int:
		exp := math.Pow(float64(val2.Value.(int)), val1.Value.(float64))
		g.Push(StackElement{Type: Float, Value: exp})
	// variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Int:
			exp := int(math.Pow(float64(g.VariableMap[val2.Value.(string)].Value.(int)), float64(g.VariableMap[val1.Value.(string)].Value.(int))))
			g.Push(StackElement{Type: Int, Value: exp})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Float:
			exp := math.Pow(g.VariableMap[val2.Value.(string)].Value.(float64), g.VariableMap[val1.Value.(string)].Value.(float64))
			g.Push(StackElement{Type: Float, Value: exp})
		// one is an int and the other is a float
		case g.VariableMap[val1.Value.(string)].Type == Int && g.VariableMap[val2.Value.(string)].Type == Float:
			exp := math.Pow(float64(g.VariableMap[val2.Value.(string)].Value.(int)), g.VariableMap[val1.Value.(string)].Value.(float64))
			g.Push(StackElement{Type: Float, Value: exp})
		case g.VariableMap[val1.Value.(string)].Type == Float && g.VariableMap[val2.Value.(string)].Type == Int:
			exp := math.Pow(g.VariableMap[val2.Value.(string)].Value.(float64), float64(g.VariableMap[val1.Value.(string)].Value.(int)))
			g.Push(StackElement{Type: Float, Value: exp})
		default:
			return errors.New("ERROR: cannot perform EXP_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch {
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Int:
			exp := int(math.Pow(float64(val2.Value.(int)), float64(g.VariableMap[val1.Value.(string)].Value.(int))))
			g.Push(StackElement{Type: Int, Value: exp})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Float:
			exp := math.Pow(val2.Value.(float64), g.VariableMap[val1.Value.(string)].Value.(float64))
			g.Push(StackElement{Type: Float, Value: exp})
		// one is an int and the other is a float
		case g.VariableMap[val1.Value.(string)].Type == Int && val2.Type == Float:
			exp := math.Pow(float64(val2.Value.(int)), g.VariableMap[val1.Value.(string)].Value.(float64))
			g.Push(StackElement{Type: Float, Value: exp})
		case g.VariableMap[val1.Value.(string)].Type == Float && val2.Type == Int:
			exp := math.Pow(val2.Value.(float64), float64(g.VariableMap[val1.Value.(string)].Value.(int)))
			g.Push(StackElement{Type: Float, Value: exp})
		default:
			return errors.New("ERROR: cannot perform EXP_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch {
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Int:
			exp := int(math.Pow(float64(g.VariableMap[val2.Value.(string)].Value.(int)), float64(val1.Value.(int))))
			g.Push(StackElement{Type: Int, Value: exp})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Float:
			exp := math.Pow(g.VariableMap[val2.Value.(string)].Value.(float64), val1.Value.(float64))
			g.Push(StackElement{Type: Float, Value: exp})
		// one is an int and the other is a float
		case g.VariableMap[val2.Value.(string)].Type == Int && val1.Type == Float:
			exp := math.Pow(float64(g.VariableMap[val2.Value.(string)].Value.(int)), val1.Value.(float64))
			g.Push(StackElement{Type: Float, Value: exp})
		case g.VariableMap[val2.Value.(string)].Type == Float && val1.Type == Int:
			exp := math.Pow(g.VariableMap[val2.Value.(string)].Value.(float64), float64(val1.Value.(int)))
			g.Push(StackElement{Type: Float, Value: exp})
		default:
			return errors.New("ERROR: cannot perform EXP_OP on different types")
		}
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
	case val.Type == Float:
		g.Push(StackElement{Type: Float, Value: val.Value.(float64) + 1})
	// variable increment
	case val.Type == Identifier:
		_, exists := g.VariableMap[val.Value.(string)]

		if !exists {
			return fmt.Errorf("ERROR: variable %v has not been declared", val.Value.(string))
		}

		switch {
		case g.VariableMap[val.Value.(string)].Type == Int:
			incVal := g.VariableMap[val.Value.(string)].Value.(int) + 1
			temp := g.VariableMap[val.Value.(string)]
			temp.Value = incVal
			g.VariableMap[val.Value.(string)] = temp
		case g.VariableMap[val.Value.(string)].Type == Float:
			incVal := g.VariableMap[val.Value.(string)].Value.(float64) + 1
			temp := g.VariableMap[val.Value.(string)]
			temp.Value = incVal
			g.VariableMap[val.Value.(string)] = temp
		default:
			return errors.New("ERROR: cannot perform INC_OP on different types")
		}
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
	case val.Type == Float:
		g.Push(StackElement{Type: Float, Value: val.Value.(float64) - 1})
	// variable decrement
	case val.Type == Identifier:
		_, exists := g.VariableMap[val.Value.(string)]

		if !exists {
			return fmt.Errorf("ERROR: variable %v has not been declared", val.Value.(string))
		}

		switch {
		case g.VariableMap[val.Value.(string)].Type == Int:
			decVal := g.VariableMap[val.Value.(string)].Value.(int) - 1
			temp := g.VariableMap[val.Value.(string)]
			temp.Value = decVal
			g.VariableMap[val.Value.(string)] = temp
		case g.VariableMap[val.Value.(string)].Type == Float:
			decVal := g.VariableMap[val.Value.(string)].Value.(float64) - 1
			temp := g.VariableMap[val.Value.(string)]
			temp.Value = decVal
			g.VariableMap[val.Value.(string)] = temp
		default:
			return errors.New("ERROR: cannot perform DEC_OP on different types")
		}
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
	case val.Type == Bool:
		g.Push(StackElement{Type: Bool, Value: !val.Value.(bool)})
	// variable negation
	case val.Type == Identifier:
		_, exists := g.VariableMap[val.Value.(string)]

		if !exists {
			return fmt.Errorf("ERROR: variable %v has not been declared", val.Value.(string))
		}

		switch g.VariableMap[val.Value.(string)].Type {
		case Int:
			negVal := -g.VariableMap[val.Value.(string)].Value.(int)
			temp := g.VariableMap[val.Value.(string)]
			temp.Value = negVal
			g.VariableMap[val.Value.(string)] = temp
		case Float:
			negVal := -g.VariableMap[val.Value.(string)].Value.(float64)
			temp := g.VariableMap[val.Value.(string)]
			temp.Value = negVal
			g.VariableMap[val.Value.(string)] = temp
		case Bool:
			negVal := !g.VariableMap[val.Value.(string)].Value.(bool)
			temp := g.VariableMap[val.Value.(string)]
			temp.Value = negVal
			g.VariableMap[val.Value.(string)] = temp
		default:
			return errors.New("ERROR: cannot perform NEG_OP on different types")
		}
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
	if val1.Type != Bool && val1.Type != Identifier || val2.Type != Bool && val2.Type != Identifier {
		return errors.New("ERROR: cannot perform AND_OP on non boolean types")
	}

	switch {
	case val1.Value.(bool) && val2.Value.(bool):
		g.Push(StackElement{Type: Bool, Value: true})
	// using variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type {
			switch g.VariableMap[val1.Value.(string)].Type {
			case Bool:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(bool) && g.VariableMap[val2.Value.(string)].Value.(bool)})
			default:
				return errors.New("ERROR: cannot perform AND_OP on non boolean types")
			}
		} else {
			return errors.New("ERROR: cannot perform AND_OP on non boolean types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch g.VariableMap[val1.Value.(string)].Type {
		case Bool:
			g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(bool) && val2.Value.(bool)})
		default:
			return errors.New("ERROR: cannot perform AND_OP on non boolean types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch g.VariableMap[val2.Value.(string)].Type {
		case Bool:
			g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(bool) && val1.Value.(bool)})
		default:
			return errors.New("ERROR: cannot perform AND_OP on non boolean types")
		}
	default:
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
	if val1.Type != Bool && val1.Type != Identifier || val2.Type != Bool && val2.Type != Identifier {
		return errors.New("ERROR: cannot perform OR_OP on non boolean types")
	}

	switch {
	case val1.Value.(bool) || val2.Value.(bool):
		g.Push(StackElement{Type: Bool, Value: true})
	// using variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type {
			switch g.VariableMap[val1.Value.(string)].Type {
			case Bool:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(bool) || g.VariableMap[val2.Value.(string)].Value.(bool)})
			default:
				return errors.New("ERROR: cannot perform OR_OP on non boolean types")
			}
		} else {
			return errors.New("ERROR: cannot perform OR_OP on non boolean types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch g.VariableMap[val1.Value.(string)].Type {
		case Bool:
			g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(bool) || val2.Value.(bool)})
		default:
			return errors.New("ERROR: cannot perform OR_OP on non boolean types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		switch g.VariableMap[val2.Value.(string)].Type {
		case Bool:
			g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(bool) || val1.Value.(bool)})
		default:
			return errors.New("ERROR: cannot perform OR_OP on non boolean types")
		}
	default:
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
	// variable negation
	case val1.Type == Identifier:
		_, exists := g.VariableMap[val1.Value.(string)]

		if !exists {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		switch g.VariableMap[val1.Value.(string)].Type {
		case Bool:
			g.Push(StackElement{Type: Bool, Value: !g.VariableMap[val1.Value.(string)].Value.(bool)})
		case Int:
			g.Push(StackElement{Type: Int, Value: g.VariableMap[val1.Value.(string)].Value.(int) * -1})
		case Float:
			g.Push(StackElement{Type: Float, Value: g.VariableMap[val1.Value.(string)].Value.(float64) * -1})
		}
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

	// if types are not equal, return false
	if val1.Type != val2.Type {
		g.Push(StackElement{Type: Bool, Value: false})
		return nil
	}

	switch {
	case val1.Type == Int && val2.Type == Int:
		g.Push(StackElement{Type: Bool, Value: val1.Value == val2.Value})
	case val1.Type == Float && val2.Type == Float:
		g.Push(StackElement{Type: Bool, Value: val1.Value == val2.Value})
	case val1.Type == String && val2.Type == String:
		g.Push(StackElement{Type: Bool, Value: val1.Value == val2.Value})
	case val1.Type == Bool && val2.Type == Bool:
		g.Push(StackElement{Type: Bool, Value: val1.Value == val2.Value})
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type {
			switch g.VariableMap[val1.Value.(string)].Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(int) == g.VariableMap[val2.Value.(string)].Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(float64) == g.VariableMap[val2.Value.(string)].Value.(float64)})
			case String:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(string) == g.VariableMap[val2.Value.(string)].Value.(string)})
			case Bool:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(bool) == g.VariableMap[val2.Value.(string)].Value.(bool)})
			}
		} else {
			g.Push(StackElement{Type: Bool, Value: false})
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == val2.Type {
			switch val2.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(int) == val2.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(float64) == val2.Value.(float64)})
			case String:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(string) == val2.Value.(string)})
			case Bool:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Value.(bool) == val2.Value.(bool)})
			}
		} else {
			g.Push(StackElement{Type: Bool, Value: false})
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val2.Value.(string)].Type == val1.Type {
			switch val1.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) == val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) == val1.Value.(float64)})
			case String:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(string) == val1.Value.(string)})
			case Bool:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(bool) == val1.Value.(bool)})
			}
		} else {
			g.Push(StackElement{Type: Bool, Value: false})
		}
	default:
		g.Push(StackElement{Type: Bool, Value: false})
	}
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

	switch {
	case val1.Type == val2.Type:
		g.Push(StackElement{Type: Bool, Value: true})
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type})
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		g.Push(StackElement{Type: Bool, Value: g.VariableMap[val1.Value.(string)].Type == val2.Type})
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Type == val1.Type})
	default:
		g.Push(StackElement{Type: Bool, Value: false})
	}

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
	// using variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type {
			switch g.VariableMap[val1.Value.(string)].Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) > g.VariableMap[val1.Value.(string)].Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) > g.VariableMap[val1.Value.(string)].Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == val2.Type {
			switch val2.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) > val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) > val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val2.Value.(string)].Type == val1.Type {
			switch val1.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) > val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) > val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_OP on different types")
		}
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
	// using variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type {
			switch g.VariableMap[val1.Value.(string)].Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) < g.VariableMap[val1.Value.(string)].Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) < g.VariableMap[val1.Value.(string)].Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform LS_THAN_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == val2.Type {
			switch val2.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) < val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) < val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform LS_THAN_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val2.Value.(string)].Type == val1.Type {
			switch val1.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) < val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) < val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform LS_THAN_OP on different types")
		}
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
	// using variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type {
			switch g.VariableMap[val1.Value.(string)].Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) >= g.VariableMap[val1.Value.(string)].Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) >= g.VariableMap[val1.Value.(string)].Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_EQ_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == val2.Type {
			switch val2.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) >= val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) >= val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_EQ_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val2.Value.(string)].Type == val1.Type {
			switch val1.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) >= val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) >= val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_EQ_OP on different types")
		}
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
	// using variables
	case val1.Type == Identifier && val2.Type == Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == g.VariableMap[val2.Value.(string)].Type {
			switch g.VariableMap[val1.Value.(string)].Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) <= g.VariableMap[val1.Value.(string)].Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) <= g.VariableMap[val1.Value.(string)].Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_EQ_OP on different types")
		}
	// val1 is a variable and val2 is not
	case val1.Type == Identifier && val2.Type != Identifier:
		_, exists1 := g.VariableMap[val1.Value.(string)]

		if !exists1 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val1.Value.(string))
		}

		if g.VariableMap[val1.Value.(string)].Type == val2.Type {
			switch val2.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) <= val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) <= val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_EQ_OP on different types")
		}
	// val2 is a variable and val1 is not
	case val1.Type != Identifier && val2.Type == Identifier:
		_, exists2 := g.VariableMap[val2.Value.(string)]

		if !exists2 {
			return fmt.Errorf("ERROR: variable %v has not been declared", val2.Value.(string))
		}

		if g.VariableMap[val2.Value.(string)].Type == val1.Type {
			switch val1.Type {
			case Int:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(int) <= val1.Value.(int)})
			case Float:
				g.Push(StackElement{Type: Bool, Value: g.VariableMap[val2.Value.(string)].Value.(float64) <= val1.Value.(float64)})
			}
		} else {
			return errors.New("ERROR: cannot perform GT_THAN_EQ_OP on different types")
		}
	default:
		return errors.New("ERROR: cannot perform LS_THAN_EQ_OP on different types")
	}
	return nil
}

func (g *Gorth) VarAssign() error {
	val1, err := g.Pop()

	if err != nil {
		return errors.New("ERROR: stack is empty, cannot assign from an empty stack")
	}

	// this should be the variable on the stack
	val2, err := g.Pop()

	if err != nil {
		return errors.New("ERROR: stack is empty, cannot assign from an empty stack")
	}

	// if this value is not an identifier
	if val2.Type != Identifier {
		return errors.New("ERROR: cannot assign a value to a non-variable")
	}

	// if the value is an identifier
	// check if the variable has been declared
	// if it has, update the value
	variable, exists := g.VariableMap[val2.Value.(string)]

	if !exists {
		return fmt.Errorf("ERROR: variable %v has not been declared on the stack", val2.Value.(string))
	}

	if g.VariableMap[val2.Value.(string)].Const {
		return fmt.Errorf("ERROR: variable %v is a constant and cannot be reassigned", val2.Value.(string))
	}

	// change the value of the variable in the variable map
	switch variable.Type {
	case Int:
		if val1.Type == Int {
			g.VariableMap[val2.Value.(string)] = Variable{Type: Int, Value: val1.Value.(int), Name: val2.Value.(string), Const: false}
		} else {
			return errors.New("ERROR: cannot assign a non-integer value to an integer variable")
		}
	case Float:
		if val1.Type == Float {
			g.VariableMap[val2.Value.(string)] = Variable{Type: Float, Value: val1.Value.(float64), Name: val2.Value.(string), Const: false}
		} else {
			return errors.New("ERROR: cannot assign a non-float value to a float variable")
		}
	case Bool:
		if val1.Type == Bool {
			g.VariableMap[val2.Value.(string)] = Variable{Type: Bool, Value: val1.Value.(bool), Name: val2.Value.(string), Const: false}
		} else {
			return errors.New("ERROR: cannot assign a non-boolean value to a boolean variable")
		}
	case String:
		if val1.Type == String {
			g.VariableMap[val2.Value.(string)] = Variable{Type: String, Value: val1.Value.(string), Name: val2.Value.(string), Const: false}
		} else {
			return errors.New("ERROR: cannot assign a non-string value to a string variable")
		}
	default:
		return errors.New("ERROR: cannot assign a value to a non-variable")
	}

	return nil
}

func (g *Gorth) PrintStack() {
	fmt.Printf("Program stack: %v\n", g.ExecStack)
}

func (g *Gorth) ExecuteProgram(program []StackElement) error {
	for _, op := range program {
		if g.DebugMode {
			fmt.Println("Current operation: " + fmt.Sprintf("%v", op.Type == Operator))
			fmt.Println("Current Stack: ", g.ExecStack)
		}

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
			case VAR_ASSIGN_OP:
				err := g.VarAssign()
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
	program, variables, err := Tokenize(strings.Join(lines, " "))

	if err != nil {
		panic(err)
	}

	// create a new gorth instance
	g := NewGorth(debugMode, strictMode)

	g.VariableMap = variables

	if g.DebugMode {
		fmt.Println("Variables: ", g.VariableMap)
		fmt.Println("Program: ", program)
	}

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
