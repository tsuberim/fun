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
		if filename == "lsp" {
			internal.LSPServer()
			return
		}

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

		lines = append(lines, line)

		// Check if the last line ends with backslash (incomplete)
		if strings.TrimSpace(line) != "" && strings.HasSuffix(strings.TrimSpace(line), "\\") {
			// Remove the backslash and continue reading
			lines[len(lines)-1] = strings.TrimSuffix(strings.TrimSpace(line), "\\")
			// Change prompt to indicate continuation
			rl.SetPrompt(": ")
			continue
		}

		// Reset prompt to normal
		rl.SetPrompt("> ")

		// Submit the input if we have any lines
		if len(lines) > 0 {
			input := strings.Join(lines, "\n")
			result := evaluateInput(program, input)
			if result != "" {
				fmt.Println(result)
			}
			lines = nil
		}
	}
}

func evaluateInput(program *internal.Program, input string) string {
	source := []byte(input)
	mod, err := program.Run(source, internal.InlineModule)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	result := mod.Pretty(0)
	if strings.TrimSpace(result) == "" {
		return ""
	}
	return result
}
