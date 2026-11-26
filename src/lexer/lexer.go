package lexer

import (
	_ "errors"
	"fmt"
	"regexp"
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


func Lex(filepath string) ([]Lexeme, error) {
	return lexOpt2(filepath)
}
