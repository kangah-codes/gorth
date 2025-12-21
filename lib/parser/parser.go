package parser

import (
	"fmt"
	"gorth/lexer"
)

type Parser struct {
	lexer           *lexer.Lexer
	currentToken    lexer.TokenType
	currentLiteral  string
	currentPosition lexer.Position
	peekToken       lexer.TokenType
	peekLiteral     string
	peekPosition    lexer.Position
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	// read next 2 tokens to init current and peek
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) currentTokenIs(t lexer.TokenType) bool {
	return p.currentToken == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken == t
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.currentLiteral = p.peekLiteral
	p.currentPosition = p.peekPosition

	lexedToken, err := p.lexer.NextToken()
	if err != nil {
		panic("not implemented")
	}

	p.peekPosition = lexedToken.Pos
	p.peekToken = lexedToken.Type
	p.peekLiteral = lexedToken.Literal
}

func (p *Parser) Parse() (*Program, error) {
	program := &Program{
		Statements: []Node{},
	}

	// check if current token is not EOF
	for !p.currentTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program, nil
}

func (p *Parser) parseStatement() (Node, error) {
	switch p.currentToken {
	case lexer.INT:
		return p.parseIntLiteral()
	case lexer.STRING:
		return p.parseStrLiteral()
	case lexer.BOOL:
		return p.parseBoolLiteral()
	case lexer.FLOAT:
		return p.parseFloatLiteral()
	case lexer.VARIABLE:
		return p.parseIdentLiteral()
	case lexer.PTR:
		return p.parsePtr()
	case lexer.PTR_DEREF:
		return p.parsePtrDeref()
	// binary ops
	case lexer.OP_PLUS, lexer.OP_MINUS, lexer.OP_MULTIPLY, lexer.OP_DIVIDE, lexer.OP_POWER, lexer.OP_MODULO:
		return p.parseBinaryOp()
	// comparison ops
	case lexer.EQ, lexer.NEQ, lexer.GT, lexer.LT, lexer.GTE, lexer.LTE:
		return p.parseComparisonOp()
	// logical ops
	case lexer.OP_AND, lexer.OP_OR:
		return p.parseLogicalOp()
	// unary ops
	case lexer.OP_NOT, lexer.OP_INC, lexer.OP_DEC, lexer.OP_DUMP: // adding dump op because even though it prints its unary & modifies stack
		return p.parseUnaryOp()
	// stack operations
	case lexer.OP_DROP, lexer.OP_SWAP, lexer.OP_DUP, lexer.OP_OVER, lexer.OP_ROT, lexer.OP_DEL:
		return p.parseStackOp()
	// print ops
	case lexer.OP_PRINT, lexer.OP_PRINTLN:
		return p.parsePrintOp()
	case lexer.ASSIGN:
		return p.parseAssignment()
	case lexer.ILLEGAL:
		return nil, fmt.Errorf("illegal token at line %d, col %d: %s",
			p.currentPosition.Line, p.currentPosition.Column, p.currentLiteral)
	// array literal
	case lexer.LBRACKET:
		return p.parseArray()
	default:
		return nil, fmt.Errorf("unexpected token %s at line %d, col %d",
			lexer.TokenMap[p.currentToken], p.currentPosition.Line, p.currentPosition.Column)
	}
}

// literal parsing
func (p *Parser) parseIntLiteral() (Node, error) {
	return &IntLiteral{
		Value:    p.currentLiteral,
		Position: p.currentPosition,
		Token:    p.currentToken,
	}, nil
}

func (p *Parser) parseFloatLiteral() (Node, error) {
	return &FloatLiteral{
		Value:    p.currentLiteral,
		Position: p.currentPosition,
		Token:    p.currentToken,
	}, nil
}

func (p *Parser) parseStrLiteral() (Node, error) {
	return &StringLiteral{
		Value:    p.currentLiteral,
		Position: p.currentPosition,
		Token:    p.currentToken,
	}, nil
}

func (p *Parser) parseBoolLiteral() (Node, error) {
	return &BoolLiteral{
		Value:    p.currentLiteral,
		Position: p.currentPosition,
		Token:    p.currentToken,
	}, nil
}

func (p *Parser) parseIdentLiteral() (Node, error) {
	return &Identifier{
		Value:    p.currentLiteral,
		Position: p.currentPosition,
		Token:    p.currentToken,
	}, nil
}

// operator parsing
func (p *Parser) parseUnaryOp() (Node, error) {
	return &UnaryExpression{
		Operator: p.currentToken,
		Position: p.currentPosition,
		// node will be filled in by interpreter based on stack value
	}, nil
}

func (p *Parser) parseBinaryOp() (Node, error) {
	return &BinaryExpression{
		Operator: p.currentToken,
		Position: p.currentPosition,
		// nodes will be filled in by interpreter based on stack values
	}, nil
}

func (p *Parser) parseComparisonOp() (Node, error) {
	return &BinaryExpression{
		Operator: p.currentToken,
		Position: p.currentPosition,
		// nodes will be filled in by interpreter based on stack values
	}, nil
}

func (p *Parser) parseLogicalOp() (Node, error) {
	return &BinaryExpression{
		Operator: p.currentToken,
		Position: p.currentPosition,
		// nodes will be filled in by interpreter based on stack values
	}, nil
}

func (p *Parser) parseStackOp() (Node, error) {
	return &StackOp{
		Operation: p.currentToken,
		Pos:       p.currentPosition,
	}, nil
}

func (p *Parser) parsePrintOp() (Node, error) {
	return &PrintStmt{
		Kind: p.currentToken,
		Pos:  p.currentPosition,
		// Value will be filled in by interpreter based on stack value
	}, nil
}

// assignment and special ops
func (p *Parser) parseAssignment() (Node, error) {
	return &StackOp{
		Operation: lexer.ASSIGN,
		Pos:       p.currentPosition,
	}, nil
}

func (p *Parser) parsePtr() (Node, error) {
	return &PointerExpr{
		Operation: lexer.PTR,
		Value: &IntLiteral{
			Value:    p.currentLiteral[1:], // remove the * prefix
			Position: p.currentPosition,
		},
		Pos: p.currentPosition,
	}, nil
}

func (p *Parser) parsePtrDeref() (Node, error) {
	return &PointerExpr{
		Operation: lexer.PTR_DEREF,
		Value: &IntLiteral{
			Value:    p.currentLiteral[1:], // remove the @ prefix
			Position: p.currentPosition,
		},
		Pos: p.currentPosition,
	}, nil
}

func (p *Parser) parseArray() (Node, error) {
	elements := []Node{}
	pos := p.currentPosition

	// Skip opening bracket [
	p.nextToken()

	// Parse elements until we hit closing bracket
	for !p.currentTokenIs(lexer.RBRACKET) && !p.currentTokenIs(lexer.EOF) {
		// Parse the element
		elem, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)

		// Move to next token
		p.nextToken()

		// If it's a comma, skip it
		if p.currentTokenIs(lexer.COMMA) {
			p.nextToken()
		} else if !p.currentTokenIs(lexer.RBRACKET) {
			return nil, fmt.Errorf("expected ',' or ']' in array but got %s at line %d, col %d",
				lexer.TokenMap[p.currentToken], p.currentPosition.Line, p.currentPosition.Column)
		}
	}

	if !p.currentTokenIs(lexer.RBRACKET) {
		return nil, fmt.Errorf("expected ']' to close array at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	return &ArrayLiteral{
		Elements: elements,
		Pos:      pos,
	}, nil
}
