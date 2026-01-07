package parser

import (
	"fmt"
	"gorth/lexer"
	"reflect"
	"strings"
)

type Node interface {
	node()
	String() string
}

// Pretty prints the AST tree structure with indentation
func PrintAST(node Node, indent int) string {
	prefix := strings.Repeat("  ", indent)
	var result strings.Builder

	switch n := node.(type) {
	case *Program:
		result.WriteString(prefix + "Program:\n")
		for _, stmt := range n.Statements {
			result.WriteString(PrintAST(stmt, indent+1))
		}

	case *IntLiteral:
		fmt.Fprintf(&result, "%s%s (value: %s)\n", prefix, n.Token, n.Value)

	case *FloatLiteral:
		fmt.Fprintf(&result, "%s%s (value: %s)\n", prefix, n.Token, n.Value)

	case *StringLiteral:
		fmt.Fprintf(&result, "%s%s (value: %q)\n", prefix, n.Token, n.Value)

	case *BoolLiteral:
		fmt.Fprintf(&result, "%s%s (value: %s)\n", prefix, n.Token, n.Value)

	case *NullLiteral:
		fmt.Fprintf(&result, "%s%s\n", prefix, n.Token)

	case *Identifier:
		fmt.Fprintf(&result, "%s%s (name: %s)\n", prefix, n.Token, n.Value)

	case *BinaryExpression:
		fmt.Fprintf(&result, "%sBinaryOp (%s):\n", prefix, lexer.TokenMap[n.Operator])
		if n.Left != nil {
			result.WriteString(prefix + "  Left:\n")
			result.WriteString(PrintAST(n.Left, indent+2))
		} else {
			result.WriteString(prefix + "  Left: <from stack>\n")
		}
		if n.Right != nil {
			result.WriteString(prefix + "  Right:\n")
			result.WriteString(PrintAST(n.Right, indent+2))
		} else {
			result.WriteString(prefix + "  Right: <from stack>\n")
		}

	case *UnaryExpression:
		fmt.Fprintf(&result, "%sUnaryOp (%s):\n", prefix, lexer.TokenMap[n.Operator])
		if n.Operand != nil {
			result.WriteString(prefix + "  Operand:\n")
			result.WriteString(PrintAST(n.Operand, indent+2))
		} else {
			result.WriteString(prefix + "  Operand: <from stack>\n")
		}

	case *StackOp:
		fmt.Fprintf(&result, "%sStackOp (%s)\n", prefix, lexer.TokenMap[n.Operation])

	case *Assignment:
		fmt.Fprintf(&result, "%sAssignment:\n", prefix)
		result.WriteString(prefix + "  Target:\n")
		if n.Target != nil {
			result.WriteString(PrintAST(n.Target, indent+2))
		} else {
			result.WriteString(prefix + "    <missing target>\n")
		}

	case *IOStmt:
		fmt.Fprintf(&result, "%sPrint (%s):\n", prefix, lexer.TokenMap[n.Kind])
		if n.Value != nil {
			result.WriteString(prefix + "  Value:\n")
			result.WriteString(PrintAST(n.Value, indent+2))
		}

	case *ArrayLiteral:
		fmt.Fprintf(&result, "%sArray (%d elements):\n", prefix, n.Len())
		for i, elem := range n.Elements {
			fmt.Fprintf(&result, "%s  [%d]:\n", prefix, i)
			result.WriteString(PrintAST(elem, indent+2))
		}

	case *VarDeclaration:
		fmt.Fprintf(&result, "%sVarDecl (name: %s)\n", prefix, n.Name)

	case *ConstDeclaration:
		if n.Value == nil {
			fmt.Fprintf(&result, "%sConstDeclTarget (name: %s)\n", prefix, n.Name)
		} else {
			fmt.Fprintf(&result, "%sConstDecl (name: %s):\n", prefix, n.Name)
			result.WriteString(prefix + "  Value:\n")
			result.WriteString(PrintAST(n.Value, indent+2))
		}

	case *BreakStmt:
		fmt.Fprintf(&result, "%sBreak\n", prefix)

	case *ContinueStmt:
		fmt.Fprintf(&result, "%sContinue\n", prefix)

	case *Procedure:
		fmt.Fprintf(&result, "%sProcedure (name: %s)\n", prefix, n.Name)
		result.WriteString(prefix + "  Params:\n")
		for _, p := range n.Parameters {
			result.WriteString(prefix + "  Param: " + p.Name + "\n")
		}

	default:
		fmt.Fprintf(&result, "%sUnknown node type: %T\n", prefix, n)
	}

	return result.String()
}

