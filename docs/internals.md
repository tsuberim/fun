# Internals

This document describes the internal implementation of the Fun programming language.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Parser Implementation](#parser-implementation)
- [Type System](#type-system)
- [Evaluation](#evaluation)
- [Module System](#module-system)
- [LSP Implementation](#lsp-implementation)
- [Performance Considerations](#performance-considerations)

## Architecture Overview

The Fun language implementation consists of several key components:

```
fun/
├── main.go              # Entry point and REPL
├── internal/            # Core language implementation
│   ├── program.go       # Program execution and module management
│   ├── expr.go          # Abstract syntax tree (AST)
│   ├── type.go          # Type system and inference
│   ├── value.go         # Runtime values and evaluation
│   ├── env.go           # Environment management
│   └── lsp.go           # Language Server Protocol
├── tree-sitter-fun/     # Parser implementation
│   ├── grammar.js       # Tree-sitter grammar
│   └── bindings/        # Language bindings
└── vscode/              # VS Code extension
```

### Execution Flow

1. **Parsing**: Source code → Tree-sitter AST → Fun AST
2. **Type Checking**: AST → Type inference → Type errors
3. **Evaluation**: AST + Environment → Runtime values
4. **Output**: Pretty-printed results

## Parser Implementation

### Tree-sitter Grammar

The parser uses Tree-sitter for robust parsing with error recovery.

```javascript
// grammar.js - Key grammar rules
module.exports = grammar({
  name: "fun",
  rules: {
    source_file: $ => $._inner_block,
    _expr: $ => choice($.int, $.str, $.var, $.app, $.lam, $.rec, $.when, $.list),
    int: $ => /\d+/,
    str: $ => seq('`', repeat(choice($.lit_str, seq('{', $._expr, '}'))), '`'),
    app: $ => seq($._expr, '(', sep($._expr, ','), ')'),
    lam: $ => seq('\\', sep($.var, ','), '->', $._expr),
    when: $ => seq('when', $._expr, 'is', sep1($.when_clause, ';'), optional(seq('else', $._expr))),
    // ... more rules
  }
});
```

### AST Construction

The parser converts Tree-sitter nodes to Fun AST nodes:

```go
// expr.go - AST node types
type Expr interface {
    expr()
    Pretty(indent int) string
}

type Int struct {
    Value int
}

type App struct {
    Fn   Expr
    Args []Expr
}

type Lam struct {
    Params []string
    Body   Expr
}

type When struct {
    Value   Expr
    Options []WhenClause
    Else    Expr
}
```

### Parsing Process

```go
// program.go - Parsing workflow
func (p *Program) Run(source []byte, importPath string) (*Module, error) {
    // 1. Parse with Tree-sitter
    tree := p.parser.Parse(source, nil)
    node := tree.RootNode()
    
    // 2. Convert to Fun AST
    expr, err := fromNode(node, source)
    if err != nil {
        return nil, errors.WithMessagef(err, "failed to parse module: %s", importPath)
    }
    
    // 3. Type check
    _, t, err := p.inferer.Infer(expr, p.env.Types())
    if err != nil {
        return nil, errors.WithMessagef(err, "failed to infer module type: %s", importPath)
    }
    
    // 4. Evaluate
    val, err := p.evaluator.Eval(expr, p.env.Values())
    if err != nil {
        return nil, errors.WithMessagef(err, "failed to evaluate module: %s", importPath)
    }
    
    return &Module{ImportPath: importPath, Expr: expr, Val: val, Type: generalize(t)}, nil
}
```

## Type System

### Hindley-Milner Type Inference

Fun implements the Hindley-Milner type system with automatic type inference.

#### Type Representation

```go
// type.go - Type system
type Type interface {
    typ()
    freeVars() *strset.Set
    apply(subst *Subst) Type
    Pretty(i int) string
}

type TypeVar struct {
    Name string
}

type TypeCons struct {
    Name string
    Args []Type
}

type TypeRec struct {
    Entries map[string]Type
    RestVar *TypeVar
    Union   bool
}

type Scheme struct {
    Forall []string  // Universal quantification
    Type   Type
}
```

#### Type Inference Algorithm

```go
// type.go - Inference implementation
func (i *Inferrer) Infer(expr Expr, env *TypeEnv) (subst *Subst, typ Type, err error) {
    switch expr := expr.(type) {
    case *Int:
        return &Subst{Subst: map[string]Type{}}, &TypeCons{Name: "Int"}, nil
        
    case *Var:
        scheme, has := env.Types[expr.Name]
        if !has {
            return nil, nil, errors.Errorf("undefined variable: %s", expr.Name)
        }
        return &Subst{Subst: map[string]Type{}}, i.instantiate(scheme), nil
        
    case *App:
        // Infer function type
        subst1, fnType, err := i.Infer(expr.Fn, env)
        if err != nil {
            return nil, nil, err
        }
        
        // Infer argument types
        var argTypes []Type
        for _, arg := range expr.Args {
            subst2, argType, err := i.Infer(arg, env)
            if err != nil {
                return nil, nil, err
            }
            argTypes = append(argTypes, argType)
        }
        
        // Unify function type with argument types
        resultType := i.freshVar()
        expectedType := &TypeCons{Name: "Lam", Args: append(argTypes, resultType)}
        
        subst3, err := i.unify(fnType, expectedType)
        if err != nil {
            return nil, nil, err
        }
        
        // Compose substitutions
        finalSubst := subst3.compose(subst2).compose(subst1)
        return finalSubst, resultType.apply(finalSubst), nil
    }
    // ... more cases
}
```

#### Type Unification

```go
// type.go - Unification algorithm
func (i *Inferrer) unify(t1, t2 Type) (*Subst, error) {
    switch t1 := t1.(type) {
    case *TypeVar:
        if t1.freeVars().Has(t1.Name) {
            return nil, errors.Errorf("occurs check failed")
        }
        return &Subst{Subst: map[string]Type{t1.Name: t2}}, nil
        
    case *TypeCons:
        if t2Cons, ok := t2.(*TypeCons); ok {
            if t1.Name != t2Cons.Name || len(t1.Args) != len(t2Cons.Args) {
                return nil, errors.Errorf("type mismatch")
            }
            
            subst := &Subst{Subst: map[string]Type{}}
            for j, arg1 := range t1.Args {
                arg2 := t2Cons.Args[j]
                newSubst, err := i.unify(arg1.apply(subst), arg2.apply(subst))
                if err != nil {
                    return nil, err
                }
                subst = newSubst.compose(subst)
            }
            return subst, nil
        }
    }
    return nil, errors.Errorf("cannot unify types")
}
```

## Evaluation

### Runtime Values

```go
// value.go - Runtime value representation
type Val interface {
    val()
    Pretty(indent int) string
}

type Int struct {
    Value int
}

type LitStr struct {
    Value string
}

type ListVal struct {
    Items []Val
}

type RecVal struct {
    Entries map[string]Val
}

type Closure struct {
    Params []string
    Env    map[string]Val
    Body   Expr
}

type Builtin struct {
    Name string
    Impl func(args []Val) (Val, error)
}
```

### Evaluation Strategy

Fun uses call-by-value evaluation with lexical scoping.

```go
// value.go - Evaluation implementation
func (e *Evaluator) Eval(expr Expr, env map[string]Val) (Val, error) {
    switch expr := expr.(type) {
    case *Int:
        return expr, nil
        
    case *Var:
        val, has := env[expr.Name]
        if !has {
            return nil, errors.Errorf("undefined variable: %s", expr.Name)
        }
        return val, nil
        
    case *Lam:
        return &Closure{
            Params: expr.Params,
            Env:    env,
            Body:   expr.Body,
        }, nil
        
    case *App:
        // Evaluate function
        fn, err := e.Eval(expr.Fn, env)
        if err != nil {
            return nil, err
        }
        
        // Handle builtin functions
        if builtin, ok := fn.(*Builtin); ok {
            var args []Val
            for _, arg := range expr.Args {
                val, err := e.Eval(arg, env)
                if err != nil {
                    return nil, err
                }
                args = append(args, val)
            }
            return builtin.Impl(args)
        }
        
        // Handle closures
        clos, ok := fn.(*Closure)
        if !ok {
            return nil, errors.Errorf("cannot apply non closure value")
        }
        
        if len(expr.Args) != len(clos.Params) {
            return nil, errors.Errorf("invalid number of arguments")
        }
        
        // Create new environment with arguments
        newEnv := maps.Clone(clos.Env)
        for i, arg := range expr.Args {
            val, err := e.Eval(arg, env)
            if err != nil {
                return nil, err
            }
            newEnv[clos.Params[i]] = val
        }
        
        return e.Eval(clos.Body, newEnv)
    }
    // ... more cases
}
```

### Built-in Functions

```go
// env.go - Standard environment
func NewStdEnv(program *Program) *Env {
    return &Env{
        Items: map[string]Item{
            "+": {
                Type: &Scheme{Type: &TypeCons{Name: "Lam", Args: []Type{
                    &TypeCons{Name: "Int"}, &TypeCons{Name: "Int"}, &TypeCons{Name: "Int"},
                }}},
                Val: &Builtin{
                    Name: "+",
                    Impl: func(args []Val) (Val, error) {
                        sum := 0
                        for _, arg := range args {
                            i, ok := arg.(*Int)
                            if !ok {
                                return nil, errors.Errorf("invalid sum value type")
                            }
                            sum += i.Value
                        }
                        return &Int{Value: sum}, nil
                    },
                },
            },
            "fix": {
                Type: &Scheme{Forall: []string{"a"}, Type: &TypeCons{Name: "Lam", Args: []Type{
                    &TypeCons{Name: "Lam", Args: []Type{&TypeVar{Name: "a"}, &TypeVar{Name: "a"}}},
                    &TypeVar{Name: "a"},
                }}},
                Val: &Builtin{
                    Name: "fix",
                    Impl: func(args []Val) (Val, error) {
                        // Fixed-point combinator implementation
                        // ... implementation details
                    },
                },
            },
        },
    }
}
```

## Module System

### Module Loading

```go
// program.go - Module management
func (p *Program) Import(importPath string) (*Module, error) {
    // Check if already loaded
    mod, has := p.Modules[importPath]
    if has {
        return mod, nil
    }
    
    // Load and parse module
    importedMod, err := p.importModule(importPath)
    if err != nil {
        return nil, err
    }
    
    // Cache module
    p.Modules[importPath] = importedMod
    return importedMod, nil
}

func (p *Program) importModule(importPath string) (*Module, error) {
    // Resolve path relative to current module
    fullPath := path.Join(path.Dir(p.currentImportPath()), importPath)
    
    // Read source file
    source, err := os.ReadFile(fullPath)
    if err != nil {
        return nil, errors.WithMessagef(err, "failed to read module `%s`", importPath)
    }
    
    // Parse and evaluate module
    return p.Run(source, importPath)
}
```

### Import Resolution

```go
// program.go - Import stack management
type Program struct {
    // ... other fields
    importStack []string  // Track import chain for relative paths
}

func (p *Program) currentImportPath() string {
    if len(p.importStack) == 0 {
        return ""
    }
    return p.importStack[len(p.importStack)-1]
}

func (p *Program) Run(source []byte, importPath string) (*Module, error) {
    if importPath != InlineModule {
        p.importStack = append(p.importStack, importPath)
        defer func() {
            p.importStack = p.importStack[:len(p.importStack)-1]
        }()
    }
    
    // ... parsing, type checking, evaluation
}
```

## LSP Implementation

### Language Server

```go
// lsp.go - LSP server implementation
func LSPServer() {
    logs.Init(log.New(os.Stderr, "", 0777))
    
    server := lsp.NewServer(&lsp.Options{
        CompletionProvider: &defines.CompletionOptions{
            TriggerCharacters: &[]string{"."},
        },
    })
    
    // Hover information
    server.OnHover(func(ctx context.Context, req *defines.HoverParams) (result *defines.Hover, err error) {
        // Extract position and provide type information
        return &defines.Hover{
            Contents: defines.MarkupContent{
                Kind:  defines.MarkupKindPlainText,
                Value: "Type information here",
            },
        }, nil
    })
    
    // Code completion
    server.OnCompletion(func(ctx context.Context, req *defines.CompletionParams) (result *[]defines.CompletionItem, err error) {
        // Provide completion suggestions
        return &[]defines.CompletionItem{{
            Label:      "suggestion",
            Kind:       &defines.CompletionItemKindText,
            InsertText: &"suggestion",
        }}, nil
    })
    
    server.Run()
}
```

### VS Code Extension

```typescript
// extension.ts - VS Code extension
export function activate(context: ExtensionContext) {
    const serverOptions: Executable = {
        transport: TransportKind.stdio,
        command: '/path/to/fun',
        args: ['lsp']
    };
    
    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'fun' }]
    };
    
    client = new LanguageClient(
        'funLanguageServer',
        'Fun Language Server',
        serverOptions,
        clientOptions
    );
    
    client.start();
}
```

## Performance Considerations

### Memory Management

- **AST Sharing**: Reuse AST nodes where possible
- **Environment Cloning**: Use efficient environment copying
- **Type Caching**: Cache inferred types to avoid recomputation

### Evaluation Optimization

- **Lazy Evaluation**: Consider lazy evaluation for large data structures
- **Tail Call Optimization**: Optimize recursive functions
- **Built-in Functions**: Use native implementations for common operations

### Parser Performance

- **Tree-sitter**: Leverages incremental parsing
- **Error Recovery**: Robust parsing with error recovery
- **Memory Efficient**: Minimal memory overhead during parsing

### Type Inference

- **Constraint Solving**: Efficient unification algorithm
- **Type Generalization**: Smart generalization to avoid over-specialization
- **Type Caching**: Cache type schemes to avoid recomputation

## Future Improvements

### Planned Features

1. **Garbage Collection**: Implement proper memory management
2. **Optimization Passes**: Add compiler optimizations
3. **Standard Library**: Expand built-in functions
4. **Error Recovery**: Better error messages and recovery
5. **Performance Profiling**: Add performance measurement tools

### Architecture Enhancements

1. **Modular Compiler**: Separate parsing, type checking, and evaluation
2. **Bytecode VM**: Consider bytecode interpretation for better performance
3. **JIT Compilation**: Just-in-time compilation for hot code paths
4. **Parallel Evaluation**: Parallel evaluation of independent expressions

The internal implementation provides a solid foundation for the Fun programming language with room for future enhancements and optimizations. 