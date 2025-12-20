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

func (l *Lexer) peekChar() rune {
	r, err := l.readChar()
	if err != nil {
		return 0
	}

	if r != 0 {
		if err := l.reader.UnreadRune(); err != nil {
			// TODO: use gorth error here
			panic(err)
		}
	}

	return r
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		// TODO: change this from panic
		panic(err)
	}

	// go back a column
	l.position.Column--
}

// reads a pointer from the stack
func (l *Lexer) readPtr() Token {
	panic("Not implemented")
}

// dereferences a pointer from the stack
func (l *Lexer) readPtrDeref() Token {
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
func (l *Lexer) readIdent() (string, TokenType) {
	var lit string

	// Start with current character
	lit = string(l.char)

	// Read all valid identifier characters
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.char = 0
				return lit, l.classifyIdent(lit)
			}
			panic(err)
		}

		l.position.Column++

		// Stop reading when we hit a non-letter character
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			l.char = r
			return lit, l.classifyIdent(lit)
		}

		lit += string(r)
	}
}

// read variable names from a sequence of chars
func (l *Lexer) readVariable() (string, TokenType) {
	var literal string

	literal = string(l.char)

	r, _ := l.readChar()
	l.position.Column++

	fmt.Println(literal, r)

	if !unicode.IsLetter(r) {
		return literal, ILLEGAL
	}

	literal += string(r)

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.char = 0
				return literal, VARIABLE
			}

			panic(err)
		}

		l.position.Column++

		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			l.char = r
			return literal, VARIABLE
		}

		literal += string(r)
	}
}

// reads strings from a sequence of chars starting with "
func (l *Lexer) readString() (string, TokenType) {
	var literal string

	// skip opening quote
	l.char, _ = l.readChar()
	l.position.Column++

	for {
		if l.char == 0 {
			panic(fmt.Errorf("unterminated string at line %d column %d", l.position.Line, l.position.Column))
		}

		// closing quotes
		if l.char == '"' {
			break
		}

		if l.char == '\n' {
			// TODO: language design, should we allow multiline strings like these? maybe not
			// l.jumpToNextLine()
			panic(fmt.Errorf("unterminated string before newline at line %d column %d", l.position.Line, l.position.Column))
		} else {
			l.position.Column++
		}

		literal += string(l.char)
		l.char, _ = l.readChar()
	}

	return literal, STRING
}

// reads numbers from a sequence of runes
func (l *Lexer) readNumber() (string, TokenType) {
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
				return literal, tokenType
			}

			// TODO: use better error
			panic(err)
		}

		l.position.Column++

		if unicode.IsSymbol(r) || unicode.IsPunct(r) && r != '.' && len(literal) > 0 {
			panic(fmt.Errorf("error: invalid token %v at line %v col %v", string(r), l.position.Line, l.position.Column))
		}

		if unicode.IsDigit(r) {
			literal += string(r)
		} else if IsDecimal(r) {
			if tokenType == INT {
				tokenType = FLOAT
				literal += string(r)
			} else {
				panic(fmt.Errorf("unexpected decimal point at line %d column %d", position.Line, position.Column))
			}
		} else {
			l.char = r
			return literal, tokenType
		}
	}
}

func (l *Lexer) NextToken() Token {
	// skip whitespace
	for l.char == ' ' || l.char == '\t' || l.char == '\r' || l.char == '\n' {
		if l.char == '\n' {
			l.jumpToNextLine()
		} else {
			l.position.Column++
		}
		l.char, _ = l.readChar()
	}

	tok := Token{Pos: l.position}

	switch l.char {
	case 0:
		tok.Type = EOF
		return tok
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
		// check if value after is a number, meaning we are using a pointer
		if unicode.IsDigit(l.peekChar()) {
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
		// check if we're doing equality
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok.Type = EQ
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = ASSIGN
			tok.Literal = string(l.char)
		}
	case '!':
		// check if we're doing negation or equality check
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok.Type = NEQ
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = OP_NOT
			tok.Literal = string(l.char)
		}
	case '>':
		// check if we're doing equality
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok.Type = GTE
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = GT
			tok.Literal = string(l.char)
		}
	case '<':
		// check if we're doing equality
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok.Type = LTE
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = LT
			tok.Literal = string(l.char)
		}
	case '&':
		if l.peekChar() == '&' {
			char := l.char
			l.readChar()
			tok.Type = OP_AND
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.char)
		}
	case '|':
		if l.peekChar() == '|' {
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
		tok.Literal, tok.Type = l.readString()
		// consume closing quote
		l.char, _ = l.readChar()
		l.position.Column++
		return tok
	case '$':
		tok.Literal, tok.Type = l.readVariable()
	default:
		if unicode.IsLetter(l.char) {
			lit, t := l.readIdent()
			tok.Literal = lit
			tok.Type = t
			return tok
		} else if unicode.IsDigit(l.char) {
			lit, t := l.readNumber()
			tok.Literal = lit
			tok.Type = t
			return tok
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.char)
		}
	}

	l.char, _ = l.readChar()
	l.position.Column++
	return tok
}
