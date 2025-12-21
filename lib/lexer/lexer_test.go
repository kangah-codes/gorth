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
			peeked, err := lexer.peekChar()

			if err != nil {
				t.Errorf("peekChar() threw unexpected error: %v", err)
			}

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
			lit, tokType, err := lexer.readIdent()

			if err != nil {
				t.Errorf("readIdent() returned error unexpectedly %v", err)
			}

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

			lit, tokType, err := lexer.readVariable()

			if err != nil {
				t.Errorf("readVariable() returned unexpected error %v", err)
			}

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
		})
	}
}

func TestReadSinglelineString(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedLiteral   string
		expectedTokenType TokenType
		wantError         bool
		errorContains     string
	}{
		{
			name:              "valid empty string",
			input:             `""`,
			expectedLiteral:   "",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid simple string",
			input:             `"hello"`,
			expectedLiteral:   "hello",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with spaces",
			input:             `"hello world"`,
			expectedLiteral:   "hello world",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with numbers",
			input:             `"test123"`,
			expectedLiteral:   "test123",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with symbols",
			input:             `"hello!@#$%^&*()"`,
			expectedLiteral:   "hello!@#$%^&*()",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with unicode",
			input:             `"αβγδε"`,
			expectedLiteral:   "αβγδε",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with tabs",
			input:             "\"hello\tworld\"",
			expectedLiteral:   "hello\tworld",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "unterminated string at EOF",
			input:             `"hello`,
			expectedLiteral:   "hello",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string at line",
		},
		{
			name:              "unterminated string with content at EOF",
			input:             `"hello world`,
			expectedLiteral:   "hello world",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string at line",
		},
		{
			name:              "string with newline",
			input:             "\"hello\nworld\"",
			expectedLiteral:   "hello",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string before newline",
		},
		{
			name:              "string ending with newline",
			input:             "\"hello\n",
			expectedLiteral:   "hello",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string before newline",
		},
		{
			name:              "empty string with newline",
			input:             "\"\n",
			expectedLiteral:   "",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string before newline",
		},
		{
			name:              "valid string followed by text",
			input:             `"hello"world`,
			expectedLiteral:   "hello",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "string with single character",
			input:             `"a"`,
			expectedLiteral:   "a",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "string with multiple quotes inside not supported",
			input:             `"say "hello""`,
			expectedLiteral:   "say ",
			expectedTokenType: STRING,
			wantError:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			lit, tokType, err := lexer.readSinglelineString()

			if tt.wantError {
				if err == nil {
					t.Errorf("readSinglelineString() expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("readSinglelineString() error = %v, want error containing %s", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("readSinglelineString() unexpected error = %v", err)
				}
			}

			if lit != tt.expectedLiteral {
				t.Errorf("readSinglelineString() literal = %s, want %s", lit, tt.expectedLiteral)
			}

			if tokType != tt.expectedTokenType {
				t.Errorf("readSinglelineString() tokenType = %s, want %s", tokType, tt.expectedTokenType)
			}
		})
	}
}