type Program struct {
	Statements []Node
}

func (p *Program) node() {}
func (p *Program) String() string {
	var result strings.Builder
	for _, stmt := range p.Statements {
		result.WriteString(stmt.String() + "\n")
	}
	return result.String()
}

type IntLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *IntLiteral) node() {}
func (i *IntLiteral) String() string {
	return fmt.Sprintf("%-15s %-15s %s", i.Token, i.Value, i.Position)
}

type FloatLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *FloatLiteral) node() {}
func (i *FloatLiteral) String() string {
	return fmt.Sprintf("%-15s %-15s %s", i.Token, i.Value, i.Position)
}

type StringLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *StringLiteral) node() {}
func (i *StringLiteral) String() string {
	return fmt.Sprintf("%-15s %-15s %s", i.Token, i.Value, i.Position)
}

type NullLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *NullLiteral) node() {}
func (i *NullLiteral) String() string {
	return fmt.Sprintf("%-15s %-15s %s", i.Token, i.Value, i.Position)
}

type BoolLiteral struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *BoolLiteral) node() {}
func (i *BoolLiteral) String() string {
	return fmt.Sprintf("%-15s %-15s %s", i.Token, i.Value, i.Position)
}

type Identifier struct {
	Value    string
	Token    lexer.TokenType
	Position lexer.Position
}

func (i *Identifier) node() {}
func (i *Identifier) String() string {
	return fmt.Sprintf("%-15s %-15s %s", i.Token, i.Value, i.Position)
}

type BinaryExpression struct {
	Left     Node
	Operator lexer.TokenType
	Right    Node
	Position lexer.Position
}

func (b *BinaryExpression) node() {}
func (b *BinaryExpression) String() string {
	return fmt.Sprintf("%-15s %-15s %s", b.Operator, lexer.TokenMap[b.Operator], b.Position)
}

type UnaryExpression struct {
	Operator lexer.TokenType
	Operand  Node
	Position lexer.Position
}

func (u *UnaryExpression) node() {}
func (u *UnaryExpression) String() string {
	return fmt.Sprintf("%-15s %-15s %s", u.Operator, lexer.TokenMap[u.Operator], u.Position)
}

// Stack operations
type StackOp struct {
	Operation lexer.TokenType
	Pos       lexer.Position
}

func (s *StackOp) node() {}
func (s *StackOp) String() string {
	return fmt.Sprintf("%-15s %-15s %s", s.Operation, lexer.TokenMap[s.Operation], s.Pos)
}

type Parameter struct {
	Name string
	Pos  lexer.Position
}

// function declaration
type Procedure struct {
	Name       string
	Parameters []Parameter
	Body       []Node
	Pos        lexer.Position
}

func (p *Procedure) node() {}
func (p *Procedure) String() string {
	return p.Name
}

// Variable assignment
type Assignment struct {
	Target Node
	Pos    lexer.Position
}

func (a *Assignment) node() {}
func (a *Assignment) String() string {
	var target string
	switch t := a.Target.(type) {
	case *Identifier:
		target = t.Value
	case *VarDeclaration:
		target = "VAR " + t.Name
	case *ConstDeclaration:
		target = "CONST " + t.Name
	default:
		target = "<unknown>"
	}

	return fmt.Sprintf("%-15s %-15s %s", lexer.OP_ASSIGN, target, a.Pos)
}

// Print/Println/Dump
type IOStmt struct {
	Kind  lexer.TokenType // DUMP_OP
	Value Node
	Pos   lexer.Position
}

func (p *IOStmt) node() {}
func (p *IOStmt) String() string {
	return fmt.Sprintf("%-15s %-15s %s", p.Kind, "", p.Pos)
}

