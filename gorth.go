package main

import (
	"errors"
	"fmt"
	"math"
	"os"
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
	DMP_OP
)

const (
	MAX_STACK_SIZE = 999_999
	DEBUG_MODE     = false
	STRICT_MODE    = false
)

type StackElement struct {
	Value       int
	IsProgramOp bool
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
func (g *Gorth) DumpStack() {
	val, err := g.Pop()

	// check if pop op returned any errors
	if err != nil {
		fmt.Println(err.Error())
	}

	// can only dump values and not operations from stack. a value should always be the last on the stack
	if val.IsProgramOp {
		fmt.Println("ERROR: Values must be the last in the stack")
		os.Exit(1)
	}

	// everything else
	fmt.Println(val.Value)
}

func (g *Gorth) Sim(program []StackElement) error {
	// check if program length exceeds max stack size
	if len(program) > MAX_STACK_SIZE {
		return errors.New("ERROR: Program size exceeds max stack size")
	}

	// iterate through the program
	for _, op := range program {
		if op.IsProgramOp {
			switch op.Value {
			case ADD_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE ADD_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				b, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: a.Value + b.Value, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER ADD_OP: %v\n", g.ProgramStack)
				}
			case SUB_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE SUB_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				b, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: b.Value - a.Value, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER SUB_OP: %v\n", g.ProgramStack)
				}
			case MUL_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE MUL_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				b, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: a.Value * b.Value, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER MUL_OP: %v\n", g.ProgramStack)
				}
			case DIV_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE DIV_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				b, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: b.Value / a.Value, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER DIV_OP: %v\n", g.ProgramStack)
				}
			case EXP_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE EXP_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				b, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				result := int(math.Pow(float64(b.Value), float64(a.Value)))
				g.Push(StackElement{Value: result, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER EXP_OP: %v\n", g.ProgramStack)
				}
			case INC_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE INC_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: a.Value + 1, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER INC_OP: %v\n", g.ProgramStack)
				}
			case DEC_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE DEC_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: a.Value - 1, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER DEC_OP: %v\n", g.ProgramStack)
				}
			case NEG_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE NEG_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: a.Value * -1, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER NEG_OP: %v\n", g.ProgramStack)
				}
			case MOD_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE MOD_OP: %v\n", g.ProgramStack)
				}
				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				b, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(StackElement{Value: b.Value % a.Value, IsProgramOp: false})
				if g.DebugMode {
					fmt.Printf("AFTER MOD_OP: %v\n", g.ProgramStack)
				}
			case DUP_OP:
				if len(g.ProgramStack) < 1 {
					return errors.New("ERROR: Stack is empty, cannot duplicate")
				}

				if g.DebugMode {
					fmt.Printf("BEFORE DUPOP: %v\n", g.ProgramStack)
				}

				var last_element = g.ProgramStack[len(g.ProgramStack)-1]

				// we can only duplicate values and not operations
				if last_element.IsProgramOp {
					return errors.New("ERROR: Cannot duplicate operators")
				}

				// push the last element to the stack
				g.Push(StackElement{Value: last_element.Value, IsProgramOp: last_element.IsProgramOp})

				if g.DebugMode {
					fmt.Printf("AFTER DUPOP: %v\n", g.ProgramStack)
				}
			case DMP_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE DMP_OP: %v\n", g.ProgramStack)
				}
				g.DumpStack()
				if g.DebugMode {
					fmt.Printf("AFTER DMP_OP: %v\n", g.ProgramStack)
				}
			case SWP_OP:
				if g.DebugMode {
					fmt.Printf("BEFORE SWP_OP: %v\n", g.ProgramStack)
				}

				a, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				b, err := g.Pop()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				g.Push(a)
				g.Push(b)

				if g.DebugMode {
					fmt.Printf("AFTER SWP_OP: %v\n", g.ProgramStack)
				}

			default:
				fmt.Println("ERROR: Unknown operation on stack")
				os.Exit(1)
			}
		} else {
			if g.DebugMode {
				fmt.Printf("BEFORE VAL: %v\n", g.ProgramStack)
			}
			g.Push(StackElement{Value: op.Value, IsProgramOp: false})
			if g.DebugMode {
				fmt.Printf("AFTER VAL: %v\n", g.ProgramStack)
			}
		}
	}

	if g.StrictMode {
		// check if the stack is empty
		// in strict mode, all operands on the stack must be consumed
		return assert(len(g.ProgramStack) == 0, "ERROR: Program stack has unused elements")
	}

	return nil
}

func main() {
	var now = time.Now()
	var g = Gorth{
		ProgramStack: []StackElement{},
		DebugMode:    false,
		StrictMode:   false,
	}

	var program = []StackElement{
		{Value: 10, IsProgramOp: false},
		{Value: 20, IsProgramOp: false},
		{Value: SWP_OP, IsProgramOp: true},
		{Value: DMP_OP, IsProgramOp: true},
		{Value: DMP_OP, IsProgramOp: true},
	}

	err := g.Sim(program)

	var after_sim = time.Now()

	if err != nil {
		fmt.Println("Program simulation failed")
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Program simulation completed in %v seconds\n", after_sim.Sub(now).Seconds())
	}
}
