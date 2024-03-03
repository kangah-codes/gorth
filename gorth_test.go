package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type TestCase []struct {
	stack       []StackElement
	variableMap map[string]Variable
	expected    []StackElement
	expectedErr error
	title       string
}

func TestNewGorth(t *testing.T) {
	debugMode := true
	strictMode := false

	g := NewGorth(debugMode, strictMode)

	if len(g.ExecStack) != 0 {
		t.Errorf("Expected empty execution stack, but got %d elements", len(g.ExecStack))
	}

	if g.DebugMode != debugMode {
		t.Errorf("Expected debug mode to be %v, but got %v", debugMode, g.DebugMode)
	}

	if g.StrictMode != strictMode {
		t.Errorf("Expected strict mode to be %v, but got %v", strictMode, g.StrictMode)
	}

	if g.MaxStackSize != MAX_STACK_SIZE {
		t.Errorf("Expected max stack size to be %d, but got %d", MAX_STACK_SIZE, g.MaxStackSize)
	}
}

func TestReadGorthFile(t *testing.T) {
	filename := "testfile.txt" // Replace with the actual filename for testing

	// Create a temporary test file
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(filename)
	defer file.Close()

	// Write test data to the file
	testData := []string{
		"line 1",
		"line 2",
		"# comment",
		"line 3",
	}
	for _, line := range testData {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			t.Fatalf("Failed to write to test file: %v", err)
		}
	}

	// Call the function under test
	lines, err := ReadGorthFile(filename)
	if err != nil {
		t.Fatalf("Failed to read Gorth file: %v", err)
	}

	// Verify the result
	expectedLines := []string{"line 1", "line 2", "line 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, but got %d lines", len(expectedLines), len(lines))
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Expected line %d to be %q, but got %q", i+1, expectedLines[i], line)
		}
	}
}
func TestTokenize(t *testing.T) {
	testCases := []struct {
		input       string
		expected    []StackElement
		expectedMap map[string]Variable
		expectedErr error
	}{
		{
			input: "/x 10 def",
			expected: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			expectedMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
			},
			expectedErr: nil,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		tokens, variables, err := Tokenize(tc.input)

		if err != tc.expectedErr {
			t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
		}

		if !reflect.DeepEqual(tokens, tc.expected) {
			t.Errorf("Expected tokens: %v, but got: %v", tc.expected, tokens)
		}

		if !reflect.DeepEqual(variables, tc.expectedMap) {
			t.Errorf("Expected variables: %v, but got: %v", tc.expectedMap, variables)
		}
	}
}
func TestPush(t *testing.T) {
	g := NewGorth(false, false)
	g.MaxStackSize = 5

	// Test pushing elements onto the stack
	err := g.Push(StackElement{Type: Identifier, Value: "x"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = g.Push(StackElement{Type: Int, Value: 10})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test stack overflow
	for i := 0; i < g.MaxStackSize-2; i++ {
		err = g.Push(StackElement{Type: Int, Value: i})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	err = g.Push(StackElement{Type: Int, Value: 100})
	if err == nil {
		t.Error("Expected stack overflow error, but got nil")
	} else if err.Error() != "ERROR: stack overflow" {
		t.Errorf("Expected error message: %q, but got: %q", "ERROR: stack overflow", err.Error())
	}
}
func TestPop(t *testing.T) {
	g := NewGorth(false, false)

	// Test popping from an empty stack
	_, err := g.Pop()
	if err == nil {
		t.Error("Expected error: cannot pop from an empty stack, but got nil")
	} else if err.Error() != "ERROR: cannot pop from an empty stack" {
		t.Errorf("Expected error message: %q, but got: %q", "ERROR: cannot pop from an empty stack", err.Error())
	}

	// Test popping from a non-empty stack
	g.ExecStack = append(g.ExecStack, StackElement{Type: Identifier, Value: "x"})
	g.ExecStack = append(g.ExecStack, StackElement{Type: Int, Value: 10})

	val, err := g.Pop()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedVal := StackElement{Type: Int, Value: 10}
	if val != expectedVal {
		t.Errorf("Expected popped value: %v, but got: %v", expectedVal, val)
	}

	if len(g.ExecStack) != 1 {
		t.Errorf("Expected stack length to be 1, but got: %d", len(g.ExecStack))
	}
}
func TestDrop(t *testing.T) {
	g := NewGorth(false, false)

	// Test dropping from an empty stack
	err := g.Drop()
	if err == nil {
		t.Error("Expected error: cannot drop from an empty stack, but got nil")
	} else if err.Error() != "ERROR: cannot drop from an empty stack" {
		t.Errorf("Expected error message: %q, but got: %q", "ERROR: cannot drop from an empty stack", err.Error())
	}

	// Test dropping from a non-empty stack
	g.ExecStack = append(g.ExecStack, StackElement{Type: Identifier, Value: "x"})
	g.ExecStack = append(g.ExecStack, StackElement{Type: Int, Value: 10})

	err = g.Drop()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(g.ExecStack) != 1 {
		t.Errorf("Expected stack length to be 1, but got: %d", len(g.ExecStack))
	}
}

func TestDump(t *testing.T) {
	g := NewGorth(false, false)

	// Test dumping an integer value
	g.ExecStack = append(g.ExecStack, StackElement{Type: Int, Value: 10})

	oldStdout := os.Stdout // Keep a reference to the original stdout
	r, w, _ := os.Pipe()   // Create a pipe to capture stdout
	os.Stdout = w          // Replace stdout with the write end of the pipe

	// Capture the output by reading from the read end of the pipe
	capturedOutput := make(chan string)
	go func() {
		out, _ := ioutil.ReadAll(r)
		capturedOutput <- string(out)
	}()

	err := g.Dump()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	w.Close()             // Close the write end of the pipe to unblock the reader
	os.Stdout = oldStdout // Restore the original stdout

	// Check the captured output
	expectedOutput := "10\n"
	actualOutput := <-capturedOutput
	if actualOutput != expectedOutput {
		t.Errorf("Expected output: %q, but got: %q", expectedOutput, actualOutput)
	}

	// Similar procedure for other test cases
}
func TestRot(t *testing.T) {
	g := NewGorth(false, false)

	// Test with less than 3 elements on stack
	err := g.Rot()
	expectedErr := "ERROR: at least 3 elements need to be on stack to perform ROT_OP"
	if err == nil {
		t.Error("Expected error: ", expectedErr)
	} else if err.Error() != expectedErr {
		t.Errorf("Expected error: %q, but got: %q", expectedErr, err.Error())
	}

	// Test with 3 elements on stack
	g.ExecStack = append(g.ExecStack, StackElement{Type: Int, Value: 1})
	g.ExecStack = append(g.ExecStack, StackElement{Type: Int, Value: 2})
	g.ExecStack = append(g.ExecStack, StackElement{Type: Int, Value: 3})

	err = g.Rot()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedStack := []StackElement{
		{Type: Int, Value: 2},
		{Type: Int, Value: 3},
		{Type: Int, Value: 1},
	}
	if !reflect.DeepEqual(g.ExecStack, expectedStack) {
		t.Errorf("Expected stack: %v, but got: %v", expectedStack, g.ExecStack)
	}
}
func TestPeek(t *testing.T) {
	g := NewGorth(false, false)

	// Test peeking from an empty stack
	_, err := g.Peek()
	expectedErr := "ERROR: cannot PEEK_OP at an empty stack"
	if err == nil {
		t.Error("Expected error: ", expectedErr)
	} else if err.Error() != expectedErr {
		t.Errorf("Expected error: %q, but got: %q", expectedErr, err.Error())
	}

	// Test peeking from a non-empty stack
	g.ExecStack = append(g.ExecStack, StackElement{Type: Identifier, Value: "x"})
	g.ExecStack = append(g.ExecStack, StackElement{Type: Int, Value: 10})

	val, err := g.Peek()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedVal := StackElement{Type: Int, Value: 10}
	if val != expectedVal {
		t.Errorf("Expected peeked value: %v, but got: %v", expectedVal, val)
	}

	if len(g.ExecStack) != 2 {
		t.Errorf("Expected stack length to be 2, but got: %d", len(g.ExecStack))
	}
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		stack       []StackElement
		variableMap map[string]Variable
		expected    []StackElement
		expectedErr error
		title       string
	}{
		// Test integer addition
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 15},
			},
			expectedErr: nil,
			title:       "Test integer addition",
		},
		// Test string concatenation
		{
			stack: []StackElement{
				{Type: String, Value: "Hello"},
				{Type: String, Value: "World"},
			},
			expected: []StackElement{
				{Type: String, Value: "WorldHello"},
			},
			expectedErr: nil,
			title:       "Test string concatenation",
		},
		// Test float addition
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
				{Type: Float, Value: 2.71},
			},
			expected: []StackElement{
				{Type: Float, Value: 5.85},
			},
			expectedErr: nil,
			title:       "Test float addition",
		},
		// Test mixed type addition (int and float)
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 7.5},
			},
			expectedErr: nil,
			title:       "Test mixed type addition (int and float)",
		},
		// Test mixed type addition (float and int)
		{
			stack: []StackElement{
				{Type: Float, Value: 2.5},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Float, Value: 7.5},
			},
			expectedErr: nil,
			title:       "Test mixed type addition (float and int)",
		},
		// Test addition with both variables
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
				"y": {Name: "y", Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 15},
			},
			expectedErr: nil,
			title:       "Test addition with both variables",
		},
		// Test addition with variable and integer
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Int, Value: 5},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 15},
			},
			expectedErr: nil,
			title:       "Test addition with variable and integer",
		},
		// Test addition with integer and variable
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 15},
			},
			expectedErr: nil,
			title:       "Test addition with integer and variable",
		},
		// Test addition with variables of different types
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
				"y": {Name: "y", Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 7.5},
			},
			expectedErr: errors.New("ERROR: cannot perform ADD_OP on different types"),
			title:       "Test addition with variables of different types",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Add()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}

