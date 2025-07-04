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

const varName = /_*[a-z][\w\d_]*/
const consName = /_*[A-Z][\w\d_]*/
const symbol = /[!@$%^&*+\-~><=]+/

module.exports = grammar({
  name: "fun",

  rules: {
    // TODO: add the actual grammar rules
    source_file: $ => $._inner_block,
    _expr: $ => choice($.int, $.str, $.var, $.sym, $.app, $.iapp, $.lam, $.rec, $.prop, $.cons, $.when, $.list, $.block),
    int: $ => /\d+/,
    lit_str: $ => /[^`{}]+/,
    str: $ => seq('`',repeat(choice($.lit_str, seq('{', $._expr, '}'))),'`'),
    var: $ => varName,
    cons_name: $ => consName,
    sym: $ => symbol,
    app: $ => prec(3, seq($._expr, '(', sep($._expr, ','), ')')),
    iapp: $ => prec.left(2,seq($._expr, $.sym, $._expr)),
    lam: $ => seq('\\', sep($.var, ','), '->', $._expr),
    rec: $ => seq('{', sep(seq($.var, ':', $._expr), ','), '}'),
    prop: $ => prec.left(3, seq($._expr, '.',$.var)),
    cons: $ => prec.left(4,seq($.cons_name, optional($._expr))),
    when: $ => prec.right(1,seq('when', $._expr, 'is', sep1(seq( $.cons_name, $.var, '->', $._expr), ';'), optional(seq('else', $._expr)))),
    list: $ => seq('[', sep($._expr, ','), ']'),
    assign: $ => seq($.var, '=', $._expr),
    bind: $ => seq($.var, '<-', $._expr),
    annot: $ => seq($.var, ':', $._type),
    import: $ => seq('import', $.var, 'from', '`', $.lit_str, '`') ,
    _decl: $ => choice($.assign, $.bind, $.annot, $.import),
    _inner_block: $ => seq(repeat(seq($._decl, choice('\n', '\\'))), $._expr),
    block: $ => prec.left(5,seq('(', $._inner_block, ')')),

    _type: $ => choice($.var, $.type_cons, $.type_rec, $.type_union),
    type_cons: $ => seq($.cons_name, optional(seq('<', sep($._type, ', '), '>'))),
    type_rec: $ => seq('{', sep(seq($.var, ':', $._type), ','), '}'),
    type_union: $ => seq('[', sep(seq($.cons_name, $._type), ','), ']'),

    _comment: _ => token(seq('#', /.*/)),
  },
  
  extras: $ => [
    $._comment,
    /[\s\f\uFEFF\u2060\u200B]|\r?\n/,
  ],

  word: $ => $.var,
});
