package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

// page 27, Pattern 1:
// Mapping Grammars to Recursive-Descent Recognizers

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
	Type int
	Text string
}

// Token types
const (
	EOF int = iota
	LBrack
	RBrack
	Name
	Comma
)

// Lexer goes through the input rune by rune and produces Tokens
type Lexer struct {
	input string // entire input
	p     int    // current position
	cur   rune   // current rune
}

func NewLexer(input string) *Lexer {
	p := 0                  // just for clarity, zero-value of int is already 0
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

func (l *Lexer) Next() (Token, error) {
	for l.cur != inputEOF {
		log.Printf("%+v\n", l)
		switch l.cur {
		case ' ', '\t', '\n', '\r':
			log.Printf("ws %+v\n", l)
			l.whitespace()
			continue
		case ',':
			log.Printf("comma %+v\n", l)
			l.consume()
			return Token{Type: Comma, Text: ","}, nil
		case '[':
			log.Printf("lbrack %+v\n", l)
			l.consume()
			return Token{Type: LBrack, Text: "["}, nil
		case ']':
			log.Printf("rbrack %+v\n", l)
			l.consume()
			return Token{Type: RBrack, Text: "]"}, nil
		default:
			log.Printf("name %+v\n", l)
			if isLetter(l.cur) {
				return l.name()
			}
			return Token{}, fmt.Errorf("invalid character: %c", l.cur)
		}
	}
	log.Printf("EOF %+v\n", l)
	return Token{Type: EOF}, nil
}

// Grammar doesn't have a lexical rule for whitespace but this functions works
// as if it was a lexical rule, just ignoring whitespaces (consumes without
// doing anything).
func (l *Lexer) whitespace() {
	for l.cur == ' ' || l.cur == '\t' || l.cur == '\n' || l.cur == '\r' {
		l.consume()
	}
}

// Lexical rule NAME. The string builder accumulates all the consecutive letters
// into a token.
func (l *Lexer) name() (Token, error) {
	s := &strings.Builder{}
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
	log.Printf("consume %+v\n", l)
}

func main() {
	ex := `[a,b,c]`
	l := NewLexer(ex)
	log.SetOutput(io.Discard)
	for tok, err := l.Next(); tok.Type != EOF; tok, err = l.Next() {
		if err != nil {
			panic(err)
		}
		fmt.Println(tok.Text)
	}
	// output:
	//[
	//a
	//,
	//b
	//,
	//c
	//]
}