func TestSub(t *testing.T) {
	testCases := []struct {
		stack       []StackElement
		variableMap map[string]Variable
		expected    []StackElement
		expectedErr error
		title       string
	}{
		// Test integer subtraction
		{
			stack: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 5},
			},
			expectedErr: nil,
			title:       "Test integer subtraction",
		},
		// Test float subtraction
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
				{Type: Float, Value: 2.0},
			},
			expected: []StackElement{
				{Type: Float, Value: 1.1400000000000001},
			},
			expectedErr: nil,
			title:       "Test float subtraction",
		},
		// Test mixed number subtraction (int and float)
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 2.5},
			},
			expectedErr: nil,
			title:       "Test mixed number subtraction (int and float)",
		},
		// Test mixed number subtraction (float and int)
		{
			stack: []StackElement{
				{Type: Float, Value: 2.5},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Float, Value: -2.5},
			},
			expectedErr: nil,
			title:       "Test mixed number subtraction (float and int)",
		},
		// Test variable subtraction
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
				"y": {Name: "y", Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 5},
			},
			expectedErr: nil,
			title:       "Test variable subtraction",
		},
		// Test variable subtraction with undeclared variable
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
			},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable y has not been declared"),
			title:       "Test variable subtraction with undeclared variable",
		},
		// Test variable subtraction with different types
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
				"y": {Name: "y", Type: Float, Value: 2.5},
			},
			expected:    []StackElement{{Type: Float, Value: 7.5}},
			expectedErr: nil,
			title:       "Test variable subtraction with different types",
		},
		// Test variable subtraction with one variable and one non-variable
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Int, Value: 5},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 5},
			},
			expectedErr: nil,
			title:       "Test variable subtraction with one variable and one non-variable",
		},
		// Test variable subtraction with undeclared variable and non-variable
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Int, Value: 5},
			},
			variableMap: map[string]Variable{
				"y": {Name: "y", Type: Int, Value: 10},
			},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable x has not been declared"),
			title:       "Test variable subtraction with undeclared variable and non-variable",
		},
		// Test variable subtraction with different types (variable and non-variable)
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Float, Value: 2.5},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
			},
			expected:    []StackElement{{Type: Float, Value: 7.5}},
			expectedErr: nil,
			title:       "Test variable subtraction with different types (variable and non-variable)",
		},
		// Test variable subtraction with different types (non-variable and variable)
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Float, Value: 2.5},
			},
			expected:    []StackElement{{Type: Float, Value: 2.5}},
			expectedErr: errors.New("ERROR: cannot perform SUB_OP on different types"),
			title:       "Test variable subtraction with different types (non-variable and variable)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Sub()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}
