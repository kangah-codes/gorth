package parser

import (
	"gorth/lexer"
	"strings"
	"testing"
)

func TestParse_DoWhileEnd_AllowsContinueInBody(t *testing.T) {
	src := `
DO
  0 DUMPLN
  CONTINUE
WHILE
  TRUE
END
`

	l := lexer.NewLexer(strings.NewReader(src))
	p := NewParser(l)
	_, err := p.Parse()
	if err != nil {
		t.Fatalf("expected parse ok, got: %v", err)
	}

}
