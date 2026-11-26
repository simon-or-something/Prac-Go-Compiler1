package lexer

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type charType = int
const (
	charTypeOperator charType = iota // + - -> . ° x..y ...
	charTypeWhitespc                 // ' ', '\n', ...
	charTypeDelim                    //  { } [ ] ( ) ... split because makes parsing easier
	charTypeNumeric                  // checks for . specifically, or an x, an 0, (or an f?)
	charTypeString                   // special parse rule, don't split on space / different type
	charTypeSymbol                   // all others
	charTypeError                    // anything not covered by symbol
)  // CharType

type dt = int
const (
	digitTypeBinary dt = iota // 0b... why not twonal, like tonal? hahahahaaaa
	//digitTypeSeximal         // 0s... base 6, it's because that's why
	digitTypeOctal            // 0o...
	digitTypeDecimal          // ...
	//digitTypeDodecimal       // 0d... i always wanted base 12
	digitTypeHexadecimal      // 0x...
) // DigitType

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

func appendAtom(charStack []rune, file_location fileloc, atoms []atom, filepath string, cType charType) ([]atom) {
	symbol := string(charStack)
	clear(charStack) // https://zetcode.com/golang/builtins-clear/
	return append(atoms, atom{cType, file_location, filepath, symbol})
}

func lexOpt2(filepath string) ([]Lexeme, error) {
	fcons, _ := os.ReadFile(filepath)

	lexemes := []Lexeme{}
	atoms := []atom{}
	charStack := []rune{}

	file_location := fileloc{-1, 0} // so i can increment at the top of the loop i guess
    curr_file_idx := 0
    FILE_LENGTH := len(fcons)

	for curr_file_idx < FILE_LENGTH {
		file_location.X += 1
		if false {
		} else /*NUMERIC*/ if unicode.IsDigit(rune(fcons[curr_file_idx])) { // so now numbers have to start with \d[.x]\d* or some bs
			base := digitTypeDecimal
			if fcons[curr_file_idx] == '0' {
				if (curr_file_idx + 2) < FILE_LENGTH {
					switch fcons[curr_file_idx + 1] {
					case 'b':
						base = digitTypeBinary
					/*
					case 's':
					    base = digitTypeSeximal
					*/
					case 'o':
						base = digitTypeOctal
					/*
					case 'd':
					    base = digitTypeDodecimal
					*/
					case 'x':
						base = digitTypeHexadecimal
					}
				} else {
					// TODO: this
					println("invalid numeric at (" + strconv.Itoa(file_location.X) + ", " + strconv.Itoa(file_location.Y) + "), but error not implemented yet")
				}
			}
			for curr_file_idx < FILE_LENGTH && (unicode.IsDigit(rune(fcons[curr_file_idx])) || strings.IndexByte("bsodx", fcons[curr_file_idx]) >= 0) {
				switch fcons[curr_file_idx] {
				}
				charStack = append(charStack, rune(fcons[curr_file_idx]))
				curr_file_idx++
				if fcons[curr_file_idx] == '"' && fcons[curr_file_idx - 1] != '\\' {
					break
				}
			}
			atoms = appendAtom(charStack, file_location, atoms, filepath, charTypeString)
			file_location.X += len(charStack)
		} else /*STRING*/ if fcons[curr_file_idx] == '"' && (curr_file_idx == 0 || fcons[curr_file_idx - 1] != '\\') {
			// i originally did the if stmts in reverse order because... idk i felt like it lol
			// because usually symbol becomes non trivial in my other approaches, so i wanted to see how good i can do this
			for curr_file_idx < FILE_LENGTH {
				charStack = append(charStack, rune(fcons[curr_file_idx]))
				curr_file_idx++
				if fcons[curr_file_idx] == '"' && fcons[curr_file_idx - 1] != '\\' {
					break
				}
			}
			atoms = appendAtom(charStack, file_location, atoms, filepath, charTypeString)
			file_location.X += len(charStack)
		} else /*SYMBOL*/ if unicode.IsOneOf([]*unicode.RangeTable{unicode.Letter, unicode.Punct}, rune(fcons[curr_file_idx])) {
			// this is [a-zA-Z_][a-zA-Z0-9_]. Note to self: if i want to not let names start with _, remove Punct from the above if stmt
			// https://www.geeksforgeeks.org/go-language/arrays-in-go/
			for curr_file_idx < FILE_LENGTH && unicode.IsOneOf([]*unicode.RangeTable{unicode.Letter, unicode.Punct, unicode.Digit}, rune(fcons[curr_file_idx])) {
				charStack = append(charStack, rune(fcons[curr_file_idx]))
				curr_file_idx++
			}
			atoms = appendAtom(charStack, file_location, atoms, filepath, charTypeSymbol)
			file_location.X += len(charStack)
		} else /*ERROR*/ {
			atoms = append(atoms, atom{charTypeError, file_location, filepath, string(fcons[curr_file_idx])})
			file_location.X += len(charStack)
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

