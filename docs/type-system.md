# Type System

Fun uses the Hindley-Milner type system, which provides automatic type inference with polymorphic types.

## Table of Contents

- [Type Inference](#type-inference)
- [Type Constructors](#type-constructors)
- [Type Variables](#type-variables)
- [Type Schemes](#type-schemes)
- [Type Annotations](#type-annotations)
- [Type Checking](#type-checking)
- [Common Type Errors](#common-type-errors)

## Type Inference

Fun automatically infers types for expressions without requiring explicit type annotations.

### Basic Inference

```fun
# Integer literal
42                    # inferred as: Int

# String literal
`hello`               # inferred as: Str

# Function application
inc(5)                # inferred as: Int (if inc : Lam<Int, Int>)

# Lambda expression
\x -> x + 1           # inferred as: Lam<Int, Int>
```

### Polymorphic Inference

```fun
# Identity function
id = \x -> x          # inferred as: Lam<a, a>

# Constant function
const = \x, y -> x    # inferred as: Lam<a, b, a>

# List operations
head = \list ->       # inferred as: Lam<List<a>, a>
    when list is
        Cons x xs -> x;
        Nil -> error "empty list"
```

## Type Constructors

### Basic Types

```fun
Int                   # Integer type
Str                   # String type
Bool                  # Boolean type (True/False)
```

### Function Types

```fun
Lam<a, b>             # Function from type a to type b
Lam<a, b, c>          # Function from a and b to c
```

### List Types

```fun
List<a>               # List containing elements of type a
List<Int>             # List of integers
List<Str>             # List of strings
```

### Record Types

```fun
{name: Str, age: Int}     # Record with named fields
{a: Int, b: Str}          # Generic record type
```

## Type Variables

Type variables represent unknown types and enable polymorphism.

### Universal Quantification

```fun
# ∀a. a → a
id : Lam<a, a>
id = \x -> x

# ∀a. ∀b. a → b → a
const : Lam<a, b, a>
const = \x, y -> x
```

### Type Variable Scope

```fun
# Different type variables in same scope
swap : Lam<a, b, {first: b, second: a}>
swap = \x, y -> {first: y, second: x}

# Same type variable used multiple times
pair : Lam<a, {first: a, second: a}>
pair = \x, y -> {first: x, second: y}
```

## Type Schemes

Type schemes represent polymorphic types with universal quantification.

### Generalization

```fun
# Type inference produces a type scheme
map : Lam<Lam<a, b>, List<a>, List<b>>
map = \f, list ->
    when list is
        Cons x xs -> Cons (f x) (map f xs);
        Nil -> Nil
```

### Instantiation

```fun
# Type scheme is instantiated for specific use
int_map = map(inc)    # instantiates to: Lam<List<Int>, List<Int>>
str_map = map(length) # instantiates to: Lam<List<Str>, List<Int>>
```

## Type Annotations

Explicit type annotations help with documentation and catch type errors early.

### Function Types

```fun
# Simple function
inc : Lam<Int, Int>
inc = \x -> x + 1

# Higher-order function
compose : Lam<Lam<b, c>, Lam<a, b>, Lam<a, c>>
compose = \f, g, x -> f(g(x))
```

### Record Types

```fun
# Person record
person : {name: Str, age: Int}
person = {name: "Alice", age: 30}

# Generic record
point : {x: a, y: a}
point = {x: 10, y: 20}
```

### List Types

```fun
# Homogeneous list
numbers : List<Int>
numbers = [1, 2, 3, 4, 5]

# List of records
people : List<{name: Str, age: Int}>
people = [
    {name: "Alice", age: 30},
    {name: "Bob", age: 25}
]
```

## Type Checking

The type checker ensures type safety at compile time.

### Type Unification

```fun
# Unification of function types
apply : Lam<Lam<a, b>, a, b>
apply = \f, x -> f(x)

# Unification of record types
get_name : Lam<{name: Str, age: Int}, Str>
get_name = \person -> person.name
```

### Type Constraints

```fun
# Equality constraint
equal : Lam<a, a, Bool>
equal = \x, y -> x == y

# Record constraint
has_name : Lam<{name: Str}, Str>
has_name = \record -> record.name
```

## Common Type Errors

### Type Mismatch

```fun
# Error: cannot apply Int to function expecting Lam<Int, Int>
inc(5, 3)             # inc expects 1 argument, got 2

# Error: cannot access field 'age' on record without that field
point.age             # point has type {x: Int, y: Int}, no 'age' field
```

### Unification Failure

```fun
# Error: cannot unify Int with Str
\x -> x + `hello`     # + expects Int, got Str

# Error: cannot unify different record types
{name: "Alice"}.age   # record has no 'age' field
```

### Polymorphic Type Errors

```fun
# Error: type variable 'a' cannot be unified with Int
id(5) + `hello`       # id returns 'a', but + expects Int

# Error: cannot unify List<a> with Int
head(42)              # head expects List<a>, got Int
```

## Advanced Type Features

### Existential Types

```fun
# Existential type for hiding implementation details
data_type : {create: Lam<Int, a>, show: Lam<a, Str>}
data_type = {
    create: \x -> x,
    show: \x -> `{x}`
}
```

### Higher-Kinded Types

```fun
# Type constructor for Maybe
data Maybe a = Just a | Nothing

# Higher-order type constructor
map_maybe : Lam<Lam<a, b>, Maybe<a>, Maybe<b>>
map_maybe = \f, maybe ->
    when maybe is
        Just x -> Just (f x);
        Nothing -> Nothing
```

### Type Families

```fun
# Type-level functions
type family ListLength a where
    ListLength (List a) = Int
    ListLength a = Int

# Usage in type annotations
length : Lam<List<a>, ListLength (List<a>)>
length = \list ->
    when list is
        Cons x xs -> 1 + length(xs);
        Nil -> 0
```

## Type Inference Algorithm

The type inference algorithm works in three phases:

1. **Constraint Generation**: Generate type constraints from expressions
2. **Constraint Solving**: Solve constraints using unification
3. **Generalization**: Generalize types to type schemes

### Example

```fun
# Step 1: Generate constraints
map = \f, list -> ...  # f: a → b, list: List<a>, result: List<b>

# Step 2: Solve constraints
# map: (a → b) → List<a> → List<b>

# Step 3: Generalize
# map: ∀a. ∀b. (a → b) → List<a> → List<b>
```

## Best Practices

### Type Annotations

```fun
# Use type annotations for public APIs
public_api : Lam<Int, Str>
public_api = \x -> `{x}`

# Use type annotations for complex functions
complex_function : Lam<List<Int>, List<Int>>
complex_function = \list -> map(inc, list)
```

### Type Safety

```fun
# Use Maybe for optional values
safe_head : Lam<List<a>, Maybe<a>>
safe_head = \list ->
    when list is
        Cons x xs -> Just x;
        Nil -> Nothing

# Use Result for operations that can fail
safe_divide : Lam<Int, Int, Maybe<Int>>
safe_divide = \x, y ->
    when y == 0 is
        True -> Nothing;
        False -> Just (x / y)
```

### Polymorphic Programming

```fun
# Write generic functions when possible
id : Lam<a, a>
id = \x -> x

# Avoid unnecessary type constraints
length : Lam<List<a>, Int>  # works for any list type
length = \list ->
    when list is
        Cons x xs -> 1 + length(xs);
        Nil -> 0
```

The type system ensures that well-typed programs cannot have runtime type errors, providing strong guarantees about program correctness. 