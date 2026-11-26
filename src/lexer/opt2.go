package lexer

import (
	"os"
	"strings"
	"unicode"
)

type charType = int
const (
	//LexemeNull     LexemeType = iota
	charTypeOperator charType = iota // + - -> . ° x..y ...
	charTypeWhitespc                 // ' ', '\n', ...
	charTypeDelim                    //  { } [ ] ( ) ... split because makes parsing easier
	charTypeNumeric                  // checks for . specifically, or an x, an o, (or an f?)
	charTypeString                   // special parse rule, don't split on space / different type
	charTypeSymbol                   // all others
	charTypeError                    // anything not covered by symbol
)  // CharType

// yes this is similar to Lexeme
// because this is something that can be replaced (lexOpt2), and lexemes may later contain other attributes
// maybe like spell check or something for treesitter?
// ive never written a treesitter parser so idk yet
// https://github.com/llvm/llvm-project/blob/main/clang/include/clang/Lex/Token.h
// this is how llvm does it, and soon they may diverge anyhow (and to keep them separate because they have different goals)
// Lexeme: be the final output to pass into the Parser / CST Generator
// Atom: be an intermediate step to the lexemes, and to try a new approach
type atom struct {
	type_    charType
	src_loc  fileloc
	src_file string
	raw_text string
}

// https://www.golinuxcloud.com/go-can-functions-be-passed-as-parameters/
// i considered doing the above, but decided against it: i don't want to have my bad habits from c carry over into go
// but it could work if i do something called "sich auf die bullscheiße einlassen"

func appendAtom(charStack []rune, splitted []string, curr_file_idx int, atoms []atom, filepath string) ([]atom) {
	symbol := string(charStack)
	file_location, _ := denormaliseFileLoc(splitted, curr_file_idx)
	clear(charStack) // https://zetcode.com/golang/builtins-clear/
	return append(atoms, atom{charTypeSymbol, file_location, filepath, symbol})
}

func lexOpt2(filepath string) ([]Lexeme, error) {
	fcons, _ := os.ReadFile(filepath)
	splitted := strings.Split(string(fcons), "\n") // i can probably optimise this into a row / column var
	// TODO: make this a row / column var

	lexemes := []Lexeme{}
	atoms := []atom{}
	charStack := []rune{}

    curr_file_idx := 0
    FILE_LENGTH := len(fcons)

	// -1 so i can peek ahead to the next character
	for curr_file_idx < FILE_LENGTH {
		if false {
		} else if fcons[curr_file_idx] == '"' && (curr_file_idx == 0 || fcons[curr_file_idx - 1] != '\\') {
			// this is a string. i did this in reverse order because... idk i felt like it lol
			// because usually symbol becomes non trivial in my other approaches, so i wanted to see how good i can do this
			for curr_file_idx < FILE_LENGTH {
				charStack = append(charStack, rune(fcons[curr_file_idx]))
				curr_file_idx++
				if fcons[curr_file_idx] == '"' && fcons[curr_file_idx - 1] != '\\' {
					break
				}
			}
			atoms = appendAtom(charStack, splitted, curr_file_idx, atoms, filepath)
		} else if unicode.IsLetter(rune(fcons[curr_file_idx])) /* || unicode.IsPunct(rune(fcons[curr_file_idx])) */ { // names cannot start with _
			// this is a symbol [a-zA-Z][a-zA-Z0-9_]
			// https://www.geeksforgeeks.org/go-language/arrays-in-go/
			for curr_file_idx < FILE_LENGTH && unicode.IsOneOf([]*unicode.RangeTable{unicode.Letter, unicode.Punct, unicode.Digit}, rune(fcons[curr_file_idx])) {
				charStack = append(charStack, rune(fcons[curr_file_idx]))
				curr_file_idx++
			}
			atoms = appendAtom(charStack, splitted, curr_file_idx, atoms, filepath)
		} else {
			file_location, _ := denormaliseFileLoc(strings.Split(string(fcons), "\n"), curr_file_idx)
			atoms = append(atoms, atom{charTypeError, file_location, filepath, string(fcons[curr_file_idx])})
		}
		curr_file_idx++
	}

	atom_idx := 0
	for atom_idx < len(atoms) {
		atom_elem := atoms[atom_idx]
		println(atom_elem.type_, atom_elem.raw_text)
		atom_idx++
	}

    return lexemes, nil
}

