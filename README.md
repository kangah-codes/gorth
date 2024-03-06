# Gorth

<p>
    <a href="https://go.dev">
        <img src="https://img.shields.io/badge/Language-Golang-077D9C" alt="Go" />
    </a>
    <img src="badge.svg" alt="Coverage" />
</p>

## Description

Gorth is a Go implementation of a simple stack-based language. It is a work in progress and is not yet feature complete. Heavily inspired by Porth.

## Features

-   Basic stack operations: push, pop, swap, duplicate
-   Arithmetic operations: add, subtract, multiply, divide
-   Comparison and logical operations
-   Variable assignment and manipulation
-   Control flow: conditional branches and loops (Coming soon!)
-   Functions or subroutines (Coming soon!)
-   Basic input/output (Coming soon!)
-   Error handling (Coming soon!)
-   Memory management (Coming soon!)
-   Standard library (Coming soon!)

## Exhaustive List of operations

-   Push (Primitive Op) - Pushes a value onto the stack
-   Peek (Primitive Op) - Peeks at the top value on the stack
-   Pop (Primitive Op) - Pops a value off the stack
-   Swap (swap) - Swaps the top two values on the stack
-   Dup (dup) - Duplicates the top value on the stack
-   Rot (rot) - Rotates the top three values on the stack
-   Drop (rot) - Drops the top value off the stack
-   Print (print) - Prints the top value on the stack
-   Dump (dump) - Dumps the top value on the stack
-   Modulo (%) - Supports int only (a b %) cos I'm lazy
-   Increment (++) - Increments the top value on the stack by 1
-   Decrement (--) - Decrements the top value on the stack by 1

## Usage

Have go installed, and run `go run gorth.go`. To test, run `go test`.

Run a gorth file
`go run gorth.go ./hello.gorth -d -s`

`-d` is for debug mode, `-s` is for strict mode.

In strict mode, all values on the stack have to be consumed by the end of the program, otherwise you get a big fat error. In non-strict mode, the stack can have values left over. Don't ask me why 🤷🏾‍♀️

Gorth supports various built-in operators, including arithmetic operations (+, -, \*, /, %, ^), logical operations (&&, ||, !), and comparison operations (==, !=, ===). Also (>=, <=, >, <.)

## Examples

### Hello World

```gorth
"Hello, World!" print drop

# or this, if in strict mode

"Hello, World!" dump
```

### Simple Arithmetic

```gorth
1 2 + print drop
```

### Logical Operations

```gorth
1 2 != print drop
```

## Contributing

Idk make a pr or something
