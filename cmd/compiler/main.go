package main

import (
	"os"

	"github.com/simon-or-something/Prac-Go-Compiler1/src/driver"
	"github.com/simon-or-something/Prac-Go-Compiler1/src/lexer"
)

func main() {
	lang_opts, _ := driver.RunDriver(os.Args)
	lexemes, _ := lexer.Lex(lang_opts.Infile)
	for i := range len(lexemes) {
		lexer.PrintLexeme(lexemes[i])
		_ = i
	}
}
