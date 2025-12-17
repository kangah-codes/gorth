package lexer

type TokenType string

const (
	// Special tokens
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	// Literals
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"
	BOOL   TokenType = "BOOL"

	// Identifiers
	IDENT TokenType = "IDENT"

	// Types
	TYPE_INT   TokenType = "TYPE_INT"
	TYPE_STR   TokenType = "TYPE_STR"
	TYPE_BOOL  TokenType = "TYPE_BOOL"
	TYPE_FLOAT TokenType = "TYPE_FLOAT"
	TYPE_PTR   TokenType = "TYPE_PTR"
	TYPE_ARR   TokenType = "TYPE_ARR"

	// keywords
	PLUS     TokenType = "PLUS"     // +
	MINUS    TokenType = "MINUS"    // -
	MULTIPLY TokenType = "MULTIPLY" // *
	DIVIDE   TokenType = "DIVIDE"   // /
	POWER    TokenType = "POWER"    // ^
	MODULO   TokenType = "MODULO"   // %

	// Comparison
	EQ  TokenType = "EQ"  // ==
	NEQ TokenType = "NEQ" // !=
	GT  TokenType = "GT"  // >
	LT  TokenType = "LT"  // <
	GTE TokenType = "GTE" // >=
	LTE TokenType = "LTE" // <=

	// Logical
	AND TokenType = "AND" // &&
	OR  TokenType = "OR"  // ||
	NOT TokenType = "NOT" // !

	// Assignment
	ASSIGN TokenType = "ASSIGN" // =

	// Keywords
	CONST   TokenType = "CONST"
	PROC    TokenType = "PROC"
	ENDPROC TokenType = "ENDPROC"
	IN      TokenType = "IN"
	RETURN  TokenType = "RETURN"

	// Stack operations
	DROP TokenType = "DROP"
	SWAP TokenType = "SWAP"
	DUP  TokenType = "DUP"
	OVER TokenType = "OVER"
	ROT  TokenType = "ROT"
	DEL  TokenType = "DEL"
	INC  TokenType = "INC"
	DEC  TokenType = "DEC"

	// I/O
	PRINT   TokenType = "PRINT"
	PRINTLN TokenType = "PRINTLN"
	DUMP    TokenType = "DUMP"

	// Pointer operations
	DEREF TokenType = "DEREF" // @
	PTR   TokenType = "PTR"   // *n format

	// Delimiters
	LBRACKET TokenType = "LBRACKET" // [
	RBRACKET TokenType = "RBRACKET" // ]

	// Comment

)

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

var keywords = map[string]TokenType{
	"const":   CONST,
	"proc":    PROC,
	"endproc": ENDPROC,
	"in":      IN,
	"return":  RETURN,
	"true":    BOOL,
	"false":   BOOL,
	"int":     TYPE_INT,
	"str":     TYPE_STR,
	"bool":    TYPE_BOOL,
	"float":   TYPE_FLOAT,
	"ptr":     TYPE_PTR,
	"arr":     TYPE_ARR,
	"drop":    DROP,
	"swap":    SWAP,
	"dup":     DUP,
	"over":    OVER,
	"rot":     ROT,
	"del":     DEL,
	"inc":     INC,
	"dec":     DEC,
	"print":   PRINT,
	"println": PRINTLN,
	"dump":    DUMP,

	"+":  PLUS,
	"-":  MINUS,
	"*":  MULTIPLY,
	"/":  DIVIDE,
	"^":  POWER,
	"%":  MODULO,
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	"<":  LT,
	">=": GTE,
	"<=": LTE,
	"&&": AND,
	"||": OR,
	"!":  NOT,
	"=":  ASSIGN,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// LookupOperator checks if an identifier is an operator
func LookupOperator(op string) TokenType {
	if tok, ok := keywords[op]; ok {
		return tok
	}

	return ILLEGAL
}

// IsOperator checks if a string is an operator
func IsOperator(s string) bool {
	_, ok := keywords[s]
	return ok && s != "true" && s != "false"
}

func IsDecimal(c rune) bool {
	return c == '.'
}

func (t TokenType) TokenIsOperator() bool {
	keywords := map[TokenType]bool{
		PLUS: true, MINUS: true, MULTIPLY: true, DIVIDE: true,
		POWER: true, MODULO: true, EQ: true, NEQ: true,
		GT: true, LT: true, GTE: true, LTE: true,
		AND: true, OR: true, NOT: true,
	}
	return keywords[t]
}

// IsKeyword checks if a token is a keyword
func (t TokenType) TokenIsKeyword() bool {
	for _, keyword := range keywords {
		if t == keyword {
			return true
		}
	}
	return false
}
