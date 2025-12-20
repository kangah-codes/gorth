package lexer

import (
	"strings"
	"testing"
)

func TestNewLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantChar rune
	}{
		{
			name:     "empty input",
			input:    "",
			wantChar: 0,
		},
		{
			name:     "single character",
			input:    "a",
			wantChar: 'a',
		},
		{
			name:     "multiple characters",
			input:    "hello",
			wantChar: 'h',
		},
		{
			name:     "starts with number",
			input:    "123",
			wantChar: '1',
		},
		{
			name:     "starts with symbol",
			input:    "+test",
			wantChar: '+',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			if lexer == nil {
				t.Fatal("NewLexer returned nil")
			}

			if lexer.char != tt.wantChar {
				t.Errorf("NewLexer() char = %c, want %c", lexer.char, tt.wantChar)
			}

			if lexer.position.Line != 1 {
				t.Errorf("NewLexer() position.Line = %d, want 1", lexer.position.Line)
			}

			if lexer.position.Column != 1 {
				t.Errorf("NewLexer() position.Column = %d, want 1", lexer.position.Column)
			}

			if lexer.reader == nil {
				t.Error("NewLexer() reader is nil")
			}
		})
	}
}

func TestJumpToNextLine(t *testing.T) {
	tests := []struct {
		name         string
		initialLine  int
		initialCol   int
		expectedLine int
		expectedCol  int
	}{
		{
			name:         "jump from line 1",
			initialLine:  1,
			initialCol:   5,
			expectedLine: 2,
			expectedCol:  1,
		},
		{
			name:         "jump from line 10",
			initialLine:  10,
			initialCol:   15,
			expectedLine: 11,
			expectedCol:  1,
		},
		{
			name:         "jump from column 0",
			initialLine:  5,
			initialCol:   0,
			expectedLine: 6,
			expectedCol:  1,
		},
		{
			name:         "jump from large line number",
			initialLine:  100,
			initialCol:   50,
			expectedLine: 101,
			expectedCol:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader("test")
			lexer := NewLexer(reader)

			// Set initial position
			lexer.position.Line = tt.initialLine
			lexer.position.Column = tt.initialCol

			// Call jumpToNextLine
			lexer.jumpToNextLine()

			// Check results
			if lexer.position.Line != tt.expectedLine {
				t.Errorf("jumpToNextLine() Line = %d, want %d", lexer.position.Line, tt.expectedLine)
			}

			if lexer.position.Column != tt.expectedCol {
				t.Errorf("jumpToNextLine() Column = %d, want %d", lexer.position.Column, tt.expectedCol)
			}
		})
	}
}
func TestReadChar(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantRunes  []rune
		wantErrors []bool
	}{
		{
			name:       "empty input",
			input:      "",
			wantRunes:  []rune{0},
			wantErrors: []bool{false},
		},
		{
			name:       "single character",
			input:      "a",
			wantRunes:  []rune{'a', 0},
			wantErrors: []bool{false, false},
		},
		{
			name:       "multiple characters",
			input:      "abc",
			wantRunes:  []rune{'a', 'b', 'c', 0},
			wantErrors: []bool{false, false, false, false},
		},
		{
			name:       "unicode characters",
			input:      "αβγ",
			wantRunes:  []rune{'α', 'β', 'γ', 0},
			wantErrors: []bool{false, false, false, false},
		},
		{
			name:       "mixed characters and numbers",
			input:      "a1b2",
			wantRunes:  []rune{'a', '1', 'b', '2', 0},
			wantErrors: []bool{false, false, false, false, false},
		},
		{
			name:       "special characters",
			input:      "!@#",
			wantRunes:  []rune{'!', '@', '#', 0},
			wantErrors: []bool{false, false, false, false},
		},
		{
			name:       "newline character",
			input:      "a\nb",
			wantRunes:  []rune{'a', '\n', 'b', 0},
			wantErrors: []bool{false, false, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			// First check the current character already read by NewLexer
			if len(tt.wantRunes) > 0 {
				if lexer.char != tt.wantRunes[0] {
					t.Errorf("NewLexer() char = %c (U+%04X), want %c (U+%04X)", lexer.char, lexer.char, tt.wantRunes[0], tt.wantRunes[0])
				}
			}

			// Then read subsequent characters
			for i := 1; i < len(tt.wantRunes); i++ {
				gotRune, err := lexer.readChar()

				if (err != nil) != tt.wantErrors[i] {
					t.Errorf("readChar() error = %v, wantError %v", err, tt.wantErrors[i])
					return
				}

				if gotRune != tt.wantRunes[i] {
					t.Errorf("readChar() = %c (U+%04X), want %c (U+%04X)", gotRune, gotRune, tt.wantRunes[i], tt.wantRunes[i])
				}
			}
		})
	}
}

