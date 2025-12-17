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
			expectedChar:   'a',
			expectedColumn: 1,
		},
		{
			name:           "backup after reading multiple characters",
			input:          "hello",
			setupSteps:     3,
			expectedChar:   'l',
			expectedColumn: 2,
		},
		{
			name:           "backup with numbers",
			input:          "123",
			setupSteps:     2,
			expectedChar:   '2',
			expectedColumn: 1,
		},
		{
			name:           "backup with symbols",
			input:          "+=*",
			setupSteps:     1,
			expectedChar:   '+',
			expectedColumn: 1,
		},
		{
			name:           "backup with unicode",
			input:          "αβγ",
			setupSteps:     2,
			expectedChar:   'β',
			expectedColumn: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lexer := NewLexer(reader)

			// Read setupSteps characters to advance position
			for i := 0; i < tt.setupSteps; i++ {
				lexer.char, _ = lexer.readChar()
				lexer.position.Column++
			}

			// Call backup
			lexer.backup()

			// Verify position was decremented
			if lexer.position.Column != tt.expectedColumn {
				t.Errorf("backup() Column = %d, want %d", lexer.position.Column, tt.expectedColumn)
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
