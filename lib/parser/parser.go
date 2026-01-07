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
	loopDepth       int
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

	if err := validateLoopControl(program); err != nil {
		return nil, err
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
		return p.parseConstDefinition()
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
	case lexer.OP_DUMP, lexer.OP_DUMPLN:
		return p.parseIOOP()
	case lexer.OP_ASSIGN:
		return p.parseAssignment()
	case lexer.ILLEGAL:
		return nil, fmt.Errorf("illegal token at line %d, col %d: %s",
			p.currentPosition.Line, p.currentPosition.Column, p.currentLiteral)
	// array literal
	case lexer.LSQUARE_BRACKET:
		return p.parseArray()
	case lexer.PROC:
		return p.parseProcedure()
	case lexer.DO:
		return p.parseDoStmt()
	case lexer.WHILE:
		return p.parseWhileStmt()
	case lexer.BREAK:
		return &BreakStmt{Pos: p.currentPosition}, nil
	case lexer.CONTINUE:
		return &ContinueStmt{Pos: p.currentPosition}, nil
	case lexer.CALL:
		return p.parseProcCall()
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
	// CONST constName
	// NOTE: This parses a CONST *target* (used by ':='), not a full const definition.
	pos := p.currentPosition
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if !p.currentTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected identifier after CONST at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	constName := p.currentLiteral

	return &ConstDeclaration{
		Name: constName,
		Pos:  pos,
		// Value intentionally omitted for inline target form
	}, nil
}

func (p *Parser) parseConstDefinition() (Node, error) {
	// CONST constName value
	pos := p.currentPosition
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if !p.currentTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected identifier after CONST at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	constName := p.currentLiteral

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if p.currentTokenIs(lexer.EOF) {
		return nil, fmt.Errorf("CONST %s requires an assigned value at declaration (e.g. CONST %s 10) at line %d, col %d",
			constName, constName, pos.Line, pos.Column)
	}

	var value Node
	switch p.currentToken {
	case lexer.INT:
		value = &IntLiteral{Value: p.currentLiteral, Position: p.currentPosition, Token: p.currentToken}
	case lexer.FLOAT:
		value = &FloatLiteral{Value: p.currentLiteral, Position: p.currentPosition, Token: p.currentToken}
	case lexer.STRING:
		value = &StringLiteral{Value: p.currentLiteral, Position: p.currentPosition, Token: p.currentToken}
	case lexer.BOOL:
		value = &BoolLiteral{Value: p.currentLiteral, Position: p.currentPosition, Token: p.currentToken}
	case lexer.NULL:
		value = &NullLiteral{Value: p.currentLiteral, Position: p.currentPosition, Token: p.currentToken}
	case lexer.LSQUARE_BRACKET:
		arr, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		value = arr
	default:
		return nil, fmt.Errorf("expected literal value after CONST %s at line %d, col %d, got %s",
			constName, p.currentPosition.Line, p.currentPosition.Column, lexer.TokenMap[p.currentToken])
	}

	return &ConstDeclaration{
		Name:  constName,
		Value: value,
		Pos:   pos,
	}, nil
}

// assignment and special ops
func (p *Parser) parseAssignment() (Node, error) {
	// current token is OP_ASSIGN :=
	pos := p.currentPosition

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	var target Node
	switch p.currentToken {
	case lexer.VAR:
		// inline declaration target
		t, err := p.parseVarDeclaration()
		if err != nil {
			return nil, err
		}
		target = t
	case lexer.CONST:
		// inline const declaration target
		t, err := p.parseConstDeclaration()
		if err != nil {
			return nil, err
		}
		target = t
	case lexer.IDENT:
		target = &Identifier{
			Value:    p.currentLiteral,
			Position: p.currentPosition,
			Token:    p.currentToken,
		}
	default:
		return nil, fmt.Errorf("expected identifier, VAR, or CONST after ':=' at line %d, col %d, got %s",
			p.currentPosition.Line, p.currentPosition.Column, lexer.TokenMap[p.currentToken])
	}

	return &Assignment{
		Target: target,
		Pos:    pos,
	}, nil
}

func (p *Parser) parseArray() (Node, error) {
	elements := []Node{}
	pos := p.currentPosition

	// Skip opening bracket [
	p.nextToken()

	// Parse elements until we hit closing bracket
	for !p.currentTokenIs(lexer.RSQUARE_BRACKET) && !p.currentTokenIs(lexer.EOF) {
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
		} else if !p.currentTokenIs(lexer.RSQUARE_BRACKET) {
			return nil, fmt.Errorf("expected ',' or ']' in array but got %s at line %d, col %d",
				lexer.TokenMap[p.currentToken], p.currentPosition.Line, p.currentPosition.Column)
		}
	}

	if !p.currentTokenIs(lexer.RSQUARE_BRACKET) {
		return nil, fmt.Errorf("expected ']' to close array at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	return &ArrayLiteral{
		Elements: elements,
		Pos:      pos,
	}, nil
}

