package main

import (
	"fmt"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	"strconv"
)

type Expr interface {
	expr()
}

func (e *Int) expr()    {}
func (e *LitStr) expr() {}
func (e *Str) expr()    {}
func (e *Var) expr()    {}
func (e *App) expr()    {}
func (e *Lam) expr()    {}
func (e *Rec) expr()    {}
func (e *Prop) expr()   {}
func (e *Cons) expr()   {}
func (e *When) expr()   {}
func (e *List) expr()   {}
func (e *Block) expr()  {}

type Int struct {
	Value int
}

type LitStr struct {
	Value string
}

type Str struct {
	Parts []Expr
}

type Var struct {
	Name string
}

type App struct {
	Fn   Expr
	Args []Expr
}

type Lam struct {
	Params []string
	Body   Expr
}

type RecEntry struct {
	Prop  string
	Value Expr
}

type Rec struct {
	Entries []RecEntry
}

type Prop struct {
	Parent Expr
	Prop   string
}

type Cons struct {
	Name    string
	Payload Expr // may be nil
}

type WhenClause struct {
	ConsName    string
	Payload     string
	Consequence Expr
}

type When struct {
	Value   Expr
	Options []WhenClause
	Else    Expr
}

type List struct {
	Items []Expr
}

type Assign struct {
	Name  string
	Value Expr
}

type Block struct {
	Assignments []Assign
	Result      Expr
}

func fromTree(node *tree_sitter.Node, source []byte) (Expr, error) {
	switch node.GrammarName() {
	case "int":
		value, err := strconv.Atoi(node.Utf8Text(source))
		if err != nil {
			return nil, err
		}
		return &Int{Value: value}, nil
	case "lit_str":
		return &LitStr{Value: node.Utf8Text(source)}, nil
	case "str":
		cursor := node.Walk()

		var exprs []Expr
		for _, child := range node.NamedChildren(cursor) {
			expr, err := fromTree(&child, source)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, expr)
		}

		return &Str{Parts: exprs}, nil
	case "var":
		name := node.Utf8Text(source)
		return &Var{Name: name}, nil
	}
	return nil, fmt.Errorf("invalid node type %s", node.GrammarName())
}
