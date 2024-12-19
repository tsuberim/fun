package main

import (
	"fmt"
	"github.com/samber/lo"
	"maps"
	"sort"
	"strings"
)
import "github.com/scylladb/go-set/strset"

// lowercase to avoid clash with user-defined cons names
const intConsName = "int"
const strConsName = "str"
const lambdaConsName = "lam"
const listConsName = "list"

type Type interface {
	typ()
	freeVars() *strset.Set
	apply(subst *Subst) Type
	Pretty(i int) string
}

func (t *TypeCons) typ() {}
func (t *TypeVar) typ()  {}

type TypeCons struct {
	Name string
	Args []Type
}

func (t *TypeCons) Pretty(indent int) string {
	if len(t.Args) == 0 {
		return t.Name
	}

	var args []string
	for _, arg := range t.Args {
		args = append(args, arg.Pretty(indent))
	}

	return fmt.Sprintf("%s<%s>", t.Name, strings.Join(args, ", "))
}

func (t *TypeCons) freeVars() *strset.Set {
	result := strset.New()
	for _, arg := range t.Args {
		result.Merge(arg.freeVars())
	}
	return result
}

func (t *TypeCons) apply(subst *Subst) Type {
	return &TypeCons{
		Name: t.Name,
		Args: lo.Map(t.Args, func(item Type, index int) Type {
			return item.apply(subst)
		}),
	}
}

type TypeVar struct {
	Name string
}

func (t *TypeVar) Pretty(indent int) string {
	return t.Name
}

func (t *TypeVar) freeVars() *strset.Set {
	return strset.New(t.Name)
}

func (t *TypeVar) apply(subst *Subst) Type {
	if t, has := subst.Subst[t.Name]; has {
		return t
	}

	return t
}

type Scheme struct {
	Forall []string
	Type   Type
}

func (s *Scheme) Pretty(indent int) string {
	if len(s.Forall) == 0 {
		return s.Type.Pretty(indent)
	}
	return fmt.Sprintf("∀%s. %s", strings.Join(s.Forall, ", "), s.Type.Pretty(indent))
}

func (s *Scheme) apply(subst *Subst) *Scheme {
	limitedSubst := maps.Clone(subst.Subst)
	for _, param := range s.Forall {
		delete(limitedSubst, param)
	}

	return &Scheme{
		Forall: s.Forall,
		Type:   s.Type.apply(&Subst{Subst: limitedSubst}),
	}
}

func (s *Scheme) freeVars() *strset.Set {
	result := s.Type.freeVars()
	result.Remove(s.Forall...)
	return result
}

type Subst struct {
	Subst map[string]Type
}

func (s *Subst) compose(subst *Subst) *Subst {
	result := &Subst{Subst: map[string]Type{}}
	for name, t := range subst.Subst {
		result.Subst[name] = t
	}
	for name, t := range s.Subst {
		result.Subst[name] = t.apply(subst)
	}
	return result
}

func bind(tvarName string, t Type) (*Subst, error) {
	if tvar, ok := t.(*TypeVar); ok && tvarName == tvar.Name {
		return &Subst{Subst: map[string]Type{}}, nil
	}

	if t.freeVars().Has(tvarName) {
		return nil, fmt.Errorf("infinite recursive type")
	}

	return &Subst{Subst: map[string]Type{tvarName: t}}, nil

}

func unify(t1, t2 Type) (*Subst, error) {
	if tvar, ok := t1.(*TypeVar); ok {
		return bind(tvar.Name, t2)
	}

	if tvar, ok := t2.(*TypeVar); ok {
		return bind(tvar.Name, t1)
	}

	cons1 := t1.(*TypeCons)
	cons2 := t2.(*TypeCons)

	if cons1.Name != cons2.Name || len(cons1.Args) != len(cons2.Args) {
		return nil, fmt.Errorf("incompatible types %s ~!~ %s", cons1.Pretty(0), cons2.Pretty(0))
	}

	subst := &Subst{Subst: map[string]Type{}}
	for i := range len(cons1.Args) {
		s, err := unify(cons1.Args[i].apply(subst), cons2.Args[i].apply(subst))
		if err != nil {
			return nil, err
		}

		subst = subst.compose(s)
	}

	return subst, nil
}

func generalize(t Type) *Scheme {
	forall := t.freeVars().List()
	sort.Strings(forall)
	return &Scheme{
		Forall: forall,
		Type:   t,
	}
}

type Inferrer struct {
	varCount int
}

func NewInferrer() *Inferrer {
	return &Inferrer{varCount: 0}
}

func (i *Inferrer) freshVar() *TypeVar {
	current := i.varCount
	i.varCount++
	return &TypeVar{Name: fmt.Sprintf("t%d", current)}
}

