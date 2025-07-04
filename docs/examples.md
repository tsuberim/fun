# Examples and Tutorials

This document provides practical examples and tutorials for learning Fun programming.

## Table of Contents

- [Getting Started](#getting-started)
- [Basic Examples](#basic-examples)
- [Intermediate Examples](#intermediate-examples)
- [Advanced Examples](#advanced-examples)
- [Common Patterns](#common-patterns)
- [Complete Programs](#complete-programs)

## Getting Started

### Hello World

```fun
# hello.fun
`Hello, World!`
```

Run it:
```bash
./fun hello.fun
# Output: Hello, World!
```

### Simple Calculator

```fun
# calculator.fun
add = \x, y -> x + y
subtract = \x, y -> x - y
multiply = \x, y -> x * y

# Test the functions
result1 = add(5, 3)      # 8
result2 = subtract(10, 4) # 6
result3 = multiply(2, 7)  # 14

{
    add: add,
    subtract: subtract,
    multiply: multiply,
    results: [result1, result2, result3]
}
```

## Basic Examples

### Mathematical Functions

```fun
# math.fun
# Basic arithmetic functions

inc = \x -> x + 1
dec = \x -> x - 1
double = \x -> x * 2
square = \x -> x * x

# Test functions
test_inc = inc(5)        # 6
test_dec = dec(10)       # 9
test_double = double(4)  # 8
test_square = square(3)  # 9

{
    inc: inc,
    dec: dec,
    double: double,
    square: square,
    tests: [test_inc, test_dec, test_double, test_square]
}
```

### String Manipulation

```fun
# strings.fun
# String operations and templates

name = "Alice"
age = 30

# String templates
greeting = `Hello, {name}!`
info = `{name} is {age} years old`
calculation = `{age} + 5 = {age + 5}`

# String concatenation (using templates)
full_name = `{name} Smith`
message = `{greeting} {info}`

{
    greeting: greeting,
    info: info,
    calculation: calculation,
    full_name: full_name,
    message: message
}
```

### Boolean Logic

```fun
# logic.fun
# Boolean operations and conditional logic

# Boolean values
is_true = True
is_false = False

# Equality comparisons
equal_test = 5 == 5      # True
unequal_test = 3 == 7    # False

# Conditional logic
check_number = \x ->
    when x == 0 is
        True -> "zero";
        False -> when x > 0 is
            True -> "positive";
            False -> "negative"

# Test the function
test1 = check_number(0)   # "zero"
test2 = check_number(5)   # "positive"
test3 = check_number(-3)  # "negative"

{
    equal_test: equal_test,
    unequal_test: unequal_test,
    check_number: check_number,
    tests: [test1, test2, test3]
}
```

## Intermediate Examples

### List Operations

```fun
# lists.fun
# Working with lists

# List construction
empty_list = []
numbers = [1, 2, 3, 4, 5]
mixed = [1, `hello`, True]

# List operations using pattern matching
head = \list ->
    when list is
        Cons x xs -> x;
        Nil -> error "empty list"

tail = \list ->
    when list is
        Cons x xs -> xs;
        Nil -> error "empty list"

length = \list ->
    when list is
        Cons x xs -> 1 + length(xs);
        Nil -> 0

# Test list operations
first = head(numbers)     # 1
rest = tail(numbers)      # [2, 3, 4, 5]
count = length(numbers)   # 5

{
    head: head,
    tail: tail,
    length: length,
    first: first,
    rest: rest,
    count: count
}
```

### Records and Data Structures

```fun
# records.fun
# Working with records

# Define record types
Person : {name: Str, age: Int}
Point : {x: Int, y: Int}

# Create records
alice = {name: "Alice", age: 30}
bob = {name: "Bob", age: 25}
origin = {x: 0, y: 0}
point1 = {x: 10, y: 20}

# Record access
alice_name = alice.name   # "Alice"
alice_age = alice.age     # 30
point_x = point1.x        # 10
point_y = point1.y        # 20

# Functions that work with records
get_name = \person -> person.name
get_age = \person -> person.age
is_adult = \person -> person.age >= 18

# Test functions
alice_is_adult = is_adult(alice)  # True
bob_is_adult = is_adult(bob)      # True

{
    get_name: get_name,
    get_age: get_age,
    is_adult: is_adult,
    alice_is_adult: alice_is_adult,
    bob_is_adult: bob_is_adult
}
```

### Recursion with Fix

```fun
# recursion.fun
# Using the fix combinator for recursion

# Factorial function
factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(n - 1)
)

# Fibonacci function
fibonacci = fix(\rec -> \n ->
    when n == 0 is
        True -> 0;
        False -> when n == 1 is
            True -> 1;
            False -> rec(n - 1) + rec(n - 2)
)

# Sum range function
sum_range = fix(\rec -> \n ->
    when n == 0 is
        True -> 0;
        False -> n + rec(n - 1)
)

# Test recursive functions
fact_5 = factorial(5)     # 120
fib_10 = fibonacci(10)    # 55
sum_100 = sum_range(100)  # 5050

{
    factorial: factorial,
    fibonacci: fibonacci,
    sum_range: sum_range,
    fact_5: fact_5,
    fib_10: fib_10,
    sum_100: sum_100
}
```

## Advanced Examples

### Higher-Order Functions

```fun
# higher_order.fun
# Functions that take or return functions

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

# Compose function
compose = \f, g, x -> f(g(x))

# Test higher-order functions
numbers = [1, 2, 3, 4, 5]
inc = \x -> x + 1
double = \x -> x * 2
is_even = \x -> x % 2 == 0

mapped = map(inc, numbers)        # [2, 3, 4, 5, 6]
filtered = filter(is_even, numbers) # [2, 4]
summed = fold(add, 0, numbers)    # 15
composed = compose(inc, double, 5) # 11

{
    map: map,
    filter: filter,
    fold: fold,
    compose: compose,
    mapped: mapped,
    filtered: filtered,
    summed: summed,
    composed: composed
}
```

### Pattern Matching with Custom Types

```fun
# pattern_matching.fun
# Advanced pattern matching examples

# Maybe type simulation
safe_head = \list ->
    when list is
        Cons x xs -> Just x;
        Nil -> Nothing

safe_tail = \list ->
    when list is
        Cons x xs -> Just xs;
        Nil -> Nothing

# Result type simulation
safe_divide = \x, y ->
    when y == 0 is
        True -> Error "division by zero";
        False -> Success (x / y)

# Tree structure
process_tree = \tree ->
    when tree is
        Leaf value -> value;
        Node left right -> process_tree(left) + process_tree(right)

# Test pattern matching
test_list = [1, 2, 3]
empty_list = []

head_result = safe_head(test_list)    # Just 1
tail_result = safe_tail(test_list)    # Just [2, 3]
empty_head = safe_head(empty_list)    # Nothing

div_success = safe_divide(10, 2)      # Success 5
div_error = safe_divide(10, 0)        # Error "division by zero"

{
    safe_head: safe_head,
    safe_tail: safe_tail,
    safe_divide: safe_divide,
    head_result: head_result,
    tail_result: tail_result,
    empty_head: empty_head,
    div_success: div_success,
    div_error: div_error
}
```

### Module System Example

```fun
# math_utils.fun
# Mathematical utilities module

# Basic math functions
inc = \x -> x + 1
dec = \x -> x - 1
abs = \x -> when x >= 0 is
    True -> x;
    False -> -x

# Statistical functions
mean = \list ->
    let
        sum = fold(add, 0, list)
        count = length(list)
    in
        when count == 0 is
            True -> 0;
            False -> sum / count

# Export functions
{
    inc: inc,
    dec: dec,
    abs: abs,
    mean: mean
}
```

```fun
# main.fun
# Main program using the math module

import math from `./math_utils`

# Use imported functions
result1 = math.inc(5)           # 6
result2 = math.dec(10)          # 9
result3 = math.abs(-7)          # 7
result4 = math.mean([1, 2, 3, 4, 5])  # 3

# Combine functions
combined = math.inc(math.abs(-5))  # 6

{
    result1: result1,
    result2: result2,
    result3: result3,
    result4: result4,
    combined: combined
}
```

## Common Patterns

### Error Handling

```fun
# error_handling.fun
# Common error handling patterns

# Maybe pattern for optional values
safe_parse_int = \str ->
    when str is
        "0" -> Just 0;
        "1" -> Just 1;
        "2" -> Just 2;
        _ -> Nothing

# Result pattern for operations that can fail
safe_operation = \x, y ->
    when y == 0 is
        True -> Error "cannot divide by zero";
        False -> Success (x / y)

# Chain operations with error handling
process_data = \input ->
    let
        parsed = safe_parse_int(input)
    in
        when parsed is
            Just value -> Success value;
            Nothing -> Error "invalid input"

# Test error handling
test1 = safe_parse_int("1")     # Just 1
test2 = safe_parse_int("abc")   # Nothing
test3 = safe_operation(10, 2)   # Success 5
test4 = safe_operation(10, 0)   # Error "cannot divide by zero"

{
    safe_parse_int: safe_parse_int,
    safe_operation: safe_operation,
    process_data: process_data,
    test1: test1,
    test2: test2,
    test3: test3,
    test4: test4
}
```

### State Management

```fun
# state.fun
# Managing state in a functional way

# Counter with state
create_counter = \initial ->
    {
        value: initial,
        increment: \counter -> {counter | value: counter.value + 1},
        decrement: \counter -> {counter | value: counter.value - 1},
        get_value: \counter -> counter.value
    }

# Test counter
counter = create_counter(0)
counter1 = counter.increment(counter)
counter2 = counter.increment(counter1)
value = counter.get_value(counter2)

# Bank account simulation
create_account = \initial_balance ->
    {
        balance: initial_balance,
        deposit: \account, amount -> {account | balance: account.balance + amount},
        withdraw: \account, amount ->
            when account.balance >= amount is
                True -> {account | balance: account.balance - amount};
                False -> account,
        get_balance: \account -> account.balance
    }

# Test bank account
account = create_account(100)
account1 = account.deposit(account, 50)
account2 = account.withdraw(account1, 30)
final_balance = account.get_balance(account2)

{
    create_counter: create_counter,
    create_account: create_account,
    counter_value: value,
    account_balance: final_balance
}
```

## Complete Programs

### Simple Calculator Program

```fun
# calculator.fun
# Complete calculator program

# Define operations
add = \x, y -> x + y
subtract = \x, y -> x - y
multiply = \x, y -> x * y
divide = \x, y -> when y == 0 is
    True -> error "division by zero";
    False -> x / y

# Calculator function
calculate = \op, x, y ->
    when op is
        "add" -> add(x, y);
        "subtract" -> subtract(x, y);
        "multiply" -> multiply(x, y);
        "divide" -> divide(x, y);
        _ -> error "unknown operation"

# Test calculations
result1 = calculate("add", 5, 3)      # 8
result2 = calculate("subtract", 10, 4) # 6
result3 = calculate("multiply", 2, 7)  # 14
result4 = calculate("divide", 15, 3)   # 5

# Export calculator
{
    add: add,
    subtract: subtract,
    multiply: multiply,
    divide: divide,
    calculate: calculate,
    results: [result1, result2, result3, result4]
}
```

### List Processing Program

```fun
# list_processor.fun
# Complete list processing program

# List utilities
map = \f, list ->
    when list is
        Cons x xs -> Cons (f x) (map f xs);
        Nil -> Nil

filter = \pred, list ->
    when list is
        Cons x xs -> when pred(x) is
            True -> Cons x (filter pred xs);
            False -> filter pred xs;
        Nil -> Nil

fold = \f, init, list ->
    when list is
        Cons x xs -> fold f (f init x) xs;
        Nil -> init

# List operations
sum = fold(add, 0)
product = fold(multiply, 1)
length = fold(\acc, x -> acc + 1, 0)

# Predicates
is_even = \x -> x % 2 == 0
is_positive = \x -> x > 0
is_negative = \x -> x < 0

# Process a list
numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

doubled = map(double, numbers)
evens = filter(is_even, numbers)
positives = filter(is_positive, numbers)
total = sum(numbers)
count = length(numbers)

# Export functions and results
{
    map: map,
    filter: filter,
    fold: fold,
    sum: sum,
    product: product,
    length: length,
    doubled: doubled,
    evens: evens,
    positives: positives,
    total: total,
    count: count
}
```

### Data Processing Program

```fun
# data_processor.fun
# Complete data processing program

# Data structures
Person : {name: Str, age: Int, city: Str}

# Sample data
people = [
    {name: "Alice", age: 30, city: "New York"},
    {name: "Bob", age: 25, city: "Boston"},
    {name: "Charlie", age: 35, city: "New York"},
    {name: "Diana", age: 28, city: "Chicago"}
]

# Data processing functions
get_names = map(\person -> person.name, people)
get_ages = map(\person -> person.age, people)
get_cities = map(\person -> person.city, people)

# Filtering functions
adults = filter(\person -> person.age >= 18, people)
new_yorkers = filter(\person -> person.city == "New York", people)
young_people = filter(\person -> person.age < 30, people)

# Aggregation functions
average_age = mean(get_ages)
total_people = length(people)
adult_count = length(adults)

# Export results
{
    people: people,
    get_names: get_names,
    get_ages: get_ages,
    get_cities: get_cities,
    adults: adults,
    new_yorkers: new_yorkers,
    young_people: young_people,
    average_age: average_age,
    total_people: total_people,
    adult_count: adult_count
}
```

These examples demonstrate the power and expressiveness of the Fun programming language. Start with the basic examples and work your way up to the more advanced patterns as you become comfortable with the language. 