// Array literal
type ArrayLiteral struct {
	Elements []Node
	Pos      lexer.Position
}

func (a *ArrayLiteral) node()    {}
func (a *ArrayLiteral) Len() int { return len(a.Elements) }
func (a *ArrayLiteral) Homo() bool {
	if len(a.Elements) <= 1 {
		return true
	}

	refType := reflect.TypeOf(a.Elements[0])
	for i := range a.Elements {
		if reflect.TypeOf(a.Elements[i]) != refType {
			return false
		}
	}

	return true
}
func (a *ArrayLiteral) String() string {
	return fmt.Sprintf("%-15s %-15s %s", "ARR", fmt.Sprintf("%d elements", a.Len()), a.Pos)
}

// Variable declaration (VAR name)
type VarDeclaration struct {
	Name string
	Pos  lexer.Position
}

func (v *VarDeclaration) node() {}
func (v *VarDeclaration) String() string {
	return fmt.Sprintf("%-15s %-15s %s", "VAR", v.Name, v.Pos)
}

// Constant declaration (CONST name value)
type ConstDeclaration struct {
	Name  string
	Value Node
	Pos   lexer.Position
}

func (c *ConstDeclaration) node() {}
func (c *ConstDeclaration) String() string {
	return fmt.Sprintf("%-15s %-15s %s", "CONST", c.Name, c.Pos)
}

type DoStmt struct {
	Condition Node
}

type IfStmt struct {
	Condition  []Node
	ThenBranch []Node
	ElseBranch []Node
	Position   lexer.Position
}

func (i *IfStmt) node() {}
func (i *IfStmt) String() string {
	return "IF statement"
}

type WhileStmt struct {
	Condition []Node
	Body      []Node
	Position  lexer.Position
}

func (i *WhileStmt) node() {}
func (i *WhileStmt) String() string {
	return "WHILE statement"
}

type BreakStmt struct {
	Pos lexer.Position
}

func (b *BreakStmt) node() {}
func (b *BreakStmt) String() string {
	return fmt.Sprintf("%-15s %-15s %s", "BREAK", "", b.Pos)
}

type ContinueStmt struct {
	Pos lexer.Position
}

func (c *ContinueStmt) node() {}
func (c *ContinueStmt) String() string {
	return fmt.Sprintf("%-15s %-15s %s", "CONTINUE", "", c.Pos)
}

