package parser

import (
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)


type ParsedCommand struct {
	Original    string
	Canonical   string
	CommandName string 
}

func Parse(cmdStr string) (*ParsedCommand, error) {
	parser := syntax.NewParser(syntax.KeepComments(false))
	file, err := parser.Parse(strings.NewReader(cmdStr), "")
	if err != nil {
		return nil, err
	}

	cmdName := extractPrimaryCommand(file)
	canonicalizeAST(file)

	printer := syntax.NewPrinter()
	var buf strings.Builder
	err = printer.Print(&buf, file)
	if err != nil {
		return nil, err
	}

	return &ParsedCommand{
		Original:    cmdStr,
		Canonical:   strings.TrimSpace(buf.String()),
		CommandName: cmdName,
	}, nil
}

func wordToString(w *syntax.Word) string {
	var sb strings.Builder
	for _, part := range w.Parts {
		switch x := part.(type) {
		case *syntax.Lit:
			if x.Value != "\\" { 
				sb.WriteString(x.Value)
			}
		case *syntax.SglQuoted: 
			sb.WriteString(x.Value) 
		case *syntax.DblQuoted:
			for _, dp := range x.Parts {
				if dl, ok := dp.(*syntax.Lit); ok {
					sb.WriteString(dl.Value)
				}
			}
		}
	}
	return sb.String()
}
func extractPrimaryCommand(node syntax.Node) string {
	var name string
	syntax.Walk(node, func(n syntax.Node) bool {
		if name != "" {
			return false 
		}
		if x, ok := n.(*syntax.CallExpr); ok {
			if len(x.Args) > 0 {
				cleanName := wordToString(x.Args[0])
				if cleanName != "" {
					name = filepath.Base(cleanName)
				}
			}
			return false
		}
		return true
	})
	return name
}

func canonicalizeAST(node syntax.Node) {
	syntax.Walk(node, func(n syntax.Node) bool {
		if x, ok := n.(*syntax.CallExpr); ok {
			for i, arg := range x.Args {
				if i == 0 {
					continue 
				}

				hasSubshell := false
				for _, part := range arg.Parts {
					if _, isCmd := part.(*syntax.CmdSubst); isCmd {
						hasSubshell = true
						break
					}
				}

				
				if hasSubshell {
					continue
				}

				cleanStr := wordToString(arg)
				if !strings.HasPrefix(cleanStr, "-") && cleanStr != "" {
					
					arg.Parts = []syntax.WordPart{
						&syntax.Lit{Value: "<ARG>"},
					}
				}
			}
		}
		return true
	})
}