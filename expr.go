package main

import (
	"fmt"
	"github.com/scylladb/go-set/strset"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	"strconv"
	"strings"
)

func dent(count int, str string) string {
	return strings.Repeat("\t", count) + str
}

type Expr interface {
	expr()
	Pretty(indent int) string
}

func (i *Int) expr()    {}
func (s *LitStr) expr() {}
func (s *Str) expr()    {}
func (v *Var) expr()    {}
func (a *App) expr()    {}
func (l *Lam) expr()    {}
func (r *Rec) expr()    {}
func (p *Prop) expr()   {}
func (c *Cons) expr()   {}
func (w *When) expr()   {}
func (l *List) expr()   {}
func (b *Block) expr()  {}

type Int struct {
	Value int
}

func (i *Int) Pretty(indent int) string {
	return dent(indent, strconv.Itoa(i.Value))
}

type LitStr struct {
	Value string
}

func (s *LitStr) Pretty(indent int) string {
	return dent(indent, s.Value)
}

type Str struct {
	Parts []Expr
}

func (s *Str) Pretty(indent int) string {
	var parts []string
	for _, part := range s.Parts {
		parts = append(parts, part.Pretty(0))
	}

	return dent(indent, fmt.Sprintf("`%s`", strings.Join(parts, "")))
}

type Var struct {
	Name     string
	IsSymbol bool
}

func (v *Var) Pretty(indent int) string {
	return dent(indent, v.Name)
}

type App struct {
	Fn   Expr
	Args []Expr
}

func (a *App) Pretty(indent int) string {
	if fn, ok := a.Fn.(*Var); ok && fn.IsSymbol && len(a.Args) == 2 {
		return dent(indent, fmt.Sprintf("%s %s %s", a.Args[0].Pretty(indent), a.Fn.Pretty(indent), a.Args[1].Pretty(indent)))
	}

	var args []string
	for _, arg := range a.Args {
		args = append(args, arg.Pretty(indent))
	}
	return dent(indent, fmt.Sprintf("%s(%s)", a.Fn.Pretty(indent), strings.Join(args, ", ")))
}

type Lam struct {
	Params []string
	Body   Expr
}

func (l *Lam) Pretty(indent int) string {
	return dent(indent, fmt.Sprintf("\\%s -> %s", strings.Join(l.Params, ", "), l.Body.Pretty(indent)))
}

type RecEntry struct {
	Prop  string
	Value Expr
}

func (r *RecEntry) Pretty(indent int) string {
	return dent(indent, fmt.Sprintf("%s: %s", r.Prop, r.Value.Pretty(0)))
}

type Rec struct {
	Entries []RecEntry
}

func (r *Rec) Pretty(indent int) string {
	var entries []string
	for _, entry := range r.Entries {
		entries = append(entries, entry.Pretty(indent))
	}
	return dent(indent, fmt.Sprintf("{%s}", strings.Join(entries, ", ")))
}

type Prop struct {
	Parent Expr
	Prop   string
}

func (p *Prop) Pretty(indent int) string {
	return dent(indent, fmt.Sprintf("%s.%s", p.Parent.Pretty(indent), p.Prop))
}

type Cons struct {
	Name    string
	Payload Expr // may be nil
}

func (c *Cons) Pretty(indent int) string {
	if c.Payload == nil {
		return dent(indent, c.Name)
	}
	return fmt.Sprintf("%s %s", c.Name, c.Payload.Pretty(indent))
}

type WhenClause struct {
	ConsName    string
	Payload     string
	Consequence Expr
}

func (w *WhenClause) Pretty(indent int) string {
	if w.Payload == "" {
		return dent(indent, fmt.Sprintf("%s -> %s", w.ConsName, w.Consequence.Pretty(indent)))
	}
	return dent(indent, fmt.Sprintf("%s %s -> %s", w.ConsName, w.Payload, w.Consequence.Pretty(indent)))
}

type When struct {
	Value   Expr
	Options []WhenClause
	Else    Expr // may be nil
}

func (w *When) Pretty(indent int) string {
	var options []string
	for _, option := range w.Options {
		options = append(options, option.Pretty(indent))
	}
	if w.Else == nil {
		return fmt.Sprintf("when %s is %s", w.Value.Pretty(indent), strings.Join(options, "; "))
	}

	return fmt.Sprintf("when %s is %s else %s", w.Value.Pretty(indent), strings.Join(options, "; "), w.Else.Pretty(indent))
}

type List struct {
	Items []Expr
}

