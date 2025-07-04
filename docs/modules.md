# Modules

Fun's module system allows you to organize code across multiple files and share functionality between programs.

## Table of Contents

- [Module Basics](#module-basics)
- [Creating Modules](#creating-modules)
- [Importing Modules](#importing-modules)
- [Exporting Values](#exporting-values)
- [Module Paths](#module-paths)
- [Module Structure](#module-structure)
- [Best Practices](#best-practices)

## Module Basics

Every `.fun` file is a module. A module can contain:

- **Function definitions**: `inc = \x -> x + 1`
- **Type annotations**: `factorial : Lam<Int, Int>`
- **Import statements**: `import lib from \`./lib\``
- **Export record**: The final expression is exported

## Creating Modules

### Simple Module

```fun
# math.fun
inc = \x -> x + 1
dec = \x -> x - 1
double = \x -> x * 2

# Export all functions as a record
{
    inc: inc,
    dec: dec,
    double: double
}
```

### Module with Type Annotations

```fun
# types.fun
# Define types for better documentation
inc : Lam<Int, Int>
inc = \x -> x + 1

dec : Lam<Int, Int>
dec = \x -> x - 1

# Export with type information
{
    inc: inc,
    dec: dec
}
```

### Module with Imports

```fun
# advanced.fun
import basic from `./basic`

# Use imported functions
factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(basic.dec(n))
)

# Export new functions
{
    factorial: factorial,
    # Re-export imported functions
    inc: basic.inc,
    dec: basic.dec
}
```

## Importing Modules

### Basic Import

```fun
# main.fun
import math from `./math`

# Use imported functions
result = math.inc(5)    # result: 6
sum = math.inc(10) + math.dec(5)  # result: 14
```

### Import with Alias

```fun
# main.fun
import utils from `./math`

# Use with different name
result = utils.inc(5)   # same as math.inc(5)
```

### Multiple Imports

```fun
# main.fun
import math from `./math`
import types from `./types`

# Use functions from different modules
result1 = math.inc(5)
result2 = types.dec(10)
```

### Nested Imports

```fun
# main.fun
import advanced from `./advanced`

# Use functions that were re-exported
result = advanced.factorial(5)  # uses advanced.factorial
inc_result = advanced.inc(3)    # uses math.inc via re-export
```

## Exporting Values

### Record Exports

The most common way to export values is using a record:

```fun
# utils.fun
inc = \x -> x + 1
dec = \x -> x - 1
square = \x -> x * x

# Export as record
{
    inc: inc,
    dec: dec,
    square: square
}
```

### Single Value Export

You can export a single value directly:

```fun
# constants.fun
PI = 3.14159
E = 2.71828

# Export single value
PI
```

### Function Export

```fun
# factorial.fun
factorial : Lam<Int, Int>
factorial = fix(\rec -> \n ->
    when n == 0 is
        True -> 1;
        False -> n * rec(n - 1)
)

# Export function directly
factorial
```

### Selective Export

```fun
# library.fun
# Internal helper functions (not exported)
internal_helper = \x -> x * 2

# Public API functions
public_inc = \x -> x + 1
public_dec = \x -> x - 1

# Export only public functions
{
    inc: public_inc,
    dec: public_dec
}
```

## Module Paths

### Relative Paths

```fun
# Import from same directory
import math from `./math`

# Import from parent directory
import shared from `../shared`

# Import from subdirectory
import utils from `./utils/helpers`

# Import from sibling directory
import common from `../common/lib`
```

### File Extensions

```fun
# .fun extension is optional
import math from `./math`        # looks for math.fun
import math from `./math.fun`    # explicit extension

# Any file extension is supported
import data from `./data.txt`    # non-standard extension
import config from `./config`    # no extension
```

### Directory Traversal

```fun
# Multiple levels up
import root from `../../root`

# Deep nesting
import deep from `./nested/very/deep/module`

# Absolute paths (relative to project root)
import global from `/src/global`
```

## Module Structure

### Complete Module Example

```fun
# user.fun

# 1. Imports
import math from `./math`
import types from `./types`

# 2. Type annotations
User : {name: Str, age: Int}
create_user : Lam<Str, Int, User>
get_name : Lam<User, Str>

# 3. Function definitions
create_user = \name, age -> {name: name, age: age}

get_name = \user -> user.name

is_adult = \user -> user.age >= 18

# 4. Internal helper functions
validate_age = \age -> age >= 0

# 5. Export record
{
    create_user: create_user,
    get_name: get_name,
    is_adult: is_adult,
    # Re-export useful functions from other modules
    inc: math.inc,
    dec: math.dec
}
```

### Module Organization

```fun
# Recommended module structure:

# 1. Module documentation
# user.fun - User management functions

# 2. Imports
import math from `./math`

# 3. Type definitions
User : {name: Str, age: Int}

# 4. Constants
MIN_AGE = 0
MAX_AGE = 150

# 5. Public functions
create_user = \name, age -> {name: name, age: age}

# 6. Private helper functions
validate_age = \age -> age >= MIN_AGE && age <= MAX_AGE

# 7. Export record
{
    create_user: create_user,
    MIN_AGE: MIN_AGE,
    MAX_AGE: MAX_AGE
}
```

## Best Practices

### Module Naming

```fun
# Use descriptive names
math.fun              # Mathematical functions
user_management.fun   # User-related functions
string_utils.fun      # String utilities

# Avoid generic names
utils.fun             # Too generic
lib.fun               # Too generic
```

### Export Strategy

```fun
# Export only what's needed
# Good: Selective export
{
    public_function: public_function,
    public_constant: public_constant
}

# Avoid: Export everything
{
    public_function: public_function,
    internal_helper: internal_helper,  # Should be private
    debug_function: debug_function     # Should be private
}
```

### Import Organization

```fun
# Group imports logically
# Standard library imports
import math from `./math`
import string from `./string`

# Third-party imports
import external from `./external/lib`

# Local imports
import user from `./user`
import auth from `./auth`
```

### Module Dependencies

```fun
# Avoid circular dependencies
# Good: Hierarchical structure
math.fun              # No dependencies
user.fun              # Depends on math.fun
app.fun               # Depends on user.fun

# Bad: Circular dependency
# user.fun imports from app.fun
# app.fun imports from user.fun
```

### Type Safety

```fun
# Use type annotations for public APIs
public_function : Lam<Int, Str>
public_function = \x -> `{x}`

# Internal functions can rely on type inference
internal_helper = \x -> x * 2
```

### Documentation

```fun
# Document your modules
# user.fun
# User management module
# Provides functions for creating and manipulating user records

import math from `./math`

# Create a new user with name and age
create_user : Lam<Str, Int, {name: Str, age: Int}>
create_user = \name, age -> {name: name, age: age}

# Check if user is an adult
is_adult : Lam<{name: Str, age: Int}, Bool>
is_adult = \user -> user.age >= 18

{
    create_user: create_user,
    is_adult: is_adult
}
```

## Common Patterns

### Utility Module

```fun
# utils.fun
# Common utility functions

# String utilities
to_upper = \str -> `{str}`  # Placeholder for actual implementation
to_lower = \str -> `{str}`  # Placeholder for actual implementation

# List utilities
head = \list ->
    when list is
        Cons x xs -> x;
        Nil -> error "empty list"

tail = \list ->
    when list is
        Cons x xs -> xs;
        Nil -> error "empty list"

{
    to_upper: to_upper,
    to_lower: to_lower,
    head: head,
    tail: tail
}
```

### Configuration Module

```fun
# config.fun
# Application configuration

# Database settings
DB_HOST = "localhost"
DB_PORT = 5432
DB_NAME = "myapp"

# Application settings
APP_NAME = "Fun App"
APP_VERSION = "1.0.0"
DEBUG_MODE = True

{
    DB_HOST: DB_HOST,
    DB_PORT: DB_PORT,
    DB_NAME: DB_NAME,
    APP_NAME: APP_NAME,
    APP_VERSION: APP_VERSION,
    DEBUG_MODE: DEBUG_MODE
}
```

### Type Module

```fun
# types.fun
# Common type definitions

# User types
User : {name: Str, age: Int}
Admin : {name: Str, age: Int, permissions: List<Str>}

# Result types
Result : {success: Bool, data: a, error: Str}

# Export type constructors (if supported)
{
    User: User,
    Admin: Admin,
    Result: Result
}
```

The module system provides a clean way to organize code and share functionality across your Fun programs. Use it to create reusable libraries and maintainable codebases. 