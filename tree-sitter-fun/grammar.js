/**
 * @file A purely fun-ctional language with amazing performance
 * @author Matan Tsuberi <tsuberim@gmail.com>
 * @license MIT
 */

/// <reference types="tree-sitter-cli/dsl" />
// @ts-check

module.exports = grammar({
  name: "fun",

  rules: {
    // TODO: add the actual grammar rules
    source_file: $ => "hello"
  }
});