func TestMul(t *testing.T) {
	testCases := []struct {
		stack       []StackElement
		variableMap map[string]Variable
		expected    []StackElement
		expectedErr error
		title       string
	}{
		// Test integer multiplication
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 50},
			},
			expectedErr: nil,
			title:       "Test integer multiplication",
		},
		// Test float multiplication
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
				{Type: Float, Value: 2.0},
			},
			expected: []StackElement{
				{Type: Float, Value: 6.28},
			},
			expectedErr: nil,
			title:       "Test float multiplication",
		},
		// Test mixed number multiplication (int and float)
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 12.5},
			},
			expectedErr: nil,
			title:       "Test mixed number multiplication (int and float)",
		},
		// Test mixed number multiplication (float and int)
		{
			stack: []StackElement{
				{Type: Float, Value: 2.5},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Float, Value: 12.5},
			},
			expectedErr: nil,
			title:       "Test mixed number multiplication (float and int)",
		},
		// Test variable multiplication
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
				"y": {Name: "y", Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 50},
			},
			expectedErr: nil,
			title:       "Test variable multiplication",
		},
		// Test variable multiplication with undeclared variable
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable y has not been declared"),
			title:       "Test variable multiplication with undeclared variable",
		},
		// Test variable multiplication with different types
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
				"y": {Name: "y", Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 12.5},
			},
			expectedErr: nil,
			title:       "Test variable multiplication with different types",
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Mul()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	testCases := []struct {
		stack       []StackElement
		variableMap map[string]Variable
		expected    []StackElement
		expectedErr error
		title       string
	}{
		// Test integer division
		{
			stack: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 2},
			},
			expectedErr: nil,
			title:       "Test integer division",
		},
		// Test float division
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
				{Type: Float, Value: 2.0},
			},
			expected: []StackElement{
				{Type: Float, Value: 1.57},
			},
			expectedErr: nil,
			title:       "Test float division",
		},
		// Test mixed number division (int and float)
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 2.0},
			},
			expectedErr: nil,
			title:       "Test mixed number division (int and float)",
		},
		// Test mixed number division (float and int)
		{
			stack: []StackElement{
				{Type: Float, Value: 2.5},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Float, Value: 0.5},
			},
			expectedErr: nil,
			title:       "Test mixed number division (float and int)",
		},
		// Test variable division
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
				"y": {Name: "y", Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 2},
			},
			expectedErr: nil,
			title:       "Test variable division",
		},
		// Test variable division with undeclared variable
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable y has not been declared"),
			title:       "Test variable division with undeclared variable",
		},
		// Test variable division with different types
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
				"y": {Name: "y", Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 2.0},
			},
			expectedErr: nil,
			title:       "Test variable division with different types",
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Div()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}
func TestMod(t *testing.T) {
	var testCases = []struct {
		stack       []StackElement
		expected    []StackElement
		expectedErr error
		title       string
	}{
		// Test integer modulo
		{
			stack: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 0},
			},
			expectedErr: nil,
			title:       "Test integer modulo",
		},
		// Test float modulo
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
				{Type: Float, Value: 2.0},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot perform MOD_OP on different types"),
			title:       "Test float modulo",
		},
		// Test mixed number modulo (int and float)
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Float, Value: 2.5},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot perform MOD_OP on different types"),
			title:       "Test mixed number modulo (int and float)",
		},
		// Test mixed number modulo (float and int)
		{
			stack: []StackElement{
				{Type: Float, Value: 2.5},
				{Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot perform MOD_OP on different types"),
			title:       "Test mixed number modulo (float and int)",
		},
		// Add more test cases as needed
		{
			stack: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 3},
			},
			expected: []StackElement{
				{Type: Int, Value: 1},
			},
			expectedErr: nil,
			title:       "Test integer modulo",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 0},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot divide by zero"),
			title:       "Test integer modulo with divisor 0",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 5},
				{Type: Int, Value: 3},
			},
			expected: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 2},
			},
			expectedErr: nil,
			title:       "Test integer modulo with more than 2 elements on stack",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack

			err := g.Mod()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}

