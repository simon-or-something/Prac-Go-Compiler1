package lexer

import (
	_ "errors"
	"os"
	"regexp"
	"strings"
)

type LexemeType int

// i typedef types so the type itself carries information about the use case
type LexemeLength int // this is a trick i do in mock-professional C
// the issue is that Go seems to have nominal typing
// this is *safer* for types, but that means the trick doesn't work here

// ill use regex based matching, see 'BraCkish' for a C implementation for a LISP
const (
	//LexemeNull     LexemeType = iota
	LexemeOperator LexemeType = iota // + - -> . ° x..y ...
	LexemeKeyword                    // word bro, word
	LexemeWhiteSpc                   // part of me wants to make operator precedence based on whitespace
	LexemeDelim                      //  { } [ ] ( ) ... split because makes parsing easier
	LexemeNumeric                    // checks for one . specifically, or an f
	LexemeString                     // special parse rule, don't split on space / different type
	LexemeComment
	LexemeSymbol // user variables, follows a scheme
	// LexemeBool // Booleans don't exist, https://en.wikipedia.org/wiki/Church_encoding
	//LexemeLiteral // not this and instead ...Numeric and ...String because this makes parsing easier down the line
	LexemeEOF
	LexemeError // that's useful. not a pattern itself but useful for parsing
) // LexemeType

