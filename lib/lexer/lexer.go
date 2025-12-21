package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Lexer struct {
	position Position
	reader   *bufio.Reader
	char     rune
}

type StackElement struct {
	Type     TokenType
	Value    string
	Position Position
}

func NewLexer(reader io.Reader) *Lexer {
	l := &Lexer{
		position: NewPosition(1, 1),
		reader:   bufio.NewReader(reader),
	}
	l.char, _ = l.readChar()
	return l
}

func (l *Lexer) jumpToNextLine() {
	l.position.Line++
	l.position.Column = 1
}

func (l *Lexer) jumpToNextColumn() {
	l.position.Column++
}

func (l *Lexer) readChar() (rune, error) {
	r, _, err := l.reader.ReadRune()

	if err != nil {
		if err == io.EOF {
			return 0, nil
		}

		return 0, err
	}

	return r, nil
}

func (l *Lexer) peekChar() (rune, error) {
	r, err := l.readChar()
	if err != nil {
		return 0, nil
	}

	if r != 0 {
		if err := l.reader.UnreadRune(); err != nil {
			return r, err
		}
	}

	return r, nil
}

func (l *Lexer) backup() error {
	if err := l.reader.UnreadRune(); err != nil {
		return err
	}

	// go back a column
	l.position.Column--
	return nil
}

// reads a pointer from the stack
func (l *Lexer) readPtr() (Token, error) {
	panic("Not implemented")
}

// dereferences a pointer from the stack
func (l *Lexer) readPtrDeref() (Token, error) {
	panic("Not implemented")
}

func (l *Lexer) classifyIdent(lit string) TokenType {
	// Check for boolean literals first
	if lit == "true" || lit == "false" {
		return BOOL
	}

	// Check if it's a keyword/operator
	if tokType, ok := keywords[lit]; ok {
		return tokType
	}

	// means its most likely illegal
	return ILLEGAL
}

// read identifiers
func (l *Lexer) readIdent() (string, TokenType, error) {
	var lit string

	// Start with current character
	lit = string(l.char)

	// Read all valid identifier characters
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.char = 0
				return lit, l.classifyIdent(lit), nil
			}
			return "", ILLEGAL, err
		}

		l.jumpToNextColumn()

		// Stop reading when we hit a non-letter character
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			l.char = r
			return lit, l.classifyIdent(lit), nil
		}

		lit += string(r)
	}
}

// read variable names from a sequence of chars
func (l *Lexer) readVariable() (string, TokenType, error) {
	var literal string

	literal = string(l.char)

	r, err := l.readChar()
	if err != nil {
		return literal, ILLEGAL, err
	}

	l.jumpToNextColumn()

	if !unicode.IsLetter(r) {
		return literal, ILLEGAL, nil
	}

	literal += string(r)

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.char = 0
				return literal, VARIABLE, nil
			}

			return "", ILLEGAL, err
		}

		l.jumpToNextColumn()

		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			l.char = r
			return literal, VARIABLE, nil
		}

		literal += string(r)
	}
}

// reads strings from a sequence of chars starting with "
func (l *Lexer) readSinglelineString() (string, TokenType, error) {
	var literal string
	var err error

	// skip opening quote
	l.char, err = l.readChar()
	if err != nil {
		return literal, ILLEGAL, err
	}

	l.jumpToNextColumn()

	for {
		if l.char == 0 {
			// TODO: use custom error here
			return literal, ILLEGAL, fmt.Errorf("unterminated string at line %d column %d", l.position.Line, l.position.Column)
		}

		// handle escape sequences
		if l.char == '\\' {
			l.jumpToNextColumn()
			l.char, err = l.readChar()
			if err != nil {
				return literal, ILLEGAL, err
			}

			// handle common escape sequences
			switch l.char {
			case 'n':
				literal += "\n"
			case 't':
				literal += "\t"
			case 'r':
				literal += "\r"
			case '"':
				literal += "\""
			case '\\':
				literal += "\\"
			default:
				// for any other character, just include it literally
				literal += string(l.char)
			}
			l.char, _ = l.readChar()
			l.jumpToNextColumn()
			continue
		}

		// closing quotes
		if l.char == '"' {
			break
		}

		if l.char == '\n' {
			return literal, ILLEGAL, fmt.Errorf("unterminated string before newline at line %d column %d", l.position.Line, l.position.Column)
		} else {
			l.jumpToNextColumn()
		}

		literal += string(l.char)
		l.char, _ = l.readChar()
	}

	return literal, STRING, nil
}

// reads strings from a sequence of chars starting with `
func (l *Lexer) readMultilineString() (string, TokenType, error) {
	var literal string
	var err error

	l.char, err = l.readChar()
	if err != nil {
		return literal, ILLEGAL, err
	}

	l.jumpToNextColumn()

	for {
		if l.char == 0 {
			return "", ILLEGAL, fmt.Errorf("unterminated string at line %d column %d", l.position.Line, l.position.Column)
		}

		// handle escape sequences
		if l.char == '\\' {
			l.jumpToNextColumn()
			l.char, err = l.readChar()
			if err != nil {
				return literal, ILLEGAL, err
			}

			// handle common escape sequences
			switch l.char {
			case 'n':
				literal += "\n"
			case 't':
				literal += "\t"
			case 'r':
				literal += "\r"
			case '`':
				literal += "`"
			case '\\':
				literal += "\\"
			default:
				// for any other character, just include it literally
				literal += string(l.char)
			}

			if l.char == '\n' {
				l.jumpToNextLine()
			} else {
				l.jumpToNextColumn()
			}

			l.char, _ = l.readChar()
			continue
		}

		// closing quotes
		if l.char == '`' {
			break
		}

		if l.char == '\n' {
			l.jumpToNextLine()
		} else {
			l.jumpToNextColumn()
		}

		literal += string(l.char)
		l.char, _ = l.readChar()
	}

	return literal, STRING, nil
}

