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

func (p *Parser) nextToken() error {
	p.currentToken = p.peekToken
	p.currentLiteral = p.peekLiteral
	p.currentPosition = p.peekPosition

	lexedToken, err := p.lexer.NextToken()
	if err != nil {
		return err
	}

	p.peekPosition = lexedToken.Pos
	p.peekToken = lexedToken.Type
	p.peekLiteral = lexedToken.Literal

	return nil
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
	case lexer.VAR:
		return p.parseVarDeclaration()
	case lexer.CONST:
		return p.parseConstDeclaration()
	case lexer.IDENT:
		return p.parseIdentLiteral()
	case lexer.NULL:
		return p.parseNullLiteral()
	// binary ops
	case lexer.OP_PLUS, lexer.OP_SUBTRACT, lexer.OP_MULTIPLY, lexer.OP_DIVIDE, lexer.OP_POWER,
		lexer.OP_MODULO, lexer.OP_EQ, lexer.OP_NEQ, lexer.OP_GT, lexer.OP_LT, lexer.OP_GTE,
		lexer.OP_LTE, lexer.OP_AND, lexer.OP_OR:
		return p.parseBinaryOp()
	// unary ops
	case lexer.OP_NOT, lexer.OP_INC, lexer.OP_DEC:
		return p.parseUnaryOp()
	// stack operations
	case lexer.OP_DROP, lexer.OP_SWAP, lexer.OP_DUP, lexer.OP_OVER, lexer.OP_ROT, lexer.OP_DEL, lexer.OP_CLEAR, lexer.OP_PICK:
		return p.parseStackOp()
	// io operations
	case lexer.OP_DUMP:
		return p.parseIOOP()
	case lexer.OP_ASSIGN:
		return p.parseAssignment()
	case lexer.ILLEGAL:
		return nil, fmt.Errorf("illegal token at line %d, col %d: %s",
			p.currentPosition.Line, p.currentPosition.Column, p.currentLiteral)
	// array literal
	case lexer.LBRACKET:
		return p.parseArray()
	case lexer.IF:
		return p.parseIfStmt()
	case lexer.WHILE:
		return p.parseWhileStmt()
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

func (p *Parser) parseNullLiteral() (Node, error) {
	return &NullLiteral{
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

func (p *Parser) parseStackOp() (Node, error) {
	return &StackOp{
		Operation: p.currentToken,
		Pos:       p.currentPosition,
	}, nil
}

func (p *Parser) parseIOOP() (Node, error) {
	return &IOStmt{
		Kind: p.currentToken,
		Pos:  p.currentPosition,
		// Value will be filled in by interpreter based on stack value
	}, nil
}

// Variable and constant declarations
func (p *Parser) parseVarDeclaration() (Node, error) {
	// VAR variableName
	pos := p.currentPosition
	p.nextToken() // move to identifier

	if !p.currentTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected identifier after VAR at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	varName := p.currentLiteral

	return &VarDeclaration{
		Name: varName,
		Pos:  pos,
	}, nil
}

func (p *Parser) parseConstDeclaration() (Node, error) {
	// CONST constName value
	pos := p.currentPosition
	p.nextToken() // move to identifier

	if !p.currentTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected identifier after CONST at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	constName := p.currentLiteral
	p.nextToken() // move to value

	// Parse the value (could be any literal)
	var value Node
	var err error

	switch p.currentToken {
	case lexer.INT:
		value, err = p.parseIntLiteral()
	case lexer.FLOAT:
		value, err = p.parseFloatLiteral()
	case lexer.STRING:
		value, err = p.parseStrLiteral()
	case lexer.BOOL:
		value, err = p.parseBoolLiteral()
	case lexer.NULL:
		value, err = p.parseNullLiteral()
	default:
		return nil, fmt.Errorf("expected literal value after const name at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	if err != nil {
		return nil, err
	}

	return &ConstDeclaration{
		Name:  constName,
		Value: value,
		Pos:   pos,
	}, nil
}

// assignment and special ops
func (p *Parser) parseAssignment() (Node, error) {
	// Current token is OP_ASSIGN (->)
	pos := p.currentPosition

	// Next token MUST be an identifier (variable name)
	p.nextToken()

	if !p.currentTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected identifier after '->' at line %d, col %d, got %s",
			p.currentPosition.Line, p.currentPosition.Column, lexer.TokenMap[p.currentToken])
	}

	varName := p.currentLiteral

	return &Assignment{
		Name: varName,
		Pos:  pos,
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

func (p *Parser) parseIfStmt() (Node, error) {
	pos := p.currentPosition

	p.nextToken()

	// then branch should be a series of nodes
	thenBranch := []Node{}
	for !p.currentTokenIs(lexer.ELSE) && !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		// add the statement to the thenbranch if if exists
		if stmt != nil {
			thenBranch = append(thenBranch, stmt)
		}

		// advance the parser
		p.nextToken()
	}

	// parse optional else
	var elseBranch []Node
	if p.currentTokenIs(lexer.ELSE) {
		// consume else
		p.nextToken()

		for !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}

			// add the statement to the thenbranch if if exists
			if stmt != nil {
				elseBranch = append(elseBranch, stmt)
			}

			// advance the parser
			p.nextToken()
		}
	}

	// expect an end
	if !p.currentTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	return &IfStmt{
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
		Position:   pos,
	}, nil

}

func (p *Parser) parseWhileStmt() (Node, error) {
	pos := p.currentPosition

	p.nextToken()

	body := []Node{}
	for !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			body = append(body, stmt)
		}

		p.nextToken()
	}

	if !p.currentTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	return &WhileStmt{
		Body:     body,
		Position: pos,
	}, nil
}
