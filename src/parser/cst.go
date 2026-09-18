package parser

import "github.com/simon-or-something/Prac-Go-Compiler1/src/lexer"

type comment struct {
	content string
	lexeme lexer.Lexeme
}
type MultilineComment struct { comment }
type LineComment struct { comment }

type operator struct {
}
type MemberAccess struct { operator } // this can be a binary operator but this is AST stuff
type NullStmt struct { operator }
type XorOperator struct { operator }

