package runtime

import "gorth/lexer"

type GorthRuntime struct {
	stack []lexer.StackElement
}