func TestExp(t *testing.T) {
	var testCases = []struct {
		stack       []StackElement
		expected    []StackElement
		expectedErr error
		title       string
	}{
		// Test integer exponentiation
		{
			stack: []StackElement{
				{Type: Int, Value: 2},
				{Type: Int, Value: 3},
			},
			expected: []StackElement{
				{Type: Int, Value: 8},
			},
			expectedErr: nil,
			title:       "Test integer exponentiation",
		},
		// Test float exponentiation
		{
			stack: []StackElement{
				{Type: Float, Value: 2.0},
				{Type: Float, Value: 3.0},
			},
			expected: []StackElement{
				{Type: Float, Value: 8.0},
			},
			expectedErr: nil,
			title:       "Test float exponentiation",
		},
		// Test mixed number exponentiation (int and float)
		{
			stack: []StackElement{
				{Type: Int, Value: 2},
				{Type: Float, Value: 3.0},
			},
			expected: []StackElement{
				{Type: Float, Value: 8.0},
			},
			expectedErr: nil,
			title:       "Test mixed number exponentiation (int and float)",
		},
		// Test mixed number exponentiation (float and int)
		{
			stack: []StackElement{
				{Type: Float, Value: 2.0},
				{Type: Int, Value: 3},
			},
			expected: []StackElement{
				{Type: Float, Value: 8.0},
			},
			expectedErr: nil,
			title:       "Test mixed number exponentiation (float and int)",
		},
		// Add more test cases as needed
		{
			stack: []StackElement{
				{Type: Int, Value: 2},
				{Type: Int, Value: 3},
				{Type: Int, Value: 4},
			},
			expected: []StackElement{
				{Type: Int, Value: 2},
				{Type: Int, Value: 81},
			},
			expectedErr: nil,
			title:       "Test integer exponentiation with more than 2 elements on stack",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 2},
				{Type: Int, Value: 0},
			},
			expected: []StackElement{
				{Type: Int, Value: 1},
			},
			expectedErr: nil,
			title:       "Test integer exponentiation with exponent 0",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 2},
				{Type: Int, Value: -3},
			},
			expected: []StackElement{
				{Type: Int, Value: 0},
			},
			expectedErr: nil,
			title:       "Test integer exponentiation with negative exponent",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 2},
				{Type: Int, Value: 3},
				{Type: Int, Value: 4},
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 2},
				{Type: Int, Value: 3},
				{Type: Int, Value: 1024},
			},
			expectedErr: nil,
			title:       "Test integer exponentiation with more than 3 elements on stack",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack

			err := g.Exp()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}

