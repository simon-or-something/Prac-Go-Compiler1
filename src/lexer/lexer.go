package lexer

import (
	_ "errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type LexemeType int

// i typedef types so the type itself carries information about the use case
type LexemeLength int // this is a trick i do in mock-professional C
// the issue is that Go seems to have nominal typing
// this is *safer* for types, but that means the trick doesn't work here

type fileloc struct {
	X, Y int
} // vec2<int>

type Lexeme struct {
	Ltype    LexemeType
	Src_loc  fileloc // struct {x, y int} // oh lol this works
	Src_file string
	Raw_text string // stupid of me to not include that wtf lmao
	// is this my first time??
}

// ill use regex based matching, see 'BraCkish' for a C implementation for a LISP
const (
	//LexemeNull     LexemeType = iota
	LexemeComment  LexemeType = iota
	LexemeOperator            // + - -> . ° x..y ...
	LexemeKeyword             // word bro, word
	LexemeWhiteSpc            // part of me wants to make operator precedence based on whitespace
	LexemeDelim               //  { } [ ] ( ) ... split because makes parsing easier
	LexemeNumeric             // checks for one . specifically, or an f
	LexemeString              // special parse rule, don't split on space / different type
	LexemeSymbol              // user variables, follows a scheme
	// LexemeBool     // Booleans don't exist, https://en.wikipedia.org/wiki/Church_encoding
	//LexemeLiteral  // not this and instead ...Numeric and ...String because this makes parsing easier down the line
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
// I WROTE THIS COMMENT BUT I DIDNT FOLLOW IT
var LexemePatterns = [][]string{
	{
		"\\/\\/.*\n?", // according to regex101, this _should_ work
		"\\/\\*.*?\\*\\/", // https://regex101.com/
	}, // LexemeComment.
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
		regexp.QuoteMeta("\\"),  //     nyi: lambdas
		regexp.QuoteMeta("`"),   //     nyi: homoiconicity
		regexp.QuoteMeta("'"),   // char
		regexp.QuoteMeta("#"),   //     nyi: to raw bytes
	}, // LexemeOperator
	{
		"nil",        // null, active sentinel value
		"func",       // for functions like in go
		"class",      //     nyi: classes, PODs by default
		"import",     //     nyi: self explanatory
		"pure",       //     nyi: optimisation directive
		"static",     //     nyi: non local storage declaration
		"debug",      // prints the AST / relevant info at that location
		"assert",     // crashes on bool expression being false
		"switch",     //     nyi: hash search for conditions
		"case",       //     nyi: individual cases for the hash search
		"defer",      //     nyi: execution at the end of scope
		"continue",   //     nyi: syntax sugar for ~goto~ loop / if not applies:
		"delete",     //     nyi: storage deletion
		"false",      // desugars to 0
		"true",       // desugars to 1
		"for",        // init + condition + step based loop
		"while",      //     nyi: condition based loop, can be chained with do
		"until",      //     nyi: reverse while but can be chained with do
		"do",         //     nyi: execution + check based loop
		"break",      //     nyi: stop loop execution
		"goto",       //     nyi: deprecated (XDXD)
		"lngjmp",     //     nyi: goto + explicit side entry into different scopes
		"srtjmp",     //     nyi: goto + same scope / global scope
		"if",         // branching
		"when",       //     nyi: variable watching
		"throw",      // toss an error yo
		"try",        // try but allow me to baseball huh an error
		"catch",      // baseball huh an error
		"return",     // gives back a value to the callsite
		"processing", //     nyi: desugars into return + lngjmp back into the function; yield
		"constexpr",  //     nyi: this expression can be evaluated at compiletime
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
		"[0-9]*(\\.)?[0-9]+", // https://regex101.com/
	}, // LexemeNumeric
	{
		"\"(.+(\\\")?)*\"", // https://regex101.com/
	}, // LexemeString
	{
		"[_a-zA-Z][_a-zA-Z0-9]*",
	}, // LexemeSymbol
} // [LexemeType][RegexPattern]

func PrintLexeme(lexeme Lexeme) {
	//fmt.Printf("%s @ (%d, %d) e %s, '%s'\n", lexer.LexemeTypeAsString(lexeme.Ltype), lexeme.Src_loc.X, lexeme.Src_loc.Y, lexeme.Src_file, lexeme.Raw_text)
	fmt.Printf("%s '%s'\n", LexemeTypeAsString(lexeme.Ltype), lexeme.Raw_text)
}

func LexemeTypeAsString(lt LexemeType) string {
	switch lt {
	case LexemeComment:
		return "LexemeComment"
	case LexemeOperator:
		return "LexemeOperator"
	case LexemeKeyword:
		return "LexemeKeyword"
	case LexemeWhiteSpc:
		return "LexemeWhiteSpc"
	case LexemeDelim:
		return "LexemeDelim"
	case LexemeNumeric:
		return "LexemeNumeric"
	case LexemeString:
		return "LexemeString"
	case LexemeSymbol:
		return "LexemeSymbol"
	// case LexemeBool:
	// return "LexemeBool"
	// case LexemeLiteral:
	// return "LexemeLiteral"
	case LexemeEOF:
		return "LexemeEOF"
	case LexemeError:
		return "LexemeError"
	default:
		return "something went SERIOUSLY wrong"
	}
}

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

// a normalised file location is the raw byte index
// denormalised includes line and column
// it is cheaper to just pass the array instead of recomputing it each time
func denormaliseFileLoc(fcons []string, loc int) (fileloc, error) {
	line := 0
	// stupid mistake
	for file_idx := 0; file_idx < len(fcons) && file_idx > len(fcons[file_idx]); file_idx++ {
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