// reads numbers from a sequence of runes
func (l *Lexer) readNumber() (string, TokenType, error) {
	position := l.position
	var literal string
	var tokenType TokenType = INT

	// Start with current character
	literal = string(l.char)

	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				l.char = 0
				return literal, tokenType, nil
			}

			return literal, ILLEGAL, err
		}

		l.jumpToNextColumn()

		if unicode.IsDigit(r) {
			literal += string(r)
		} else if IsDecimal(r) {
			if tokenType == INT {
				tokenType = FLOAT
				literal += string(r)
			} else {
				return literal, ILLEGAL, fmt.Errorf("unexpected decimal point at line %d column %d", position.Line, position.Column)
			}
		} else {
			l.char = r
			return literal, tokenType, nil
		}
	}
}

func (l *Lexer) NextToken() (Token, error) {
	// skip whitespace
	for l.char == ' ' || l.char == '\t' || l.char == '\r' || l.char == '\n' {
		if l.char == '\n' {
			l.jumpToNextLine()
		} else {
			l.jumpToNextColumn()
		}
		l.char, _ = l.readChar()
	}

	tok := Token{Pos: l.position}

	switch l.char {
	case 0:
		tok.Type = EOF
		return tok, nil
	case '#':
		// skip everything until end of line
		for {
			r, _ := l.readChar()
			if r == '\n' || r == 0 {
				l.char = r
				break
			}
		}
		return l.NextToken()
	case '+':
		tok.Type = OP_PLUS
		tok.Literal = string(l.char)
	case '-':
		tok.Type = OP_MINUS
		tok.Literal = string(l.char)
	case '*':
		tok.Literal = string(l.char)
		peek, err := l.peekChar()
		if err != nil {
			return Token{}, nil
		}
		// check if value after is a number, meaning we are using a pointer
		if unicode.IsDigit(peek) {
			return l.readPtr()
		}
		// else it's just a multiply
		tok.Type = OP_MULTIPLY
	case '/':
		tok.Type = OP_DIVIDE
		tok.Literal = string(l.char)
	case '^':
		tok.Type = OP_POWER
		tok.Literal = string(l.char)
	case '%':
		tok.Type = OP_MODULO
		tok.Literal = string(l.char)
	case '@':
		return l.readPtrDeref()
	case '=':
		peek, err := l.peekChar()
		if err != nil {
			return Token{}, nil
		}
		// check if we're doing equality
		if peek == '=' {
			char := l.char
			l.readChar()
			tok.Type = EQ
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = ASSIGN
			tok.Literal = string(l.char)
		}
	case '!':
		peek, err := l.peekChar()
		if err != nil {
			return Token{}, nil
		}
		// check if we're doing negation or equality check
		if peek == '=' {
			l.readChar()
			tok.Type = NEQ
			tok.Literal = "!="
		} else {
			tok.Type = OP_NOT
			tok.Literal = string(l.char)
		}
	case '>':
		peek, err := l.peekChar()
		if err != nil {
			return Token{}, nil
		}
		// check if we're doing equality
		if peek == '=' {
			l.readChar()
			tok.Type = GTE
			tok.Literal = ">="
		} else {
			tok.Type = GT
			tok.Literal = string(l.char)
		}
	case '<':
		peek, err := l.peekChar()
		if err != nil {
			return Token{}, nil
		}
		// check if we're doing equality
		if peek == '=' {
			l.readChar()
			tok.Type = LTE
			tok.Literal = "<="
		} else {
			tok.Type = LT
			tok.Literal = string(l.char)
		}
	case '&':
		peek, err := l.peekChar()
		if err != nil {
			return Token{}, nil
		}
		if peek == '&' {
			char := l.char
			l.readChar()
			tok.Type = OP_AND
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.char)
		}
	case '|':
		peek, err := l.peekChar()
		if err != nil {
			return Token{}, nil
		}
		if peek == '|' {
			char := l.char
			// consume next char
			l.readChar()
			tok.Type = OP_OR
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.char)
		}
	case '[':
		tok.Type = LBRACKET
		tok.Literal = string(l.char)
	case ']':
		tok.Type = RBRACKET
		tok.Literal = string(l.char)
	case ',':
		tok.Type = COMMA
		tok.Literal = string(l.char)
	case '"':
		lit, tokType, err := l.readSinglelineString()
		if err != nil {
			return Token{}, err
		}

		tok.Literal, tok.Type = lit, tokType
		// consume closing quote
		l.char, _ = l.readChar()
		l.jumpToNextColumn()
		return tok, nil
	case '$':
		lit, tokType, err := l.readVariable()
		if err != nil {
			return Token{}, err
		}

		tok.Literal, tok.Type = lit, tokType
	default:
		if unicode.IsLetter(l.char) {
			lit, tokType, err := l.readIdent()
			if err != nil {
				return Token{}, err
			}
			tok.Literal = lit
			tok.Type = tokType
			return tok, nil
		} else if unicode.IsDigit(l.char) {
			lit, t, err := l.readNumber()
			if err != nil {
				return Token{}, err
			}
			tok.Literal = lit
			tok.Type = t
			return tok, nil
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.char)
		}
	}

	l.char, _ = l.readChar()
	l.jumpToNextColumn()
	return tok, nil
}
