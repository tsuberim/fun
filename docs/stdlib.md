# Standard Library

The Fun standard library provides essential functions and types for common programming tasks.

## Table of Contents

- [Arithmetic Functions](#arithmetic-functions)
- [Comparison Functions](#comparison-functions)
- [Recursion Functions](#recursion-functions)
- [Built-in Types](#built-in-types)
- [Type Constructors](#type-constructors)

## Arithmetic Functions

### Addition (`+`)

```fun
+ : Lam<Int, Int, Int>
```

Adds two or more integers.

```fun
5 + 3                    # result: 8
10 + 20 + 30             # result: 60
```

**Examples:**
```fun
# Sum a list of numbers
sum = \list ->
    when list is
        Cons x xs -> x + sum(xs);
        Nil -> 0

sum([1, 2, 3, 4, 5])     # result: 15
```

### Subtraction (`-`)

```fun
- : Lam<Int, Int, Int>
```

Subtracts integers. The first argument is subtracted by all subsequent arguments.

```fun
10 - 3                   # result: 7
20 - 5 - 3               # result: 12 (20 - 5 - 3)
```

**Examples:**
```fun
# Decrement function
dec = \x -> x - 1

# Calculate difference
difference = \a, b -> a - b
```

## Comparison Functions

### Equality (`==`)

```fun
== : Lam<a, a, Bool>
```

Compares two values for equality. Works with any type that supports equality comparison.

```fun
5 == 5                   # result: True
5 == 3                   # result: False
`hello` == `hello`       # result: True
`hello` == `world`       # result: False
```

**Examples:**
```fun
# Check if list contains element
contains = \list, elem ->
    when list is
        Cons x xs -> when x == elem is
            True -> True;
            False -> contains(xs, elem);
        Nil -> False

contains([1, 2, 3], 2)   # result: True
contains([1, 2, 3], 5)   # result: False
```

## Recursion Functions

### Fixed-Point Combinator (`fix`)

```fun
fix : Lam<Lam<a, a>, a>
```

The Y combinator that enables recursion in a pure functional language.

```fun
# Factorial using fix
factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(n - 1)
)

factorial(5)              # result: 120
```

**Examples:**
```fun
# Fibonacci using fix
fibonacci = fix(\rec -> \n ->
    when n == 0 is
        True -> 0;
        False -> when n == 1 is
            True -> 1;
            False -> rec(n - 1) + rec(n - 2)
)

fibonacci(10)             # result: 55

# Sum range using fix
sum_range = fix(\rec -> \n ->
    when n == 0 is
        True -> 0;
        False -> n + rec(n - 1)
)

sum_range(100)            # result: 5050
```

## Built-in Types

### Boolean Values

```fun
True                     # Boolean true value
False                    # Boolean false value
```

**Usage:**
```fun
# Boolean logic
is_even = \n -> n % 2 == 0
is_positive = \n -> n > 0

# Conditional logic
when is_even(4) is
    True -> "even";
    False -> "odd"
```

### Integer Type (`Int`)

```fun
42                       # Positive integer
-17                      # Negative integer
0                        # Zero
```

**Properties:**
- Unbounded precision (limited by available memory)
- Supports all arithmetic operations
- Can be compared for equality

### String Type (`Str`)

```fun
`hello world`            # Simple string
`Hello {name}!`          # String template
`Value: {x + y}`         # String with expression
```

**String Templates:**
```fun
name = "Alice"
greeting = `Hello {name}!`    # result: "Hello Alice!"

x = 5
y = 3
result = `{x} + {y} = {x + y}`  # result: "5 + 3 = 8"
```

## Type Constructors

### Function Type (`Lam`)

```fun
Lam<a, b>                # Function from a to b
Lam<a, b, c>             # Function from a and b to c
```

**Examples:**
```fun
# Function type annotations
inc : Lam<Int, Int>
inc = \x -> x + 1

add : Lam<Int, Int, Int>
add = \x, y -> x + y

# Higher-order function
compose : Lam<Lam<b, c>, Lam<a, b>, Lam<a, c>>
compose = \f, g, x -> f(g(x))
```

### List Type (`List`)

```fun
List<a>                  # List containing elements of type a
```

**List Constructors:**
```fun
Cons x xs                # Non-empty list with head x and tail xs
Nil                      # Empty list
```

**Examples:**
```fun
# List construction
empty = []
numbers = [1, 2, 3, 4, 5]

# List operations
head = \list ->
    when list is
        Cons x xs -> x;
        Nil -> error "empty list"

tail = \list ->
    when list is
        Cons x xs -> xs;
        Nil -> error "empty list"
```

### Record Type

```fun
{field1: a, field2: b}   # Record with named fields
```

**Examples:**
```fun
# Record construction
person = {name: "Alice", age: 30}
point = {x: 10, y: 20}

# Record access
person.name              # result: "Alice"
point.x                  # result: 10

# Record type annotation
person : {name: Str, age: Int}
person = {name: "Bob", age: 25}
```

## Common Patterns

### List Processing

```fun
# Map function
map = \f, list ->
    when list is
        Cons x xs -> Cons (f x) (map f xs);
        Nil -> Nil

# Filter function
filter = \pred, list ->
    when list is
        Cons x xs -> when pred(x) is
            True -> Cons x (filter pred xs);
            False -> filter pred xs;
        Nil -> Nil

# Fold function
fold = \f, init, list ->
    when list is
        Cons x xs -> fold f (f init x) xs;
        Nil -> init
```

### Error Handling

```fun
# Maybe type for optional values
data Maybe a = Just a | Nothing

# Safe division
safe_divide = \x, y ->
    when y == 0 is
        True -> Nothing;
        False -> Just (x / y)

# Safe head
safe_head = \list ->
    when list is
        Cons x xs -> Just x;
        Nil -> Nothing
```

### Recursion Patterns

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

# Mutual recursion
even = \n ->
    when n == 0 is
        True -> True;
        False -> odd(n - 1)

odd = \n ->
    when n == 0 is
        True -> False;
        False -> even(n - 1)
```

## Best Practices

### Use Built-in Functions

```fun
# Prefer built-in + over manual addition
sum = fold (+) 0         # Use built-in + function

# Use == for equality comparisons
contains = \list, elem ->
    any (\x -> x == elem) list
```

### Leverage Type Inference

```fun
# Let the compiler infer types when possible
id = \x -> x             # inferred as: Lam<a, a>

# Add type annotations for clarity
public_api : Lam<Int, Str>
public_api = \x -> `{x}`
```

### Use Pattern Matching

```fun
# Prefer pattern matching over manual checks
safe_head = \list ->
    when list is
        Cons x xs -> Just x;
        Nil -> Nothing

# Instead of manual null checks
unsafe_head = \list ->
    when list == [] is
        True -> error "empty list";
        False -> head(list)
```

The standard library provides the foundation for building complex applications in Fun. Combine these functions with custom types and higher-order functions to create powerful abstractions. 