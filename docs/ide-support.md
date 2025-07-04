# IDE Support

Fun provides excellent IDE support through a VS Code extension and Language Server Protocol (LSP) implementation.

## Table of Contents

- [VS Code Extension](#vs-code-extension)
- [Language Server Protocol](#language-server-protocol)
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)

## VS Code Extension

The Fun VS Code extension provides a complete development environment for the Fun programming language.

### Features

- **Syntax Highlighting**: Color-coded syntax for Fun code
- **IntelliSense**: Code completion and suggestions
- **Error Reporting**: Real-time syntax and type error detection
- **Hover Information**: Type information on hover
- **Go to Definition**: Navigate to function and variable definitions
- **Symbol Search**: Find symbols across your codebase

### Extension Structure

```
vscode/
├── client/              # VS Code extension client
│   ├── src/
│   │   └── extension.ts # Main extension code
│   ├── package.json     # Extension manifest
│   └── tsconfig.json    # TypeScript configuration
├── package.json         # Workspace configuration
└── README.md           # Extension documentation
```

## Language Server Protocol

The Fun LSP server provides language intelligence features that work with any LSP-compatible editor.

### LSP Features

- **Syntax Analysis**: Parse and validate Fun code
- **Type Checking**: Real-time type inference and error detection
- **Code Completion**: Suggest functions, variables, and types
- **Hover Information**: Display type information and documentation
- **Error Diagnostics**: Report syntax and type errors
- **Symbol Information**: Provide symbol details and locations

### Starting the LSP Server

```bash
# Start LSP server directly
./fun lsp

# The server runs on stdio and communicates with the client
```

## Features

### Syntax Highlighting

The extension provides syntax highlighting for:

- **Keywords**: `when`, `is`, `else`, `import`, `from`
- **Operators**: `+`, `-`, `==`, `->`
- **Literals**: Numbers, strings, booleans
- **Functions**: Lambda expressions, function calls
- **Types**: Type annotations, type constructors
- **Comments**: Single-line comments with `#`

### IntelliSense

#### Code Completion

```fun
# Type to see suggestions
inc = \x -> x + 1
inc(  # Shows parameter hints
```

#### Hover Information

```fun
# Hover over functions to see types
inc = \x -> x + 1
# Hover over 'inc' shows: Lam<Int, Int>

# Hover over variables
name = "Alice"
# Hover over 'name' shows: Str
```

#### Go to Definition

```fun
# Right-click on function name
import math from `./math`
math.inc(5)  # Go to definition of 'inc' in math.fun
```

### Error Reporting

#### Syntax Errors

```fun
# Missing semicolon
when x == 0 is
    True -> 1
    False -> 2  # Error: missing semicolon

# Invalid syntax
inc = \x -> x +  # Error: incomplete expression
```

#### Type Errors

```fun
# Type mismatch
inc = \x -> x + 1
inc(5, 3)  # Error: too many arguments

# Undefined variable
result = undefined_var  # Error: undefined variable
```

#### Import Errors

```fun
# Missing module
import missing from `./missing`  # Error: module not found

# Invalid import path
import lib from `invalid/path`   # Error: cannot read module
```

## Installation

### Building the Extension

1. **Build the Fun executable**:
   ```bash
   go build
   ```

2. **Open the extension directory**:
   ```bash
   cd vscode
   code .  # Opens VS Code in the extension directory
   ```

3. **Install dependencies**:
   ```bash
   npm install
   ```

4. **Run the extension**:
   - Press `F5` in VS Code
   - This opens a new Extension Development Host window
   - Open a `.fun` file to test the extension

### Installing in VS Code

1. **Package the extension**:
   ```bash
   cd vscode
   vsce package
   ```

2. **Install the extension**:
   - In VS Code, go to Extensions (Ctrl+Shift+X)
   - Click "Install from VSIX..."
   - Select the generated `.vsix` file

### Manual Installation

1. **Copy the extension**:
   ```bash
   cp -r vscode ~/.vscode/extensions/fun-language
   ```

2. **Restart VS Code**:
   - Close and reopen VS Code
   - The extension should be available

## Configuration

### Extension Settings

Add these settings to your VS Code settings:

```json
{
  "fun.languageServerPath": "/path/to/your/fun",
  "fun.enableTypeChecking": true,
  "fun.showTypeInformation": true,
  "fun.enableCompletion": true
}
```

### Language Configuration

The extension includes language configuration for:

- **File associations**: `.fun` files are recognized as Fun code
- **Comments**: `#` for single-line comments
- **Brackets**: Automatic bracket matching
- **Indentation**: Smart indentation for Fun syntax

### Keybindings

Default keybindings:

- `Ctrl+Space`: Trigger suggestions
- `F12`: Go to definition
- `Shift+F12`: Find all references
- `Ctrl+Shift+O`: Go to symbol in file

## Troubleshooting

### Common Issues

#### Extension Not Loading

```bash
# Check if the Fun executable is built
ls -la fun

# Verify the extension path in settings
# Make sure fun.languageServerPath points to the correct executable
```

#### LSP Server Not Starting

```bash
# Test the LSP server manually
./fun lsp

# Check for error messages
# Common issues:
# - Executable not found
# - Permission denied
# - Missing dependencies
```

#### No Syntax Highlighting

1. **Check file association**:
   - Make sure `.fun` files are associated with Fun language
   - Right-click on a `.fun` file → "Open With" → "Fun"

2. **Reload the extension**:
   - Press `Ctrl+Shift+P`
   - Type "Developer: Reload Window"

#### No IntelliSense

1. **Check LSP connection**:
   - Open Output panel (View → Output)
   - Select "Fun Language Server" from dropdown
   - Look for connection errors

2. **Verify file structure**:
   - Make sure your `.fun` files are valid
   - Check for syntax errors that might prevent analysis

### Debug Mode

Enable debug mode for detailed logging:

```json
{
  "fun.debug": true,
  "fun.trace": "verbose"
}
```

### Manual LSP Testing

Test the LSP server manually:

```bash
# Start server
./fun lsp

# In another terminal, send LSP messages
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":123,"rootUri":"file:///tmp","capabilities":{}}}' | ./fun lsp
```

## Advanced Configuration

### Custom LSP Settings

Configure LSP features:

```json
{
  "fun.lsp": {
    "completion": {
      "triggerCharacters": [".", "("]
    },
    "hover": {
      "enabled": true
    },
    "diagnostics": {
      "enabled": true,
      "delay": 1000
    }
  }
}
```

### Workspace Settings

Project-specific settings in `.vscode/settings.json`:

```json
{
  "fun.languageServerPath": "./fun",
  "fun.workspaceRoot": ".",
  "fun.includePaths": ["./src", "./lib"]
}
```

### Extension Development

For developers working on the extension:

```bash
# Install development dependencies
cd vscode
npm install

# Run tests
npm test

# Build extension
npm run compile

# Package for distribution
vsce package
```

## Best Practices

### File Organization

```
project/
├── src/
│   ├── main.fun
│   ├── math.fun
│   └── utils.fun
├── .vscode/
│   └── settings.json
└── fun  # Executable
```

### Extension Usage

1. **Use type annotations** for better IntelliSense
2. **Organize code in modules** for better navigation
3. **Use consistent naming** for better completion
4. **Enable error reporting** to catch issues early

### Performance Tips

1. **Keep files small** for faster analysis
2. **Use modules** to organize large codebases
3. **Avoid circular dependencies** that slow down analysis
4. **Use type annotations** to reduce inference time

The IDE support provides a professional development experience for Fun programming, with features comparable to mainstream programming languages. 