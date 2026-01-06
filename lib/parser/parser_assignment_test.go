package parser

import (
	"gorth/lexer"
	"strings"
	"testing"
)

func TestAssignmentTargets(t *testing.T) {
	input := "1000 := VAR x\nVAR y\n10 := y\n"
	p := NewParser(lexer.NewLexer(strings.NewReader(input)))
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if got := len(program.Statements); got != 5 {
		// expected: INT, ASSIGN, VAR, INT, ASSIGN
		t.Fatalf("expected 5 statements, got %d", got)
	}

	assign1, ok := program.Statements[1].(*Assignment)
	if !ok {
		t.Fatalf("expected stmt[1] to be *Assignment, got %T", program.Statements[1])
	}
	if _, ok := assign1.Target.(*VarDeclaration); !ok {
		t.Fatalf("expected first assignment target to be *VarDeclaration, got %T", assign1.Target)
	}

	assign2, ok := program.Statements[4].(*Assignment)
	if !ok {
		t.Fatalf("expected stmt[4] to be *Assignment, got %T", program.Statements[4])
	}
	if _, ok := assign2.Target.(*Identifier); !ok {
		t.Fatalf("expected second assignment target to be *Identifier, got %T", assign2.Target)
	}
}
