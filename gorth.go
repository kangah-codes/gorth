package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// lets start with basic operations, add, minus, multiply, divide
const (
	// Arithmetic operations
	ADD_OP = iota
	SUB_OP
	MUL_OP
	DIV_OP
	MOD_OP
	EXP_OP
	INC_OP
	DEC_OP
	NEG_OP
	// Stack manipulation operations
	SWP_OP
	DUP_OP
	DRP_OP
)

const (
	MAX_STACK_SIZE = 999_999
	DEBUG_MODE     = false
	STRICT_MODE    = false
)

const (
	// StackElement types
	IntType = iota
	OpType
)

type StackElement struct {
	Value     int
	ValueType int
}

// Gorth program
type Gorth struct {
	ProgramStack []StackElement
	DebugMode    bool
	StrictMode   bool // in strict mode, the entire stack does not need to be consumed
}

// helper functions
func assert(condition bool, msg string) error {
	if !condition {
		return errors.New(msg)
	}

	return nil
}

// Create a program stack with a default size
func NewProgramStack(debug_mode bool, strict_mode bool) *Gorth {
	return &Gorth{
		ProgramStack: []StackElement{},
		DebugMode:    debug_mode,
		StrictMode:   strict_mode,
	}
}

// Push operation to the program stack
func (g *Gorth) Push(val StackElement) error {
	if len(g.ProgramStack) > MAX_STACK_SIZE {
		return errors.New("ERROR: Stack is full, cannot push element")
	}

	g.ProgramStack = append(g.ProgramStack, val)
	return nil
}

// print the program stack
func (g *Gorth) PrintStack() {
	for _, element := range g.ProgramStack {
		fmt.Println(element)
	}
}

// Pop operation from the program stack
func (g *Gorth) Pop() (StackElement, error) {
	if len(g.ProgramStack) < 1 {
		return StackElement{}, errors.New("ERROR: Stack is empty, cannot pop from stack")
	}

	// get the last element from the stack
	val := g.ProgramStack[len(g.ProgramStack)-1]

	// remove the last element from the stack
	g.ProgramStack = g.ProgramStack[:len(g.ProgramStack)-1]

	return val, nil
}

// Dump operation which pops off the last element on the stack to stdout
func (g *Gorth) Drop() {
	val, err := g.Pop()

	// check if pop op returned any errors
	if err != nil {
		fmt.Println(err.Error())
	}

	// can only dump values and not operations from stack. a value should always be the last on the stack
	if val.ValueType == OpType {
		fmt.Println("ERROR: Values must be the last in the stack")
		os.Exit(1)
	}

	// everything else
	fmt.Println(val.Value)
}

// Peek operation which returns the last element on the stack
func (g *Gorth) Peek() (StackElement, error) {
	if len(g.ProgramStack) < 1 {
		return StackElement{}, errors.New("ERROR: Stack is empty, cannot peek from stack")
	}

	// get the last element from the stack
	val := g.ProgramStack[len(g.ProgramStack)-1]

	return val, nil
}

// Function to print the top element of the stack
func (g *Gorth) Print() {
	val, err := g.Peek()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(val.Value)
}

