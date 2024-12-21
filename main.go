package main

import (
	"fmt"
	"github.com/maxott/go-repl"
	"log"
	"os"
)

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
	"-": {
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
	program, err := NewProgram()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		filename := os.Args[1]
		source, err := os.ReadFile(filename)
		if err != nil {
			println(err.Error())
		}

		mod, err := program.Run(source, RootModule)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		println(mod.Pretty(0))
	} else {
		r := repl.NewRepl(NewReplHandler(program))
		err := r.Loop()
		if err != nil {
			panic(err)
		}
	}
}

type ReplHandler struct{ program *Program }

func NewReplHandler(program *Program) *ReplHandler {
	return &ReplHandler{program: program}
}

func (r *ReplHandler) Prompt() string {
	return ">"
}

func (r *ReplHandler) Eval(buffer string) string {
	source := []byte(buffer)
	mod, err := r.program.Run(source, RootModule)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	return mod.Pretty(0)
}

func (r *ReplHandler) Tab(buffer string) string {
	// TODO: use tree-sitter to complete?
	return ""
}
