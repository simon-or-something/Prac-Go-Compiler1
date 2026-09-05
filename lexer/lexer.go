package lexer

import (
	"fmt"
	"os"
	"strings"
)

//
type lexemeType int
const (
		LexemeNull lexemeType = iota
		LexemeOperator // { } + - -> ( )
		LexemeSymbol // word bro, word
		LexemeWhiteSpc
		LexemeNumeric // checks for . or f
		LexemeString // special parse rule, dont split on space / different type
)

type fileloc struct {
		x, y int
}

type Lexeme struct {
		ltype lexemeType
		src_loc fileloc // struct {x, y int} // oh lol this works
		src_file string
}

// https://stackoverflow.com/questions/21743841/how-to-avoid-annoying-error-declared-and-not-used
func Lex(filepath string) ([]Lexeme, error) {
		// it is possible to read it in chunks and pipeline it asynchronously
		// noted on roadmap, not bothered to make it overkill for a v0.0.1
		fconts, _ := os.ReadFile(filepath)
		lexemes := []Lexeme{}
		//lexemes = append(lexemes, Lexeme{LexemeNull, fileloc{3, 3}, "hello"})
		//lines := strings.Split(string(fconts), "\n")
		//fmt.Println(lines)
		for idx := range fconts {
		}
		return lexemes, nil
}
