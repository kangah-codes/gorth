package runtime

import (
	"fmt"
	"gorth/lexer"
	"gorth/parser"
	"math"
	"strconv"
)

type GorthRuntime struct {
	dataStack   *Stack
	returnStack *Stack
	variables   map[string]Value
	words       map[string]*Word
	memory      []Value
	halted      bool
}

func NewRuntime() *GorthRuntime {
	return &GorthRuntime{
		dataStack:   &Stack{items: make([]Value, 0), max: 999_999_999},
		returnStack: &Stack{items: make([]Value, 0), max: 999_999_999},
		variables:   make(map[string]Value),
		words:       make(map[string]*Word),
		memory:      make([]Value, 0),
		halted:      false,
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

func (r *GorthRuntime) executeNode(node parser.Node) error {
	switch n := node.(type) {
	case *parser.IntLiteral:
		return r.execIntLiteral(n)
	case *parser.FloatLiteral:
		return r.execFloatLiteral(n)
	case *parser.StringLiteral:
		return r.execStrLiteral(n)
	case *parser.NullLiteral:
		return r.execNullLiteral()
	case *parser.BoolLiteral:
		return r.execBoolLiteral(n)
	case *parser.UnaryExpression:
		return r.execUnaryOp(n)
	case *parser.BinaryExpression:
		return r.execBinaryOp(n)
	case *parser.IOStmt:
		return r.execIOStmt(n)
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

	// check if its user defined
	// TODO: implement

	return fmt.Errorf("undefined identifier: %s", n.Value)
}

func (r *GorthRuntime) execIOStmt(n *parser.IOStmt) error {
	val, err := r.dataStack.Pop()
	if err != nil {
		return err
	}

	switch n.Kind {
	case lexer.OP_DUMP:
		fmt.Print(r.valueToString(val))
	}

	return nil
}

func (r *GorthRuntime) execUnaryOp(n *parser.UnaryExpression) error {
	operand, err := r.dataStack.Pop()
	if err != nil {
		return fmt.Errorf("unary op %s: %v", lexer.TokenMap[n.Operator], err)
	}

	var result Value

	switch n.Operator {
	case lexer.OP_NOT:
		result, err = r.not(operand)
	case lexer.OP_INC:
		result, err = r.inc(operand)
	case lexer.OP_DEC:
		result, err = r.dec(operand)
	default:
		return fmt.Errorf("unknown operator: %s", lexer.TokenMap[n.Operator])
	}

	return r.dataStack.Push(result)
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

func (r *GorthRuntime) not(op Value) (Value, error) {
	if op.Type != TYPE_BOOL {
		return Value{}, fmt.Errorf("cannot use not operation on non-boolean type: %s", op.Type)
	}

	return Value{Type: TYPE_BOOL, Data: !op.Data.(bool)}, nil
}

func (r *GorthRuntime) inc(op Value) (Value, error) {
	if op.Type != TYPE_INT {
		return Value{}, fmt.Errorf("cannot use inc operation on non-int type: %s", op.Type)
	}

	return Value{Type: TYPE_INT, Data: op.Data.(int) + 1}, nil
}

func (r *GorthRuntime) dec(op Value) (Value, error) {
	if op.Type != TYPE_INT {
		return Value{}, fmt.Errorf("cannot use dec operation on non-int type: %s", op.Type)
	}

	return Value{Type: TYPE_INT, Data: op.Data.(int) - 1}, nil
}

func (r *GorthRuntime) add(left, right Value) (Value, error) {
	isLeftNumeric := left.Type == TYPE_INT || left.Type == TYPE_FLOAT
	isRightNumeric := right.Type == TYPE_INT || right.Type == TYPE_FLOAT

	if !isLeftNumeric || !isRightNumeric {
		return Value{}, fmt.Errorf(
			"cannot add %s and %s",
			left.Type, right.Type,
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
			left.Type, right.Type,
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
			left.Type, right.Type,
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
			left.Type, right.Type,
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
			left.Type, right.Type,
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
			left.Type, right.Type,
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
		return fmt.Sprintf("%d", v.Data.(int))
	case TYPE_FLOAT:
		return fmt.Sprintf("%f", v.Data.(float64))
	case TYPE_STR:
		return v.Data.(string)
	case TYPE_BOOL:
		return fmt.Sprintf("%t", v.Data.(bool))
	case TYPE_NULL:
		return "null"
	default:
		return fmt.Sprintf("%v", v.Data)
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

	fmt.Println("\nUser-defined Words:")
	if len(r.words) == 0 {
		fmt.Println("  <none>")
	} else {
		for name := range r.words {
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
