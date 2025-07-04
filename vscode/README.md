# Fun Language Extension for VS Code

This extension provides syntax highlighting and language support for the Fun programming language.

## Features

- Syntax highlighting for Fun language files (.fun)
- Document symbol provider for outline view
- Tree-sitter based parsing
- Auto-closing brackets and parentheses
- Comment support (# for line comments)
- Indentation rules

## Installation

1. Clone this repository
2. Navigate to the `vscode` directory
3. Run `npm install` to install dependencies
4. Run `npm run compile` to build the extension
5. Press F5 in VS Code to launch the extension in a new window

## Usage

- Open any `.fun` file to see syntax highlighting
- Use the outline view to see document symbols
- Comments start with `#`
- Strings use backticks `` ` ``
- Function definitions use `\` for lambda syntax
- Pattern matching uses `when` and `is` keywords

## Development

- `npm run compile` - Compile TypeScript
- `npm run watch` - Watch for changes and recompile
- `npm run lint` - Run ESLint

## Tree-sitter Integration

This extension uses the tree-sitter grammar from the parent directory to provide accurate parsing and syntax highlighting. 