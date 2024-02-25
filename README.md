# Gorth

Gorth is a Go implementation of a simple Forth-like language. It is a work in progress and is not yet feature complete. Heavily inspired by Porth, may even copy some of his code tbh.

## Description

This project is a stack-based programming language implemented in Go. It provides a simple and efficient way to write programs using stack operations.

## Features

-   Basic stack operations: push, pop, swap, duplicate
-   Arithmetic operations: add, subtract, multiply, divide
-   Comparison and logical operations
-   Control flow: conditional branches and loops
-   Functions or subroutines
-   Basic input/output
-   Error handling
-   Memory management (optional)
-   Concurrency (optional)
-   Debugging tools

## Usage

Have go installed, and run `go run gorth.go`. To test, run `go test`.

Run a gorth file
`go run gorth.go ./hello.gorth -d -s`

`-d` is for debug mode, `-s` is for strict mode.

In strict mode, all values on the stack have to be consumed by the end of the program, otherwise you get a big fat error. In non-strict mode, the stack can have values left over. Don't ask me why 🤷🏾‍♀️

Gorth supports various built-in operators, including arithmetic operations (+, -, \*, /, %, ^), logical operations (&&, ||, !), and comparison operations (==, !=, ===). Working on adding >, <, >=, <=.

## Examples

### Hello World

```gorth
"Hello, World!" print drop

# or

"Hello, World!" dump
```

### Simple Arithmetic

```gorth
1 2 + print drop # 3
```

### Logical Operations

```gorth
1 2 != print drop # true
```

## Contributing

If you're open to contributions, provide instructions on how potential contributors can help improve your project.