func (l *List) Pretty(indent int) string {
	var items []string
	for _, item := range l.Items {
		items = append(items, item.Pretty(indent))
	}
	return dent(indent, fmt.Sprintf("[%s]", strings.Join(items, ", ")))
}

type Declaration interface {
	Pretty(indent int) string
	decl()
}

func (a *Assignment) decl()     {}
func (a *TypeAnnotation) decl() {}
func (a *Import) decl()         {}

type Assignment struct {
	Name  string
	Value Expr
}

func (a *Assignment) Pretty(indent int) string {
	return dent(indent, fmt.Sprintf("%s = %s", a.Name, a.Value.Pretty(indent)))
}

type TypeAnnotation struct {
	Name   string
	Scheme *Scheme
}

func (a *TypeAnnotation) Pretty(indent int) string {
	return dent(indent, fmt.Sprintf("%s : %s", a.Name, a.Scheme.Pretty(indent)))
}

type Import struct {
	Name string
	Path string
}

func (i *Import) Pretty(indent int) string {
	return fmt.Sprintf("import %s from %s", i.Name, i.Path)
}

type Block struct {
	Decs   []Declaration
	Result Expr
}

func (b *Block) Pretty(indent int) string {
	var declarations []string
	for _, declaration := range b.Decs {
		declarations = append(declarations, declaration.Pretty(indent))
	}
	return dent(indent, fmt.Sprintf("(%s;%s)", strings.Join(declarations, ";"), b.Result.Pretty(indent)))
}

func fromNode(node *tree_sitter.Node, source []byte) (Expr, error) {
	if node.HasError() {
		return nil, fmt.Errorf("parse error")
	}
	if node == nil {
		return nil, nil
	}

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
			expr, err := fromNode(&child, source)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, expr)
		}

		return &Str{Parts: exprs}, nil
	case "var":
		name := node.Utf8Text(source)
		return &Var{Name: name}, nil
	case "sym":
		name := node.Utf8Text(source)
		return &Var{Name: name, IsSymbol: true}, nil
	case "app":
		first, err := fromNode(node.NamedChild(0), source)
		if err != nil {
			return nil, err
		}

		var args []Expr
		for i := uint(1); i < node.NamedChildCount(); i++ {
			child := node.NamedChild(i)
			expr, err := fromNode(child, source)
			if err != nil {
				return nil, err
			}

			args = append(args, expr)
		}

		return &App{Fn: first, Args: args}, nil
	case "iapp":
		a, err := fromNode(node.NamedChild(0), source)
		if err != nil {
			return nil, err
		}
		op, err := fromNode(node.NamedChild(1), source)
		if err != nil {
			return nil, err
		}
		b, err := fromNode(node.NamedChild(2), source)
		if err != nil {
			return nil, err
		}

		return &App{Fn: op, Args: []Expr{a, b}}, nil
	case "lam":
		var exprs []Expr
		for _, child := range node.NamedChildren(node.Walk()) {
			expr, err := fromNode(&child, source)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, expr)
		}

		last := exprs[len(exprs)-1]
		var params []string
		for _, param := range exprs[:len(exprs)-1] {
			v, ok := param.(*Var)
			if !ok {
				return nil, fmt.Errorf("unexpected param expression type %t", param)
			}
			params = append(params, v.Name)
		}

		return &Lam{
			Params: params,
			Body:   last,
		}, nil
	case "rec":
		var entires []RecEntry

		names := strset.New()
		prop := true
		entry := RecEntry{}
		for _, child := range node.NamedChildren(node.Walk()) {
			expr, err := fromNode(&child, source)
			if err != nil {
				return nil, err
			}

			if prop {
				p, ok := expr.(*Var)
				if !ok {
					return nil, fmt.Errorf("unexpected record lhs expression type %t", expr)
				}

				entry.Prop = p.Name
				if names.Has(p.Name) {
					return nil, fmt.Errorf("duplicate record prop name: %s", p.Name)
				}
				names.Add(p.Name)
			} else {
				entry.Value = expr
				entires = append(entires, entry)
				entry = RecEntry{}
			}

			prop = !prop
		}

		return &Rec{Entries: entires}, nil
	case "prop":
		parent, err := fromNode(node.NamedChild(0), source)
		if err != nil {
			return nil, err
		}

		prop, err := fromNode(node.NamedChild(1), source)
		if err != nil {
			return nil, err
		}

		v, ok := prop.(*Var)
		if !ok {
			return nil, fmt.Errorf("unexpected property expression type %t", prop)
		}

		return &Prop{
			Parent: parent,
			Prop:   v.Name,
		}, nil
	case "lhs":
		expr, err := fromNode(node.NamedChild(0), source)
		if err != nil {
			return nil, err
		}

		prop, err := fromNode(node.NamedChild(1), source)
		if err != nil {
			return nil, err
		}

		v, ok := prop.(*Var)
		if !ok {
			return nil, fmt.Errorf("unexpected lhs expression type %t", prop)
		}

		return &Prop{
			Parent: expr,
			Prop:   v.Name,
		}, nil
	case "cons":
		consName := node.NamedChild(0).Utf8Text(source)
		payload, err := fromNode(node.NamedChild(1), source)
		if err != nil {
			return nil, err
		}
		return &Cons{Name: consName, Payload: payload}, nil
	case "when":
		value, err := fromNode(node.NamedChild(0), source)
		if err != nil {
			return nil, err
		}

		count := node.NamedChildCount()

		cases := strset.New()
		var options []WhenClause
		i := uint(1)
		for i+3 <= count {
			cons := node.NamedChild(i).Utf8Text(source)
			if cases.Has(cons) {
				return nil, fmt.Errorf("duplciate when clause cons name: %s", cons)
			}
			cases.Add(cons)

			payloadVar, err := fromNode(node.NamedChild(i+1), source)
			if err != nil {
				return nil, err
			}
			v, ok := payloadVar.(*Var)
			if !ok {
				return nil, fmt.Errorf("unexpected when payload expression type %t", payloadVar)
			}

			consequence, err := fromNode(node.NamedChild(i+2), source)
			if err != nil {
				return nil, err
			}

			options = append(options, WhenClause{
				ConsName:    cons,
				Payload:     v.Name,
				Consequence: consequence,
			})
			i += 3
		}

		var elseConsequence Expr
		if i < node.NamedChildCount() {
			expr, err := fromNode(node.NamedChild(i), source)
			if err != nil {
				return nil, err
			}

			elseConsequence = expr
		}

		return &When{
			Value:   value,
			Options: options,
			Else:    elseConsequence,
		}, nil

	case "list":
		var exprs []Expr
		for _, child := range node.NamedChildren(node.Walk()) {
			expr, err := fromNode(&child, source)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, expr)
		}
		return &List{Items: exprs}, nil
	case "annot":
	case "block", "source_file":
		var declarations []Declaration
		children := node.NamedChildren(node.Walk())
		for _, child := range children[:len(children)-1] {
			switch child.GrammarName() {
			case "assign":
				assign, err := assignFromNode(&child, source)
				if err != nil {
					return nil, err
				}
				declarations = append(declarations, assign)
			case "annot":
				annot, err := annotFromNode(&child, source)
				if err != nil {
					return nil, err
				}
				declarations = append(declarations, annot)
			case "import":
				importDec, err := importFromNode(&child, source)
				if err != nil {
					return nil, err
				}
				declarations = append(declarations, importDec)
			default:
				return nil, fmt.Errorf("unexpected declaration name %s", child.GrammarName())
			}
		}

		last, err := fromNode(&children[len(children)-1], source)
		if err != nil {
			return nil, err
		}

		return &Block{Decs: declarations, Result: last}, nil
	}
	return nil, fmt.Errorf("invalid node type %s", node.GrammarName())
}

