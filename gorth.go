package main

import (
	"errors"
	"fmt"
	"math"
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

	// Print operation
	PRINT_OP
)

type ValueType int

const (
	Int ValueType = iota
	String
	Bool
	Operator
)

type StackElement struct {
	Type  ValueType
	Value interface{} // Use interface{} to support both int and string values
}

type Gorth struct {
	ProgramStack []StackElement
	DebugMode    bool
	StrictMode   bool
}

func NewGorth(debugMode, strictMode bool) *Gorth {
	return &Gorth{
		ProgramStack: []StackElement{},
		DebugMode:    debugMode,
		StrictMode:   strictMode,
	}
}

// Push adds a value to the program stack.
// It appends the given value to the ProgramStack slice.
func (g *Gorth) Push(val StackElement) {
	g.ProgramStack = append(g.ProgramStack, val)
}

// Pop removes and returns the top element from the program stack.
// If the stack is empty, it returns an empty StackElement and an error.
func (g *Gorth) Pop() (StackElement, error) {
	if len(g.ProgramStack) < 1 {
		return StackElement{}, errors.New("ERROR: cannot pop from an empty stack")
	}
	val := g.ProgramStack[len(g.ProgramStack)-1]
	g.ProgramStack = g.ProgramStack[:len(g.ProgramStack)-1]
	return val, nil
}

// Drop removes the top element from the program stack.
// It panics if the stack is empty.
func (g *Gorth) Drop() {
	if len(g.ProgramStack) < 1 {
		panic("ERROR: cannot drop from an empty stack")
	}
	g.Pop()
}

// Peek returns the top element of the program stack without removing it.
// It returns an error if the stack is empty.
func (g *Gorth) Peek() (StackElement, error) {
	if len(g.ProgramStack) < 1 {
		return StackElement{}, errors.New("ERROR: cannot peek at an empty stack")
	}
	return g.ProgramStack[len(g.ProgramStack)-1], nil
}

// Print prints the top element of the Gorth stack.
// It checks the type of the top element and prints its value accordingly.
// If the top element is not a printable type, it panics with an error message.
func (g *Gorth) Print() {
	val, err := g.Peek()
	if err != nil {
		panic(err)
	}
	switch val.Type {
	case Int:
		fmt.Println(val.Value)
	case String:
		fmt.Println(val.Value)
	default:
		panic("ERROR: top element is not a printable type")
	}
}

// Add performs addition or concatenation operation on the top two elements of the stack.
// If both elements are integers, it performs integer addition.
// If both elements are strings, it concatenates them.
// Returns an error if the operation cannot be performed due to different types.
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
	case val1.Type == Int && val2.Type == Int:
		// If both are integers, perform integer addition
		sum := val1.Value.(int) + val2.Value.(int)
		g.Push(StackElement{Type: Int, Value: sum})
	case val1.Type == String && val2.Type == String:
		// If both are strings, concatenate them
		concat := val2.Value.(string) + val1.Value.(string)
		g.Push(StackElement{Type: String, Value: concat})
	default:
		return errors.New("ERROR: cannot perform ADD_OP on different types")
	}
	return nil
}

// Sub subtracts the top two elements from the stack and pushes the result back onto the stack.
// It returns an error if the operation cannot be performed.
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
	case val1.Type == Int && val2.Type == Int:
		sub := val2.Value.(int) - val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: sub})
	case (val1.Type == String && val2.Type == Int) || (val1.Type == Int && val2.Type == String):
		return errors.New("ERROR: cannot perform SUB_OP on different types")
	case (val1.Type == String && val2.Type == String):
		return errors.New("ERROR: cannot perform SUB_OP on strings")
	default:
		return errors.New("ERROR: cannot perform SUB_OP on unknown types")
	}

	return nil
}