func TestPeekChar(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedPeek    rune
		expectedCurrent rune
	}{
		{
			name:            "peek at EOF",
			input:           "",
			expectedPeek:    0,
			expectedCurrent: 0,
		},
		{
			name:            "peek single character",
			input:           "a",
			expectedPeek:    0,
			expectedCurrent: 'a',
		},
		{
			name:            "peek second character",
			input:           "ab",
			expectedPeek:    'b',
			expectedCurrent: 'a',
		},
		{
			name:            "peek multiple characters",
			input:           "hello",
			expectedPeek:    'e',
			expectedCurrent: 'h',
		},
		{
			name:            "peek number",
			input:           "123",
			expectedPeek:    '2',
			expectedCurrent: '1',
		},
		{
			name:            "peek symbol",
			input:           "+=",
			expectedPeek:    '=',
			expectedCurrent: '+',
		},
		{
			name:            "peek unicode character",
			input:           "αβ",
			expectedPeek:    'β',
			expectedCurrent: 'α',
		},
		{
			name:            "peek newline",
			input:           "a\n",
			expectedPeek:    '\n',
			expectedCurrent: 'a',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			// Store original position and character
			originalPos := lexer.position

			// Peek the next character
			peeked := lexer.peekChar()

			// Verify peek result
			if peeked != tt.expectedPeek {
				t.Errorf("peekChar() = %c (U+%04X), want %c (U+%04X)", peeked, peeked, tt.expectedPeek, tt.expectedPeek)
			}

			// Verify current character unchanged
			if lexer.char != tt.expectedCurrent {
				t.Errorf("after peekChar() current char = %c, want %c", lexer.char, tt.expectedCurrent)
			}

			// Verify position unchanged
			if lexer.position != originalPos {
				t.Errorf("after peekChar() position = %v, want %v", lexer.position, originalPos)
			}

			// Verify that the next readChar() returns the same character as peek
			if tt.expectedPeek != 0 {
				nextChar, err := lexer.readChar()
				if err != nil {
					t.Errorf("readChar() after peek failed: %v", err)
				}
				if nextChar != tt.expectedPeek {
					t.Errorf("readChar() after peek = %c, want %c", nextChar, tt.expectedPeek)
				}
			}
		})
	}
}

func TestBackup(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		setupSteps     int // number of characters to read before backup
		expectedChar   rune
		expectedColumn int
	}{
		{
			name:           "backup after reading one character",
			input:          "ab",
			setupSteps:     1,
			expectedChar:   'b',
			expectedColumn: 1,
		},
		{
			name:           "backup after reading multiple characters",
			input:          "hello",
			setupSteps:     3,
			expectedChar:   'l',
			expectedColumn: 3,
		},
		{
			name:           "backup with numbers",
			input:          "123",
			setupSteps:     2,
			expectedChar:   '3',
			expectedColumn: 2,
		},
		{
			name:           "backup with symbols",
			input:          "+=*",
			setupSteps:     1,
			expectedChar:   '=',
			expectedColumn: 1,
		},
		{
			name:           "backup with unicode",
			input:          "αβγ",
			setupSteps:     2,
			expectedChar:   'γ',
			expectedColumn: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			// Read setupSteps characters to advance position
			for i := 0; i < tt.setupSteps; i++ {
				lexer.char, _ = lexer.readChar()
			}

			// Store the current column before backup
			colBeforeBackup := lexer.position.Column

			// Call backup
			lexer.backup()

			// Verify position was decremented by 1
			if lexer.position.Column != colBeforeBackup-1 {
				t.Errorf("backup() Column = %d, want %d", lexer.position.Column, colBeforeBackup-1)
			}

			// Verify we can read the backed up character
			nextChar, err := lexer.readChar()
			if err != nil {
				t.Errorf("readChar() after backup failed: %v", err)
			}
			if nextChar != tt.expectedChar {
				t.Errorf("readChar() after backup = %c, want %c", nextChar, tt.expectedChar)
			}

			// Verify line number is unchanged
			if lexer.position.Line != 1 {
				t.Errorf("backup() changed Line to %d, want 1", lexer.position.Line)
			}
		})
	}
}

