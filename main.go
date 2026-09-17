package main

import (
	"fmt"
	"os"

	"github.com/simon-or-something/Prac-Go-Compiler1/driver"
	"github.com/simon-or-something/Prac-Go-Compiler1/lexer"
)

func printLexeme(lexeme lexer.Lexeme) {
	//fmt.Printf("%s @ (%d, %d) e %s, '%s'\n", lexer.LexemeTypeAsString(lexeme.Ltype), lexeme.Src_loc.X, lexeme.Src_loc.Y, lexeme.Src_file, lexeme.Raw_text)
	fmt.Printf("%s '%s'\n", lexer.LexemeTypeAsString(lexeme.Ltype), lexeme.Raw_text)
}

func main() {
	lang_opts, _ := driver.RunDriver(os.Args)
	lexemes, _ := lexer.Lex(lang_opts.Infile)
	for i := range len(lexemes) {
		printLexeme(lexemes[i])
		_ = i
	}
}
