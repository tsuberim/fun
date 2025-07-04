package internal

import (
	"maps"
	"os"

	"github.com/pkg/errors"
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

func taskType(result Type, errType Type) *TypeCons {
	return &TypeCons{
		Name: taskConsName,
		Args: []Type{
			result,
			errType,
		},
	}
}

func lamType(args ...Type) *TypeCons {
	return &TypeCons{
		Name: lambdaConsName,
		Args: args,
	}
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
					Impl: func(e *Evaluator, args []Val) (Val, error) {
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
					Impl: func(e *Evaluator, args []Val) (Val, error) {
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
					Impl: func(e *Evaluator, args []Val) (Val, error) {
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
					Impl: func(e *Evaluator, args []Val) (Val, error) {
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
			"flat_map": {
				Type: &Scheme{
					Forall: []string{"a", "b", "e"},
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							taskType(&TypeVar{Name: "a"}, &TypeVar{Name: "e"}),
							lamType(&TypeVar{Name: "a"}, taskType(&TypeVar{Name: "b"}, &TypeVar{Name: "e"})),
							taskType(&TypeVar{Name: "b"}, &TypeVar{Name: "e"}),
						},
					},
				},
				Val: &Builtin{
					Name: "flat_map",
					Impl: func(e *Evaluator, args []Val) (Val, error) {
						task := args[0]
						mapper := args[1]
						return &Builtin{
							Name: "flat_map_thunk",
							Impl: func(e *Evaluator, args []Val) (Val, error) {
								// execute first task
								res, err := e.evalFn(task, nil)
								if err != nil {
									return nil, err
								}
								// map the task result to consequent task
								outTask, err := e.evalFn(mapper, []Val{res})
								if err != nil {
									return nil, err
								}
								// execute the consequent task
								out, err := e.evalFn(outTask, nil)
								if err != nil {
									return nil, err
								}

								return out, nil
							},
						}, nil
					},
				},
			},
			"ok": {
				Type: &Scheme{
					Forall: []string{"a", "rest"},
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							&TypeVar{Name: "a"},
							taskType(&TypeVar{Name: "a"}, &TypeRec{
								Entries: map[string]Type{},
								RestVar: &TypeVar{Name: "rest"},
								Union:   true,
							}),
						},
					},
				},
				Val: &Builtin{
					Name: "ok",
					Impl: func(e *Evaluator, args []Val) (Val, error) {
						if len(args) != 1 {
							return nil, errors.Errorf("expecting 1 arguments, got %d", len(args))
						}
						arg := args[0]
						return &Builtin{
							Name: "ok_thunk",
							Impl: func(e *Evaluator, args []Val) (Val, error) {
								return arg, nil
							},
						}, nil
					},
				},
			},
			"err": {
				Type: &Scheme{
					Forall: []string{"rest"},
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							&TypeCons{
								Name: strConsName,
								Args: nil,
							},
							taskType(neverType, &TypeRec{
								Entries: map[string]Type{"Err": &TypeCons{
									Name: strConsName,
									Args: nil,
								}},
								RestVar: &TypeVar{Name: "rest"},
								Union:   true,
							}),
						},
					},
				},
				Val: &Builtin{
					Name: "err",
					Impl: func(e *Evaluator, args []Val) (Val, error) {
						if len(args) != 1 {
							return nil, errors.Errorf("expecting 1 arguments, got %d", len(args))
						}
						arg := args[0]
						str, ok := arg.(*LitStr)
						if !ok {
							return nil, errors.Errorf("expecting type string literal")
						}
						return &Builtin{
							Name: "err_thunk",
							Impl: func(e *Evaluator, args []Val) (Val, error) {
								return nil, errors.Errorf(str.Value)
							},
						}, nil
					},
				},
			},
			"write": {
				Type: &Scheme{
					Forall: []string{"rest"},
					Type: &TypeCons{
						Name: lambdaConsName,
						Args: []Type{
							&TypeCons{
								Name: strConsName,
								Args: nil,
							},
							&TypeCons{
								Name: strConsName,
								Args: nil,
							},
							taskType(
								unitType,
								&TypeRec{
									Entries: map[string]Type{"Err": &TypeCons{
										Name: strConsName,
										Args: nil,
									}},
									RestVar: &TypeVar{Name: "rest"},
									Union:   true,
								},
							),
						},
					},
				},
				Val: &Builtin{
					Name: "write",
					Impl: func(e *Evaluator, args []Val) (Val, error) {
						if len(args) != 2 {
							return nil, errors.Errorf("expecting 2 arguments, got %d", len(args))
						}
						filename, ok := args[0].(*LitStr)
						if !ok {
							return nil, errors.Errorf("invalid filename type %t", args[1])
						}

						content, ok := args[1].(*LitStr)
						if !ok {
							return nil, errors.Errorf("invalid content type %t", args[0])
						}

						return &Builtin{
							Name: "write_thunk",
							Impl: func(e *Evaluator, args []Val) (Val, error) {
								err := os.WriteFile(filename.Value, []byte(content.Value), 0644)
								if err != nil {
									return nil, err
								}
								return unitVal, nil
							},
						}, nil
					},
				},
			},
		},
	}
}
