package lexer

import (
	_ "errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type LexemeType int
type LexemeLength int // this is a trick i do in mock-professional C
// i typedef types so the type itself carries information about the use case

// ill use regex based matching, see 'BraCkish' for a C implementation for a LISP
const (
	LexemeNull     LexemeType = iota
	LexemeOperator            // + - -> . ° x..y ...
	LexemeSymbol              // word bro, word
	LexemeWhiteSpc            // part of me wants to make operator precedence based on whitespace
	LexemeDelim               //  { } [ ] ( ) ... split because makes parsing easier
	LexemeNumeric             // checks for one . specifically, or an f
	LexemeString              // special parse rule, don't split on space / different type
	LexemeComment
	// LexemeBool // Booleans don't exist, https://en.wikipedia.org/wiki/Church_encoding
	//LexemeLiteral // not this and instead ...Numeric and ...String because this makes parsing easier down the line
	LexemeError // that's useful. not a pattern itself but useful for parsing
)

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
var lexeme_patterns = [][]string {
	{
		"->",  //     nyi: copy | proof / implication
		"<-",  //     nyi: storage | suggestion
		">>=", //     nyi: bind
		"=<<", //     nyi: deferred application
		".",   // member access
		",",   //     nyi: structural assignment
		"_",   //     nyi: catch all
		"°",   //     nyi: after operator
		";",   // null statement
		"?",   //     nyi: ternary prefix
		":",   //     nyi: ternary suffix
		"^",   // xor operator
		"~",   // negation
		"|",   // chain | logical or
		"&",   // reference | logical and
		"||",  // bitwise or
		"&&",  // bitwise and
		"==",  // equality
		"<",   // lesser
		">",   // greater
		">=",  // lesser or equal
		"<=",  // greater or equal
		"=",   // assignment
		"+",   // addition
		"-",   // negation / subtraction
		"*",   // multiplication
		"/",   // division
		"%",   // modulation
		"\\",  //     nyi: idk yet
		"`",   //     nyi: homoiconicity
		"'",   // char
		"#",   //     nyi: to raw bytes
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
		"processing", //     nyi: desugars into return + lngjmp back into the function; yield
		"constexpr",  //     nyi: this expression can be evaluated at compiletime
	}, // LexemeSymbol
	nil,   // TODO: nil means there is an ad hoc pattern... that is incredibly bad but idk how else
	{
		"(",
		")",
		"[",
		"]",
		"{",
		"}",
	}, // LexemeDelim
	nil,
	nil,
	nil,
}

type fileloc struct {
	x, y int
}

type Lexeme struct {
	ltype    LexemeType
	src_loc  fileloc // struct {x, y int} // oh lol this works
	src_file string
}

func getNextLexeme(fcons string, idx int) (LexemeType, LexemeLength, error) {
	for patterns := range len(lexeme_patterns) {
		if lexeme_patterns[patterns] == nil {
			pattern := "$"
			switch patterns {
			case int(LexemeWhiteSpc): {
				pattern = "\\s"
			}
			case int(LexemeNumeric): {
				pattern = "[0-9]*(?.)[0-9]+"
			}
			case int(LexemeString): {
				pattern = "\".*\""
			}
			}
			re, _ := regexp.Compile(pattern)
			// for this to work I would need to find the length of the subexpr
			// but Go doesn't seem to have a regex function which contains this
			// https://whalelogic.io/posts/references/regex-go-guide/
			// 
		} else {
			for pattern := range len(lexeme_patterns[patterns]) {
				// https://zetcode.com/golang/regexp-quotemeta/
				re := regexp.MustCompile(regexp.QuoteMeta(lexeme_patterns[patterns][pattern]))
				// https://stackoverflow.com/questions/28886616/convert-array-to-slice-in-go
				// https://www.geeksforgeeks.org/go-language/strings-in-golang/
				// len somehow didn't work lul, but it was a casting issue
				if re.MatchString(fcons[idx:]) {
					return LexemeType(patterns), LexemeLength(len(lexeme_patterns[patterns][pattern])), nil
				}
			}
		}
	}
	return LexemeNull, 0, nil
}

// https://gobyexample.com/if-else

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
