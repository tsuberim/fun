# REPL Guide

The Fun REPL (Read-Eval-Print Loop) provides an interactive development environment for experimenting with code and testing functions.

## Table of Contents

- [Starting the REPL](#starting-the-repl)
- [Basic Usage](#basic-usage)
- [REPL Features](#repl-features)
- [Working with Modules](#working-with-modules)
- [Debugging](#debugging)
- [Tips and Tricks](#tips-and-tricks)

## Starting the REPL

### Basic Startup

```bash
# Start the REPL
./fun

# You'll see the prompt
>
```

### REPL Prompt

The REPL shows a simple `>` prompt and waits for input:

```bash
$ ./fun
>
```

## Basic Usage

### Simple Expressions

```fun
> 42
42 : Int

> `hello world`
hello world : Str

> 5 + 3
8 : Int

> True
True : Bool
```

### Function Definitions

```fun
> inc = \x -> x + 1
<closure> : Lam<Int, Int>

> inc(5)
6 : Int

> double = \x -> x * 2
<closure> : Lam<Int, Int>

> double(10)
20 : Int
```

### Complex Expressions

```fun
> factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(n - 1)
)
<closure> : Lam<Int, Int>

> factorial(5)
120 : Int
```

### Pattern Matching

```fun
> safe_head = \list ->
    when list is
        Cons x xs -> Just x;
        Nil -> Nothing
<closure> : Lam<List<a>, Maybe<a>>

> safe_head([1, 2, 3])
Just 1 : Maybe<Int>

> safe_head([])
Nothing : Maybe<a>
```

## REPL Features

### Type Information

The REPL automatically shows type information for all expressions:

```fun
> 42
42 : Int

> \x -> x + 1
<closure> : Lam<Int, Int>

> [1, 2, 3]
[1, 2, 3] : List<Int>
```

### Error Reporting

The REPL provides clear error messages:

```fun
> inc(5, 3)
Error: invalid number of arguments for function

> undefined_variable
Error: undefined variable: undefined_variable

> 5 + `hello`
Error: cannot apply non closure value of type Int
```

### Multi-line Input

The REPL supports multi-line expressions:

```fun
> factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(n - 1)
)
<closure> : Lam<Int, Int>
```

### Expression History

The REPL maintains a history of expressions (though not currently displayed):

```fun
> 1
1 : Int

> 2
2 : Int

> 3
3 : Int
# Previous expressions are remembered internally
```

## Working with Modules

### Importing in REPL

```fun
> import math from `./math`
<record> : {inc: Lam<Int, Int>, dec: Lam<Int, Int>}

> math.inc(5)
6 : Int

> math.dec(10)
9 : Int
```

### Creating Modules in REPL

```fun
> # Define functions
> inc = \x -> x + 1
<closure> : Lam<Int, Int>

> dec = \x -> x - 1
<closure> : Lam<Int, Int>

> # Create module record
> math_module = {inc: inc, dec: dec}
{inc: <closure>, dec: <closure>} : {inc: Lam<Int, Int>, dec: Lam<Int, Int>}

> # Use the module
> math_module.inc(5)
6 : Int
```

### Testing Module Functions

```fun
> # Import and test
> import utils from `./utils`
<record> : {head: Lam<List<a>, a>, tail: Lam<List<a>, List<a>>}

> utils.head([1, 2, 3])
1 : Int

> utils.tail([1, 2, 3])
[2, 3] : List<Int>
```

## Debugging

### Inspecting Values

```fun
> # Check function types
> inc
<closure> : Lam<Int, Int>

> # Check record structure
> person = {name: "Alice", age: 30}
{name: Alice, age: 30} : {name: Str, age: Int}

> person.name
Alice : Str

> person.age
30 : Int
```

### Testing Edge Cases

```fun
> # Test with edge cases
> factorial(0)
1 : Int

> factorial(1)
1 : Int

> factorial(10)
3628800 : Int

> # Test error conditions
> safe_head([])
Nothing : Maybe<a>
```

### Step-by-step Development

```fun
> # Build complex functions step by step
> add = \x, y -> x + y
<closure> : Lam<Int, Int, Int>

> add(3, 4)
7 : Int

> # Test the function
> add(10, 20)
30 : Int

> # Now use it in more complex expressions
> add(add(1, 2), add(3, 4))
10 : Int
```

## Tips and Tricks

### Quick Testing

```fun
> # Test arithmetic
> 1 + 2 * 3
7 : Int

> # Test string templates
> name = "Alice"
Alice : Str

> `Hello {name}!`
Hello Alice! : Str

> # Test boolean logic
> 5 == 5
True : Bool

> 3 == 7
False : Bool
```

### Function Composition

```fun
> # Define simple functions
> inc = \x -> x + 1
<closure> : Lam<Int, Int>

> double = \x -> x * 2
<closure> : Lam<Int, Int>

> # Compose functions
> compose = \f, g, x -> f(g(x))
<closure> : Lam<Lam<b, c>, Lam<a, b>, Lam<a, c>>

> # Test composition
> compose(inc, double, 5)
11 : Int  # inc(double(5)) = inc(10) = 11
```

### List Operations

```fun
> # Create lists
> numbers = [1, 2, 3, 4, 5]
[1, 2, 3, 4, 5] : List<Int>

> # Test list functions
> head = \list ->
    when list is
        Cons x xs -> x;
        Nil -> error "empty list"
<closure> : Lam<List<a>, a>

> head(numbers)
1 : Int

> # Test with empty list
> head([])
Error: empty list
```

### Type Exploration

```fun
> # Explore type inference
> id = \x -> x
<closure> : Lam<a, a>

> id(42)
42 : Int

> id(`hello`)
hello : Str

> id(True)
True : Bool

> # Test polymorphic functions
> const = \x, y -> x
<closure> : Lam<a, b, a>

> const(5, `hello`)
5 : Int

> const(`world`, 42)
world : Str
```

### Interactive Development

```fun
> # Start with simple functions
> inc = \x -> x + 1
<closure> : Lam<Int, Int>

> # Test the function
> inc(5)
6 : Int

> # Refine the function
> inc = \x -> x + 1
<closure> : Lam<Int, Int>

> # Build more complex functions
> sum_range = fix(\rec -> \n ->
    when n == 0 is
        True -> 0;
        False -> n + rec(n - 1)
)
<closure> : Lam<Int, Int>

> # Test the complex function
> sum_range(10)
55 : Int
```

## Exiting the REPL

To exit the REPL, use Ctrl+C (or Cmd+C on macOS):

```bash
> # Type Ctrl+C to exit
Goodbye!
$
```

## Best Practices

### Use Clear Variable Names

```fun
> # Good: descriptive names
> user_name = "Alice"
Alice : Str

> user_age = 30
30 : Int

> # Avoid: unclear names
> x = 42
42 : Int
```

### Test Incrementally

```fun
> # Test each part separately
> add = \x, y -> x + y
<closure> : Lam<Int, Int, Int>

> add(1, 2)
3 : Int

> # Then use in larger expressions
> add(add(1, 2), add(3, 4))
10 : Int
```

### Use Type Annotations for Clarity

```fun
> # Add type annotations for complex functions
> factorial : Lam<Int, Int>
> factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(n - 1)
)
<closure> : Lam<Int, Int>
```

The REPL is an excellent tool for learning Fun, testing functions, and developing code interactively. Use it to experiment with the language and verify your understanding of concepts. 