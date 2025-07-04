# Fun Language Guide

This guide covers the syntax and features of the Fun programming language.

## Table of Contents

- [Basic Syntax](#basic-syntax)
- [Types](#types)
- [Functions](#functions)
- [Pattern Matching](#pattern-matching)
- [Records](#records)
- [Lists](#lists)
- [Type Annotations](#type-annotations)
- [Comments](#comments)

## Basic Syntax

### Comments

```fun
# This is a single-line comment
# Comments start with # and continue to the end of the line
```

### Expressions

Every expression in Fun evaluates to a value. Expressions can be:

- **Literals**: `42`, `` `hello` ``
- **Variables**: `x`, `myVariable`
- **Function calls**: `f(x, y)`
- **Operators**: `x + y`, `a == b`

## Types

### Basic Types

#### Integers

```fun
42          # positive integer
-17         # negative integer
0           # zero
```

#### Strings

```fun
`hello world`           # simple string
`Hello {name}!`         # string template with expression
`Value: {x + y}`        # string template with computation
```

String templates use curly braces `{}` to embed expressions.

## Functions

### Lambda Expressions

```fun
# Single parameter
inc = \x -> x + 1

# Multiple parameters
add = \x, y -> x + y

# No parameters (unit function)
const = \ -> 42
```

### Function Application

```fun
# Direct application
inc(5)                  # result: 6
add(3, 4)               # result: 7

# Operator syntax for infix operators
3 + 4                   # same as +(3, 4)
10 == 5                 # same as ==(10, 5)
```

### Higher-Order Functions

```fun
# Function that takes a function as parameter
apply_twice = \f, x -> f(f(x))

# Usage
apply_twice(inc, 5)     # result: 7 (inc(inc(5)))
```

## Pattern Matching

### When Expressions

Pattern matching is done with `when` expressions:

```fun
when value is
    Just x -> x;
    Nothing -> 0
else 42
```

### Constructor Patterns

```fun
# Define a data type
data Maybe a = Just a | Nothing

# Pattern match on constructors
safe_head = \list ->
    when list is
        Cons head tail -> Just head;
        Nil -> Nothing

# Pattern match with payload binding
process_result = \result ->
    when result is
        Success value -> value;
        Error msg -> 0
```

### Nested Patterns

```fun
# Pattern match on nested structures
process_nested = \data ->
    when data is
        Just (Cons x xs) -> x;
        Just Nil -> 0;
        Nothing -> -1
```

## Records

### Record Construction

```fun
# Simple record
person = {name: "Alice", age: 30}

# Record with expressions
point = {x: 10, y: 20, distance: sqrt(x*x + y*y)}
```

### Record Access

```fun
person.name              # result: "Alice"
person.age               # result: 30
point.distance           # result: computed distance
```

### Record Updates

```fun
# Create new record with updated field
older_person = {person | age: person.age + 1}
```

### Record Types

```fun
# Type annotation for records
person : {name: Str, age: Int}
person = {name: "Bob", age: 25}
```

## Lists

### List Construction

```fun
# Empty list
empty = []

# List with elements
numbers = [1, 2, 3, 4, 5]

# List with expressions
squares = [x*x | x <- [1, 2, 3, 4, 5]]
```

### List Operations

```fun
# Head and tail (pattern matching)
head = \list ->
    when list is
        Cons x xs -> x;
        Nil -> error "empty list"

# Length
length = \list ->
    when list is
        Cons x xs -> 1 + length(xs);
        Nil -> 0
```

## Type Annotations

### Explicit Types

```fun
# Function type annotation
factorial : Lam<Int, Int>
factorial = \n -> when n == 0 is
    True -> 1;
    False -> n * factorial(n - 1)

# Variable type annotation
count : Int
count = 42
```

### Type Constructors

```fun
# Generic types
id : Lam<a, a>
id = \x -> x

# Complex types
map : Lam<Lam<a, b>, List<a>, List<b>>
map = \f, list -> when list is
    Cons x xs -> Cons (f x) (map f xs);
    Nil -> Nil
```

### Type Variables

```fun
# Polymorphic functions
const : Lam<a, b, a>
const = \x, y -> x

# Type variables are lowercase
swap : Lam<a, b, {first: b, second: a}>
swap = \x, y -> {first: y, second: x}
```

## Built-in Functions

### Arithmetic

```fun
# Addition
5 + 3                    # result: 8

# Subtraction  
10 - 4                   # result: 6

# Comparison
5 == 5                   # result: True
3 == 7                   # result: False
```

### Fixed-Point Combinator

```fun
# Recursion using fix
factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(n - 1)
)
```

## Best Practices

### Naming Conventions

```fun
# Variables and functions: lowercase with underscores
my_function = \x -> x + 1
user_name = "Alice"

# Constructors: start with uppercase
data Color = Red | Green | Blue

# Type variables: single lowercase letters
id : Lam<a, a>
```

### Code Organization

```fun
# Group related functions together
# Math utilities
inc = \x -> x + 1
dec = \x -> x - 1
double = \x -> x * 2

# Export as record
{
    inc: inc,
    dec: dec,
    double: double
}
```

### Error Handling

```fun
# Use Maybe types for optional values
safe_divide = \x, y ->
    when y == 0 is
        True -> Nothing;
        False -> Just (x / y)

# Use Result types for operations that can fail
parse_int = \str ->
    when str is
        "0" -> Success 0;
        "1" -> Success 1;
        _ -> Error "invalid number"
```

## Common Patterns

### List Processing

```fun
# Map over list
map = \f, list ->
    when list is
        Cons x xs -> Cons (f x) (map f xs);
        Nil -> Nil

# Filter list
filter = \pred, list ->
    when list is
        Cons x xs -> when pred(x) is
            True -> Cons x (filter pred xs);
            False -> filter pred xs;
        Nil -> Nil
```

### Recursion

```fun
# Tail recursion with accumulator
factorial_tail = \n ->
    let
        fact_acc = \n, acc ->
            when n == 0 is
                True -> acc;
                False -> fact_acc(n - 1, n * acc)
    in
        fact_acc(n, 1)
```

This guide covers the essential syntax and features of Fun. For more advanced topics, see the [Type System](type-system.md) and [Examples](examples.md) documentation. 