// SimulateStack shows what the stack looks like as the program executes
func SimulateStack(program *Program) string {
	var result strings.Builder
	stack := []string{}

	result.WriteString("Stack Simulation:\n")
	result.WriteString("=================\n\n")

	for i, stmt := range program.Statements {
		fmt.Fprintf(&result, "Step %d: ", i+1)

		switch n := stmt.(type) {
		case *IntLiteral:
			stack = append(stack, n.Value)
			fmt.Fprintf(&result, "PUSH %s\n", n.Value)

		case *FloatLiteral:
			stack = append(stack, n.Value)
			fmt.Fprintf(&result, "PUSH %s\n", n.Value)

		case *StringLiteral:
			stack = append(stack, fmt.Sprintf("%q", n.Value))
			fmt.Fprintf(&result, "PUSH %q\n", n.Value)

		case *BoolLiteral:
			stack = append(stack, n.Value)
			fmt.Fprintf(&result, "PUSH %s\n", n.Value)

		case *Identifier:
			stack = append(stack, n.Value)
			fmt.Fprintf(&result, "PUSH %s\n", n.Value)

		case *Assignment:
			if len(stack) < 1 {
				fmt.Fprintf(&result, "ERROR: ASSIGN requires 1 value on stack, stack has %d\n", len(stack))
				break
			}
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			target := "<unknown>"
			switch t := n.Target.(type) {
			case *Identifier:
				target = t.Value
			case *VarDeclaration:
				target = "VAR " + t.Name
			}
			fmt.Fprintf(&result, "ASSIGN %s %s\n", target, val)

		case *VarDeclaration:
			fmt.Fprintf(&result, "DECLARE VAR %s\n", n.Name)

		case *ConstDeclaration:
			if n.Value == nil {
				fmt.Fprintf(&result, "DECLARE CONST (missing value) %s\n", n.Name)
			} else {
				fmt.Fprintf(&result, "DECLARE CONST %s\n", n.Name)
			}

		case *BinaryExpression:
			if len(stack) < 2 {
				fmt.Fprintf(&result, "ERROR: %s requires 2 operands, stack has %d\n",
					lexer.TokenMap[n.Operator], len(stack))
				break
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			resultVal := fmt.Sprintf("(%s %s %s)", left, lexer.TokenMap[n.Operator], right)
			stack = append(stack, resultVal)
			fmt.Fprintf(&result, "%s: pop %s, pop %s, push %s\n",
				lexer.TokenMap[n.Operator], right, left, resultVal)

		case *UnaryExpression:
			if len(stack) < 1 {
				fmt.Fprintf(&result, "ERROR: %s requires 1 operand, stack is empty\n",
					lexer.TokenMap[n.Operator])
				break
			}
			operand := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			resultVal := fmt.Sprintf("(%s %s)", lexer.TokenMap[n.Operator], operand)
			stack = append(stack, resultVal)
			fmt.Fprintf(&result, "%s: pop %s, push %s\n",
				lexer.TokenMap[n.Operator], operand, resultVal)

		case *StackOp:
			switch n.Operation {
			case lexer.OP_DUP:
				if len(stack) < 1 {
					result.WriteString("ERROR: DUP requires at least 1 item\n")
					break
				}
				top := stack[len(stack)-1]
				stack = append(stack, top)
				fmt.Fprintf(&result, "DUP: duplicate %s\n", top)

			case lexer.OP_CLEAR:
				if len(stack) < 1 {
					result.WriteString("ERROR: CLEAR requires at least 1 item\n")
					break
				}

				stack = stack[:0]
				fmt.Fprintf(&result, "CLEAR")

			case lexer.OP_DROP:
				if len(stack) < 1 {
					result.WriteString("ERROR: DROP requires at least 1 item\n")
					break
				}
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				fmt.Fprintf(&result, "DROP: remove %s\n", top)

			case lexer.OP_SWAP:
				if len(stack) < 2 {
					result.WriteString("ERROR: SWAP requires at least 2 items\n")
					break
				}
				a := stack[len(stack)-2]
				b := stack[len(stack)-1]
				stack[len(stack)-2] = b
				stack[len(stack)-1] = a
				fmt.Fprintf(&result, "SWAP: %s <-> %s\n", a, b)

			case lexer.OP_OVER:
				if len(stack) < 2 {
					result.WriteString("ERROR: OVER requires at least 2 items\n")
					break
				}
				second := stack[len(stack)-2]
				stack = append(stack, second)
				fmt.Fprintf(&result, "OVER: copy %s to top\n", second)

			case lexer.OP_ROT:
				if len(stack) < 3 {
					result.WriteString("ERROR: ROT requires at least 3 items\n")
					break
				}
				// Rotate top 3: a b c -> b c a
				a := stack[len(stack)-3]
				b := stack[len(stack)-2]
				c := stack[len(stack)-1]
				stack[len(stack)-3] = b
				stack[len(stack)-2] = c
				stack[len(stack)-1] = a
				fmt.Fprintf(&result, "ROT: %s %s %s -> %s %s %s\n", a, b, c, b, c, a)

			default:
				fmt.Fprintf(&result, "%s\n", lexer.TokenMap[n.Operation])
			}

		case *IOStmt:
			if len(stack) < 1 {
				result.WriteString("ERROR: DUMP requires 1 item\n")
				break
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			fmt.Fprintf(&result, "DUMP: output %s (pop)\n", top)

		case *BreakStmt:
			fmt.Fprintf(&result, "BREAK\n")

		case *ContinueStmt:
			fmt.Fprintf(&result, "CONTINUE\n")

		default:
			fmt.Fprintf(&result, "%T\n", stmt)
		}

		// Show current stack state
		fmt.Fprintf(&result, "         Stack: [%s]\n\n", strings.Join(stack, ", "))
	}

	fmt.Fprintf(&result, "Final stack: [%s]\n", strings.Join(stack, ", "))
	return result.String()
}
