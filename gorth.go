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
	DUMP_OP

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
	case val1.Type == Int && val2.Type == Int:
		sum := val1.Value.(int) + val2.Value.(int)
		g.Push(StackElement{Type: Int, Value: sum})
	case val1.Type == String && val2.Type == String:
		concat := val2.Value.(string) + val1.Value.(string)
		g.Push(StackElement{Type: String, Value: concat})
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
	case val1.Type == Int && val2.Type == Int:
		sub := val2.Value.(int) - val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: sub})
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
	case val1.Type == Int && val2.Type == Int:
		mul := val1.Value.(int) * val2.Value.(int)
		g.Push(StackElement{Type: Int, Value: mul})
	case val1.Type == String && val2.Type == Int:
		str := val1.Value.(string)
		num := val2.Value.(int)
		var concat string
		for i := 0; i < num; i++ {
			concat += str
		}
		g.Push(StackElement{Type: String, Value: concat})
	case val1.Type == Int && val2.Type == String:
		str := val2.Value.(string)
		num := val1.Value.(int)
		var concat string
		for i := 0; i < num; i++ {
			concat += str
		}
		g.Push(StackElement{Type: String, Value: concat})
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
	case val1.Type == Int && val2.Type == Int:
		div := val2.Value.(int) / val1.Value.(int)
		g.Push(StackElement{Type: Int, Value: div})
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
	case val1.Type == Int && val2.Type == Int:
		exp := int(math.Pow(float64(val2.Value.(int)), float64(val1.Value.(int))))
		g.Push(StackElement{Type: Int, Value: exp})
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

func main() {
	start := time.Now()

	g := NewGorth(true, true)

	program := []StackElement{
		// perform this (3 * 2) + 1
		{Value: 3, Type: Int},
		{Value: 2, Type: Int},
		{Value: MUL_OP, Type: Operator},
		{Value: 1, Type: Int},
		{Value: ADD_OP, Type: Operator},
		{Value: DUMP_OP, Type: Operator},
		// now print hello world
		{Value: "hello ", Type: String},
		{Value: "world", Type: String},
		{Value: ADD_OP, Type: Operator},
		{Value: PRINT_OP, Type: Operator},
		{Value: DROP_OP, Type: Operator},
	}

	if g.DebugMode {
		fmt.Printf("Program stack at start of execution\n\t%v\n", g.ExecStack)
	}

	err := g.ExecuteProgram(program)

	end := time.Now()

	if err != nil {
		fmt.Println("Program simulation failed")
		fmt.Println(err)
	} else {
		fmt.Printf("Program simulation completed in %v seconds\n", end.Sub(start).Seconds())
	}
}
