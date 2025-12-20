package parser

import "gorth/lexer"

type Node interface {
	node()
	String() string
}

type Program struct {
	Statements []Node
}

func (p *Program) node()          {}
func (p *Program) String() string { return "Program" }

type IntLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *IntLiteral) node()          {}
func (i *IntLiteral) String() string { return i.Value }

type FloatLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *FloatLiteral) node()          {}
func (i *FloatLiteral) String() string { return i.Value }

type StringLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *StringLiteral) node()          {}
func (i *StringLiteral) String() string { return i.Value }

type BoolLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *BoolLiteral) node()          {}
func (i *BoolLiteral) String() string { return i.Value }

type Identifier struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *Identifier) node()          {}
func (i *Identifier) String() string { return i.Value }

type BinaryExpression struct {
	Left     Node
	Operator lexer.TokenType
	Right    Node
	Position lexer.Position
}

func (b *BinaryExpression) node()          {}
func (b *BinaryExpression) String() string { return "BinaryExpression" }

type UnaryExpression struct {
	Operator lexer.TokenType
	Operand  Node
	Position lexer.Position
}

func (u *UnaryExpression) node()          {}
func (u *UnaryExpression) String() string { return "UnaryExpression" }

// Stack operations
type StackOp struct {
	Operation lexer.TokenType
	Pos       lexer.Position
}

func (s *StackOp) node()          {}
func (s *StackOp) String() string { return "StackOp" }

// Variable assignment
type Assignment struct {
	Name    string
	Value   Node
	IsConst bool
	Pos     lexer.Position
}

func (a *Assignment) node()          {}
func (a *Assignment) String() string { return "Assignment" }

// Print/Println/Dump
type PrintStmt struct {
	Kind  lexer.TokenType // PRINT_OP, PRINTLN_OP, or DUMP_OP
	Value Node
	Pos   lexer.Position
}

func (p *PrintStmt) node()          {}
func (p *PrintStmt) String() string { return "PrintStmt" }

// Pointer operations
type PointerExpr struct {
	Operation lexer.TokenType // PTR or DEREF_OP
	Value     Node
	Pos       lexer.Position
}

func (p *PointerExpr) node()          {}
func (p *PointerExpr) String() string { return "PointerExpr" }

// Array literal
type ArrayLiteral struct {
	Elements []Node
	Pos      lexer.Position
}

func (a *ArrayLiteral) node()          {}
func (a *ArrayLiteral) String() string { return "ArrayLiteral" }