// Mul performs the multiplication operation on the top two elements of the stack.
// If both elements are integers, it multiplies them and pushes the result onto the stack.
// If one element is a string and the other is an integer, it repeats the string onto the stack the number of times specified by the integer.
// If both elements are strings, it returns an error indicating that the multiplication operation is not supported on strings.
// If the types of the top two elements are not supported, it returns an error indicating that the multiplication operation is not supported on the given types.
// Returns an error if there are not enough elements on the stack or if any other error occurs.
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
	case val1.Type == Int && val2.Type == Int:
		mul := val1.Value.(int) * val2.Value.(int)
		g.Push(StackElement{Type: Int, Value: mul})
	case (val1.Type == String && val2.Type == Int) || (val1.Type == Int && val2.Type == String):
		// do like Python and repeat the string onto the stack the number of times of the integer
		if val1.Type == String {
			str := val1.Value.(string)
			num := val2.Value.(int)
			for i := 0; i < num; i++ {
				g.Push(StackElement{Type: String, Value: str})
			}
		} else {
			str := val2.Value.(string)
			num := val1.Value.(int)
			for i := 0; i < num; i++ {
				g.Push(StackElement{Type: String, Value: str})
			}
		}
	case (val1.Type == String && val2.Type == String):
		return errors.New("ERROR: cannot perform MUL_OP on strings")
	default:
		return errors.New("ERROR: cannot perform MUL_OP on unknown types")
	}

	return nil
}

// Div performs division operation on the top two elements of the stack.
// It pops two elements from the stack and performs division based on their types.
// If both elements are integers, it performs integer division.
// If one element is a string and the other is an integer, it returns an error.
// If both elements are strings, it returns an error.
// If the types are unknown, it returns an error.
// The result of the division operation is pushed back onto the stack.
// Returns an error if there are not enough elements on the stack or if the division operation is not possible.
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
	case val1.Type == Int && val2.Type == Int:
		div := val2.Value.(int) / val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: div})
	case (val1.Type == String && val2.Type == Int) || (val1.Type == Int && val2.Type == String):
		return errors.New("ERROR: cannot perform DIV_OP on different types")
	case (val1.Type == String && val2.Type == String):
		return errors.New("ERROR: cannot perform DIV_OP on strings")
	default:
		return errors.New("ERROR: cannot perform DIV_OP on unknown types")
	}

	return nil
}

// Mod performs the modulo operation on the top two elements of the stack.
// It pops two elements from the stack and checks their types.
// If both elements are integers, it calculates the modulo and pushes the result back onto the stack.
// If one element is a string and the other is an integer, or if both elements are strings,
// it returns an error indicating that the modulo operation is not supported for those types.
// If the types are unknown, it returns an error indicating that the modulo operation is not supported for unknown types.
// Returns an error if there are not enough elements on the stack or if an error occurs during popping.
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
	case val1.Type == Int && val2.Type == Int:
		mod := val2.Value.(int) % val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: mod})
	case (val1.Type == String && val2.Type == Int) || (val1.Type == Int && val2.Type == String):
		return errors.New("ERROR: cannot perform MOD_OP on different types")
	case (val1.Type == String && val2.Type == String):
		return errors.New("ERROR: cannot perform MOD_OP on strings")
	default:
		return errors.New("ERROR: cannot perform MOD_OP on unknown types")
	}

	return nil
}

// Exp performs the exponentiation operation on the top two elements of the stack.
// It pops two elements from the stack and checks their types.
// If both elements are integers, it calculates the result of the exponentiation operation and pushes the result back onto the stack.
// If one element is a string and the other is an integer, or if both elements are strings, it returns an error indicating that the operation is not supported.
// If the types of the elements are unknown, it returns an error indicating that the operation is not supported.
// Returns an error if there are not enough elements on the stack or if an error occurs during the operation.
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
	case val1.Type == Int && val2.Type == Int:
		exp := int(math.Pow(float64(val2.Value.(int)), float64(val1.Value.(int))))
		g.Push(StackElement{Type: Int, Value: exp})
	case (val1.Type == String && val2.Type == Int) || (val1.Type == Int && val2.Type == String):
		return errors.New("ERROR: cannot perform EXP_OP on different types")
	case (val1.Type == String && val2.Type == String):
		return errors.New("ERROR: cannot perform EXP_OP on strings")
	default:
		return errors.New("ERROR: cannot perform EXP_OP on unknown types")
	}

	return nil
}

// Inc increments the value at the top of the stack by 1.
// If the value is an integer, it will be incremented by 1.
// If the value is a string, an error will be returned.
// If the value is of an unknown type, an error will be returned.
// Returns an error if the stack is empty or if an error occurs during the operation.
func (g *Gorth) Inc() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	switch {
	case val.Type == Int:
		g.Push(StackElement{Type: Int, Value: val.Value.(int) + 1})
	case val.Type == String:
		return errors.New("ERROR: cannot perform INC_OP on strings")
	default:
		return errors.New("ERROR: cannot perform INC_OP on unknown types")
	}

	return nil
}

