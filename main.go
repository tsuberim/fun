package main

import (
	"fmt"
	"fun/internal"
	"github.com/maxott/go-repl"
	"log"
	"os"
)

func main() {
	program, err := internal.NewProgram()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		filename := os.Args[1]
		source, err := os.ReadFile(filename)
		if err != nil {
			println(err.Error())
		}

		mod, err := program.Run(source, filename)
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

type ReplHandler struct{ program *internal.Program }

func NewReplHandler(program *internal.Program) *ReplHandler {
	return &ReplHandler{program: program}
}

func (r *ReplHandler) Prompt() string {
	return ">"
}

func (r *ReplHandler) Eval(buffer string) string {
	source := []byte(buffer)
	mod, err := r.program.Run(source, internal.InlineModule)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	return mod.Pretty(0)
}

func (r *ReplHandler) Tab(buffer string) string {
	// TODO: use tree-sitter to complete?
	return ""
}
