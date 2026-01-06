# Gorth Language Support (VS Code)

Syntax highlighting and basic language configuration for `.gorth` files.

## Highlighted (currently implemented)

- Comments: `# ...`
- Strings: `"..."` (with escapes like `\\n`, `\\t`, `\\"`)
- Numbers: integers and floats
- Keywords: `CONST`, `VAR`, `IF`, `ELSE`, `END`, `WHILE`
- Stack ops: `DROP`, `SWAP`, `DUP`, `OVER`, `ROT`, `DEL`, `INC`, `DEC`, `CLEAR`, `PICK`
- I/O: `DUMP`, `DUMPLN`
- Operators: `+ - * / % ^`, comparisons (`== != >= <= > <`), logical (`&& || !`), assignment (`=` and `:=`)
- Literals: booleans (`true`/`false` or `TRUE`/`FALSE`), `NULL`
- Arrays: `[ ... ]` with `,` separators

## Install (from source)

- Close all VS Code windows.
- Remove any existing installed copies of this extension (duplicate extension IDs can cause VS Code to load the older one):
  - macOS/Linux: remove folders matching `~/.vscode/extensions/kangah-codes.gorth-language-support*`
  - Windows: remove folders matching `%USERPROFILE%\\.vscode\\extensions\\kangah-codes.gorth-language-support*`
- Copy this folder into your VS Code extensions directory:
  - macOS/Linux:
    - `cp -R editor-support/vscode ~/.vscode/extensions/kangah-codes.gorth-language-support-0.1.0`
  - Windows (PowerShell):
    - `Copy-Item -Recurse -Force editor-support\\vscode $env:USERPROFILE\\.vscode\\extensions\\kangah-codes.gorth-language-support-0.1.0`
- Start VS Code and run **Developer: Reload Window**.

If `DUMPLN` still isn’t highlighted:

- Confirm the file is in **Gorth** language mode (bottom-right language picker).
- Run **Developer: Inspect Editor Tokens and Scopes**, click `DUMPLN`, and verify it has scope `keyword.other.io.gorth`.

## Develop

- Open `editor-support/vscode` in VS Code
- Press `F5` to launch an Extension Development Host
