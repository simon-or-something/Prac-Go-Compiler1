package lexer

import (
	"os"
	"regexp"
	"strings"
)

func getNextLexeme(fcons string, idx int) (LexemeType, LexemeLength, error) {
	// honestly, i like Gos decision to not include while loops in the language
	// even if that's a little backwards, for -> while -> goto (simply and usually)

	// this also simplifies control flow, there is no function calling a function calling a function
	// each function does one thing: this gets the next lexeme, including its "type"
	// before, i would need to find the length of the subexpr
	// https://whalelogic.io/posts/references/regex-go-guide/
	// but Go doesn't seem to have a regex function which contains this
	// so instead i specify that it should find the string _at the start_
	if idx >= len(fcons) {
		return LexemeEOF, 0, nil
	}
	for lex_type_idx := range len(LexemePatterns) {
		for pattern := range len(LexemePatterns[lex_type_idx]) {
			// https://zetcode.com/golang/regexp-quotemeta/
			// besides, this was a bug lul. i did regexp.QuoteMeta twice
			re := regexp.MustCompile("^" + LexemePatterns[lex_type_idx][pattern])
			// https://stackoverflow.com/questions/28886616/convert-array-to-slice-in-go
			// https://www.geeksforgeeks.org/go-language/strings-in-golang/
			// len somehow didn't work lul, but it was a casting issue
			// https://pkg.go.dev/regexp#Regexp.FindString
			if str := re.FindStringSubmatch(fcons[idx:]); str != nil {
				// so FindStringSubmatch returns []string, not string
				// which was cause for headache
				return LexemeType(lex_type_idx), LexemeLength(len(str[0])), nil
			}
			// https://gobyexample.com/if-else
			// oh and ternary isn't a thing either
			// i guess that stops people from writing silly stuff
			// hah! just watch me
			// let it be known i hhhhhhhhhhate Java with a passion
			// but somehow i don't mind the restrictiveness of Go
			// maybe because im older and sooooo much more mature
		}
	}
	// what genuinely appals me is the fact that unary op postfix increment is legitimate syntax sugar
	// for a single statement. it doesn't "return" anything like it does in C
	// in C you can "imagine" it as (offset += 1, offset), note the comma operator
	// in the previous attempt, i had a for loop that checks for the type of the lexeme by checking the next char
	// until a different lexeme type was returned. while this works conceptually, that isn't really necessary here
	// we (you; dear voice in my head; and i) have a useful regex engine now
	// i had another bug in the iteration that made it return negative numbers btw:
	// why?
	// idx is 0. i did idx - offset to find the length. So of course its negative
	// but offset _is_ the length
	// Ladies and gentlemen, a mathematician

	return LexemeError, 1, nil
}

// https://stackoverflow.com/questions/21743841/how-to-avoid-annoying-error-declared-and-not-used
func lexOpt1(filepath string) ([]Lexeme, error) {
	// it is possible to read it in chunks and pipeline it asynchronously
	// noted on roadmap, not bothered to make it overkill for a v0.0.1
	fcons, _ := os.ReadFile(filepath)
	lines := strings.Split(string(fcons), "\n")
	lexemes := []Lexeme{} // https://go.dev/wiki/SliceTricks, https://pkg.go.dev/golang.org/x/exp/slices

	file_idx := 0
	FILE_LENGTH := len(fcons)
	for file_idx < FILE_LENGTH {
		// now i am getting OOB errors
		lexeme_type, lexeme_length, _ := getNextLexeme(string(fcons), file_idx)
		location, _ := denormaliseFileLoc(lines, file_idx)
		// https://go.dev/tour/moretypes/7e
		// slices btw are indices like in python, not lengths. idk why i assumed otherwis
		lexemes = append(lexemes, Lexeme{lexeme_type,
			location,
			filepath,
			string(fcons)[file_idx : file_idx+int(lexeme_length)],
		})
		file_idx += int(lexeme_length)
	}

	return lexemes, nil
}
