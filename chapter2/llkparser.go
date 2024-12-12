package main

import (
	"fmt"
)

// page 41, Pattern 3:
// LL(k) Recursive-Descent Parser

// Grammar to be parsed (ANTLR syntax):
//
// grammar NestedNameList;
// list     : '[' elements ']' ;       // match bracketed list
// elements : element (',' element)* ; // match comma-separated list
// element  : NAME '=' NAME            // match assignment such as a=b
//			| NAME
//			| list
//			;
// NAME     : ('a'..'z'|'A'..'Z')+ ;   // NAME is sequence of >=1 letter

// We need two state variables to keep track of the parse state: an input token
// stream and a lookahead circular buffer. To report parse errors we could
// panic, but here we'll just use a variable to track it, though this isn't the
// optimal solution (it only reports the last error and does not stop the
// parser).
type LLkParser struct {
	input *Lexer
	buf   []Token // circular lookahead buffer
	k     int     // how many lookahead symbols (length of the buffer)
	pos   int     // circular index of next token position to fill
	err   error
}

func NewLLkParser(l *Lexer, k int) *LLkParser {
	buf := make([]Token, k)
	p := &LLkParser{input: l, buf: buf, k: k}

	// initialize the buffer with first k tokens
	for range k {
		p.consume()
	}
	return p
}

func (p *LLkParser) list() {
	p.match(LBrack)
	p.elements()
	p.match(RBrack)
}

func (p *LLkParser) elements() {
	p.element()
	for p.lookahead(1).Type == Comma {
		p.match(Comma)
		p.element()
	}
}

// element needs 2 lookahead tokens to make a decision on whether it's an
// assignment or not.
func (p *LLkParser) element() {
	first, second := p.lookahead(1), p.lookahead(2)

	if first.Type == Name && second.Type == Equals {
		p.match(Name)
		p.match(Equals)
		p.match(Name)
	} else if first.Type == Name {
		p.match(Name)
	} else if first.Type == LBrack {
		p.list()
	} else {
		p.err = fmt.Errorf("%w: expecting name or list, found %+v", SyntaxError, p.lookahead(1).Type)
	}
}

// lookahead returns the nth next Token in the buffer. This kind of method is
// often called `peek()`
func (p *LLkParser) lookahead(n int) Token {
	index := (p.pos + n - 1) % p.k
	return p.buf[index]
}

// match checks if the current lookahead token if of the type we're looking for.
// Goes to the next token if it is or reports an error if it isn't.
func (p *LLkParser) match(typ TokenType) {
	if p.lookahead(1).Type == typ {
		// go to next token
		p.consume()
	} else {
		p.err = fmt.Errorf("%w: expecting %v, got %v", SyntaxError, typ, p.lookahead(1).Type)
	}
}

func (p *LLkParser) consume() {
	tok, err := p.input.Next()

	p.buf[p.pos] = tok
	// add 1 until we reach k, then wraps around to 0
	p.pos = (p.pos + 1) % p.k

	// if at the end of token input stream don't assign to err otherwise we
	// overwrite the last error
	if tok.Type != EOF {
		p.err = err
	}
}