// TODO: nyi = not yet implemented
// so apparently:
/*
var arr = [][]string {
    {"hello"},
    {"deez nuts", "balls"},
}
*/
// is valid, legal, perfectly acceptable Go
// honourable mentions for keywords: `chain` (resolved by |), `as` (cast), `blank` (;)
// some of these keywords are horribly cursed, I am well aware of this
// but you cannot stop me
// MARK: it is absolutely imperative that the LexemeType corresponds to the index in lexeme_patterns
var lexeme_patterns = [][]string{
	{
		regexp.QuoteMeta("->"),  //     nyi: copy | proof / implication
		regexp.QuoteMeta("<-"),  //     nyi: storage | suggestion
		regexp.QuoteMeta(">>="), //     nyi: bind
		regexp.QuoteMeta("=<<"), //     nyi: deferred application
		regexp.QuoteMeta("."),   // member access
		regexp.QuoteMeta(","),   //     nyi: structural assignment
		regexp.QuoteMeta("_"),   //     nyi: catch all
		regexp.QuoteMeta("°"),   //     nyi: after operator
		regexp.QuoteMeta(";"),   // null statement
		regexp.QuoteMeta("?"),   //     nyi: ternary prefix
		regexp.QuoteMeta(":"),   //     nyi: ternary suffix
		regexp.QuoteMeta("^"),   // xor operator
		regexp.QuoteMeta("~"),   // negation
		regexp.QuoteMeta("|"),   // chain | logical or
		regexp.QuoteMeta("&"),   // reference | logical and
		regexp.QuoteMeta("||"),  // bitwise or
		regexp.QuoteMeta("&&"),  // bitwise and
		regexp.QuoteMeta("=="),  // equality
		regexp.QuoteMeta("<"),   // lesser
		regexp.QuoteMeta(">"),   // greater
		regexp.QuoteMeta(">="),  // lesser or equal
		regexp.QuoteMeta("<="),  // greater or equal
		regexp.QuoteMeta("="),   // assignment
		regexp.QuoteMeta("+"),   // addition
		regexp.QuoteMeta("-"),   // negation / subtraction
		regexp.QuoteMeta("*"),   // multiplication
		regexp.QuoteMeta("/"),   // division
		regexp.QuoteMeta("%"),   // modulation
		regexp.QuoteMeta("\\"),  //     nyi: idk yet
		regexp.QuoteMeta("`"),   //     nyi: homoiconicity
		regexp.QuoteMeta("'"),   // char
		regexp.QuoteMeta("#"),   //     nyi: to raw bytes
	}, // LexemeOperator
	{
		regexp.QuoteMeta("nil"),        // null, active sentinel value
		regexp.QuoteMeta("func"),       // for functions like in go
		regexp.QuoteMeta("class"),      //     nyi: classes, PODs by default
		regexp.QuoteMeta("import"),     //     nyi: self explanatory
		regexp.QuoteMeta("pure"),       //     nyi: optimisation directive
		regexp.QuoteMeta("static"),     //     nyi: non local storage declaration
		regexp.QuoteMeta("debug"),      // prints the AST / relevant info at that location
		regexp.QuoteMeta("assert"),     // crashes on bool expression being false
		regexp.QuoteMeta("switch"),     //     nyi: hash search for conditions
		regexp.QuoteMeta("case"),       //     nyi: individual cases for the hash search
		regexp.QuoteMeta("defer"),      //     nyi: execution at the end of scope
		regexp.QuoteMeta("continue"),   //     nyi: syntax sugar for ~goto~ loop / if not applies:
		regexp.QuoteMeta("delete"),     //     nyi: storage deletion
		regexp.QuoteMeta("false"),      // desugars to 0
		regexp.QuoteMeta("true"),       // desugars to 1
		regexp.QuoteMeta("for"),        // init + condition + step based loop
		regexp.QuoteMeta("while"),      //     nyi: condition based loop, can be chained with do
		regexp.QuoteMeta("until"),      //     nyi: reverse while but can be chained with do
		regexp.QuoteMeta("do"),         //     nyi: execution + check based loop
		regexp.QuoteMeta("break"),      //     nyi: stop loop execution
		regexp.QuoteMeta("goto"),       //     nyi: deprecated (XDXD)
		regexp.QuoteMeta("lngjmp"),     //     nyi: goto + explicit side entry into different scopes
		regexp.QuoteMeta("srtjmp"),     //     nyi: goto + same scope / global scope
		regexp.QuoteMeta("if"),         // branching
		regexp.QuoteMeta("when"),       //     nyi: variable watching
		regexp.QuoteMeta("throw"),      // toss an error yo
		regexp.QuoteMeta("try"),        // try but allow me to baseball huh an error
		regexp.QuoteMeta("catch"),      // baseball huh an error
		regexp.QuoteMeta("processing"), //     nyi: desugars into return + lngjmp back into the function; yield
		regexp.QuoteMeta("constexpr"),  //     nyi: this expression can be evaluated at compiletime
	}, // LexemeKeyword
	{
		"\\s",
	}, // LexemeWhiteSpc
	{
		regexp.QuoteMeta("("),
		regexp.QuoteMeta(")"),
		regexp.QuoteMeta("["),
		regexp.QuoteMeta("]"),
		regexp.QuoteMeta("{"),
		regexp.QuoteMeta("}"),
	}, // LexemeDelim
	{
		"[0-9]*(?.)[0-9]+",
	}, // LexemeNumeric
	{
		"\".*\"",
	}, // LexemeString
	{
		"\\/\\/.*?\n",
		"\\/\\*.*?\\*\\/",
	}, // LexemeComment
	{
		"[_a-zA-Z][_a-zA-Z0-9]*",
	}, // LexemeSymbol
} // [LexemeType][RegexPattern]

type fileloc struct {
	x, y int
} // vec2<int>

type Lexeme struct {
	ltype    LexemeType
	src_loc  fileloc // struct {x, y int} // oh lol this works
	src_file string
}