func TestInc(t *testing.T) {
	var testCases = TestCase{
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 6},
			},
			expectedErr: nil,
			title:       "Test incrementing an integer",
		},
		// test on a string
		{
			stack: []StackElement{
				{Type: String, Value: "Hello"},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot perform INC_OP on different types"),
			title:       "Test incrementing a string",
		},
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
			},
			expected: []StackElement{
				{Type: Float, Value: 4.140000000000001},
			},
			expectedErr: nil,
			title:       "Test incrementing a float",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: nil,
			title:       "Test incrementing a variable",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable x has not been declared"),
			title:       "Test incrementing a variable of non-integer type",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"y": {Name: "y", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable x has not been declared"),
			title:       "Test incrementing an undeclared variable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Inc()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}

func TestDec(t *testing.T) {
	var testCases = TestCase{
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: 4},
			},
			expectedErr: nil,
			title:       "Test decrementing an integer",
		},
		// test on a string
		{
			stack: []StackElement{
				{Type: String, Value: "Hello"},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot perform DEC_OP on different types"),
			title:       "Test decrementing a string",
		},
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
			},
			expected: []StackElement{
				{Type: Float, Value: 2.14},
			},
			expectedErr: nil,
			title:       "Test decrementing a float",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: nil,
			title:       "Test decrementing a variable",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable x has not been declared"),
			title:       "Test decrementing a variable of non-integer type",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"y": {Name: "y", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable x has not been declared"),
			title:       "Test decrementing an undeclared variable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Dec()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}

func TestNeg(t *testing.T) {
	var testCases = TestCase{
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Int, Value: -5},
			},
			expectedErr: nil,
			title:       "Test negating an integer",
		},
		// test on a string
		{
			stack: []StackElement{
				{Type: String, Value: "Hello"},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot perform NEG_OP on different types"),
			title:       "Test negating a string",
		},
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
			},
			expected: []StackElement{
				{Type: Float, Value: -3.14},
			},
			expectedErr: nil,
			title:       "Test negating a float",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: nil,
			title:       "Test negating a variable",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable x has not been declared"),
			title:       "Test negating a variable of non-integer type",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
			},
			variableMap: map[string]Variable{
				"y": {Name: "y", Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: fmt.Errorf("ERROR: variable x has not been declared"),
			title:       "Test negating an undeclared variable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Neg()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	var testCases = TestCase{
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 10},
				{Type: Int, Value: 5},
			},
			expectedErr: nil,
			title:       "Test swapping two integers",
		},
		{
			stack: []StackElement{
				{Type: Float, Value: 3.14},
				{Type: Float, Value: 2.0},
			},
			expected: []StackElement{
				{Type: Float, Value: 2.0},
				{Type: Float, Value: 3.14},
			},
			expectedErr: nil,
			title:       "Test swapping two floats",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Float, Value: 2.5},
			},
			expected: []StackElement{
				{Type: Float, Value: 2.5},
				{Type: Int, Value: 5},
			},
			expectedErr: nil,
			title:       "Test swapping an integer and a float",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
				"y": {Name: "y", Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Identifier, Value: "y"},
				{Type: Identifier, Value: "x"},
			},
			expectedErr: nil,
			title:       "Test swapping two variables",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Identifier, Value: "y"},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
			},
			// the swap should still happen, but I know the error would throw when you try to use the undeclared var
			expected: []StackElement{
				{Type: Identifier, Value: "y"},
				{Type: Identifier, Value: "x"},
			},
			expectedErr: nil,
			title:       "Test swapping a variable and an undeclared variable",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Int, Value: 5},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 10},
			},
			expected: []StackElement{
				{Type: Int, Value: 5},
				{Type: Identifier, Value: "x"},
			},
			expectedErr: nil,
			title:       "Test swapping an undeclared variable and a variable",
		},
		{
			stack: []StackElement{
				{Type: Identifier, Value: "x"},
				{Type: Float, Value: 2.5},
			},
			variableMap: map[string]Variable{
				"x": {Name: "x", Type: Int, Value: 5},
			},
			expected: []StackElement{
				{Type: Float, Value: 2.5},
				{Type: Identifier, Value: "x"},
			},
			expectedErr: nil,
			title:       "Test swapping a variable and a non-variable",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
			},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot pop from an empty stack"),
			title:       "Test swapping with only one element on stack",
		},
		{
			stack:       []StackElement{},
			expected:    []StackElement{},
			expectedErr: errors.New("ERROR: cannot pop from an empty stack"),
			title:       "Test swapping with no elements on stack",
		},
		{
			stack: []StackElement{
				{Type: Int, Value: 5},
				{Type: Int, Value: 10},
				{Type: Int, Value: 15},
			},
			expected: []StackElement{
				{Type: Int, Value: 5},
				{Type: Int, Value: 15},
				{Type: Int, Value: 10},
			},
			expectedErr: nil,
			title:       "Test swapping with more than 2 elements on stack",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			g := NewGorth(false, false)
			g.ExecStack = tc.stack
			g.VariableMap = tc.variableMap

			err := g.Swap()
			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %q, but got: %q", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(g.ExecStack, tc.expected) {
				t.Errorf("Expected stack: %v, but got: %v", tc.expected, g.ExecStack)
			}
		})
	}
}