func (i *Inferrer) instantiate(scheme *Scheme) Type {
	subst := &Subst{Subst: map[string]Type{}}
	for _, param := range scheme.Forall {
		subst.Subst[param] = i.freshVar()
	}
	return scheme.Type.apply(subst)
}

type TypeEnv struct {
	Types map[string]*Scheme
}

func (e *TypeEnv) apply(subst *Subst) *TypeEnv {
	result := &TypeEnv{Types: map[string]*Scheme{}}
	for name, scheme := range e.Types {
		result.Types[name] = scheme.apply(subst)
	}
	return result
}

func (e *TypeEnv) extend(name string, scheme *Scheme) *TypeEnv {
	cloned := maps.Clone(e.Types)
	cloned[name] = scheme
	return &TypeEnv{Types: cloned}
}

func (i *Inferrer) Infer(expr Expr, env *TypeEnv) (subst *Subst, typ Type, err error) {
	subst = &Subst{Subst: map[string]Type{}}

	switch expr := expr.(type) {
	case *Int:
		return subst, &TypeCons{
			Name: intConsName,
			Args: nil,
		}, nil
	case *LitStr:
		return subst, &TypeCons{
			Name: strConsName,
			Args: nil,
		}, nil
	case *Str:
		for _, part := range expr.Parts {
			env = env.apply(subst)

			s, t, err := i.Infer(part, env)
			if err != nil {
				return nil, nil, err
			}

			subst = subst.compose(s)

			s, err = unify(t, &TypeCons{
				Name: "str",
				Args: nil,
			})
			if err != nil {
				return nil, nil, err
			}

			subst = subst.compose(s)
		}

		return subst, &TypeCons{
			Name: "str",
			Args: nil,
		}, nil
	case *Var:
		scheme, has := env.Types[expr.Name]
		if !has {
			return nil, nil, fmt.Errorf("unbound variable %s", expr.Name)
		}

		t := i.instantiate(scheme)
		return &Subst{Subst: map[string]Type{}}, t, nil
	case *Lam:
		newEnv := &TypeEnv{Types: maps.Clone(env.Types)}

		var args []Type
		for _, param := range expr.Params {
			fresh := i.freshVar()
			newEnv.Types[param] = &Scheme{
				Forall: nil,
				Type:   fresh,
			}
			args = append(args, fresh)
		}

		s, t, err := i.Infer(expr.Body, newEnv)
		if err != nil {
			return nil, nil, err
		}

		subst = subst.compose(s)
		args = append(args, t)

		result := &TypeCons{
			Name: lambdaConsName,
			Args: args,
		}

		return subst, result.apply(subst), nil
	case *App:
		resultVar := i.freshVar()

		var args []Type
		for _, arg := range expr.Args {
			env = env.apply(subst)

			s, t, err := i.Infer(arg, env)
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)
			args = append(args, t)
		}
		env = env.apply(subst)

		s, t, err := i.Infer(expr.Fn, env)
		if err != nil {
			return nil, nil, err
		}

		subst = subst.compose(s)

		args = append(args, resultVar)
		s, err = unify(t.apply(subst), &TypeCons{
			Name: lambdaConsName,
			Args: args,
		})
		if err != nil {
			return nil, nil, err
		}
		subst = subst.compose(s)

		return subst, resultVar.apply(subst), nil
	case *List:
		var fresh Type = i.freshVar()
		for _, item := range expr.Items {
			fresh = fresh.apply(subst)
			env = env.apply(subst)

			s, t, err := i.Infer(item, env)
			if err != nil {
				return nil, nil, err
			}

			subst = subst.compose(s)

			s, err = unify(t, fresh)
			if err != nil {
				return nil, nil, err
			}

			subst = subst.compose(s)
		}

		itemType := fresh.apply(subst)

		return subst, &TypeCons{
			Name: listConsName,
			Args: []Type{itemType},
		}, nil
	case *Rec:
	case *Prop:
	case *Cons:
		s, t, err := i.Infer(expr.Payload, env)
		if err != nil {
			return nil, nil, err
		}
		subst = subst.compose(s)
		return subst, &TypeCons{Name: expr.Name, Args: []Type{t.apply(subst)}}, nil
	case *When:
	case *Block:
		for _, assignment := range expr.Assignments {
			env = env.apply(subst)
			s, t, err := i.Infer(assignment.Value, env)
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)
			env.extend(assignment.Name, generalize(t))
		}
		env = env.apply(subst)

		s, t, err := i.Infer(expr.Result, env)
		if err != nil {
			return nil, nil, err
		}
		subst = subst.compose(s)
		return subst, t.apply(subst), nil
	}

	return nil, nil, fmt.Errorf("invalid expression type: %T", expr)
}
