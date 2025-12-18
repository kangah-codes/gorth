# Gorth Language Support for VS Code

Syntax highlighting and language support for the Gorth programming language.

## Features

- Syntax highlighting for `.gorth` files
- Comment support (line comments with `#`)
- Bracket matching and auto-closing
- Support for all Gorth language constructs:
  - Keywords: `const`, `proc`, `endproc`, `in`, `return`
  - Types: `int`, `str`, `bool`, `float`, `ptr`, `arr`
  - Stack operations: `drop`, `swap`, `dup`, `over`, `rot`, `del`, `inc`, `dec`
  - I/O operations: `print`, `println`, `dump`
  - Operators: arithmetic, comparison, logical
  - Variables (prefixed with `$`)
  - Arrays with `[]` syntax
  - Pointers with `@` and `*n` syntax

## Installation

### From Source

1. Clone this repository
2. Copy the `language-support` folder to your VS Code extensions directory:
   - **macOS/Linux**: `~/.vscode/extensions/`
   - **Windows**: `%USERPROFILE%\.vscode\extensions\`
3. Restart VS Code

### Development

1. Open the `language-support` folder in VS Code
2. Press `F5` to launch a new VS Code window with the extension loaded
3. Open any `.gorth` file to test the syntax highlighting

## Customization

You can customize the colors by adding this to your `settings.json`:

```json
{
  "editor.tokenColorCustomizations": {
    "textMateRules": [
      {
        "scope": "keyword.control.gorth",
        "settings": {
          "foreground": "#C586C0",
          "fontStyle": "bold"
        }
      },
      {
        "scope": "storage.type.gorth",
        "settings": {
          "foreground": "#4EC9B0"
        }
      },
      {
        "scope": "keyword.other.stack.gorth",
        "settings": {
          "foreground": "#DCDCAA"
        }
      },
      {
        "scope": "keyword.other.io.gorth",
        "settings": {
          "foreground": "#4FC1FF"
        }
      }
    ]
  }
}
```

## License

See the main Gorth repository for license information.