// Dec decrements the value at the top of the stack by 1.
// If the value is an integer, it will be decremented by 1.
// If the value is a string, an error will be returned.
// If the value is of an unknown type, an error will be returned.
// Returns an error if the stack is empty or if an error occurs during the operation.
func (g *Gorth) Dec() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	switch {
	case val.Type == Int:
		g.Push(StackElement{Type: Int, Value: val.Value.(int) - 1})
	case val.Type == String:
		return errors.New("ERROR: cannot perform DEC_OP on strings")
	default:
		return errors.New("ERROR: cannot perform DEC_OP on unknown types")
	}

	return nil
}

// Neg performs the negation operation on the top element of the stack.
// If the top element is an integer, it negates the value.
// If the top element is a string, it returns an error indicating that negation is not supported for strings.
// If the top element is of an unknown type, it returns an error indicating that negation is not supported for unknown types.
// Returns an error if the stack is empty or if there is an error while popping the element from the stack.
func (g *Gorth) Neg() error {
	val, err := g.Pop()
	if err != nil {
		return err
	}

	switch {
	case val.Type == Int:
		g.Push(StackElement{Type: Int, Value: -val.Value.(int)})
	case val.Type == String:
		return errors.New("ERROR: cannot perform NEG_OP on strings")
	default:
		return errors.New("ERROR: cannot perform NEG_OP on unknown types")
	}

	return nil
}

// Swap swaps the top two values on the stack.
// It pops two values from the stack, and then pushes them back in reverse order.
// If there are less than two values on the stack, it returns an error.
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

// Dup duplicates the top value of the stack and pushes it onto the stack.
// It returns an error if the stack is empty.
func (g *Gorth) Dup() error {
	val, err := g.Peek()
	if err != nil {
		return err
	}

	g.Push(val)

	return nil
}

func (g *Gorth) PrintStack() {
	fmt.Printf("Program stack: %v\n", g.ProgramStack)
}

// ExecuteProgram executes a program by iterating over each operation in the program slice.
// If the operation is an operator, it performs the corresponding action on the Gorth instance.
// If the operation is not an operator, it pushes the operation onto the program stack.
// If an error occurs during execution, it is returned.
// If the Gorth instance is in strict mode, it checks if there are any unconsumed elements on the stack at the end of execution.
// If the Gorth instance is in debug mode, it prints the program stack at the end of execution.
// Returns nil if execution is successful.
func (g *Gorth) ExecuteProgram() error {
	for _, op := range g.ProgramStack {
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
			case PRINT_OP:
				g.Print()
			case DROP_OP:
				g.Drop()
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
			default:
				return fmt.Errorf("unknown operation: %v", op.Value)
			}
		} else {
			g.Push(op)
		}
	}

	// if we are in strict mode, make sure the entire stack is consumed
	if g.StrictMode {
		if len(g.ProgramStack) > 0 {
			return fmt.Errorf("ERROR: unconsumed elements remain on the stack\n\t%v", g.ProgramStack)
		}
	}

	// if we are in debug mode, print the stack
	if g.DebugMode {
		fmt.Printf("Program stack at end of execution\n\t%v\n", g.ProgramStack)
	}

	return nil
}

// main is the entry point of the program.
// It initializes a new Gorth interpreter, executes a program,
// and measures the time it takes to complete the simulation.
func main() {
	start := time.Now()

	g := NewGorth(true, true)

	g.ProgramStack = []StackElement{
		{Type: Int, Value: 10},
		{Type: Int, Value: 20},
		{Type: Operator, Value: ADD_OP},
		{Type: String, Value: "Hello, "},
		{Type: String, Value: "world!"},
		{Type: Operator, Value: ADD_OP},
		{Type: Operator, Value: PRINT_OP},
		{Type: Operator, Value: DROP_OP},
		{Type: Operator, Value: PRINT_OP},
		{Type: Operator, Value: DROP_OP},
	}

	err := g.ExecuteProgram()

	end := time.Now()

	if err != nil {
		fmt.Println("Program simulation failed")
		fmt.Println(err)
	} else {
		fmt.Printf("Program simulation completed in %v seconds\n", end.Sub(start).Seconds())
	}
}
