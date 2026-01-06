# Gorth Language Support (VS Code)

Syntax highlighting and basic language configuration for `.gorth` files.

## Highlighted (currently implemented)

- Comments: `# ...`
- Strings: `"..."` (with escapes like `\\n`, `\\t`, `\\"`)
- Numbers: integers and floats
- Keywords: `CONST`, `VAR`, `IF`, `ELSE`, `END`, `WHILE`
- Stack ops: `DROP`, `SWAP`, `DUP`, `OVER`, `ROT`, `DEL`, `INC`, `DEC`, `CLEAR`, `PICK`
- I/O: `DUMP`
- Operators: `+ - * / % ^`, comparisons (`== != >= <= > <`), logical (`&& || !`), assignment (`=` and `:=`)
- Literals: booleans (`true`/`false` or `TRUE`/`FALSE`), `NULL`
- Arrays: `[ ... ]` with `,` separators

## Install (from source)

- Copy `editor-support/vscode` into your VS Code extensions directory:
  - macOS/Linux: `~/.vscode/extensions/`
  - Windows: `%USERPROFILE%\\.vscode\\extensions\\`

## Develop

- Open `editor-support/vscode` in VS Code
- Press `F5` to launch an Extension Development Host
