/**
 * @file A purely fun-ctional language with amazing performance
 * @author Matan Tsuberi <tsuberim@gmail.com>
 * @license MIT
 */

/// <reference types="tree-sitter-cli/dsl" />
// @ts-check

function sep1(e, del) {
  return seq(repeat(seq(e, del)), e, optional(del))
}

function sep(e, del) {
  return optional(sep1(e, del))
}

const varName = /[a-z][\w\d_]*/
const consName = /[A-Z][\w\d_]*/
const symbol = /[!@#$%^&*+\-~]+/

module.exports = grammar({
  name: "fun",

  rules: {
    // TODO: add the actual grammar rules
    source_file: $ => $._inner_block,
    _expr: $ => choice($.int, $.str, $.var, $.sym, $.app, $.iapp, $.lam, $.rec, $.prop, $.cons, $.when, $.list, $.block),
    int: $ => /\d+/,
    str: $ => seq('`',repeat(choice(/[^`{}]+/, seq('{', $._expr, '}'))),'`'),
    var: $ => varName,
    sym: $ => symbol,
    app: $ => prec(2, seq($._expr, '(', sep($._expr, ','), ')')),
    iapp: $ => prec.left(3,seq($._expr, $.sym, $._expr)),
    lam: $ => seq('\\', sep($.var, ','), '->', $._expr),
    rec: $ => seq('{', sep(seq($.var, ':', $._expr), ','), '}'),
    prop: $ => prec.left(3, seq($._expr, '.',$.var)),
    cons: $ => prec.left(4,seq(consName, optional($._expr))),
    when: $ => prec.right(6,seq('when', $._expr, 'is', sep1(seq(consName, optional($.var), '->', $._expr), ';'),optional(seq(';', 'else', $._expr)))),
    list: $ => seq('[', sep($._expr, ','), ']'),
    _inner_block: $ => seq(repeat(seq($.var, '=', $._expr, '\n')), $._expr),
    block: $ => prec.left(5,seq('(', $._inner_block, ')')),
  },

  words: ['when']
});
