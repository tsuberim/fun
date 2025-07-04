# Fun Programming Language Documentation

Welcome to the Fun programming language! Fun is a purely functional programming language with type inference, built in Go with Tree-sitter parsing and LSP support.

## Quick Navigation

- **[Language Guide](language-guide.md)** - Learn the Fun language syntax and features
- **[Type System](type-system.md)** - Understanding Hindley-Milner type inference
- **[Standard Library](stdlib.md)** - Built-in functions and types
- **[Modules](modules.md)** - Code organization and imports
- **[REPL Guide](repl.md)** - Interactive development environment
- **[IDE Support](ide-support.md)** - VS Code extension and LSP features
- **[Examples](examples.md)** - Code examples and tutorials
- **[Internals](internals.md)** - Language implementation details

## What is Fun?

Fun is a functional programming language designed for:

- **Type Safety**: Hindley-Milner type inference with polymorphic types
- **Functional Programming**: Lambda expressions, pattern matching, immutable data
- **Developer Experience**: REPL, LSP support, VS Code integration
- **Performance**: Tree-sitter parsing, efficient evaluation

## Key Features

### Functional Programming
- Lambda expressions: `\x -> x + 1`
- Pattern matching with `when` expressions
- Immutable data structures
- Higher-order functions

### Type System
- Automatic type inference
- Polymorphic types
- Type annotations when needed
- Compile-time type checking

### Developer Tools
- Interactive REPL
- Language Server Protocol
- VS Code extension
- Syntax highlighting

## Getting Started

1. **Install**: Clone and build the project
2. **Try REPL**: Run `./fun` for interactive development
3. **Run Programs**: `./fun examples/basic/main`
4. **IDE Setup**: Install the VS Code extension

## Example

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

## Community

- **Author**: Matan Tsuberi <tsuberim@gmail.com>
- **License**: MIT
- **Repository**: [GitHub](https://github.com/your-repo/fun)

---

Start with the [Language Guide](language-guide.md) to learn Fun programming!

