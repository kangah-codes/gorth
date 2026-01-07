package runtime

import (
	"testing"
)

func TestStack_Push(t *testing.T) {
	tests := []struct {
		name        string
		stack       *Stack
		value       Value
		wantErr     bool
		expectedLen int
	}{
		{
			name:        "push to empty stack",
			stack:       &Stack{items: []Value{}, max: 0},
			value:       Value{Type: TYPE_INT, Data: 42},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:        "push to stack with items",
			stack:       &Stack{items: []Value{{Type: TYPE_INT, Data: 1}}, max: 0},
			value:       Value{Type: TYPE_STR, Data: "hello"},
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:        "push with no max limit",
			stack:       &Stack{items: []Value{}, max: 0},
			value:       Value{Type: TYPE_BOOL, Data: true},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:        "push within max limit",
			stack:       &Stack{items: []Value{{Type: TYPE_INT, Data: 1}}, max: 3},
			value:       Value{Type: TYPE_FLOAT, Data: 3.14},
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:        "push at max limit",
			stack:       &Stack{items: []Value{{Type: TYPE_INT, Data: 1}, {Type: TYPE_INT, Data: 2}}, max: 2},
			value:       Value{Type: TYPE_INT, Data: 3},
			wantErr:     true,
			expectedLen: 2,
		},
		{
			name:        "push exceeding max limit",
			stack:       &Stack{items: []Value{{Type: TYPE_INT, Data: 1}, {Type: TYPE_INT, Data: 2}, {Type: TYPE_INT, Data: 3}}, max: 2},
			value:       Value{Type: TYPE_INT, Data: 4},
			wantErr:     true,
			expectedLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.stack.Push(tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("Stack.Push() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.stack.items) != tt.expectedLen {
				t.Errorf("Stack.Push() resulted in length = %v, expected %v", len(tt.stack.items), tt.expectedLen)
			}

			if !tt.wantErr && len(tt.stack.items) > 0 {
				lastItem := tt.stack.items[len(tt.stack.items)-1]
				if lastItem.Type != tt.value.Type || lastItem.Data != tt.value.Data {
					t.Errorf("Stack.Push() last item = %+v, expected %+v", lastItem, tt.value)
				}
			}
		})
	}
}
func TestStack_Pop(t *testing.T) {
	tests := []struct {
		name        string
		stack       *Stack
		wantValue   Value
		wantErr     bool
		expectedLen int
	}{
		{
			name:        "pop from empty stack",
			stack:       &Stack{items: []Value{}, max: 0},
			wantValue:   Value{},
			wantErr:     true,
			expectedLen: 0,
		},
		{
			name:        "pop from stack with one item",
			stack:       &Stack{items: []Value{{Type: TYPE_INT, Data: 42}}, max: 0},
			wantValue:   Value{Type: TYPE_INT, Data: 42},
			wantErr:     false,
			expectedLen: 0,
		},
		{
			name:        "pop from stack with multiple items",
			stack:       &Stack{items: []Value{{Type: TYPE_INT, Data: 1}, {Type: TYPE_STR, Data: "hello"}}, max: 0},
			wantValue:   Value{Type: TYPE_STR, Data: "hello"},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:        "pop bool value",
			stack:       &Stack{items: []Value{{Type: TYPE_BOOL, Data: true}}, max: 0},
			wantValue:   Value{Type: TYPE_BOOL, Data: true},
			wantErr:     false,
			expectedLen: 0,
		},
		{
			name:        "pop float value",
			stack:       &Stack{items: []Value{{Type: TYPE_FLOAT, Data: 3.14}}, max: 0},
			wantValue:   Value{Type: TYPE_FLOAT, Data: 3.14},
			wantErr:     false,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.stack.Pop()

			if (err != nil) != tt.wantErr {
				t.Errorf("Stack.Pop() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.stack.items) != tt.expectedLen {
				t.Errorf("Stack.Pop() resulted in length = %v, expected %v", len(tt.stack.items), tt.expectedLen)
			}

			if !tt.wantErr {
				if value.Type != tt.wantValue.Type || value.Data != tt.wantValue.Data {
					t.Errorf("Stack.Pop() value = %+v, expected %+v", value, tt.wantValue)
				}
			}
		})
	}
}
func TestStack_Peek(t *testing.T) {
	tests := []struct {
		name      string
		stack     *Stack
		wantValue Value
		wantErr   bool
	}{
		{
			name:      "peek from empty stack",
			stack:     &Stack{items: []Value{}, max: 0},
			wantValue: Value{},
			wantErr:   true,
		},
		{
			name:      "peek from stack with one item",
			stack:     &Stack{items: []Value{{Type: TYPE_INT, Data: 42}}, max: 0},
			wantValue: Value{Type: TYPE_INT, Data: 42},
			wantErr:   false,
		},
		{
			name:      "peek from stack with multiple items",
			stack:     &Stack{items: []Value{{Type: TYPE_INT, Data: 1}, {Type: TYPE_STR, Data: "hello"}}, max: 0},
			wantValue: Value{Type: TYPE_STR, Data: "hello"},
			wantErr:   false,
		},
		{
			name:      "peek bool value",
			stack:     &Stack{items: []Value{{Type: TYPE_BOOL, Data: true}}, max: 0},
			wantValue: Value{Type: TYPE_BOOL, Data: true},
			wantErr:   false,
		},
		{
			name:      "peek float value",
			stack:     &Stack{items: []Value{{Type: TYPE_FLOAT, Data: 3.14}}, max: 0},
			wantValue: Value{Type: TYPE_FLOAT, Data: 3.14},
			wantErr:   false,
		},
		{
			name:      "peek null value",
			stack:     &Stack{items: []Value{{Type: TYPE_NULL, Data: nil}}, max: 0},
			wantValue: Value{Type: TYPE_NULL, Data: nil},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalLen := len(tt.stack.items)
			value, err := tt.stack.Peek()

			if (err != nil) != tt.wantErr {
				t.Errorf("Stack.Peek() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.stack.items) != originalLen {
				t.Errorf("Stack.Peek() modified stack length from %v to %v", originalLen, len(tt.stack.items))
			}

			if !tt.wantErr {
				if value.Type != tt.wantValue.Type || value.Data != tt.wantValue.Data {
					t.Errorf("Stack.Peek() value = %+v, expected %+v", value, tt.wantValue)
				}
			}
		})
	}
}
func TestStack_Size(t *testing.T) {
	tests := []struct {
		name     string
		stack    *Stack
		expected int
	}{
		{
			name:     "empty stack",
			stack:    &Stack{items: []Value{}, max: 0},
			expected: 0,
		},
		{
			name:     "stack with one item",
			stack:    &Stack{items: []Value{{Type: TYPE_INT, Data: 42}}, max: 0},
			expected: 1,
		},
		{
			name:     "stack with multiple items",
			stack:    &Stack{items: []Value{{Type: TYPE_INT, Data: 1}, {Type: TYPE_STR, Data: "hello"}, {Type: TYPE_BOOL, Data: true}}, max: 0},
			expected: 3,
		},
		{
			name:     "stack with max limit set",
			stack:    &Stack{items: []Value{{Type: TYPE_FLOAT, Data: 3.14}, {Type: TYPE_NULL, Data: nil}}, max: 5},
			expected: 2,
		},
		{
			name:     "stack at max capacity",
			stack:    &Stack{items: []Value{{Type: TYPE_INT, Data: 1}, {Type: TYPE_INT, Data: 2}}, max: 2},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := tt.stack.Size()

			if size != tt.expected {
				t.Errorf("Stack.Size() = %v, expected %v", size, tt.expected)
			}
		})
	}
}
