package main

import (
	"fmt"
	"strings"
)

// page 31, Pattern 2:
// LL(1) Recursive-Descent Lexer

// Implementation
//
// Rules:
// * for each rule r in grammar G, create a method r() with the same name
// * if the rule has a subrule, call the subrule method from within r()
//
// Subrules:
// * alternatives like `rule := a | b | c` become either a switch or if-else
//   sequence for branches a, b or c
// * we test whether the branches a, b or c apply by looking at the next token
// * if a subrule is optional, remove the `default` or `else` part of the branch
//
// Tokens:
// * tokens have a type (a kind)
// * all references to token of type T becomes a call to match(T)
// * match(T) is a helper function/method that consumes a token if T is the
//   current lookahead token
// * if there's a mismatch and the current token is not T, match(T) returns an
//   error
// * we define the types of tokens with each type being descring by an integer
// * we define some meta token types like EOF and INVALID_TOKEN
//
// Lexer:
// * lexical rules are functions/methods in the lexer
// * it has a `next()` function that produces the next token
// * this function looks at the next char and determines what it could be, then
//   calls the appropriate lexical rule function

// Grammar to be matched (ANTLR syntax):
//
// grammar NestedNameList;
// list     : '[' elements ']' ;       // match bracketed list
// elements : element (',' element)* ; // match comma-separated list
// element  : NAME | list ;            // element is name or nested list
// NAME     : ('a'..'z'|'A'..'Z')+ ;   // NAME is sequence of >=1 lette
//
// example strings that should be matched
// [a,b,c]
// [a,[b,c],d]

type Token struct {
	Type TokenType
	Text string
}

type TokenType int

// Token types
const (
	EOF TokenType = iota
	LBrack
	RBrack
	Name
	Comma
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case LBrack:
		return "LBrack"
	case RBrack:
		return "RBrack"
	case Name:
		return "Name"
	case Comma:
		return "Comma"
	default:
		return "Unknown"
	}
}

// Lexer goes through the input rune by rune and produces Tokens. Lexers are
// also called "scanners" or "tokenizers".
type Lexer struct {
	input string // entire input
	p     int    // current position
	cur   rune   // current rune
}

func NewLexer(input string) *Lexer {
	p := 0 // just for clarity, zero-value of int is already 0
	if input == "" {
		return &Lexer{input: input, p: p, cur: inputEOF}
	}
	cur := []rune(input)[p] // convert string to a rune slice so that indexing is rune-based and not byte-based
	return &Lexer{input: input, p: p, cur: cur}
}

// isLetter is a helper function, only recognizes ASCII letters
func isLetter(r rune) bool {
	return r >= 'a' && r < 'z' || r >= 'A' && r <= 'Z'
}

// marks the end of input, if input was a Reader we wouldn't need this as we
// have io.EOF
var inputEOF = rune(-1)

// Next return the next Token at each invocation or an error if the input cannot
// be recognized.
func (l *Lexer) Next() (Token, error) {
	for l.cur != inputEOF {
		switch l.cur {
		case ' ', '\t', '\n', '\r':
			l.consume()
			continue
		case ',':
			l.consume()
			return Token{Type: Comma, Text: ","}, nil
		case '[':
			l.consume()
			return Token{Type: LBrack, Text: "["}, nil
		case ']':
			l.consume()
			return Token{Type: RBrack, Text: "]"}, nil
		default:
			if isLetter(l.cur) {
				return l.name()
			}
			return Token{}, fmt.Errorf("invalid character: %c", l.cur)
		}
	}
	return Token{Type: EOF}, nil
}

// Lexical rule NAME. The string builder accumulates all the consecutive letters
// into a token.
func (l *Lexer) name() (Token, error) {
	var s strings.Builder
	for isLetter(l.cur) {
		s.WriteRune(l.cur)
		l.consume()
	}

	return Token{Type: Name, Text: s.String()}, nil
}

// Consume moves the current position forward by one and saves the next current
// rune. The check for input length wouldn't be this way if we were using a
// Reader-based iteration, we could just check for io.EOF instead.
func (l *Lexer) consume() {
	l.p += 1
	if l.p >= len(l.input) {
		// signals end of input as we're not using Reader
		l.cur = inputEOF
	} else {
		// saves the next rune
		l.cur = []rune(l.input)[l.p]
	}
}
