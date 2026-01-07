package parser

import (
	"gorth/lexer"
	"strings"
	"testing"
)

func TestParse_ContinueOutsideLoop_Errors(t *testing.T) {
	src := "CONTINUE"
	l := lexer.NewLexer(strings.NewReader(src))
	p := NewParser(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestParse_BreakOutsideLoop_Errors(t *testing.T) {
	src := "BREAK"
	l := lexer.NewLexer(strings.NewReader(src))
	p := NewParser(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
