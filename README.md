# Fun Programming Language

A purely functional programming language with type inference, built in Go with Tree-sitter parsing and LSP support.

## Features

- **Functional Programming**: Lambda expressions, pattern matching, and immutable data structures
- **Type System**: Hindley-Milner type inference with polymorphic types
- **Pattern Matching**: `when` expressions with constructor patterns
- **Records**: Named field data structures with type safety
- **Lists**: Immutable list data structures
- **Modules**: Import system for code organization
- **REPL**: Interactive development environment
- **LSP Support**: Language Server Protocol for IDE integration
- **Tree-sitter Grammar**: Robust parsing with syntax highlighting

## Quick Start

### Installation

```bash
git clone <repository>
cd fun
go build
```

### Running Programs

```bash
# Run a file
./fun examples/basic/main
```

# Start REPL
./fun

# Start LSP server
./fun lsp
```

### Example

```fun
import lib from `./lib`

# computes the sum of integers from 1..n
sum_range : Lam<Int, Int>
sum_range = fix(\rec -> \n ->
    when n == 0 is
        True t -> 0;
        False f -> n + rec(lib.dec(n))
)

sum_range(100) # result: 5050
```

## Language Features

### Basic Types

- **Integers**: `42`, `-17`
- **Strings**: `` `hello world` ``
- **String Templates**: `` `hello {name}` ``

### Functions

```fun
# Lambda expressions
inc = \x -> x + 1

# Function application
inc(5)  # result: 6

# Multi-parameter functions
add = \x, y -> x + y
```

### Pattern Matching

```fun
when value is
    Just x -> x;
    Nothing -> 0
else 42
```

### Records

```fun
# Record construction
person = {name: "Alice", age: 30}

# Record access
person.name  # result: "Alice"

# Record with type annotation
person : {name: Str, age: Int}
person = {name: "Bob", age: 25}
```

### Lists

```fun
# List construction
numbers = [1, 2, 3, 4, 5]

# Empty list
empty = []
```

### Type Annotations

```fun
# Explicit type annotation
factorial : Lam<Int, Int>
factorial = \n -> when n == 0 is
    True -> 1;
    False -> n * factorial(n - 1)
```

### Modules

Modules allow you to organize code across multiple files. Each file is a module that can export values and import from other modules. The `.fun` extension is the conventional choice.

#### Creating a Module

```fun
# lib.fun - a utility module
inc = \x -> x + 1
dec = \x -> x - 1

# Export values by returning a record
{
    inc: inc,
    dec: dec
}
```

#### Importing Modules

```fun
# main.fun - importing and using the lib module
import lib from `./lib`

# Use imported functions
result1 = lib.inc(5)    # result: 6
result2 = lib.dec(10)   # result: 9

# You can also import with a different name
import utils from `./lib`
result3 = utils.inc(3)  # result: 4
```

#### Module Paths

Import paths use backticks and can be:
- **Relative paths**: `./lib`, `../utils/math`
- **File extensions**: Any extension is supported, `.fun` is the convention
- **Directory traversal**: `../../shared/helpers`

#### Module Structure

A module file can contain:
- **Function definitions**: `inc = \x -> x + 1`
- **Type annotations**: `factorial : Lam<Int, Int>`
- **Import statements**: `import other from \`./other\``
- **Export record**: The final expression is exported

```fun
# math.fun - a more complex module
import lib from `./lib`

# Type annotations
factorial : Lam<Int, Int>
factorial = \n -> when n == 0 is
    True -> 1;
    False -> n * factorial(lib.dec(n))

# Export multiple functions
{
    factorial: factorial,
    inc: lib.inc,
    dec: lib.dec
}
```

## Built-in Functions

- `+` : Addition for integers
- `-` : Subtraction for integers  
- `==` : Equality comparison
- `fix` : Fixed-point combinator for recursion

## Development

### Project Structure

```
fun/
├── main.go              # Main executable
├── internal/            # Core language implementation
│   ├── program.go       # Program execution
│   ├── expr.go          # Expression AST
│   ├── type.go          # Type system
│   ├── value.go         # Runtime values
│   ├── env.go           # Environment management
│   └── lsp.go           # Language Server Protocol
├── tree-sitter-fun/     # Parser implementation
│   ├── grammar.js       # Tree-sitter grammar
│   └── bindings/        # Language bindings
├── examples/            # Example programs
└── vscode/              # VS Code extension
```

### Building

```bash
# Build main executable
go build

# Build tree-sitter parser
cd tree-sitter-fun
npm install
npx tree-sitter generate
```

### Testing

```bash
# Run tree-sitter tests
cd tree-sitter-fun
npx tree-sitter test

# Run Go tests
go test ./...
```

## IDE Support

### VS Code Extension

The project includes a VS Code extension with:

- Syntax highlighting
- Language Server Protocol support
- Code completion
- Hover information

To install the extension:

1. Build the main executable
2. Open the `vscode` folder in VS Code
3. Press F5 to run the extension

### Language Server

The LSP server provides:

- Syntax error reporting
- Type information on hover
- Code completion suggestions
- Go-to-definition support

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Author

Matan Tsuberi <tsuberim@gmail.com>
