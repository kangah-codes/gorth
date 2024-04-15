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

-   Basic stack operations: push, pop, swap, dup, rot, drop, over
-   Arithmetic operations: add, subtract, multiply, divide
-   Stack Pointers
-   Comparison and logical operations (Coming soon!)
-   Variable assignment and manipulation
-   Control flow: conditional branches and loops (Coming soon!)
-   Procedures (Coming soon!)
-   Basic input/output (Coming soon!)
-   Error handling (Coming soon!)
-   Standard library (Coming soon!)

## Exhaustive List of operations

| Operation | Description                                                                                              |
| --------- | -------------------------------------------------------------------------------------------------------- |
| `swap`    | Swaps the top two values on the stack                                                                    |
| `dup`     | Duplicates the top value on the stack                                                                    |
| `rot`     | Rotates the top three values on the stack                                                                |
| `print`   | Prints the top value on the stack                                                                        |
| `dump`    | Drops and prints the top value on the stack                                                              |
| `drop`    | Drops the top value on the stack                                                                         |
| `delete`  | Deletes a varable from mem (not to be confused with drop, which will only drop the value from the stack) |
| `over`    | Duplicates the second value on the stack to the top of the stack                                         |
| `inc`     | Increments the top value on the stack by 1                                                               |
| `dec`     | Decrements the top value on the stack by 1                                                               |
| `+`       | Adds the top 2 values on the stack                                                                       |
| `-`       | Subtracts the top 2 values on the stack                                                                  |
| `*`       | Multiplies the top 2 values on the stack                                                                 |
| `/`       | Divides the top 2 values on the stack                                                                    |
| `^`       | Raises the top value to the power of the second value                                                    |
| `%`       | Performs mod operation on top 2 values on the stack                                                      |

## Usage

Have go installed, and run `go run gorth.go`. To test, run `go test`.

Run a gorth file
`go run gorth.go -f=./hello.gorth -d -s`

`-f` is the path to the file,`-d` is for debug mode, `-s` is for strict mode.

In strict mode, all values on the stack have to be consumed by the end of the program, otherwise you get a big fat error. In non-strict mode, the stack can have values left over. Don't ask me why 🤷🏾‍♀️

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

### Assigning a variable

```gorth
pi 3.14 =

## adding type annotations is optional since the value will be inferred
pi 3.14 float =

myName "Joshua" string =
```

### Pointers

Assuming I have a stack ["Hello", "World", "!", "Gorth", "is", "cool"], and I want to print Hello world, I would have to drop all the other values on the stack, before performing a swap operation to get Hello on top, then calling print twice. But with pointers, I can reference them by their index on the execution stack, and concatenate them like this:

```gorth
"Hello" "World" "!" "Gorth" "is" "cool"

&0 " " + &1 + println
```

## Contributing

Idk make a pr or something
