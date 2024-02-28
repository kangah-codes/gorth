package main

import (
	"reflect"
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 10, Type: Int},
			{Value: 5, Type: Int},
			{Value: ADD_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 10, Type: Int},
			{Value: 5, Type: Int},
			{Value: SUB_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: MUL_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 10, Type: Int},
			{Value: 2, Type: Int},
			{Value: DIV_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 10, Type: Int},
			{Value: 3, Type: Int},
			{Value: MOD_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 2, Type: Int},
			{Value: 3, Type: Int},
			{Value: EXP_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: INC_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: DEC_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing decrement operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != 4 {
		t.Errorf("Error in decrement operation. Expected result 4, got %d", val.Value)
	}
}

func TestDupOperation(t *testing.T) {
	g := NewGorth(false, false)

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: DUP_OP, Type: Operator},
		},
	)
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

	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: SWAP_OP, Type: Operator},
		},
	)
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

	if len(g.ExecStack) != 0 {
		t.Errorf("Error in drop operation. Expected stack to be empty, got %v", g.ExecStack)
	}
}
func TestAndOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Both elements are true
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: true, Type: Bool},
			{Value: true, Type: Bool},
			{Value: AND_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing AND operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in AND operation. Expected result true, got %v", val.Value)
	}

	// Test case 2: One element is false
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: false, Type: Bool},
			{Value: true, Type: Bool},
			{Value: AND_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing AND operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in AND operation. Expected result false, got %v", val.Value)
	}

	// Test case 3: Both elements are false
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: false, Type: Bool},
			{Value: false, Type: Bool},
			{Value: AND_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing AND operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in AND operation. Expected result false, got %v", val.Value)
	}

	// Test case 4: One element is not a boolean
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: true, Type: Bool},
			{Value: AND_OP, Type: Operator},
		},
	)
	if err == nil {
		t.Errorf("Expected error performing AND operation, but got nil")
	}
}
func TestOrOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Both elements are true
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: true, Type: Bool},
			{Value: true, Type: Bool},
			{Value: OR_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing OR operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in OR operation. Expected result true, got %v", val.Value)
	}

	// Test case 2: One element is false
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: false, Type: Bool},
			{Value: true, Type: Bool},
			{Value: OR_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing OR operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in OR operation. Expected result true, got %v", val.Value)
	}

	// Test case 3: Both elements are false
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: false, Type: Bool},
			{Value: false, Type: Bool},
			{Value: OR_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing OR operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in OR operation. Expected result false, got %v", val.Value)
	}

	// Test case 4: One element is not a boolean
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: true, Type: Bool},
			{Value: OR_OP, Type: Operator},
		},
	)
	if err == nil {
		t.Errorf("Expected error performing OR operation, but got nil")
	}
}

func TestNotOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: true -> false
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: true, Type: Bool},
			{Value: NOT_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing NOT operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in NOT operation. Expected result false, got %v", val.Value)
	}

	// Test case 2: false -> true
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: false, Type: Bool},
			{Value: NOT_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing NOT operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in NOT operation. Expected result true, got %v", val.Value)
	}

	// Test case 3: non-boolean value
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: NOT_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing NOT operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != -5 {
		t.Errorf("Error in NOT operation. Expected result -5, got %v", val.Value)
	}
}

func TestEqualOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Equal integers
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 5, Type: Int},
			{Value: EQUAL_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing equal operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in equal operation. Expected result true, got %v", val.Value)
	}

	// Test case 2: Unequal integers
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: EQUAL_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing equal operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in equal operation. Expected result false, got %v", val.Value)
	}

	// Test case 3: Equal booleans
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: true, Type: Bool},
			{Value: true, Type: Bool},
			{Value: EQUAL_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing equal operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in equal operation. Expected result true, got %v", val.Value)
	}

	// Test case 4: Unequal booleans
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: true, Type: Bool},
			{Value: false, Type: Bool},
			{Value: EQUAL_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing equal operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in equal operation. Expected result false, got %v", val.Value)
	}

	// Test case 5: Mixed types
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: true, Type: Bool},
			{Value: EQUAL_OP, Type: Operator},
		},
	)

	if err != nil {
		t.Errorf("Error performing equal operation: %v", err)
	}

	val, err = g.Pop()

	if err != nil || val.Value != false {
		t.Errorf("Error in equal operation. Expected result false, got %v", val.Value)
	}
}
func TestNotEqualOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Equal elements
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 5, Type: Int},
			{Value: NOT_EQUAL_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing not equal operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in not equal operation. Expected result false, got %v", val.Value)
	}

	// Test case 2: Not equal elements
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: NOT_EQUAL_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing not equal operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in not equal operation. Expected result true, got %v", val.Value)
	}

	// Test case 3: Non-integer elements
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: "hello", Type: String},
			{Value: "world", Type: String},
			{Value: NOT_EQUAL_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing not equal operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in not equal operation. Expected result true, got %v", val.Value)
	}
}
func TestEqualType(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Equal types
	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: 20, Type: Int})

	err := g.EqualType()
	if err != nil {
		t.Errorf("Error in EqualType: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in EqualType. Expected result true, got %v", val.Value)
	}

	// Test case 2: Different types
	g.Push(StackElement{Value: 10, Type: Int})
	g.Push(StackElement{Value: true, Type: Bool})

	err = g.EqualType()
	if err != nil {
		t.Errorf("Error in EqualType: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in EqualType. Expected result false, got %v", val.Value)
	}

	// Test case 3: Empty stack
	err = g.EqualType()
	if err == nil {
		t.Errorf("Expected error in EqualType, but got nil")
	}
}
func TestTokenize(t *testing.T) {
	// Test case 1: Valid integer token
	input1 := "123"
	expected1 := []StackElement{{Type: Int, Value: 123}}
	tokens1, _, err := Tokenize(input1)
	if err != nil {
		t.Errorf("Error in TestTokenize: %v", err)
	}
	if !reflect.DeepEqual(tokens1, expected1) {
		t.Errorf("TestTokenize failed for input '%s'. Expected %v, got %v", input1, expected1, tokens1)
	}

	// Test case 2: Valid float token
	input2 := "3.14"
	expected2 := []StackElement{{Type: Float, Value: 3.14}}
	tokens2, _, err := Tokenize(input2)
	if err != nil {
		t.Errorf("Error in TestTokenize: %v", err)
	}
	if !reflect.DeepEqual(tokens2, expected2) {
		t.Errorf("TestTokenize failed for input '%s'. Expected %v, got %v", input2, expected2, tokens2)
	}

	// Test case 3: Valid string token
	input3 := `"hello world"`
	expected3 := []StackElement{{Type: String, Value: "hello world"}}
	tokens3, _, err := Tokenize(input3)
	if err != nil {
		t.Errorf("Error in TestTokenize: %v", err)
	}
	if !reflect.DeepEqual(tokens3, expected3) {
		t.Errorf("TestTokenize failed for input '%s'. Expected %v, got %v", input3, expected3, tokens3)
	}

	// Test case 4: Valid boolean token
	input4 := "true"
	expected4 := []StackElement{{Type: Bool, Value: true}}
	tokens4, _, err := Tokenize(input4)
	if err != nil {
		t.Errorf("Error in TestTokenize: %v", err)
	}
	if !reflect.DeepEqual(tokens4, expected4) {
		t.Errorf("TestTokenize failed for input '%s'. Expected %v, got %v", input4, expected4, tokens4)
	}

	// Test case 5: Valid operator token
	input5 := "+"
	expected5 := []StackElement{{Type: Operator, Value: ADD_OP}}
	tokens5, _, err := Tokenize(input5)
	if err != nil {
		t.Errorf("Error in TestTokenize: %v", err)
	}
	if !reflect.DeepEqual(tokens5, expected5) {
		t.Errorf("TestTokenize failed for input '%s'. Expected %v, got %v", input5, expected5, tokens5)
	}

	// Test case 6: Invalid token
	input6 := "abc"
	_, _, err = Tokenize(input6)
	if err == nil {
		t.Errorf("Expected error in TestTokenize for input '%s', but got nil", input6)
	}
}
func TestGreaterThanOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Int values, val2 > val1
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: GT_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThan operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in GreaterThan operation. Expected result false, got %v", val.Value)
	}

	// Test case 2: Int values, val2 < val1
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: GT_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThan operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in GreaterThan operation. Expected result false, got %v", val.Value)
	}

	// Test case 3: Float values, val2 > val1
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 10.5, Type: Float},
			{Value: 5.5, Type: Float},
			{Value: GT_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThan operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in GreaterThan operation. Expected result true, got %v", val.Value)
	}

	// Test case 4: Float values, val2 < val1
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5.5, Type: Float},
			{Value: 10.5, Type: Float},
			{Value: GT_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThan operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in GreaterThan operation. Expected result false, got %v", val.Value)
	}

	// Test case 5: Different types
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 5.5, Type: Float},
			{Value: GT_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThan operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in GreaterThan operation. Expected result false, got %v", val.Value)
	}
}
func TestLessThanOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Int values
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 10, Type: Int},
			{Value: 5, Type: Int},
			{Value: LS_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThan operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in LessThan operation. Expected result false, got %v", val.Value)
	}

	// Test case 2: Float values
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 3.14, Type: Float},
			{Value: 2.71, Type: Float},
			{Value: LS_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThan operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in LessThan operation. Expected result false, got %v", val.Value)
	}

	// Test case 3: Different types
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 3.14, Type: Float},
			{Value: 5, Type: Int},
			{Value: LS_THAN_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThan operation: %v", err)
	}
	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in LessThan operation. Expected result true, got %v", val.Value)
	}
}
func TestGreaterThanEqualOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Integers
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: GT_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThanEqual operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in GreaterThanEqual operation. Expected result false, got %v", val.Value)
	}

	// Test case 2: Floats
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 10.5, Type: Float},
			{Value: 5.5, Type: Float},
			{Value: GT_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThanEqual operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in GreaterThanEqual operation. Expected result true, got %v", val.Value)
	}

	// Test case 3: Mixed types
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 5.5, Type: Float},
			{Value: GT_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing GreaterThanEqual operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in GreaterThanEqual operation. Expected result false, got %v", val.Value)
	}
}
func TestLessThanEqualOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Int values, val2 <= val1
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 10, Type: Int},
			{Value: LS_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThanEqual operation: %v", err)
	}

	val, err := g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in LessThanEqual operation. Expected result true, got %v", val.Value)
	}

	// Test case 2: Int values, val2 > val1
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 10, Type: Int},
			{Value: 5, Type: Int},
			{Value: LS_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThanEqual operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in LessThanEqual operation. Expected result false, got %v", val.Value)
	}

	// Test case 3: Float values, val2 <= val1
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5.5, Type: Float},
			{Value: 10.5, Type: Float},
			{Value: LS_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThanEqual operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in LessThanEqual operation. Expected result true, got %v", val.Value)
	}

	// Test case 4: Float values, val2 > val1
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 10.5, Type: Float},
			{Value: 5.5, Type: Float},
			{Value: LS_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThanEqual operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != false {
		t.Errorf("Error in LessThanEqual operation. Expected result false, got %v", val.Value)
	}

	// Test case 5: Different types
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 5, Type: Int},
			{Value: 5.5, Type: Float},
			{Value: LS_THAN_EQ_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing LessThanEqual operation: %v", err)
	}

	val, err = g.Pop()
	if err != nil || val.Value != true {
		t.Errorf("Error in LessThanEqual operation. Expected result true, got %v", val.Value)
	}
}
func TestRotOperation(t *testing.T) {
	g := NewGorth(false, false)

	// Test case 1: Valid rotation
	err := g.ExecuteProgram(
		[]StackElement{
			{Value: 1, Type: Int}, // a
			{Value: 2, Type: Int}, // b
			{Value: 3, Type: Int}, // c
			{Value: ROT_OP, Type: Operator},
		},
	)
	if err != nil {
		t.Errorf("Error performing rotation operation: %v", err)
	}

	val1, err := g.Pop()
	if err != nil || val1.Value != 1 {
		t.Errorf("Error in rotation operation. Expected result 1, got %d", val1.Value)
	}

	val2, err := g.Pop()
	if err != nil || val2.Value != 3 {
		t.Errorf("Error in rotation operation. Expected result 3, got %d", val2.Value)
	}

	val3, err := g.Pop()
	if err != nil || val3.Value != 2 {
		t.Errorf("Error in rotation operation. Expected result 2, got %d", val3.Value)
	}

	// Test case 2: Insufficient elements on stack
	err = g.ExecuteProgram(
		[]StackElement{
			{Value: 1, Type: Int},
			{Value: ROT_OP, Type: Operator},
		},
	)
	if err == nil {
		t.Errorf("Expected error performing rotation operation, but got nil")
	}
}
