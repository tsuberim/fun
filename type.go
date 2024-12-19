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
func (t *TypeRec) typ()  {}

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

type TypeRec struct {
	Entries map[string]Type
	RestVar *TypeVar
	Union   bool
}

func (t *TypeRec) Pretty(indent int) string {
	keys := lo.Keys(t.Entries)
	sort.Strings(keys)
	var entries []string
	for _, key := range keys {
		if t.Union {
			entries = append(entries, fmt.Sprintf("%s %s", key, t.Entries[key].Pretty(indent)))
		} else {
			entries = append(entries, fmt.Sprintf("%s: %s", key, t.Entries[key].Pretty(indent)))
		}
	}

	if t.Union {
		if t.RestVar != nil {
			return fmt.Sprintf("[%s |%s]", strings.Join(entries, " | "), t.RestVar.Pretty(indent))
		}
		return fmt.Sprintf("[%s]", strings.Join(entries, " | "))
	} else {
		if t.RestVar != nil {
			return fmt.Sprintf("{%s |%s}", strings.Join(entries, ", "), t.RestVar.Pretty(indent))
		}
		return fmt.Sprintf("{%s}", strings.Join(entries, ", "))
	}
}

func (t *TypeRec) freeVars() *strset.Set {
	result := strset.New()
	for _, t := range t.Entries {
		result.Merge(t.freeVars())
	}
	if t.RestVar != nil {
		result.Add(t.RestVar.Name)
	}
	return result
}

