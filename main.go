package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/simon-or-something/Prac-Go-Compiler1/driver"
	"github.com/simon-or-something/Prac-Go-Compiler1/lexer"
)

func main() {
	lang_opts, _ := driver.RunDriver(os.Args)
	lexemes, _ := lexer.Lex(lang_opts.Infile)
	for i := range len(lexemes) {
		fmt.Println(lexemes[i])
	}
	fmt.Printf("%T\n", unicode.IsDigit)
}