func getLexemeType(fcons string, idx int) (LexemeType, error) {
	// honestly, i like Gos decision to not include while loops in the language
	// even if that's a little backwards, for -> while -> goto (simply and usually)
	if idx >= len(fcons) {
		return LexemeEOF, nil
	}
	for lex_type_idx := range len(lexeme_patterns) {
		for pattern := range len(lexeme_patterns[lex_type_idx]) {
			// https://zetcode.com/golang/regexp-quotemeta/
			re := regexp.MustCompile(regexp.QuoteMeta(lexeme_patterns[lex_type_idx][pattern]))
			// https://stackoverflow.com/questions/28886616/convert-array-to-slice-in-go
			// https://www.geeksforgeeks.org/go-language/strings-in-golang/
			// len somehow didn't work lul, but it was a casting issue
			if re.MatchString(fcons[idx:]) {
				return LexemeType(lex_type_idx), nil
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
	return LexemeError, nil
}

func getNextLexeme(fcons string, idx int) (LexemeType, LexemeLength, error) {
	// for this to work I would need to find the length of the subexpr
	// https://whalelogic.io/posts/references/regex-go-guide/
	// but Go doesn't seem to have a regex function which contains this
	INITIAL_LEXEME_TYPE, _ := getLexemeType(fcons, idx)
	offset := 0
	// i would need 2 post statements (assignment : new_type, and operation: increment offset)
	// so i wont do either lul
	// but then again, this for loop only calculates the length of the lexeme by matching the type
	// what genuinely appals me is the fact that unary op postfix increment is legitimate syntax sugar
	// for a single statement. It doesn't "return" anything like it does in C
	// in C you can "imagine" it as (offset += 1, offset), note the comma operator

	// i know this line is long but the format says so and i follow the format
	// although i would much prefer 2 spaces instead of 1 tab but tally ho lads innit
	// the first issue is my exam in 10 days
	// the second issue is this compiler isn't done
	// the third issue is the fact that my go formatter is "nOt To My ExPeCtAtIoNs" (skill issue)

	// okay somehow this returns negative numbers
	// why the f### does it return negative numbers
	// ...
	// ...
	// because offset is 0 and this returns the length, offset is the length itself
	// not the new index. Ladies and gentlemen, a mathematician
	for new_type := INITIAL_LEXEME_TYPE; new_type == INITIAL_LEXEME_TYPE && new_type != LexemeEOF; new_type, _ = getLexemeType(fcons, idx + offset) {
		offset++
	}
	// now, at offset, there will be a new lexeme, so i acknowledge this and say that the lexeme is [idx, retval)
	return INITIAL_LEXEME_TYPE, LexemeLength(offset), nil
}

// a normalised file location is the raw byte index
// denormalised includes line and column
// it is cheaper to just pass the array instead of recomputing it each time
func denormaliseFileLoc(fcons []string, loc int) (fileloc, error) {
	line := 0
	for file_idx := 0; loc > len(fcons[file_idx]); file_idx++ {
		loc -= len(fcons[file_idx])
		line++
	}
	return fileloc{loc, line}, nil
}

// https://stackoverflow.com/questions/21743841/how-to-avoid-annoying-error-declared-and-not-used
func Lex(filepath string) ([]Lexeme, error) {
	// it is possible to read it in chunks and pipeline it asynchronously
	// noted on roadmap, not bothered to make it overkill for a v0.0.1
	fcons, _ := os.ReadFile(filepath)
	lines := strings.Split(string(fcons), "\n")
	lexemes := []Lexeme{} // https://go.dev/wiki/SliceTricks, https://pkg.go.dev/golang.org/x/exp/slices
	//lexemes = append(lexemes, Lexeme{LexemeNull, fileloc{3, 3}, "hello"})

	file_idx := 0
	FILE_LENGTH := len(fcons)
	for file_idx < FILE_LENGTH {
		// now i am getting OOB errors
		lexeme_type, lexeme_length, _ := getNextLexeme(string(fcons), file_idx)
		location, _ := denormaliseFileLoc(lines, file_idx)
		println(string(fcons[file_idx:lexeme_length]))
		println("---")
		lexemes = append(lexemes, Lexeme{lexeme_type, location, string(fcons[file_idx:lexeme_length - 1])})
		file_idx += int(lexeme_length) + 1
	}

	return lexemes, nil
}
