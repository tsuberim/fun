package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"maps"
	"strings"
)

type Val interface {
	val()
	Pretty(indent int) string
}

func (i *Int) val()     {}
func (s *LitStr) val()  {}
func (l *ListVal) val() {}
func (r *RecVal) val()  {}
func (c *ConsVal) val() {}
func (c *Closure) val() {}
func (c *Builtin) val() {}

type ListVal struct {
	Items []Val
}

func (l *ListVal) Pretty(indent int) string {
	var items []string
	for _, item := range l.Items {
		items = append(items, item.Pretty(indent))
	}
	return dent(indent, fmt.Sprintf("[%s]", strings.Join(items, ", ")))
}

type RecVal struct {
	Entries map[string]Val
}

func (r *RecVal) Pretty(indent int) string {
	var entries []string
	keys := lo.Keys(r.Entries)
	for _, key := range keys {
		val := r.Entries[key]
		entries = append(entries, fmt.Sprintf("%s: %s", key, val.Pretty(indent)))
	}
	return dent(indent, fmt.Sprintf("{%s}", strings.Join(entries, ",\n")))
}

type ConsVal struct {
	Name    string
	Payload Val
}

func (c *ConsVal) Pretty(indent int) string {
	if c.Payload == nil {
		return dent(indent, c.Name)
	}
	return fmt.Sprintf("%s %s", c.Name, c.Payload.Pretty(indent))
}

type Closure struct {
	Params []string
	Env    map[string]Val
	Body   Expr
}

func (c *Closure) Pretty(indent int) string {
	return "<closure>"
}

type Builtin struct {
	Name string
	Impl func(args []Val) (Val, error)
}

func (b *Builtin) Pretty(indent int) string {
	return fmt.Sprintf("<builtin %s>", b.Name)
}

type Evaluator struct {
	program *Program
}

func NewEvaluator(program *Program) *Evaluator {
	return &Evaluator{program: program}
}

func (e *Evaluator) Eval(expr Expr, env map[string]Val) (Val, error) {
	switch expr := expr.(type) {
	case *Int:
		return expr, nil
	case *LitStr:
		return expr, nil
	case *Str:
		sum := ""
		for _, part := range expr.Parts {
			val, err := e.Eval(part, env)
			if err != nil {
				return nil, err
			}

			str, ok := val.(*LitStr)
			if !ok {
				return nil, errors.Errorf("invalid value type for string template %t", val)
			}

			sum += str.Value
		}
		return &LitStr{sum}, nil
	case *Var:
		val, has := env[expr.Name]
		if !has {
			return nil, errors.Errorf("undefined variable: %s", expr.Name)
		}

		return val, nil
	case *Lam:
		return &Closure{
			Params: expr.Params,
			Env:    env,
			Body:   expr.Body,
		}, nil
	case *App:
		fn, err := e.Eval(expr.Fn, env)
		if err != nil {
			return nil, err
		}

		if builtin, ok := fn.(*Builtin); ok {
			var args []Val
			for _, arg := range expr.Args {
				val, err := e.Eval(arg, env)
				if err != nil {
					return nil, err
				}

				args = append(args, val)
			}

			return builtin.Impl(args)
		}

		clos, ok := fn.(*Closure)
		if !ok {
			return nil, errors.Errorf("cannot apply non closure value of type %t", fn)
		}

		if len(expr.Args) != len(clos.Params) {
			return nil, errors.Errorf("invalid number of arguments for function %t", fn)
		}

		newEnv := maps.Clone(clos.Env)
		for i, arg := range expr.Args {
			val, err := e.Eval(arg, env)
			if err != nil {
				return nil, err
			}
			newEnv[clos.Params[i]] = val
		}

		return e.Eval(clos.Body, newEnv)
	case *List:
		var vals []Val
		for _, item := range expr.Items {
			val, err := e.Eval(item, env)
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
		}

		return &ListVal{Items: vals}, nil
	case *Rec:
		entries := map[string]Val{}
		for _, entry := range expr.Entries {
			val, err := e.Eval(entry.Value, env)
			if err != nil {
				return nil, err
			}

			entries[entry.Prop] = val
		}

		return &RecVal{Entries: entries}, nil
	case *Prop:
		val, err := e.Eval(expr.Parent, env)
		if err != nil {
			return nil, err
		}

		rec, ok := val.(*RecVal)
		if !ok {
			return nil, errors.Errorf("invalid value type for prop parent %t", val)
		}

		val, has := rec.Entries[expr.Prop]
		if !has {
			return nil, errors.Errorf("record does not contain prop %s", expr.Prop)
		}

		return val, nil
	case *Cons:
		val, err := e.Eval(expr.Payload, env)
		if err != nil {
			return nil, err
		}

		return &ConsVal{Name: expr.Name, Payload: val}, nil
	case *When:
		val, err := e.Eval(expr.Value, env)
		if err != nil {
			return nil, err
		}

		cons, ok := val.(*ConsVal)
		if !ok {
			return nil, errors.Errorf("invalid value type for when %t", val)
		}

		for _, clause := range expr.Options {
			if clause.ConsName != cons.Name {
				continue
			}

			return e.Eval(clause.Consequence, extend(env, clause.Payload, cons.Payload))
		}

		if expr.Else == nil {
			return nil, errors.Errorf("no when clause matches cons name %s", cons.Name)
		}

		return e.Eval(expr.Else, env)
	case *Block:
		blockEnv := maps.Clone(env)
		for _, decl := range expr.Decs {
			switch decl := decl.(type) {
			case *Assignment:
				val, err := e.Eval(decl.Value, blockEnv)
				if err != nil {
					return nil, err
				}

				blockEnv[decl.Name] = val
			case *Import:
				mod, err := e.program.Import(decl.Path)
				if err != nil {
					return nil, err
				}

				blockEnv[decl.Name] = mod.Val
			}
		}

		return e.Eval(expr.Result, blockEnv)
	}

	return nil, errors.Errorf("invalid expression type: %T", expr)
}

func extend(env map[string]Val, name string, val Val) map[string]Val {
	cloned := maps.Clone(env)
	cloned[name] = val
	return cloned
}