func TestClassifyIdent(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedToken TokenType
	}{
		// Boolean literals
		{
			name:          "classify false as a boolean",
			input:         "false",
			expectedToken: BOOL,
		},
		{
			name:          "classify true as a boolean",
			input:         "true",
			expectedToken: BOOL,
		},

		// Keywords
		{
			name:          "classify const keyword",
			input:         "const",
			expectedToken: CONST,
		},
		{
			name:          "classify proc keyword",
			input:         "proc",
			expectedToken: PROC,
		},
		{
			name:          "classify endproc keyword",
			input:         "endproc",
			expectedToken: ENDPROC,
		},
		{
			name:          "classify in keyword",
			input:         "in",
			expectedToken: IN,
		},
		{
			name:          "classify return keyword",
			input:         "return",
			expectedToken: RETURN,
		},

		// Type keywords
		{
			name:          "classify int type",
			input:         "int",
			expectedToken: TYPE_INT,
		},
		{
			name:          "classify str type",
			input:         "str",
			expectedToken: TYPE_STR,
		},
		{
			name:          "classify bool type",
			input:         "bool",
			expectedToken: TYPE_BOOL,
		},
		{
			name:          "classify float type",
			input:         "float",
			expectedToken: TYPE_FLOAT,
		},
		{
			name:          "classify ptr type",
			input:         "ptr",
			expectedToken: TYPE_PTR,
		},
		{
			name:          "classify arr type",
			input:         "arr",
			expectedToken: TYPE_ARR,
		},
		// Stack operations
		{
			name:          "classify drop operation",
			input:         "drop",
			expectedToken: OP_DROP,
		},
		{
			name:          "classify swap operation",
			input:         "swap",
			expectedToken: OP_SWAP,
		},
		{
			name:          "classify dup operation",
			input:         "dup",
			expectedToken: OP_DUP,
		},
		{
			name:          "classify over operation",
			input:         "over",
			expectedToken: OP_OVER,
		},
		{
			name:          "classify rot operation",
			input:         "rot",
			expectedToken: OP_ROT,
		},
		{
			name:          "classify del operation",
			input:         "del",
			expectedToken: OP_DEL,
		},
		{
			name:          "classify inc operation",
			input:         "inc",
			expectedToken: OP_INC,
		},
		{
			name:          "classify dec operation",
			input:         "dec",
			expectedToken: OP_DEC,
		},
		// I/O operations
		{
			name:          "classify print operation",
			input:         "print",
			expectedToken: OP_PRINT,
		},
		{
			name:          "classify println operation",
			input:         "println",
			expectedToken: OP_PRINTLN,
		},
		{
			name:          "classify dump operation",
			input:         "dump",
			expectedToken: OP_DUMP,
		},
		// Invalid identifiers
		{
			name:          "unknown identifier returns illegal",
			input:         "unknown",
			expectedToken: ILLEGAL,
		},
		{
			name:          "random string returns illegal",
			input:         "randomString",
			expectedToken: ILLEGAL,
		},
		{
			name:          "mixed case boolean returns illegal",
			input:         "True",
			expectedToken: ILLEGAL,
		},
		{
			name:          "mixed case boolean false returns illegal",
			input:         "False",
			expectedToken: ILLEGAL,
		},
		{
			name:          "empty string returns illegal",
			input:         "",
			expectedToken: ILLEGAL,
		},
		{
			name:          "single character returns illegal",
			input:         "x",
			expectedToken: ILLEGAL,
		},
		{
			name:          "number-like string returns illegal",
			input:         "123abc",
			expectedToken: ILLEGAL,
		},
	}

	testedKeywords := make(map[string]bool)
	for _, tt := range tests {
		if tt.expectedToken != ILLEGAL {
			testedKeywords[tt.input] = true
		}
	}

	// Check that all non-operator keywords are tested
	// Operators like +, -, *, etc. are tested separately in NextToken tests
	missingTests := []string{}
	for kw := range keywords {
		if len(kw) <= 2 && (kw == "+" || kw == "-" || kw == "*" || kw == "/" ||
			kw == "^" || kw == "%" || kw == "==" || kw == "!=" ||
			kw == ">" || kw == "<" || kw == ">=" || kw == "<=" ||
			kw == "&&" || kw == "||" || kw == "!" || kw == "=") {
			continue
		}

		if !testedKeywords[kw] {
			missingTests = append(missingTests, kw)
		}
	}

	if len(missingTests) > 0 {
		t.Errorf("Missing test coverage for keywords: %v\nPlease add test cases for these keywords to ensure proper coverage", missingTests)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)
			tokType := lexer.classifyIdent(tt.input)

			if tokType != tt.expectedToken {
				t.Errorf("classifyIdent() expected %s got %s", tt.expectedToken, tokType)
			}
		})
	}
}

