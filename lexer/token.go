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
	OP_PLUS     TokenType = "OP_PLUS"     // +
	OP_MINUS    TokenType = "OP_MINUS"    // -
	OP_MULTIPLY TokenType = "OP_MULTIPLY" // *
	OP_DIVIDE   TokenType = "OP_DIVIDE"   // /
	OP_POWER    TokenType = "OP_POWER"    // ^
	OP_MODULO   TokenType = "OP_MODULO"   // %

	// Comparison
	EQ  TokenType = "EQ"  // ==
	NEQ TokenType = "NEQ" // !=
	GT  TokenType = "GT"  // >
	LT  TokenType = "LT"  // <
	GTE TokenType = "GTE" // >=
	LTE TokenType = "LTE" // <=

	// Logical
	OP_AND TokenType = "OP_AND" // &&
	OP_OR  TokenType = "OP_OR"  // ||
	OP_NOT TokenType = "OP_NOT" // !

	// Assignment
	ASSIGN TokenType = "ASSIGN" // =

	// Keywords
	CONST   TokenType = "CONST"
	PROC    TokenType = "PROC"
	ENDPROC TokenType = "ENDPROC"
	IN      TokenType = "IN"
	RETURN  TokenType = "RETURN"

	// Stack operations
	OP_DROP TokenType = "OP_DROP"
	OP_SWAP TokenType = "OP_SWAP"
	OP_DUP  TokenType = "OP_DUP"
	OP_OVER TokenType = "OP_OVER"
	OP_ROT  TokenType = "OP_ROT"
	OP_DEL  TokenType = "OP_DEL"
	OP_INC  TokenType = "OP_INC"
	OP_DEC  TokenType = "OP_DEC"

	// I/O
	OP_PRINT   TokenType = "OP_PRINT"
	OP_PRINTLN TokenType = "OP_PRINTLN"
	OP_DUMP    TokenType = "OP_DUMP"

	// Pointer operations
	PTR_DEREF TokenType = "PTR_DEREF" // @
	PTR       TokenType = "PTR"       // *n format

	// Delimiters
	LBRACKET TokenType = "LBRACKET" // [
	RBRACKET TokenType = "RBRACKET" // ]
	COMMA    TokenType = "COMMA"    // ,

	// Comment

	// Variables
	VARIABLE TokenType = "VARIABLE"
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

var TokenMap = map[TokenType]string{
	EOF:         "EOF",
	ILLEGAL:     "ILLEGAL",
	INT:         "INT",
	FLOAT:       "FLOAT",
	STRING:      "STRING",
	BOOL:        "BOOL",
	IDENT:       "IDENT",
	TYPE_INT:    "TYPE_INT",
	TYPE_STR:    "TYPE_STR",
	TYPE_BOOL:   "TYPE_BOOL",
	TYPE_FLOAT:  "TYPE_FLOAT",
	TYPE_PTR:    "TYPE_PTR",
	TYPE_ARR:    "TYPE_ARR",
	OP_PLUS:     "OP_PLUS",
	OP_MINUS:    "OP_MINUS",
	OP_MULTIPLY: "OP_MULTIPLY",
	OP_DIVIDE:   "OP_DIVIDE",
	OP_POWER:    "OP_POWER",
	OP_MODULO:   "OP_MODULO",
	EQ:          "EQ",
	NEQ:         "NEQ",
	GT:          "GT",
	LT:          "LT",
	GTE:         "GTE",
	LTE:         "LTE",
	OP_AND:      "OP_AND",
	OP_OR:       "OP_OR",
	OP_NOT:      "OP_NOT",
	ASSIGN:      "ASSIGN",
	CONST:       "CONST",
	PROC:        "PROC",
	ENDPROC:     "ENDPROC",
	IN:          "IN",
	RETURN:      "RETURN",
	OP_DROP:     "OP_DROP",
	OP_SWAP:     "OP_SWAP",
	OP_DUP:      "OP_DUP",
	OP_OVER:     "OP_OVER",
	OP_ROT:      "OP_ROT",
	OP_DEL:      "OP_DEL",
	OP_INC:      "OP_INC",
	OP_DEC:      "OP_DEC",
	OP_PRINT:    "OP_PRINT",
	OP_PRINTLN:  "OP_PRINTLN",
	OP_DUMP:     "OP_DUMP",
	PTR_DEREF:   "PTR_DEREF",
	PTR:         "PTR",
	LBRACKET:    "LBRACKET",
	RBRACKET:    "RBRACKET",
	COMMA:       "COMMA",
	VARIABLE:    "VARIABLE",
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
	"drop":    OP_DROP,
	"swap":    OP_SWAP,
	"dup":     OP_DUP,
	"over":    OP_OVER,
	"rot":     OP_ROT,
	"del":     OP_DEL,
	"inc":     OP_INC,
	"dec":     OP_DEC,
	"print":   OP_PRINT,
	"println": OP_PRINTLN,
	"dump":    OP_DUMP,

	"+":  OP_PLUS,
	"-":  OP_MINUS,
	"*":  OP_MULTIPLY,
	"/":  OP_DIVIDE,
	"^":  OP_POWER,
	"%":  OP_MODULO,
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	"<":  LT,
	">=": GTE,
	"<=": LTE,
	"&&": OP_AND,
	"||": OP_OR,
	"!":  OP_NOT,
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
		OP_PLUS: true, OP_MINUS: true, OP_MULTIPLY: true, OP_DIVIDE: true,
		OP_POWER: true, OP_MODULO: true, EQ: true, NEQ: true,
		GT: true, LT: true, GTE: true, LTE: true,
		OP_AND: true, OP_OR: true, OP_NOT: true,
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