func assignFromNode(node *tree_sitter.Node, source []byte) (*Assignment, error) {
	lhs, err := fromNode(node.NamedChild(0), source)
	if err != nil {
		return nil, err
	}

	v, ok := lhs.(*Var)
	if !ok {
		return nil, fmt.Errorf("unexpected lhs expression type %t", lhs)
	}

	rhs, err := fromNode(node.NamedChild(1), source)
	if err != nil {
		return nil, err
	}

	return &Assignment{
		Name:  v.Name,
		Value: rhs,
	}, nil
}

func annotFromNode(node *tree_sitter.Node, source []byte) (*TypeAnnotation, error) {
	lhs, err := fromNode(node.NamedChild(0), source)
	if err != nil {
		return nil, err
	}

	v, ok := lhs.(*Var)
	if !ok {
		return nil, fmt.Errorf("unexpected lhs expression type %t", lhs)
	}

	typ, err := typeFromNode(node.NamedChild(1), source)
	if err != nil {
		return nil, err
	}

	return &TypeAnnotation{
		Name:   v.Name,
		Scheme: generalize(typ),
	}, nil
}

func importFromNode(node *tree_sitter.Node, source []byte) (*Import, error) {
	lhs := node.NamedChild(0).Utf8Text(source)
	importName := node.NamedChild(1).Utf8Text(source)

	return &Import{
		Name: lhs,
		Path: strings.Trim(importName, "`"),
	}, nil
}
