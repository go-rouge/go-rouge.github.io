// Command lexlist prints every lexer go-rouge registers, as YAML for the
// landing page's data directory. Generating it from the library is what keeps
// the page from claiming a language the code does not have, or omitting one it
// does.
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-rouge/rouge"
)

func main() {
	// Every tag and alias the registry answers to, discovered by asking the
	// library rather than by keeping a list here.
	seen := map[string]*lex{}
	var order []string
	for _, probe := range candidates() {
		l := rouge.FindLexer(probe)
		if l == nil {
			continue
		}
		if _, ok := seen[l.Tag()]; ok {
			continue
		}
		seen[l.Tag()] = &lex{Title: l.Title(), Tag: l.Tag(), Aliases: l.Aliases()}
		order = append(order, l.Tag())
	}
	sort.Slice(order, func(i, j int) bool {
		return strings.ToLower(seen[order[i]].Title) < strings.ToLower(seen[order[j]].Title)
	})

	var b strings.Builder
	b.WriteString("# Generated: (cd tools/lexlist && go run .) > data/languages.yaml\n")
	b.WriteString("# Every entry is a lexer go-rouge/rouge actually registers.\n")
	fmt.Fprintf(&b, "languages:\n")
	for _, tag := range order {
		l := seen[tag]
		tags := append([]string{l.Tag}, l.Aliases...)
		fmt.Fprintf(&b, "  - { title: %q, tags: %q }\n", l.Title, strings.Join(tags, ", "))
	}
	if _, err := os.Stdout.WriteString(b.String()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%d lexers\n", len(order))
}

type lex struct {
	Title   string
	Tag     string
	Aliases []string
}

// candidates lists one known tag per lexer. FindLexer resolves each to its
// canonical lexer, so a wrong guess simply yields nothing.
func candidates() []string {
	return []string{
		"plaintext", "json", "diff", "go", "sql", "css", "shell", "yaml",
		"python", "javascript", "html", "markdown", "ruby", "php", "tex",
		"java", "c", "cpp", "rust", "cobol", "fortran", "hcl", "terraform",
	}
}
