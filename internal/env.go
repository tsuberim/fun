package internal

import (
	"github.com/pkg/errors"
	"maps"
)

type Item struct {
	Val
	Type *Scheme
}

type Env struct {
	Items map[string]Item
}

func (e *Env) Values() map[string]Val {
	result := map[string]Val{}
	for k, v := range e.Items {
		result[k] = v.Val
	}
	return result
}

func (e *Env) Types() *TypeEnv {
	result := &TypeEnv{Types: map[string]*Scheme{}}
	for k, v := range e.Items {
		result.Types[k] = v.Type
	}
	return result
}

func NewStdEnv(program *Program) *Env {
	return &Env{
		Items: map[string]Item{
			"+": {
				Type: &Scheme{
					Forall: nil,
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							&TypeCons{
								Name: intConsName,
								Args: nil,
							},
							&TypeCons{
								Name: intConsName,
								Args: nil,
							},
							&TypeCons{
								Name: intConsName,
								Args: nil,
							},
						},
					},
				},
				Val: &Builtin{
					Name: "+",
					Impl: func(args []Val) (Val, error) {
						sum := 0
						for _, arg := range args {
							i, ok := arg.(*Int)
							if !ok {
								return nil, errors.Errorf("invalid sum value type %t", arg)
							}

							sum += i.Value
						}

						return &Int{Value: sum}, nil
					},
				},
			},
			"-": {
				Type: &Scheme{
					Forall: nil,
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							&TypeCons{
								Name: intConsName,
								Args: nil,
							},
							&TypeCons{
								Name: intConsName,
								Args: nil,
							},
							&TypeCons{
								Name: intConsName,
								Args: nil,
							},
						},
					},
				},
				Val: &Builtin{
					Name: "-",
					Impl: func(args []Val) (Val, error) {
						first := args[0]
						i, ok := first.(*Int)
						if !ok {
							return nil, errors.Errorf("invalid sum value type %t", first)
						}

						sum := i.Value
						for _, arg := range args[1:] {
							i, ok := arg.(*Int)
							if !ok {
								return nil, errors.Errorf("invalid sum value type %t", arg)
							}

							sum -= i.Value
						}

						return &Int{Value: sum}, nil
					},
				},
			},
			"==": {
				Type: &Scheme{
					Forall: []string{"a"},
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							&TypeVar{Name: "a"},
							&TypeVar{Name: "a"},
							boolType,
						},
					},
				},
				Val: &Builtin{
					Name: "==",
					Impl: func(args []Val) (Val, error) {
						if len(args) != 2 {
							return nil, errors.Errorf("expecting 2 arguments, got %d", len(args))
						}

						arg1 := args[0]
						arg2 := args[1]
						if arg1.Pretty(0) == arg2.Pretty(0) {
							return trueVal, nil
						} else {
							return falseVal, nil
						}
					},
				},
			},
			"fix": {
				Type: &Scheme{
					Forall: []string{"a"},
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							&TypeCons{
								Name: lambdaConsName,
								Args: []Type{
									&TypeVar{Name: "a"},
									&TypeVar{Name: "a"},
								},
							},
							&TypeVar{Name: "a"},
						},
					},
				},
				Val: &Builtin{
					Name: "fix",
					Impl: func(args []Val) (Val, error) {
						cont, ok := args[0].(*Closure)
						if !ok {
							return nil, errors.Errorf("invalid closure type %t", args[0])
						}

						if 1 != len(cont.Params) {
							return nil, errors.Errorf("invalid number of arguments for function")
						}

						newEnv := maps.Clone(cont.Env)
						result, err := program.evaluator.Eval(cont.Body, newEnv)
						newEnv[cont.Params[0]] = result
						return result, err
					},
				},
			},
		},
	}
}
