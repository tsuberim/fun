package main

import (
	"fmt"
	"fun/internal"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
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

		val, typ := program.EvalTask(mod)

		println(Pretty(val, typ, 0))
	} else {
		runReadlineRepl(program)
	}
}

func runReadlineRepl(program *internal.Program) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/fun_repl_history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start REPL:", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Println("Type .help for REPL commands.")

	var lines []string
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(lines) == 0 {
				break
			} else {
				lines = nil
				continue
			}
		} else if err != nil {
			break
		}

		if strings.TrimSpace(line) == "quit" || strings.TrimSpace(line) == "exit" {
			break
		}

		// Handle .clear command
		if strings.TrimSpace(line) == ".clear" {
			fmt.Print("\033[H\033[2J")
			lines = nil
			continue
		}

		// Handle .help command
		if strings.TrimSpace(line) == ".help" {
			fmt.Println(`REPL commands:
  .help   Show this help message
  .clear  Clear the screen
  quit    Exit the REPL
  exit    Exit the REPL

Multiline: End a line with '\' to continue input on the next line.`)
			lines = nil
			continue
		}

		lines = append(lines, line)

		// Check if the last line ends with backslash (explicit continuation)
		if strings.TrimSpace(line) != "" && strings.HasSuffix(strings.TrimSpace(line), "\\") {
			lines[len(lines)-1] = strings.TrimSuffix(strings.TrimSpace(line), "\\")
			rl.SetPrompt(": ")
			continue
		}

		rl.SetPrompt("> ")

		if len(lines) > 0 {
			input := strings.Join(lines, "\n")
			if strings.TrimSpace(input) == "" {
				lines = nil
				continue
			}
			result := evaluateInput(program, input)
			if result != "" {
				fmt.Println(result)
			}
			lines = nil
		}
	}
}

func Pretty(val internal.Val, scheme *internal.Scheme, indent int) string {
	return fmt.Sprintf("%s : %s", val.Pretty(indent), scheme.Pretty(indent))
}

func evaluateInput(program *internal.Program, input string) string {
	source := []byte(input)
	mod, err := program.Run(source, internal.InlineModule)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	val, typ := program.EvalTask(mod)

	result := Pretty(val, typ, 0)
	if strings.TrimSpace(result) == "" {
		return ""
	}
	return result
}