func TestReadMultilineString(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedLiteral   string
		expectedTokenType TokenType
		wantError         bool
		errorContains     string
	}{
		{
			name:              "valid empty string",
			input:             "``",
			expectedLiteral:   "",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid simple string",
			input:             "`hello`",
			expectedLiteral:   "hello",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with spaces",
			input:             "`hello world`",
			expectedLiteral:   "hello world",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with numbers",
			input:             "`test123`",
			expectedLiteral:   "test123",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with symbols",
			input:             "`hello!@#$%^&*()`",
			expectedLiteral:   "hello!@#$%^&*()",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with unicode",
			input:             "`αβγδε`",
			expectedLiteral:   "αβγδε",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid string with tabs",
			input:             "`hello\tworld`",
			expectedLiteral:   "hello\tworld",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid multiline string",
			input:             "`hello\nworld`",
			expectedLiteral:   "hello\nworld",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid multiline string with multiple newlines",
			input:             "`line1\nline2\nline3`",
			expectedLiteral:   "line1\nline2\nline3",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "valid multiline string with carriage returns",
			input:             "`line1\r\nline2`",
			expectedLiteral:   "line1\r\nline2",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "unterminated string at EOF",
			input:             "`hello",
			expectedLiteral:   "",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string at line",
		},
		{
			name:              "unterminated multiline string with content at EOF",
			input:             "`hello\nworld",
			expectedLiteral:   "",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string at line",
		},
		{
			name:              "unterminated empty string at EOF",
			input:             "`",
			expectedLiteral:   "",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unterminated string at line",
		},
		{
			name:              "valid string followed by text",
			input:             "`hello`world",
			expectedLiteral:   "hello",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "string with single character",
			input:             "`a`",
			expectedLiteral:   "a",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "string with backtick inside (ends at first backtick)",
			input:             "`say `hello``",
			expectedLiteral:   "say ",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "multiline string with empty lines",
			input:             "`line1\n\nline3`",
			expectedLiteral:   "line1\n\nline3",
			expectedTokenType: STRING,
			wantError:         false,
		},
		{
			name:              "string with only newlines",
			input:             "`\n\n\n`",
			expectedLiteral:   "\n\n\n",
			expectedTokenType: STRING,
			wantError:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			lit, tokType, err := lexer.readMultilineString()

			if tt.wantError {
				if err == nil {
					t.Errorf("readMultilineString() expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("readMultilineString() error = %v, want error containing %s", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("readMultilineString() unexpected error = %v", err)
				}
			}

			if lit != tt.expectedLiteral {
				t.Errorf("readMultilineString() literal = %s, want %s", lit, tt.expectedLiteral)
			}

			if tokType != tt.expectedTokenType {
				t.Errorf("readMultilineString() tokenType = %s, want %s", tokType, tt.expectedTokenType)
			}
		})
	}
}
func TestReadNumber(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedLiteral   string
		expectedTokenType TokenType
		wantError         bool
		errorContains     string
	}{
		{
			name:              "single digit integer",
			input:             "5",
			expectedLiteral:   "5",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "multiple digit integer",
			input:             "123",
			expectedLiteral:   "123",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "integer followed by space",
			input:             "42 ",
			expectedLiteral:   "42",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "integer followed by operator",
			input:             "99+",
			expectedLiteral:   "99",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "simple float",
			input:             "3.14",
			expectedLiteral:   "3.14",
			expectedTokenType: FLOAT,
			wantError:         false,
		},
		{
			name:              "float starting with zero",
			input:             "0.5",
			expectedLiteral:   "0.5",
			expectedTokenType: FLOAT,
			wantError:         false,
		},
		{
			name:              "float with multiple decimal places",
			input:             "123.456789",
			expectedLiteral:   "123.456789",
			expectedTokenType: FLOAT,
			wantError:         false,
		},
		{
			name:              "float followed by space",
			input:             "2.5 ",
			expectedLiteral:   "2.5",
			expectedTokenType: FLOAT,
			wantError:         false,
		},
		{
			name:              "integer at end of input",
			input:             "999",
			expectedLiteral:   "999",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "float at end of input",
			input:             "1.618",
			expectedLiteral:   "1.618",
			expectedTokenType: FLOAT,
			wantError:         false,
		},
		{
			name:              "zero integer",
			input:             "0",
			expectedLiteral:   "0",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "number with multiple decimal points",
			input:             "1.2.3",
			expectedLiteral:   "1.2",
			expectedTokenType: ILLEGAL,
			wantError:         true,
			errorContains:     "unexpected decimal point",
		},
		{
			name:              "number with symbol after digit",
			input:             "123@",
			expectedLiteral:   "123",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "number with punctuation after digit",
			input:             "456!",
			expectedLiteral:   "456",
			expectedTokenType: INT,
			wantError:         false,
		},
		{
			name:              "large integer",
			input:             "9876543210",
			expectedLiteral:   "9876543210",
			expectedTokenType: INT,
			wantError:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			lit, tokType, err := lexer.readNumber()

			if tt.wantError {
				if err == nil {
					t.Errorf("readNumber() expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("readNumber() error = %v, want error containing %s", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("readNumber() unexpected error = %v", err)
				}
			}

			if lit != tt.expectedLiteral {
				t.Errorf("readNumber() literal = %s, want %s", lit, tt.expectedLiteral)
			}

			if tokType != tt.expectedTokenType {
				t.Errorf("readNumber() tokenType = %s, want %s", tokType, tt.expectedTokenType)
			}
		})
	}
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []Token
	}{
		{
			name:  "EOF token",
			input: "",
			expectedTokens: []Token{
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 1}},
			},
		},
		{
			name:  "simple arithmetic operators",
			input: "+ - * / ^ %",
			expectedTokens: []Token{
				{Type: OP_PLUS, Literal: "+", Pos: Position{Line: 1, Column: 1}},
				{Type: OP_MINUS, Literal: "-", Pos: Position{Line: 1, Column: 3}},
				{Type: OP_MULTIPLY, Literal: "*", Pos: Position{Line: 1, Column: 5}},
				{Type: OP_DIVIDE, Literal: "/", Pos: Position{Line: 1, Column: 7}},
				{Type: OP_POWER, Literal: "^", Pos: Position{Line: 1, Column: 9}},
				{Type: OP_MODULO, Literal: "%", Pos: Position{Line: 1, Column: 11}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 12}},
			},
		},
		{
			name:  "comparison operators",
			input: "== != > < >= <=",
			expectedTokens: []Token{
				{Type: EQ, Literal: "==", Pos: Position{Line: 1, Column: 1}},
				{Type: NEQ, Literal: "!=", Pos: Position{Line: 1, Column: 4}},
				{Type: GT, Literal: ">", Pos: Position{Line: 1, Column: 7}},
				{Type: LT, Literal: "<", Pos: Position{Line: 1, Column: 9}},
				{Type: GTE, Literal: ">=", Pos: Position{Line: 1, Column: 11}},
				{Type: LTE, Literal: "<=", Pos: Position{Line: 1, Column: 14}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 16}},
			},
		},
		{
			name:  "logical operators",
			input: "&& || !",
			expectedTokens: []Token{
				{Type: OP_AND, Literal: "&&", Pos: Position{Line: 1, Column: 1}},
				{Type: OP_OR, Literal: "||", Pos: Position{Line: 1, Column: 4}},
				{Type: OP_NOT, Literal: "!", Pos: Position{Line: 1, Column: 7}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 8}},
			},
		},
		{
			name:  "assignment operator",
			input: "=",
			expectedTokens: []Token{
				{Type: ASSIGN, Literal: "=", Pos: Position{Line: 1, Column: 1}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 2}},
			},
		},
		{
			name:  "brackets and comma",
			input: "[ ] ,",
			expectedTokens: []Token{
				{Type: LBRACKET, Literal: "[", Pos: Position{Line: 1, Column: 1}},
				{Type: RBRACKET, Literal: "]", Pos: Position{Line: 1, Column: 3}},
				{Type: COMMA, Literal: ",", Pos: Position{Line: 1, Column: 5}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 6}},
			},
		},
		{
			name:  "string literals",
			input: `"hello" "world"`,
			expectedTokens: []Token{
				{Type: STRING, Literal: "hello", Pos: Position{Line: 1, Column: 1}},
				{Type: STRING, Literal: "world", Pos: Position{Line: 1, Column: 9}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 16}},
			},
		},
		{
			name:  "variable",
			input: "$myvar",
			expectedTokens: []Token{
				{Type: VARIABLE, Literal: "$myvar", Pos: Position{Line: 1, Column: 1}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 7}},
			},
		},
		{
			name:  "identifiers and keywords",
			input: "true false proc endproc",
			expectedTokens: []Token{
				{Type: BOOL, Literal: "true", Pos: Position{Line: 1, Column: 1}},
				{Type: BOOL, Literal: "false", Pos: Position{Line: 1, Column: 6}},
				{Type: PROC, Literal: "proc", Pos: Position{Line: 1, Column: 12}},
				{Type: ENDPROC, Literal: "endproc", Pos: Position{Line: 1, Column: 17}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 24}},
			},
		},
		{
			name:  "numbers",
			input: "123 456.789",
			expectedTokens: []Token{
				{Type: INT, Literal: "123", Pos: Position{Line: 1, Column: 1}},
				{Type: FLOAT, Literal: "456.789", Pos: Position{Line: 1, Column: 5}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 12}},
			},
		},
		{
			name:  "comments are skipped",
			input: "# this is a comment\ntrue",
			expectedTokens: []Token{
				{Type: BOOL, Literal: "true", Pos: Position{Line: 2, Column: 1}},
				{Type: EOF, Literal: "", Pos: Position{Line: 2, Column: 5}},
			},
		},
		{
			name:  "whitespace handling",
			input: "   \t\r\n   +   \n   -   ",
			expectedTokens: []Token{
				{Type: OP_PLUS, Literal: "+", Pos: Position{Line: 2, Column: 4}},
				{Type: OP_MINUS, Literal: "-", Pos: Position{Line: 3, Column: 4}},
				{Type: EOF, Literal: "", Pos: Position{Line: 3, Column: 8}},
			},
		},
		{
			name:  "single ampersand (illegal)",
			input: "&",
			expectedTokens: []Token{
				{Type: ILLEGAL, Literal: "&", Pos: Position{Line: 1, Column: 1}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 2}},
			},
		},
		{
			name:  "single pipe (illegal)",
			input: "|",
			expectedTokens: []Token{
				{Type: ILLEGAL, Literal: "|", Pos: Position{Line: 1, Column: 1}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 2}},
			},
		},
		{
			name:  "illegal characters",
			input: "~` £",
			expectedTokens: []Token{
				{Type: ILLEGAL, Literal: "~", Pos: Position{Line: 1, Column: 1}},
				{Type: ILLEGAL, Literal: "`", Pos: Position{Line: 1, Column: 2}},
				{Type: ILLEGAL, Literal: "£", Pos: Position{Line: 1, Column: 4}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 5}},
			},
		},
		{
			name:  "complex expression",
			input: "$x $y + 42 ==",
			expectedTokens: []Token{
				{Type: VARIABLE, Literal: "$x", Pos: Position{Line: 1, Column: 1}},
				{Type: VARIABLE, Literal: "$y", Pos: Position{Line: 1, Column: 4}},
				{Type: OP_PLUS, Literal: "+", Pos: Position{Line: 1, Column: 7}},
				{Type: INT, Literal: "42", Pos: Position{Line: 1, Column: 9}},
				{Type: EQ, Literal: "==", Pos: Position{Line: 1, Column: 12}},
				{Type: EOF, Literal: "", Pos: Position{Line: 1, Column: 13}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			for i, expected := range tt.expectedTokens {
				token, err := lexer.NextToken()
				if err != nil {
					t.Errorf("NextToken() returned unexpected error: %v", err)
					continue
				}

				if token.Type != expected.Type {
					t.Errorf("Token %v: Type = %s, want %s", tt.expectedTokens[i], token.Type, expected.Type)
				}

				if token.Literal != expected.Literal {
					t.Errorf("Token %v: Literal = %s, want %s", tt.expectedTokens[i], token.Literal, expected.Literal)
				}

				if token.Pos.Line != expected.Pos.Line {
					t.Errorf("Token %v: Position.Line = %d, want %d", tt.expectedTokens[i], token.Pos.Line, expected.Pos.Line)
				}

				if token.Pos.Column != expected.Pos.Column {
					t.Errorf("Token %v: Position.Column = %d, want %d", tt.expectedTokens[i], token.Pos.Column, expected.Pos.Column)
				}
			}
		})
	}
}

