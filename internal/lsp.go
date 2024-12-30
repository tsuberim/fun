package internal

import (
	"context"
	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp"
	"github.com/TobiasYin/go-lsp/lsp/defines"
	"log"
	"os"
)

func LSPServer() {
	logs.Init(log.New(os.Stderr, "", 0777))

	server := lsp.NewServer(&lsp.Options{CompletionProvider: &defines.CompletionOptions{
		TriggerCharacters: &[]string{"."},
	}})
	server.OnHover(func(ctx context.Context, req *defines.HoverParams) (result *defines.Hover, err error) {
		logs.Println(req)
		return &defines.Hover{Contents: defines.MarkupContent{Kind: defines.MarkupKindPlainText, Value: "hello world"}}, nil
	})

	server.OnCompletion(func(ctx context.Context, req *defines.CompletionParams) (result *[]defines.CompletionItem, err error) {
		logs.Println(req)
		d := defines.CompletionItemKindText
		str := "Hello"
		return &[]defines.CompletionItem{{
			Label:      "code",
			Kind:       &d,
			InsertText: &str,
		}}, nil
	})

	server.Run()
}
