package lexer

import (
	_ "errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type lexemeType int

const (
	LexemeNull     lexemeType = iota
	LexemeOperator            // + - -> . ° x..y ...
	LexemeSymbol              // word bro, word
	LexemeWhiteSpc
	LexemeDelim   //  { } [ ] ( ) ; : ... split because makes parsing easier
	LexemeNumeric // checks for one . specifically, or an f
	LexemeString  // special parse rule, don't split on space / different type
	LexemeComment // part of me wants to make operator precedence based on whitespace
	// LexemeBool // Booleans don't exist, https://en.wikipedia.org/wiki/Church_encoding
	//LexemeLiteral // not this and instead ...Numeric and ...String because this makes parsing easier down the line
	LexemeError   // that's useful
)

type fileloc struct {
	x, y int
}

type Lexeme struct {
	ltype    lexemeType
	src_loc  fileloc // struct {x, y int} // oh lol this works
	src_file string
}

func matchCategoryFromLength(fnc func(rune) bool, idx int, file string) (int, bool, error) {
	new_idx := idx
	if fnc(rune(file[new_idx])) {
		for fnc(rune(file[new_idx])) {
			new_idx++
		}
		return new_idx - idx, true, nil
	}
	// i can't decide if i want this to be an error, because it isn't an error
	// but idk if making a sentinel value is idiosyncratic
	return 0, false, nil
}

// this return type flexibility feels very strange
func lexemeLength(idx int, file string) (int, lexemeType, error) {
	if idx >= len(file) {
		return 0, LexemeNull, nil
	}

	if length, ok, _ := matchCategoryFromLength(unicode.IsSpace, idx, file); ok {
		return length, LexemeWhiteSpc, nil
	}

	if length, ok, _ := matchCategoryFromLength(unicode.IsLetter, idx, file); ok {
		return length, LexemeSymbol, nil
	}

	/*
	LexemeNull
	LexemeOperator
	LexemeSymbol
	LexemeWhiteSpc
	LexemeDelim
	LexemeNumeric
	LexemeString
	LexemeComment
	LexemeError
	*/

	if unicode.IsDigit(rune(file[idx])) {
		new_idx := idx
		dot_count := 0
		for unicode.IsNumber(rune(file[new_idx])) || (file[new_idx] == '.' && dot_count < 1) {
			new_idx++
			// https://gobyexample.com/if-else
			if file[new_idx] == '.' {
				dot_count++
			}
		}
		return new_idx - idx, LexemeNumeric, nil
	}
	return 0, LexemeError, nil
}

// https://stackoverflow.com/questions/21743841/how-to-avoid-annoying-error-declared-and-not-used
func Lex(filepath string) ([]Lexeme, error) {
	// it is possible to read it in chunks and pipeline it asynchronously
	// noted on roadmap, not bothered to make it overkill for a v0.0.1
	fconts, _ := os.ReadFile(filepath)
	lexemes := []Lexeme{} // https://go.dev/wiki/SliceTricks, https://pkg.go.dev/golang.org/x/exp/slices
	//lexemes = append(lexemes, Lexeme{LexemeNull, fileloc{3, 3}, "hello"})
	lines := strings.Split(string(fconts), "\n")
	//fmt.Println(lines)
	for idx := range lines {
		fmt.Printf("<%s>\n", string(lines[idx]))
	}
	return lexemes, nil
}
