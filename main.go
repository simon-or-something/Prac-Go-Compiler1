package main

import (
	_ "fmt"
	"os"

	"github.com/simon-or-something/Prac-Go-Compiler1/driver"
	"github.com/simon-or-something/Prac-Go-Compiler1/lexer"
)

func main() {
		lang_opts, _ := driver.RunDriver(os.Args)
		lexer.Lex(lang_opts.Infile)
}
