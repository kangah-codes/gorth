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
	return &Lexer{
		position: NewPosition(1, 1),
		reader:   bufio.NewReader(reader),
	}
}

func (l *Lexer) jumpToNextLine() {
	l.position.Line++
	l.position.Column = 0
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
		l.backup()
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

// skips whitespace in buffer
func (l *Lexer) skipWhitespace() {
	// jump ahead if current rune is whitespace
	l.position.Column++
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

	// will be variable soon
	return ILLEGAL
}

func (l *Lexer) readIdent() (string, TokenType) {
	var lit string

	// Read all valid identifier characters
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit, l.classifyIdent(lit)
			}
			panic(err)
		}

		l.position.Column++

		// Stop reading when we hit a non-letter character
		if !unicode.IsLetter(r) {
			l.backup()
			return lit, l.classifyIdent(lit)
		}

		lit += string(r)
	}
}

// reads numbers from a sequence of runes
func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position
	var literal string
	var tokenType TokenType = INT

	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return literal, tokenType
			}

			// TODO: use better error
			panic(err)
		}

		if unicode.IsSymbol(r) || unicode.IsPunct(r) && r != '=' && len(literal) > 0 {
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
			l.backup()
			break
		}
	}

	return literal, tokenType
}

func (l *Lexer) NextToken() Token {
	tok := Token{Pos: l.position}

	switch l.char {
	case 0:
		tok.Type = EOF
	case '\n':
		tok.Type = EOF
		l.readChar()
	case '#':
		l.jumpToNextLine()
	case '+':
		tok.Type = PLUS
		tok.Literal = string(l.char)
	case '-':
		tok.Type = MINUS
		tok.Literal = string(l.char)
	case '*':
		tok.Literal = string(l.char)
		// check if value after is a number, meaning we are using a pointer
		if unicode.IsDigit(l.peekChar()) {
			return l.readPtr()
		}
		// else it's just a multiple
		tok.Type = MULTIPLY
	case '/':
		tok.Type = DIVIDE
		tok.Literal = string(l.char)
	case '^':
		tok.Type = POWER
		tok.Literal = string(l.char)
	case '%':
		tok.Type = MODULO
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
			tok.Type = NOT
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
			tok.Type = AND
			tok.Literal = string(char) + string(l.char)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.char)
		}
	case '|':
		if l.peekChar() == '|' {
			char := l.char
			l.readChar()
			tok.Type = OR
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
	case '"':
		tok.Type = STRING
		tok.Literal = string(l.char)
	default:
		if unicode.IsLetter(l.char) {
			l, t := l.readIdent()
			tok.Literal = l
			tok.Type = t
		} else if unicode.IsDigit(l.char) {
			l, t := l.readNumber()
			tok.Literal = l
			tok.Type = t
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.char)
		}
	}

	l.readChar()
	return tok
}