func (p *Parser) parseDoStmt() (Node, error) {
	// DO ... WHILE <cond> END
	// DO ... [ELSE ...] IF <cond> END
	pos := p.currentPosition

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	thenOrBody := []Node{}
	for !p.currentTokenIs(lexer.WHILE) && !p.currentTokenIs(lexer.IF) && !p.currentTokenIs(lexer.ELSE) && !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			thenOrBody = append(thenOrBody, stmt)
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
	}

	// do-while form
	if p.currentTokenIs(lexer.WHILE) {
		// consume WHILE
		if err := p.nextToken(); err != nil {
			return nil, err
		}

		cond := []Node{}
		for !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				cond = append(cond, stmt)
			}
			if err := p.nextToken(); err != nil {
				return nil, err
			}
		}

		if !p.currentTokenIs(lexer.END) {
			return nil, fmt.Errorf("expected END at line %d, col %d",
				p.currentPosition.Line, p.currentPosition.Column)
		}

		return &WhileStmt{Body: thenOrBody, Condition: cond, Position: pos}, nil
	}

	// do-if form
	var elseBranch []Node
	if p.currentTokenIs(lexer.ELSE) {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		for !p.currentTokenIs(lexer.IF) && !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				elseBranch = append(elseBranch, stmt)
			}
			if err := p.nextToken(); err != nil {
				return nil, err
			}
		}
	}

	if !p.currentTokenIs(lexer.IF) {
		return nil, fmt.Errorf("expected IF before END at line %d, col %d, got %s",
			p.currentPosition.Line, p.currentPosition.Column, lexer.TokenMap[p.currentToken])
	}

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	cond := []Node{}
	for !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			cond = append(cond, stmt)
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
	}

	if !p.currentTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	return &IfStmt{ThenBranch: thenOrBody, ElseBranch: elseBranch, Position: pos, Condition: cond}, nil
}

func (p *Parser) parseWhileStmt() (Node, error) {
	pos := p.currentPosition

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	body := []Node{}
	for !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			body = append(body, stmt)
		}

		if err := p.nextToken(); err != nil {
			return nil, err
		}
	}

	if !p.currentTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	return &WhileStmt{
		Body:      body,
		Condition: nil,
		Position:  pos,
	}, nil
}

func (p *Parser) parseProcCall() (Node, error) {
	pos := p.currentPosition

	// consume token
	p.nextToken()

	return &CallStmt{Pos: pos, ProcName: p.currentLiteral}, nil
}

func (p *Parser) parseProcedure() (Node, error) {
	pos := p.currentPosition

	// PROC add (a b)
	// 	a b +
	// END

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if !p.currentTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected procedure name after PROC at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	procName := p.currentLiteral
	// parse parameters
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if !p.currentTokenIs(lexer.LBRACKET) {
		return nil, fmt.Errorf("expected '(' after procedure name at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	params := []Parameter{}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	for !p.currentTokenIs(lexer.RBRACKET) && !p.currentTokenIs(lexer.EOF) {
		if !p.currentTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected parameter name in procedure definition at line %d, col %d",
				p.currentPosition.Line, p.currentPosition.Column)
		}
		params = append(params, Parameter{Name: p.currentLiteral, Pos: p.currentPosition})

		if err := p.nextToken(); err != nil {
			return nil, err
		}

		if p.currentTokenIs(lexer.COMMA) {
			if err := p.nextToken(); err != nil {
				return nil, err
			}
		}
	}

	if !p.currentTokenIs(lexer.RBRACKET) {
		return nil, fmt.Errorf("expected ')' after procedure parameters at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if !p.currentTokenIs(lexer.IN) {
		return nil, fmt.Errorf("expected IN after procedure parameters at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	// parse procedure body
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	body := []Node{}
	for !p.currentTokenIs(lexer.END) && !p.currentTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			body = append(body, stmt)
		}

		if err := p.nextToken(); err != nil {
			return nil, err
		}
	}

	if !p.currentTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END at line %d, col %d",
			p.currentPosition.Line, p.currentPosition.Column)
	}

	return &Procedure{
		Name:       procName,
		Parameters: params,
		Body:       body,
		Pos:        pos,
	}, nil

}

func validateLoopControl(program *Program) error {
	var walk func(n Node, depth int) error
	walk = func(n Node, depth int) error {
		switch x := n.(type) {
		case *Program:
			for _, s := range x.Statements {
				if err := walk(s, depth); err != nil {
					return err
				}
			}
		case *WhileStmt:
			for _, s := range x.Body {
				if err := walk(s, depth+1); err != nil {
					return err
				}
			}
			for _, s := range x.Condition {
				if err := walk(s, depth+1); err != nil {
					return err
				}
			}
		case *IfStmt:
			for _, s := range x.ThenBranch {
				if err := walk(s, depth); err != nil {
					return err
				}
			}
			for _, s := range x.ElseBranch {
				if err := walk(s, depth); err != nil {
					return err
				}
			}
			for _, s := range x.Condition {
				if err := walk(s, depth); err != nil {
					return err
				}
			}
		case *BreakStmt:
			if depth <= 0 {
				return fmt.Errorf("BREAK can only be used inside a WHILE loop at line %d, col %d", x.Pos.Line, x.Pos.Column)
			}
		case *ContinueStmt:
			if depth <= 0 {
				return fmt.Errorf("CONTINUE can only be used inside a WHILE loop at line %d, col %d", x.Pos.Line, x.Pos.Column)
			}
		default:
			// other nodes: ok
		}
		return nil
	}

	return walk(program, 0)
}
