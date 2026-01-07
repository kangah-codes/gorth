package runtime

import (
	"fmt"
)

type Value struct {
	Type ValueType
	Data any
}

type ValueType int

const (
	TYPE_INT ValueType = iota
	TYPE_FLOAT
	TYPE_STR
	TYPE_BOOL
	TYPE_NULL
	TYPE_ARRAY
	TYPE_WORD
)

type Stack struct {
	items []Value
	max   int
}

func (s *Stack) Push(v Value) error {
	if s.max > 0 && len(s.items) >= s.max {
		return fmt.Errorf("stack overflow")
	}

	s.items = append(s.items, v)
	return nil
}

func (s *Stack) Pop() (Value, error) {
	if len(s.items) == 0 {
		return Value{}, fmt.Errorf("stack underflow")
	}

	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, nil
}

func (s *Stack) Peek() (Value, error) {
	if len(s.items) == 0 {
		return Value{}, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack) Size() int {
	return len(s.items)
}
