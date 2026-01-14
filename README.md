# Gorth Langu

Gorth is a stack-based, postfix notation programming language inspired by Forth. It features:

## Features

- **Primitive Types:**
    - Integers, Floats, Strings, Booleans, Null
- **Arrays:**
    - Literal syntax: `[1, 2, 3]`
    - Array methods: `->sizeof`, `->reverse`, etc. (coming soon)
- **Variables and Constants:**
    - `VAR`, `CONST` keywords
    - Declaration: `VAR x`
    - Assignment: `42 := x`
    - Inline Assignment: `42 := VAR x`
    - Constants: `CONST PI 3.14`
- **Operators:**
    - Math: `+`, `-`, `*`, `/`, `^`, `%`
    - Logical: `&&`, `||`, `!`
    - Comparison: `<`, `>`, `<=`, `>=`, `==`, `!=`
- **Control Flow:**
    - `IF`, `ELSE`, `END`, `WHILE`, `DO`, `BREAK`, `CONTINUE`
- **Procedures:**
    - Define: `PROC ... ENDP`
    - Call: `CALL`, `RETURN`, `IN
- **Stack Operations:**
    - `DUMP`, `PRINT`, direct stack manipulation
- **Comments:**
    - Lines starting with `#`
- **Error Handling:**
    - Runtime errors, illegal tokens, position tracking

## Example Usage

Run a Gorth program:

```bash
go run lib/main.go examples/hello-world.gorth
```

Enable debug mode for token/AST output:

```bash
go run lib/main.go examples/hello-world.gorth -d
```

## Example Programs

- `hello-world.gorth`: Print a string
- `object-methods.gorth`: Array and string methods
- `control-flow.gorth`: Control flow constructs
- `test.gorth`: Primitives, math, logic, comparisons

See the `examples/` folder for more.

---


For more details, see the source code and examples.