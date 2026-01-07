package runtime

import (
	"errors"
	"fmt"
	"gorth/lexer"
	"gorth/parser"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var errBreakSignal = errors.New("gorth: break")
var errContinueSignal = errors.New("gorth: continue")

type Frame struct {
	locals map[string]Value
}

type GorthRuntime struct {
	dataStack   *Stack
	callStack   []*Frame
	returnStack *Stack
	variables   map[string]Value
	constants   map[string]Value
	procedures  map[string]*parser.Procedure
	memory      []Value
	halted      bool
	out         io.Writer
}

func NewRuntime() *GorthRuntime {
	return &GorthRuntime{
		dataStack:   &Stack{items: make([]Value, 0), max: 100},
		returnStack: &Stack{items: make([]Value, 0), max: 100},
		variables:   make(map[string]Value),
		constants:   make(map[string]Value),
		procedures:  make(map[string]*parser.Procedure),
		memory:      make([]Value, 0),
		halted:      false,
		out:         os.Stdout,
	}
}

func (r *GorthRuntime) Execute(program *parser.Program) error {
	for _, stmt := range program.Statements {
		if r.halted {
			break
		}

		err := r.executeNode(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *GorthRuntime) debug_CurrentStack() {
	// show prev, current and next element on stack up to nth element
	fmt.Println("Current Stack State:")

	for i := range r.dataStack.items {
		val := r.dataStack.items[i]
		pointer := " "
		if i == len(r.dataStack.items) {
			pointer = ">"
		}
		fmt.Printf("%s [%d] %s: %s\n", pointer, i, r.typeToString(val.Type), r.valueToString(val))
	}
	fmt.Println()
}

func (r *GorthRuntime) executeNode(node parser.Node) error {
	switch n := node.(type) {
	case *parser.IntLiteral:
		return r.execIntLiteral(n)
	case *parser.FloatLiteral:
		return r.execFloatLiteral(n)
	case *parser.StringLiteral:
		return r.execStrLiteral(n)
	case *parser.BoolLiteral:
		return r.execBoolLiteral(n)
	case *parser.NullLiteral:
		return r.execNullLiteral()
	case *parser.Identifier:
		return r.execIdentifier(n)
	case *parser.BinaryExpression:
		return r.execBinaryOp(n)
	case *parser.UnaryExpression:
		return r.execUnaryOp(n)
	case *parser.StackOp:
		return r.execStackOp(n)
	case *parser.ArrayLiteral:
		return r.execArrayLiteral(n)
	case *parser.IOStmt:
		return r.execIOStmt(n)
	case *parser.VarDeclaration:
		return r.execVarDeclaration(n)
	case *parser.ConstDeclaration:
		return r.execConstDeclaration(n)
	case *parser.Assignment:
		return r.execAssignment(n)
	case *parser.IfStmt:
		return r.execIfStmt(n)
	case *parser.WhileStmt:
		return r.execWhileStmt(n)
	case *parser.BreakStmt:
		return errBreakSignal
	case *parser.ContinueStmt:
		return errContinueSignal
	case *parser.Procedure:
		return r.execProcedure(n)
	default:
		return fmt.Errorf("unknown node type: %T", node)
	}
}

func (r *GorthRuntime) execIntLiteral(n *parser.IntLiteral) error {
	val, err := strconv.Atoi(n.Value)
	if err != nil {
		return err
	}

	return r.dataStack.Push(Value{Type: TYPE_INT, Data: val})
}

func (r *GorthRuntime) execNullLiteral() error {
	return r.dataStack.Push(Value{Type: TYPE_NULL, Data: nil})
}

func (r *GorthRuntime) execUnaryOp(n *parser.UnaryExpression) error {
	val, err := r.dataStack.Pop()
	if err != nil {
		return fmt.Errorf("unary op %s: %v", lexer.TokenMap[n.Operator], err)
	}

	var result Value

	switch n.Operator {
	case lexer.OP_NOT:
		if val.Type != TYPE_BOOL {
			return fmt.Errorf("NOT requires boolean, got %s", r.typeToString(val.Type))
		}
		result = Value{Type: TYPE_BOOL, Data: !val.Data.(bool)}

	case lexer.OP_INC:
		switch val.Type {
		case TYPE_INT:
			result = Value{Type: TYPE_INT, Data: val.Data.(int) + 1}
		case TYPE_FLOAT:
			result = Value{Type: TYPE_FLOAT, Data: val.Data.(float64) + 1.0}
		default:
			return fmt.Errorf("INC requires numeric type, got %s", r.typeToString(val.Type))
		}

	case lexer.OP_DEC:
		switch val.Type {
		case TYPE_INT:
			result = Value{Type: TYPE_INT, Data: val.Data.(int) - 1}
		case TYPE_FLOAT:
			result = Value{Type: TYPE_FLOAT, Data: val.Data.(float64) - 1.0}
		default:
			return fmt.Errorf("DEC requires numeric type, got %s", r.typeToString(val.Type))
		}

	default:
		return fmt.Errorf("unknown unary operator: %s", lexer.TokenMap[n.Operator])
	}

	return r.dataStack.Push(result)
}

func (r *GorthRuntime) execVarDeclaration(n *parser.VarDeclaration) error {
	// VAR just declares a variable with null value
	r.variables[n.Name] = Value{Type: TYPE_NULL, Data: nil}
	return nil
}

func (r *GorthRuntime) execConstDeclaration(n *parser.ConstDeclaration) error {
	// CONST name value - evaluate the value and store it
	err := r.executeNode(n.Value)
	if err != nil {
		return err
	}

	val, err := r.dataStack.Pop()
	if err != nil {
		return err
	}

	r.constants[n.Name] = val
	return nil
}

func (r *GorthRuntime) execStackOp(n *parser.StackOp) error {
	switch n.Operation {
	case lexer.OP_DUP:
		val, err := r.dataStack.Peek()
		if err != nil {
			return err
		}
		return r.dataStack.Push(val)

	case lexer.OP_DROP:
		_, err := r.dataStack.Pop()
		return err

	case lexer.OP_SWAP:
		b, err := r.dataStack.Pop()
		if err != nil {
			return err
		}
		a, err := r.dataStack.Pop()
		if err != nil {
			return err
		}
		r.dataStack.Push(b)
		r.dataStack.Push(a)
		return nil

	case lexer.OP_OVER:
		if r.dataStack.Size() < 2 {
			return fmt.Errorf("OVER requires 2 items on stack")
		}
		second := r.dataStack.items[r.dataStack.Size()-2]
		return r.dataStack.Push(second)

	case lexer.OP_ROT:
		if r.dataStack.Size() < 3 {
			return fmt.Errorf("ROT requires 3 items on stack")
		}
		c, _ := r.dataStack.Pop()
		b, _ := r.dataStack.Pop()
		a, _ := r.dataStack.Pop()
		r.dataStack.Push(b)
		r.dataStack.Push(c)
		r.dataStack.Push(a)
		return nil

	case lexer.OP_PICK:
		val, err := r.dataStack.Pop()
		if err != nil {
			return err
		}

		if val.Type != TYPE_INT {
			return fmt.Errorf("cannot use pick on non-integer types: got %s", r.typeToString(val.Type))
		}

		valInt := val.Data.(int)
		if valInt >= 0 && valInt < len(r.dataStack.items) {
			r.dataStack.Push(r.dataStack.items[val.Data.(int)])
		} else {
			return fmt.Errorf("stack index out of range: %d items total, got index %d", len(r.dataStack.items), valInt)
		}

		return nil

	default:
		return fmt.Errorf("unknown stack operation: %s", lexer.TokenMap[n.Operation])
	}
}

func (r *GorthRuntime) execAssignment(n *parser.Assignment) error {
	val, err := r.dataStack.Pop()
	if err != nil {
		return err
	}

	var name string
	declare := false

	switch target := n.Target.(type) {
	case *parser.Identifier:
		name = target.Value
	case *parser.VarDeclaration:
		name = target.Name
		declare = true
	case *parser.ConstDeclaration:
		name = target.Name
		declare = true
	default:
		return fmt.Errorf("invalid assignment target: %T", n.Target)
	}

	// Cannot assign to const
	if _, ok := r.constants[name]; ok {
		return fmt.Errorf("cannot reassign value to const %s", name)
	}

	if declare {
		switch n.Target.(type) {
		case *parser.VarDeclaration:
			// Inline declaration form: value := VAR x
			if _, ok := r.variables[name]; ok {
				return fmt.Errorf("variable already declared: %s", name)
			}
			r.variables[name] = val
		case *parser.ConstDeclaration:
			// Inline declaration form: value := VAR x
			if _, ok := r.constants[name]; ok {
				return fmt.Errorf("variable already declared: %s", name)
			}
			r.constants[name] = val
		default:
			return fmt.Errorf("invalid assignment target: %T", n.Target)
		}
		return nil
	}

	// Regular assignment: value := x (requires x to exist)
	if _, ok := r.variables[name]; !ok {
		return fmt.Errorf("undefined variable %s", name)
	}

	r.variables[name] = val
	return nil
}

func (r *GorthRuntime) execArrayLiteral(n *parser.ArrayLiteral) error {
	// Evaluate each element and collect values
	elements := make([]Value, len(n.Elements))
	for i, elem := range n.Elements {
		err := r.executeNode(elem)
		if err != nil {
			return err
		}
		val, err := r.dataStack.Pop()
		if err != nil {
			return err
		}
		elements[i] = val
	}

	return r.dataStack.Push(Value{Type: TYPE_ARRAY, Data: elements})
}

func (r *GorthRuntime) execStrLiteral(n *parser.StringLiteral) error {
	return r.dataStack.Push(Value{Type: TYPE_STR, Data: n.Value})
}

func (r *GorthRuntime) execBoolLiteral(n *parser.BoolLiteral) error {
	if n.Value == "TRUE" {
		return r.dataStack.Push(Value{Type: TYPE_BOOL, Data: true})
	} else {
		return r.dataStack.Push(Value{Type: TYPE_BOOL, Data: false})
	}
}

func (r *GorthRuntime) execFloatLiteral(n *parser.FloatLiteral) error {
	val, err := strconv.ParseFloat(n.Value, 64)
	if err != nil {
		return err
	}

	return r.dataStack.Push(Value{Type: TYPE_FLOAT, Data: val})
}

func (r *GorthRuntime) execIdentifier(n *parser.Identifier) error {
	// check if its a variable
	if val, ok := r.variables[n.Value]; ok {
		return r.dataStack.Push(val)
	}

	// check if its a constant
	if val, ok := r.constants[n.Value]; ok {
		return r.dataStack.Push(val)
	}

	// check if its user defined word
	// if _, ok := r.procedures[n.Value]; ok {
	// 	// TODO: execute word
	// 	return fmt.Errorf("word execution not yet implemented: %s", n.Value)
	// }

	// If identifier is undefined, push its name as a string
	// This allows it to be used for assignment: value identifier :=
	return r.dataStack.Push(Value{Type: TYPE_STR, Data: n.Value})
}

func (r *GorthRuntime) execIOStmt(n *parser.IOStmt) error {
	val, err := r.dataStack.Pop()
	if err != nil {
		return err
	}

	switch n.Kind {
	case lexer.OP_DUMP:
		fmt.Fprint(r.out, r.valueToString(val))
	case lexer.OP_DUMPLN:
		fmt.Fprintln(r.out, r.valueToString(val))
	}

	return nil
}

func (r *GorthRuntime) execBinaryOp(n *parser.BinaryExpression) error {
	// Pop two operands
	right, err := r.dataStack.Pop()
	if err != nil {
		return fmt.Errorf("binary op %s: %v", lexer.TokenMap[n.Operator], err)
	}

	left, err := r.dataStack.Pop()
	if err != nil {
		return fmt.Errorf("binary op %s: %v", lexer.TokenMap[n.Operator], err)
	}

	// Perform operation based on operator
	var result Value

	switch n.Operator {
	case lexer.OP_PLUS:
		result, err = r.add(left, right)
	case lexer.OP_SUBTRACT:
		result, err = r.subtract(left, right)
	case lexer.OP_MULTIPLY:
		result, err = r.multiply(left, right)
	case lexer.OP_DIVIDE:
		result, err = r.divide(left, right)
	case lexer.OP_POWER:
		result, err = r.power(left, right)
	case lexer.OP_EQ:
		result, err = r.equal(left, right)
	case lexer.OP_NEQ:
		result, err = r.equal(left, right)
		result.Data = !result.Data.(bool)
	case lexer.OP_GT:
		result, err = r.compareNumeric(left, right, func(a, b float64) bool {
			return a > b
		})
	case lexer.OP_GTE:
		result, err = r.compareNumeric(left, right, func(a, b float64) bool {
			return a >= b
		})
	case lexer.OP_LT:
		result, err = r.compareNumeric(left, right, func(a, b float64) bool {
			return a < b
		})
	case lexer.OP_LTE:
		result, err = r.compareNumeric(left, right, func(a, b float64) bool {
			return a <= b
		})
	case lexer.OP_AND:
		result, err = r.compareLogical(left, right, func(a, b Value) bool {
			return a.Data.(bool) && b.Data.(bool)
		})
	case lexer.OP_OR:
		result, err = r.compareLogical(left, right, func(a, b Value) bool {
			return a.Data.(bool) || b.Data.(bool)
		})
	default:
		return fmt.Errorf("unknown operator: %s", lexer.TokenMap[n.Operator])
	}

	if err != nil {
		return err
	}

	return r.dataStack.Push(result)
}

func (r *GorthRuntime) execIfStmt(n *parser.IfStmt) error {
	for _, stmt := range n.Condition {
		if err := r.executeNode(stmt); err != nil {
			return err
		}
	}

	condition, err := r.dataStack.Pop()
	if err != nil {
		return err
	}

	// condition must be boolean
	if condition.Type != TYPE_BOOL {
		return fmt.Errorf("IF condition must be boolean, got %s", r.typeToString(condition.Type))
	}

	// execute branch
	if condition.Data.(bool) {
		for _, stmt := range n.ThenBranch {
			if err := r.executeNode(stmt); err != nil {
				return err
			}
		}
	} else if len(n.ElseBranch) > 0 {
		for _, stmt := range n.ElseBranch {
			if err := r.executeNode(stmt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *GorthRuntime) execWhileStmt(n *parser.WhileStmt) error {
	for {
		// execute body
		for _, stmt := range n.Body {
			if err := r.executeNode(stmt); err != nil {
				if errors.Is(err, errContinueSignal) {
					// skip remaining body statements, proceed to condition check
					break
				}
				if errors.Is(err, errBreakSignal) {
					// exit the nearest enclosing loop immediately
					return nil
				}
				return err
			}
		}

		// execute postfix condition block
		for _, stmt := range n.Condition {
			if err := r.executeNode(stmt); err != nil {
				if errors.Is(err, errContinueSignal) {
					break
				}
				if errors.Is(err, errBreakSignal) {
					return nil
				}
				return err
			}
		}

		condition, err := r.dataStack.Pop()
		if err != nil {
			return fmt.Errorf("WHILE statement requires condition on stack: %v", err)
		}

		if condition.Type != TYPE_BOOL {
			return fmt.Errorf("WHILE condition must be boolean, got %s", r.typeToString(condition.Type))
		}

		if !condition.Data.(bool) {
			break
		}
	}

	return nil
}

func (r *GorthRuntime) execProcedure(n *parser.Procedure) error {
	r.procedures[n.Name] = n

	for _, param := range n.Parameters {
		// declare parameters as variables with null values
		r.variables[param.Name] = Value{Type: TYPE_NULL, Data: nil}
	}

	return nil
}

func (r *GorthRuntime) not(op Value) (Value, error) {
	if op.Type != TYPE_BOOL {
		return Value{}, fmt.Errorf("cannot use not operation on non-boolean type: %s", r.typeToString(op.Type))
	}

	return Value{Type: TYPE_BOOL, Data: !op.Data.(bool)}, nil
}

func (r *GorthRuntime) inc(op Value) (Value, error) {
	if op.Type != TYPE_INT {
		return Value{}, fmt.Errorf("cannot use inc operation on non-int type: %s", r.typeToString(op.Type))
	}

	return Value{Type: TYPE_INT, Data: op.Data.(int) + 1}, nil
}

func (r *GorthRuntime) dec(op Value) (Value, error) {
	if op.Type != TYPE_INT {
		return Value{}, fmt.Errorf("cannot use dec operation on non-int type: %s", r.typeToString(op.Type))
	}

	return Value{Type: TYPE_INT, Data: op.Data.(int) - 1}, nil
}

func (r *GorthRuntime) add(left, right Value) (Value, error) {
	isLeftNumeric := left.Type == TYPE_INT || left.Type == TYPE_FLOAT
	isRightNumeric := right.Type == TYPE_INT || right.Type == TYPE_FLOAT
	isLeftString := left.Type == TYPE_STR
	isRightString := right.Type == TYPE_STR

	// TODO: refactor this later
	// string concatenation
	if isLeftString && isRightString {
		return Value{Type: TYPE_STR, Data: fmt.Sprintf("%s%s", left.Data.(string), right.Data.(string))}, nil
	}

	// string + numeric
	if isLeftString && isRightNumeric {
		return Value{Type: TYPE_STR, Data: fmt.Sprintf("%s%s", left.Data.(string), r.valueToString(right))}, nil
	}

	// numeric + string
	if isLeftNumeric && isRightString {
		return Value{Type: TYPE_STR, Data: fmt.Sprintf("%s%s", r.valueToString(left), right.Data.(string))}, nil
	}

	// numeric addition
	if !isLeftNumeric || !isRightNumeric {
		return Value{}, fmt.Errorf(
			"cannot add %s and %s",
			r.typeToString(left.Type), r.typeToString(right.Type),
		)
	}

	if left.Type == TYPE_INT && right.Type == TYPE_INT {
		return Value{
			Type: TYPE_INT,
			Data: left.Data.(int) + right.Data.(int),
		}, nil
	}

	lf := r.toFloat(left)
	rf := r.toFloat(right)

	return Value{Type: TYPE_FLOAT, Data: lf + rf}, nil
}

func (r *GorthRuntime) subtract(left, right Value) (Value, error) {
	// a b - = b - a
	isLeftNumeric := left.Type == TYPE_INT || left.Type == TYPE_FLOAT
	isRightNumeric := right.Type == TYPE_INT || right.Type == TYPE_FLOAT

	if !isLeftNumeric || !isRightNumeric {
		return Value{}, fmt.Errorf(
			"cannot subtract %s and %s",
			r.typeToString(left.Type), r.typeToString(right.Type),
		)
	}

	if left.Type == TYPE_INT && right.Type == TYPE_INT {
		return Value{
			Type: TYPE_INT,
			Data: left.Data.(int) - right.Data.(int),
		}, nil
	}

	lf := r.toFloat(left)
	rf := r.toFloat(right)

	return Value{Type: TYPE_FLOAT, Data: lf - rf}, nil
}

func (r *GorthRuntime) multiply(left, right Value) (Value, error) {
	isLeftNumeric := left.Type == TYPE_INT || left.Type == TYPE_FLOAT
	isRightNumeric := right.Type == TYPE_INT || right.Type == TYPE_FLOAT

	if !isLeftNumeric || !isRightNumeric {
		return Value{}, fmt.Errorf(
			"cannot multiple %s and %s",
			r.typeToString(left.Type), r.typeToString(right.Type),
		)
	}

	if left.Type == TYPE_INT && right.Type == TYPE_INT {
		return Value{
			Type: TYPE_INT,
			Data: right.Data.(int) * left.Data.(int),
		}, nil
	}

	lf := r.toFloat(left)
	rf := r.toFloat(right)

	return Value{Type: TYPE_FLOAT, Data: rf * lf}, nil
}

func (r *GorthRuntime) divide(left, right Value) (Value, error) {
	isLeftNumeric := left.Type == TYPE_INT || left.Type == TYPE_FLOAT
	isRightNumeric := right.Type == TYPE_INT || right.Type == TYPE_FLOAT

	if !isLeftNumeric || !isRightNumeric {
		return Value{}, fmt.Errorf(
			"cannot divide %s and %s",
			r.typeToString(left.Type), r.typeToString(right.Type),
		)
	}

	lf := r.toFloat(left)
	rf := r.toFloat(right)

	return Value{Type: TYPE_FLOAT, Data: lf / rf}, nil
}

func (r *GorthRuntime) power(left, right Value) (Value, error) {
	isLeftNumeric := left.Type == TYPE_INT || left.Type == TYPE_FLOAT
	isRightNumeric := right.Type == TYPE_INT || right.Type == TYPE_FLOAT

	if !isLeftNumeric || !isRightNumeric {
		return Value{}, fmt.Errorf(
			"cannot exponentiate %s and %s",
			r.typeToString(left.Type), r.typeToString(right.Type),
		)
	}

	lf := r.toFloat(left)
	rf := r.toFloat(right)
	pow := math.Pow(lf, rf)

	if left.Type == TYPE_INT && right.Type == TYPE_INT {
		return Value{
			Type: TYPE_INT,
			Data: int(pow),
		}, nil
	}

	return Value{Type: TYPE_FLOAT, Data: pow}, nil
}

func (r *GorthRuntime) equal(left, right Value) (Value, error) {
	if isNumeric(left.Type) && isNumeric(right.Type) {
		return Value{
			Type: TYPE_BOOL,
			Data: r.toFloat(left) == r.toFloat(right),
		}, nil
	}

	if left.Type == right.Type {
		return Value{
			Type: TYPE_BOOL,
			Data: left.Data == right.Data,
		}, nil
	}

	return Value{
		Type: TYPE_BOOL,
		Data: false,
	}, nil
}

func (r *GorthRuntime) compareNumeric(
	left, right Value,
	op func(a, b float64) bool,
) (Value, error) {
	if !isNumeric(left.Type) || !isNumeric(right.Type) {
		return Value{}, fmt.Errorf("invalid numeric comparison")
	}

	return Value{
		Type: TYPE_BOOL,
		Data: op(r.toFloat(left), r.toFloat(right)),
	}, nil
}

func (r *GorthRuntime) compareLogical(
	left, right Value,
	op func(a, b Value) bool,
) (Value, error) {
	if !isBool(left.Type) || !isBool(right.Type) {
		return Value{Type: TYPE_BOOL, Data: false}, nil
	}

	return Value{
		Type: TYPE_BOOL,
		Data: op(left, right),
	}, nil
}

func (r *GorthRuntime) modulo(left, right Value) (Value, error) {
	isLeftNumeric := left.Type == TYPE_INT || left.Type == TYPE_FLOAT
	isRightNumeric := right.Type == TYPE_INT || right.Type == TYPE_FLOAT

	if !isLeftNumeric || !isRightNumeric {
		return Value{}, fmt.Errorf(
			"cannot subtract %s and %s",
			r.typeToString(left.Type), r.typeToString(right.Type),
		)
	}

	if left.Type == TYPE_INT && right.Type == TYPE_INT {
		return Value{
			Type: TYPE_INT,
			Data: left.Data.(int) % right.Data.(int),
		}, nil
	}

	lf := r.toFloat(left)
	rf := r.toFloat(right)

	return Value{Type: TYPE_FLOAT, Data: math.Mod(lf, rf)}, nil
}

func (r *GorthRuntime) toFloat(v Value) float64 {
	switch v.Type {
	case TYPE_INT:
		return float64(v.Data.(int))
	case TYPE_FLOAT:
		return v.Data.(float64)
	default:
		// TODO: might just need to make this throw an error
		return 0.0
	}
}

func (r *GorthRuntime) valueToString(v Value) string {
	switch v.Type {
	case TYPE_INT:
		return strconv.Itoa(v.Data.(int))
	case TYPE_FLOAT:
		return strconv.FormatFloat(v.Data.(float64), 'f', -1, 64)
	case TYPE_STR:
		return v.Data.(string)
	case TYPE_BOOL:
		if v.Data.(bool) {
			return "TRUE"
		}
		return "FALSE"
	case TYPE_NULL:
		return "NULL"
	case TYPE_ARRAY:
		// represent array as [elem1, elem2, ...]
		arr := v.Data.([]Value)
		var str strings.Builder
		str.WriteString("[")
		for i, elem := range arr {
			str.WriteString(r.valueToString(elem))
			if i < len(arr)-1 {
				str.WriteString(", ")
			}
		}
		str.WriteString("]")
		return str.String()
	case TYPE_WORD:
		return "<WORD>"
	default:
		return "<UNKNOWN>"
	}
}

// PrintState prints the current state of the runtime (for debugging)
func (r *GorthRuntime) PrintState() {
	fmt.Println("Data Stack:")
	if r.dataStack.Size() == 0 {
		fmt.Println("  <empty>")
	} else {
		for i := r.dataStack.Size() - 1; i >= 0; i-- {
			val := r.dataStack.items[i]
			fmt.Printf("  [%d] %s: %s\n", i, r.typeToString(val.Type), r.valueToString(val))
		}
	}

	fmt.Println("\nVariables:")
	if len(r.variables) == 0 {
		fmt.Println("  <none>")
	} else {
		for name, val := range r.variables {
			fmt.Printf("  %s = %s (%s)\n", name, r.valueToString(val), r.typeToString(val.Type))
		}
	}

	fmt.Println("\nConstants:")
	if len(r.constants) == 0 {
		fmt.Println("  <none>")
	} else {
		for name, val := range r.constants {
			fmt.Printf("  %s = %s (%s)\n", name, r.valueToString(val), r.typeToString(val.Type))
		}
	}

	fmt.Println("\nUser-defined Procedures:")
	if len(r.procedures) == 0 {
		fmt.Println("  <none>")
	} else {
		for name := range r.procedures {
			fmt.Printf("  %s\n", name)
		}
	}
}

func (r *GorthRuntime) typeToString(t ValueType) string {
	switch t {
	case TYPE_INT:
		return "int"
	case TYPE_FLOAT:
		return "float"
	case TYPE_STR:
		return "string"
	case TYPE_BOOL:
		return "bool"
	case TYPE_NULL:
		return "null"
	case TYPE_ARRAY:
		return "array"
	case TYPE_WORD:
		return "word"
	default:
		return "unknown"
	}
}

func isNumeric(t ValueType) bool {
	return t == TYPE_INT || t == TYPE_FLOAT
}

func isBool(t ValueType) bool {
	return t == TYPE_BOOL
}