func (t *TypeRec) apply(subst *Subst) Type {
	newRec := &TypeRec{Entries: map[string]Type{}, RestVar: t.RestVar, Union: t.Union}
	for name, t := range t.Entries {
		newRec.Entries[name] = t.apply(subst)
	}

	if t.RestVar != nil {
		rest := t.RestVar.apply(subst)
		if tvar, ok := rest.(*TypeVar); ok {
			newRec.RestVar = tvar
		} else if trec, ok := rest.(*TypeRec); ok {
			for name, t := range trec.Entries {
				newRec.Entries[name] = t
			}
			newRec.RestVar = trec.RestVar
		} else {
			panic("impossible type record rest type")
		}
	}

	return newRec
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

func (i *Inferrer) unify(t1, t2 Type) (*Subst, error) {
	if tvar, ok := t1.(*TypeVar); ok {
		return bind(tvar.Name, t2)
	}

	if tvar, ok := t2.(*TypeVar); ok {
		return bind(tvar.Name, t1)
	}

	if cons1, ok := t1.(*TypeCons); ok {
		cons2, ok := t2.(*TypeCons)
		if !ok {
			return nil, fmt.Errorf("incompatible types %s ~!~ %s", t1.Pretty(0), t2.Pretty(0))
		}

		return i.unifyCons(cons1, cons2)
	}

	if rec1, ok := t1.(*TypeRec); ok {
		rec2, ok := t2.(*TypeRec)
		if !ok {
			return nil, fmt.Errorf("incompatible types %s ~!~ %s", t1.Pretty(0), t2.Pretty(0))
		}

		return i.unifyRecs(rec1, rec2)
	}

	return nil, fmt.Errorf("incompatible types %s ~!~ %s", t1.Pretty(0), t2.Pretty(0))
}

func (i *Inferrer) unifyRecs(rec1 *TypeRec, rec2 *TypeRec) (*Subst, error) {
	if rec1.Union != rec2.Union {
		return nil, fmt.Errorf("incompatible types %s ~!~ %s", rec1.Pretty(0), rec2.Pretty(0))
	}
	union := rec1.Union

	keys1 := strset.New(lo.Keys(rec1.Entries)...)
	keys2 := strset.New(lo.Keys(rec2.Entries)...)
	intersection := strset.Intersection(keys1, keys2)

	subst := &Subst{Subst: map[string]Type{}}
	for _, key := range intersection.List() {
		s, err := i.unify(rec1.Entries[key], rec2.Entries[key])
		if err != nil {
			return nil, err
		}
		subst = subst.compose(s)
	}

	keys1MinusKeys2 := strset.Difference(keys1, keys2)
	keys2MinusKeys1 := strset.Difference(keys2, keys1)
	open := rec1.RestVar != nil && rec2.RestVar != nil
	assignableToT1 := keys2MinusKeys1.IsEmpty() || rec1.RestVar != nil
	assignableToT2 := keys1MinusKeys2.IsEmpty() || rec2.RestVar != nil
	fresh := i.freshVar()
	if open || (assignableToT2 && assignableToT1) {
		if rec1.RestVar != nil {
			entries2 := map[string]Type{}
			for _, key := range keys2MinusKeys1.List() {
				entries2[key] = rec2.Entries[key]
			}

			s, err := i.unify(rec1.RestVar, &TypeRec{
				Entries: entries2,
				RestVar: fresh,
				Union:   union,
			})
			if err != nil {
				return nil, err
			}
			subst = subst.compose(s)
		}

		if rec2.RestVar != nil {
			entries1 := map[string]Type{}
			for _, key := range keys1MinusKeys2.List() {
				entries1[key] = rec1.Entries[key]
			}

			s, err := i.unify(rec2.RestVar, &TypeRec{
				Entries: entries1,
				RestVar: fresh,
				Union:   union,
			})
			if err != nil {
				return nil, err
			}
			subst = subst.compose(s)
		}
	} else {
		return nil, fmt.Errorf("incompatible types %s ~!~ %s", rec1.Pretty(0), rec2.Pretty(0))
	}

	return subst, nil
}

func (i *Inferrer) unifyCons(cons1, cons2 *TypeCons) (*Subst, error) {
	if cons1.Name != cons2.Name || len(cons1.Args) != len(cons2.Args) {
		return nil, fmt.Errorf("incompatible types %s ~!~ %s", cons1.Pretty(0), cons2.Pretty(0))
	}

	subst := &Subst{Subst: map[string]Type{}}
	for idx := range len(cons1.Args) {
		s, err := i.unify(cons1.Args[idx].apply(subst), cons2.Args[idx].apply(subst))
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

			s, err = i.unify(t, &TypeCons{
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
		s, err = i.unify(t.apply(subst), &TypeCons{
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

			s, err = i.unify(t, fresh)
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
		recType := &TypeRec{
			Entries: map[string]Type{},
			RestVar: nil, // not open
			Union:   false,
		}

		for _, entry := range expr.Entries {
			env = env.apply(subst)
			s, t, err := i.Infer(entry.Value, env)
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)
			recType.Entries[entry.Prop] = t
		}

		return subst, recType.apply(subst), nil
	case *Prop:
		s, t, err := i.Infer(expr.Parent, env)
		if err != nil {
			return nil, nil, err
		}
		subst = subst.compose(s)

		resultVar := i.freshVar()
		s, err = i.unify(t, &TypeRec{
			Entries: map[string]Type{expr.Prop: resultVar},
			RestVar: i.freshVar(),
		})
		if err != nil {
			return nil, nil, err
		}
		subst = subst.compose(s)
		return subst, resultVar.apply(subst), nil
	case *Cons:
		s, t, err := i.Infer(expr.Payload, env)
		if err != nil {
			return nil, nil, err
		}
		subst = subst.compose(s)
		return subst, &TypeRec{
			Entries: map[string]Type{
				expr.Name: t.apply(subst),
			},
			RestVar: i.freshVar(),
			Union:   true,
		}, nil
	case *When:
		var resultType Type = i.freshVar()

		s, valueType, err := i.Infer(expr.Value, env)
		if err != nil {
			return nil, nil, err
		}
		subst = subst.compose(s)
		valueType = valueType.apply(subst)

		for _, clause := range expr.Options {
			env = env.apply(subst)
			fresh := i.freshVar()
			s, t, err := i.Infer(clause.Consequence, env.extend(clause.Payload, &Scheme{
				Forall: nil,
				Type:   fresh,
			}))
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)

			var restVar *TypeVar
			if expr.Else != nil {
				restVar = i.freshVar()
			}
			s, err = i.unify(valueType, &TypeRec{
				Entries: map[string]Type{clause.ConsName: fresh},
				RestVar: restVar,
				Union:   true,
			})
			if err != nil {
				return nil, nil, err
			}

			subst = subst.compose(s)
			resultType = resultType.apply(subst)
			s, err = i.unify(resultType, t.apply(subst))
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)
		}

		if expr.Else != nil {
			s, t, err := i.Infer(expr.Else, env)
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)
			t = t.apply(subst)

			s, err = i.unify(valueType, &TypeRec{
				Entries: map[string]Type{},
				RestVar: i.freshVar(), // open set
				Union:   true,
			})
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)

			resultType = resultType.apply(subst)
			s, err = i.unify(resultType, t.apply(subst))
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)
		}

		return subst, resultType.apply(subst), nil
	case *Block:
		for _, assignment := range expr.Assignments {
			env = env.apply(subst)
			s, t, err := i.Infer(assignment.Value, env)
			if err != nil {
				return nil, nil, err
			}
			subst = subst.compose(s)
			env = env.extend(assignment.Name, generalize(t))
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