func (g *Gorth) Sim() error {
	// check if program length exceeds max stack size
	if len(g.ProgramStack) > MAX_STACK_SIZE {
		return errors.New("ERROR: Program size exceeds max stack size")
	}

	// Iterate through the program
	for _, op := range g.ProgramStack {
		if op.ValueType == OpType {
			switch op.Value {
			case ADD_OP, SUB_OP, MUL_OP, DIV_OP, EXP_OP, MOD_OP, INC_OP, DEC_OP, NEG_OP, SWP_OP:
				// Collect operands from the stack
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				// For binary operations, collect the second operand
				if op.Value != INC_OP && op.Value != DEC_OP && op.Value != NEG_OP {
					b, err := g.Pop()
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}

					// Perform operation based on the operator
					switch op.Value {
					case ADD_OP:
						g.Push(StackElement{Value: a.Value + b.Value, ValueType: IntType})
					case SUB_OP:
						g.Push(StackElement{Value: b.Value - a.Value, ValueType: IntType})
					case MUL_OP:
						g.Push(StackElement{Value: a.Value * b.Value, ValueType: IntType})
					case DIV_OP:
						g.Push(StackElement{Value: b.Value / a.Value, ValueType: IntType})
					case EXP_OP:
						result := int(math.Pow(float64(b.Value), float64(a.Value)))
						g.Push(StackElement{Value: result, ValueType: IntType})
					case MOD_OP:
						g.Push(StackElement{Value: b.Value % a.Value, ValueType: IntType})
					case SWP_OP:
						g.Push(a)
						g.Push(b)
					}
				} else {
					// For unary operations, perform operation on a single operand
					switch op.Value {
					case INC_OP:
						g.Push(StackElement{Value: a.Value + 1, ValueType: IntType})
					case DEC_OP:
						g.Push(StackElement{Value: a.Value - 1, ValueType: IntType})
					case NEG_OP:
						g.Push(StackElement{Value: a.Value * -1, ValueType: IntType})
					}
				}

			case DUP_OP:
				// Duplicate the top element of the stack
				if len(g.ProgramStack) < 1 {
					fmt.Println("ERROR: Stack is empty, cannot duplicate")
					os.Exit(1)
				}

				lastElement := g.ProgramStack[len(g.ProgramStack)-1]

				// We can only duplicate values and not operations
				if lastElement.ValueType == OpType {
					fmt.Println("ERROR: Cannot duplicate operators")
					os.Exit(1)
				}

				// Push the last element to the stack
				g.Push(StackElement{Value: lastElement.Value, ValueType: lastElement.ValueType})

			case DRP_OP:
				// Drop the stack
				g.Drop()

			default:
				fmt.Println("ERROR: Unknown operation on stack")
				os.Exit(1)
			}
		} else {
			// Push value onto the stack
			g.Push(StackElement{Value: op.Value, ValueType: IntType})
		}
	}

	if g.StrictMode {
		// check if the stack is empty
		// in strict mode, all operands on the stack must be consumed
		return assert(len(g.ProgramStack) == 0, "ERROR: Program stack has unused elements")
	}

	return nil
}

// Function to parse a line and convert it into a slice of StackElement
func parseLine(line string) ([]StackElement, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil // Ignore empty lines
	}

	// Check if the line ends with a semicolon
	if !strings.HasSuffix(line, ";") {
		return nil, fmt.Errorf("ERROR: Line does not end with semicolon: %s", line)
	}

	// Remove the semicolon at the end
	line = strings.TrimSuffix(line, ";")

	// Split the line into operands and operator
	elements := strings.Split(line, " ")

	var program []StackElement

	// Parse operands
	for _, elem := range elements[:len(elements)-1] {
		// Check if the element is an integer
		if Value, err := strconv.Atoi(elem); err == nil {
			// Element is an integer operand
			program = append(program, StackElement{ValueType: IntType, Value: Value})
		} else {
			return nil, fmt.Errorf("ERROR: invalid operand: %s", elem)
		}
	}

	// Parse operator
	switch elements[len(elements)-1] {
	case "add":
		program = append(program, StackElement{Value: ADD_OP, ValueType: OpType})
	case "sub":
		program = append(program, StackElement{Value: SUB_OP, ValueType: OpType})
	case "drop":
		program = append(program, StackElement{Value: DRP_OP, ValueType: OpType})
	// Add cases for other operators as needed
	default:
		return nil, fmt.Errorf("ERROR: invalid operator: %s", elements[len(elements)-1])
	}

	return program, nil
}

// Function to read a program from a file and parse it into a slice of StackElement
func readProgramFromFile(filename string) ([]StackElement, error) {
	// check if filename ends in .gorth
	if !strings.HasSuffix(filename, ".gorth") {
		return nil, errors.New("ERROR: Invalid file extension. File must end with .gorth")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var program []StackElement
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		element, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		program = append(program, element...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return program, nil
}

func printUsage() {
	fmt.Println("Usage: gorth <filename> -d -s")
	fmt.Println("  filename: The name of the file containing the gorth program")
	fmt.Println("  -d: Enable debug mode")
	fmt.Println("  -s: Enable strict mode")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	filename := os.Args[1]
	debug_mode := len(os.Args) > 2 && os.Args[2] == "-d"
	strict_mode := len(os.Args) > 3 && os.Args[3] == "-s"

	// Read the program from the file
	program, err := readProgramFromFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var now = time.Now()
	var g = Gorth{
		ProgramStack: program,
		DebugMode:    debug_mode,
		StrictMode:   strict_mode,
	}

	err = g.Sim()

	var after_sim = time.Now()

	if err != nil {
		fmt.Println("Program simulation failed")
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Program simulation completed in %v seconds\n", after_sim.Sub(now).Seconds())
	}
}
