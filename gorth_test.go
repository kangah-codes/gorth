package main

import (
	"testing"
)

func TestPushOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 10, Type: Int})

	val, err := g.Pop()
	if err != nil {
		t.Errorf("Error popping from stack: %v", err)
	}
	if val.Value != 10 {
		t.Errorf("Expected popped value to be 10, got %d", val.Value)
	}
}

func TestAddOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: ADD_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing addition operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 15 {
		t.Errorf("Error in addition operation. Expected result 15, got %d", val.Value)
	}
}

func TestSubOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: SUB_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing subtraction operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 5 {
		t.Errorf("Error in subtraction operation. Expected result 5, got %d", val.Value)
	}
}

func TestMulOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: MUL_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing multiplication operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 50 {
		t.Errorf("Error in multiplication operation. Expected result 50, got %d", val.Value)
	}
}

func TestDivOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: 2, Type: Int})
	g.Push(StackElement{Value: DIV_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing division operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 5 {
		t.Errorf("Error in division operation. Expected result 5, got %d", val.Value)
	}
}

func TestModOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: 3, Type: Int})
	g.Push(StackElement{Value: MOD_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing modulus operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 1 {
		t.Errorf("Error in modulus operation. Expected result 1, got %d", val.Value)
	}
}

func TestExpOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 2, Type: Int})
	g.Push(StackElement{Value: 3, Type: Int})
	g.Push(StackElement{Value: EXP_OP, Type: Operator})

	g.PrintStack()

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing exponentiation operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 8 {
		t.Errorf("Error in exponentiation operation. Expected result 8, got %d", val.Value)
	}
}

func TestIncOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: INC_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing increment operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 6 {
		t.Errorf("Error in increment operation. Expected result 6, got %d", val.Value)
	}
}

func TestDecOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: DEC_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing decrement operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 4 {
		t.Errorf("Error in decrement operation. Expected result 4, got %d", val.Value)
	}
}

func TestNegOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: NEG_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing negation operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != -5 {
		t.Errorf("Error in negation operation. Expected result -5, got %d", val.Value)
	}
}

func TestDupOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: DUP_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing duplicate operation: %v", err)
	}

	val1, err := g.Pop()
	if err != nil || val1.Value != 5 {
		t.Errorf("Error in duplicate operation. Expected result 5, got %d", val1.Value)
	}

	val2, err := g.Pop()
	if err != nil || val2.Value != 5 {
		t.Errorf("Error in duplicate operation. Expected result 5, got %d", val2.Value)
	}
}

func TestSwpOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 5, Type: Int})
	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: SWAP_OP, Type: Operator})

	err := g.ExecuteProgram()
	if err != nil {
		t.Errorf("Error performing swap operation: %v", err)
	}

	val1, err := g.Pop()
	if err != nil || val1.Value != 5 {
		t.Errorf("Error in swap operation. Expected result 5, got %d", val1.Value)
	}

	val2, err := g.Pop()
	if err != nil || val2.Value != 10 {
		t.Errorf("Error in swap operation. Expected result 10, got %d", val2.Value)
	}
}

func TestDmpOperation(t *testing.T) {
	g := NewGorth(false, false)

	g.Push(StackElement{Value: 10, Type: Int})

	g.Drop()

	if len(g.ProgramStack) != 0 {
		t.Errorf("Error in drop operation. Expected stack to be empty, got %v", g.ProgramStack)
	}
}