func TestReadIdent(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedLiteral   string
		expectedTokenType TokenType
	}{
		{
			name:              "return correct literal and tokentype for boolean input",
			input:             "true",
			expectedLiteral:   "true",
			expectedTokenType: BOOL,
		},
		{
			name:              "return correct literal and tokentype for keyword input",
			input:             "dump",
			expectedLiteral:   "dump",
			expectedTokenType: OP_DUMP,
		},
		{
			name:              "return correct literal and tokentype for illegal input",
			input:             "mala",
			expectedLiteral:   "mala",
			expectedTokenType: ILLEGAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)
			lit, tokType := lexer.readIdent()

			if lit != tt.expectedLiteral {
				t.Errorf("readIdent() literal expected %s got %s", tt.expectedLiteral, lit)
			}

			if tokType != tt.expectedTokenType {
				t.Errorf("readIdent() tokenType expected %s got %s", tt.expectedTokenType, tokType)
			}
		})
	}
}

func TestReadVariable(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedLiteral   string
		expectedTokenType TokenType
		shouldPanic       bool
	}{
		{
			name:              "valid simple variable",
			input:             "$variable",
			expectedLiteral:   "$variable",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "valid variable with numbers",
			input:             "$var123",
			expectedLiteral:   "$var123",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "valid single letter variable",
			input:             "$a",
			expectedLiteral:   "$a",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "valid variable with mixed case",
			input:             "$MyVariable",
			expectedLiteral:   "$MyVariable",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "valid variable ending with number",
			input:             "$test123",
			expectedLiteral:   "$test123",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "invalid variable starting with number",
			input:             "$1variable",
			expectedLiteral:   "$",
			expectedTokenType: ILLEGAL,
		},
		{
			name:              "invalid variable starting with symbol",
			input:             "$@invalid",
			expectedLiteral:   "$",
			expectedTokenType: ILLEGAL,
		},
		{
			name:              "invalid variable with only dollar sign",
			input:             "$ ",
			expectedLiteral:   "$",
			expectedTokenType: ILLEGAL,
		},
		{
			name:              "valid variable followed by space",
			input:             "$var ",
			expectedLiteral:   "$var",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "valid variable followed by symbol",
			input:             "$var+",
			expectedLiteral:   "$var",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "valid variable at end of input",
			input:             "$test",
			expectedLiteral:   "$test",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "invalid variable with underscore",
			input:             "$var_name",
			expectedLiteral:   "$var",
			expectedTokenType: VARIABLE,
		},
		{
			name:              "valid unicode variable",
			input:             "$αβγ",
			expectedLiteral:   "$αβγ",
			expectedTokenType: VARIABLE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("readVariable() expected to panic but didn't")
					}
				}()
			}

			lit, tokType := lexer.readVariable()

			if !tt.shouldPanic {
				if lit != tt.expectedLiteral {
					t.Errorf("readVariable() literal = %s, want %s", lit, tt.expectedLiteral)
				}

				if tokType != tt.expectedTokenType {
					t.Errorf("readVariable() tokenType = %s, want %s", tokType, tt.expectedTokenType)
				}

				// verify the lexer's current character is set correctly
				if tt.expectedTokenType == VARIABLE && len(tt.input) > len(tt.expectedLiteral) {
					expectedNextChar := rune(tt.input[len(tt.expectedLiteral)])
					if lexer.char != expectedNextChar {
						t.Errorf("readVariable() left lexer.char = %c, want %c", lexer.char, expectedNextChar)
					}
				}
			}
		})
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedLiteral   string
		expectedTokenType TokenType
		shouldPanic       bool
	}{
		{
			name:              "return correct literal and tokentype for valid variable input",
			input:             `"Hello, world"`,
			expectedLiteral:   "Hello, world",
			expectedTokenType: STRING,
		},
		{
			name:              "return correct literal and tokentype for valid variable input",
			input:             `1`,
			expectedLiteral:   "Hello, world",
			expectedTokenType: STRING,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("readVariable() expected to panic but didn't")
					}
				}()
			}

			lit, tokType := lexer.readString()

			if !tt.shouldPanic {
				if lit != tt.expectedLiteral {
					t.Errorf("readString() literal expected %s got %s", tt.expectedLiteral, lit)
				}

				if tokType != tt.expectedTokenType {
					t.Errorf("readString() tokenType expected %s got %s", tt.expectedTokenType, tokType)
				}
			}
		})
	}
}
