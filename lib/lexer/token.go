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

	// Null value
	NULL TokenType = "NULL"

	// Identifiers (variable names, word names, etc.)
	IDENT TokenType = "IDENT"

	// Math operators
	OP_PLUS     TokenType = "OP_PLUS"     // +
	OP_SUBTRACT TokenType = "OP_SUBTRACT" // -
	OP_MULTIPLY TokenType = "OP_MULTIPLY" // *
	OP_DIVIDE   TokenType = "OP_DIVIDE"   // /
	OP_POWER    TokenType = "OP_POWER"    // ^
	OP_MODULO   TokenType = "OP_MODULO"   // %

	// Comparison
	OP_EQ  TokenType = "OP_EQ"  // ==
	OP_NEQ TokenType = "OP_NEQ" // !=
	OP_GT  TokenType = "OP_GT"  // >
	OP_LT  TokenType = "OP_LT"  // <
	OP_GTE TokenType = "OP_GTE" // >=
	OP_LTE TokenType = "OP_LTE" // <=

	// Logical
	OP_AND TokenType = "OP_AND" // &&
	OP_OR  TokenType = "OP_OR"  // ||
	OP_NOT TokenType = "OP_NOT" // !

	// Assignment
	OP_ASSIGN TokenType = "OP_ASSIGN" // =

	// Keywords
	CONST    TokenType = "CONST"
	VAR      TokenType = "VAR"
	PROC     TokenType = "PROC"
	IN       TokenType = "IN"
	RETURN   TokenType = "RETURN"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	END      TokenType = "END"
	WHILE    TokenType = "WHILE"
	DO       TokenType = "DO"
	BREAK    TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"
	CALL     TokenType = "CALL"

	// Stack operations
	OP_DROP  TokenType = "OP_DROP"
	OP_SWAP  TokenType = "OP_SWAP"
	OP_DUP   TokenType = "OP_DUP"
	OP_OVER  TokenType = "OP_OVER"
	OP_ROT   TokenType = "OP_ROT"
	OP_DEL   TokenType = "OP_DEL"
	OP_INC   TokenType = "OP_INC"
	OP_DEC   TokenType = "OP_DEC"
	OP_CLEAR TokenType = "OP_CLEAR"
	OP_PICK  TokenType = "OP_PICK"

	// I/O
	OP_DUMP   TokenType = "OP_DUMP"
	OP_DUMPLN TokenType = "OP_DUMPLN"

	// Delimiters
	LSQUARE_BRACKET TokenType = "LSQUARE_BRACKET" // [
	RSQUARE_BRACKET TokenType = "RSQUARE_BRACKET" // ]
	COMMA           TokenType = "COMMA"           // ,
	LCURL_BRACKET   TokenType = "LCURL_BRACKET"   // {
	RCURL_BRACKET   TokenType = "RCURL_BRACKET"   // }
	LBRACKET        TokenType = "LBRACKET"        // (
	RBRACKET        TokenType = "RBRACKET"        // )
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

var TokenMap = map[TokenType]string{
	EOF:             "EOF",
	ILLEGAL:         "ILLEGAL",
	INT:             "INT",
	FLOAT:           "FLOAT",
	STRING:          "STRING",
	BOOL:            "BOOL",
	NULL:            "NULL",
	IDENT:           "IDENT",
	OP_PLUS:         "+",
	OP_SUBTRACT:     "-",
	OP_MULTIPLY:     "*",
	OP_DIVIDE:       "/",
	OP_POWER:        "^",
	OP_MODULO:       "%",
	OP_EQ:           "==",
	OP_NEQ:          "!=",
	OP_GT:           ">",
	OP_LT:           "<",
	OP_GTE:          ">=",
	OP_LTE:          "<=",
	OP_AND:          "&&",
	OP_OR:           "||",
	OP_NOT:          "!",
	OP_ASSIGN:       ":=",
	CONST:           "CONST",
	VAR:             "VAR",
	PROC:            "PROC",
	IN:              "IN",
	RETURN:          "RETURN",
	OP_DROP:         "DROP",
	OP_SWAP:         "SWAP",
	OP_DUP:          "DUP",
	OP_OVER:         "OVER",
	OP_ROT:          "ROT",
	OP_DEL:          "DEL",
	OP_INC:          "INC",
	OP_DEC:          "DEC",
	OP_DUMP:         "DUMP",
	OP_DUMPLN:       "DUMPLN",
	OP_CLEAR:        "CLEAR",
	OP_PICK:         "PICK",
	LSQUARE_BRACKET: "[",
	RSQUARE_BRACKET: "]",
	COMMA:           ",",
	IF:              "IF",
	ELSE:            "ELSE",
	END:             "END",
	WHILE:           "WHILE",
	DO:              "DO",
	CONTINUE:        "CONTINUE",
	BREAK:           "BREAK",
	CALL:            "CALL",
}

var keywords = map[string]TokenType{
	"CONST":    CONST,
	"VAR":      VAR,
	"PROC":     PROC,
	"IN":       IN,
	"RETURN":   RETURN,
	"TRUE":     BOOL,
	"FALSE":    BOOL,
	"NULL":     NULL,
	"DROP":     OP_DROP,
	"SWAP":     OP_SWAP,
	"DUP":      OP_DUP,
	"OVER":     OP_OVER,
	"ROT":      OP_ROT,
	"DEL":      OP_DEL,
	"DUMP":     OP_DUMP,
	"DUMPLN":   OP_DUMPLN,
	"CLEAR":    OP_CLEAR,
	"PICK":     OP_PICK,
	"INC":      OP_INC,
	"DEC":      OP_DEC,
	"IF":       IF,
	"ELSE":     ELSE,
	"END":      END,
	"WHILE":    WHILE,
	"DO":       DO,
	"BREAK":    BREAK,
	"CONTINUE": CONTINUE,
	"CALL":     CALL,
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
		OP_PLUS: true, OP_SUBTRACT: true, OP_MULTIPLY: true, OP_DIVIDE: true,
		OP_POWER: true, OP_MODULO: true, OP_EQ: true, OP_NEQ: true,
		OP_GT: true, OP_LT: true, OP_GTE: true, OP_LTE: true,
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
