package main

import (
	"fmt"
	tree_sitter_fun "fun/tree-sitter-fun/bindings/go"
	"github.com/maxott/go-repl"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	"maps"
	"os"
)

var stdlib = map[string]Val{
	"+": &Builtin{
		Name: "+",
		Impl: func(args []Val) (Val, error) {
			sum := 0
			for _, arg := range args {
				i, ok := arg.(*Int)
				if !ok {
					return nil, fmt.Errorf("invalid sum value type %t", arg)
				}

				sum += i.Value
			}

			return &Int{Value: sum}, nil
		},
	},
	"-": &Builtin{
		Name: "-",
		Impl: func(args []Val) (Val, error) {
			first := args[0]
			i, ok := first.(*Int)
			if !ok {
				return nil, fmt.Errorf("invalid sum value type %t", first)
			}

			sum := i.Value
			for _, arg := range args[1:] {
				i, ok := arg.(*Int)
				if !ok {
					return nil, fmt.Errorf("invalid sum value type %t", arg)
				}

				sum -= i.Value
			}

			return &Int{Value: sum}, nil
		},
	},
	"==": &Builtin{
		Name: "==",
		Impl: func(args []Val) (Val, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("expecting 2 arguments, got %d", len(args))
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
	"fix": &Builtin{
		Name: "fix",
		Impl: func(args []Val) (Val, error) {
			cont, ok := args[0].(*Closure)
			if !ok {
				return nil, fmt.Errorf("invalid closure type %t", args[0])
			}

			if 1 != len(cont.Params) {
				return nil, fmt.Errorf("invalid number of arguments for function")
			}

			newEnv := maps.Clone(cont.Env)
			result, err := Eval(cont.Body, newEnv)
			newEnv[cont.Params[0]] = result
			return result, err
		},
	},
}

var unitType = &TypeRec{
	Entries: map[string]Type{},
	RestVar: nil,
	Union:   false,
}

var boolType = &TypeRec{
	Entries: map[string]Type{"False": unitType, "True": unitType},
	RestVar: nil,
	Union:   true,
}

var falseVal = &ConsVal{
	Name:    "False",
	Payload: nil,
}

var trueVal = &ConsVal{
	Name:    "True",
	Payload: nil,
}

var typeEnv = &TypeEnv{Types: map[string]*Scheme{
	"+": {
		Forall: nil,
		Type: &TypeCons{
			Name: lambdaConsName,
			Args: []Type{
				&TypeCons{
					Name: "int",
					Args: nil,
				},
				&TypeCons{
					Name: "int",
					Args: nil,
				},
				&TypeCons{
					Name: "int",
					Args: nil,
				},
			},
		},
	},
	"-": {
		Forall: nil,
		Type: &TypeCons{
			Name: lambdaConsName,
			Args: []Type{
				&TypeCons{
					Name: "int",
					Args: nil,
				},
				&TypeCons{
					Name: "int",
					Args: nil,
				},
				&TypeCons{
					Name: "int",
					Args: nil,
				},
			},
		},
	},
	"==": {
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
	"fix": {
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
}}

func main() {
	if len(os.Args) > 1 {
		filename := os.Args[1]
		source, err := os.ReadFile(filename)
		if err != nil {
			println(err.Error())
		}
		println(eval(source))
	} else {
		r := repl.NewRepl(&ReplHandler{})
		err := r.Loop()
		if err != nil {
			panic(err)
		}
	}
}

func eval(source []byte) string {
	parser := tree_sitter.NewParser()
	defer parser.Close()
	err := parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_fun.Language()))
	if err != nil {
		panic(err)
	}

	tree := parser.Parse(source, nil)
	node := tree.RootNode()
	expr, err := fromNode(node, source)
	if err != nil {
		return fmt.Sprintf("ParseError: %s", err)
	}

	inferrer := NewInferrer()
	_, t, err := inferrer.Infer(expr, typeEnv)
	if err != nil {
		return fmt.Sprintf("TypeError: %s", err)
	}
	scheme := generalize(t)

	val, err := Eval(expr, stdlib)
	if err != nil {
		return fmt.Sprintf("ValueError: %s", err)
	}

	return fmt.Sprintf("%s : %s", val.Pretty(0), scheme.Pretty(0))
}

type ReplHandler struct{}

func (r *ReplHandler) Prompt() string {
	return ">"
}

func (r *ReplHandler) Eval(buffer string) string {
	source := []byte(buffer)
	return eval(source)
}

func (r *ReplHandler) Tab(buffer string) string {
	// TODO: use tree-sitter to complete?
	return ""
}