func TestJumpToNextColumn(t *testing.T) {
	tests := []struct {
		name        string
		initialCol  int
		expectedCol int
	}{
		{
			name:        "jump from column 1",
			initialCol:  1,
			expectedCol: 2,
		},
		{
			name:        "jump from column 10",
			initialCol:  10,
			expectedCol: 11,
		},
		{
			name:        "jump from column 0",
			initialCol:  0,
			expectedCol: 1,
		},
		{
			name:        "jump from large column number",
			initialCol:  100,
			expectedCol: 101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader("test")
			lexer := NewLexer(reader)

			// Set initial column
			lexer.position.Column = tt.initialCol

			// Call jumpToNextColumn
			lexer.jumpToNextColumn()

			// Check result
			if lexer.position.Column != tt.expectedCol {
				t.Errorf("jumpToNextColumn() Column = %d, want %d", lexer.position.Column, tt.expectedCol)
			}

			// Verify line number is unchanged
			if lexer.position.Line != 1 {
				t.Errorf("jumpToNextColumn() changed Line to %d, want 1", lexer.position.Line)
			}
		})
	}
}

func TestNewPosition(t *testing.T) {
	tests := []struct {
		name         string
		line         int
		column       int
		expectedLine int
		expectedCol  int
	}{
		{
			name:         "create position 1,1",
			line:         1,
			column:       1,
			expectedLine: 1,
			expectedCol:  1,
		},
		{
			name:         "create position 5,10",
			line:         5,
			column:       10,
			expectedLine: 5,
			expectedCol:  10,
		},
		{
			name:         "create position 0,0",
			line:         0,
			column:       0,
			expectedLine: 0,
			expectedCol:  0,
		},
		{
			name:         "create position with large numbers",
			line:         1000,
			column:       2000,
			expectedLine: 1000,
			expectedCol:  2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition(tt.line, tt.column)

			if pos.Line != tt.expectedLine {
				t.Errorf("NewPosition() Line = %d, want %d", pos.Line, tt.expectedLine)
			}

			if pos.Column != tt.expectedCol {
				t.Errorf("NewPosition() Column = %d, want %d", pos.Column, tt.expectedCol)
			}
		})
	}

}
