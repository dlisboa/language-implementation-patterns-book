package main

import (
	"errors"
	"fmt"
)

// page 36, Pattern 3:
// LL(1) Recursive-Descent Parser

// Grammar to be parsed (ANTLR syntax):
//
// grammar NestedNameList;
// list     : '[' elements ']' ;       // match bracketed list
// elements : element (',' element)* ; // match comma-separated list
// element  : NAME | list ;            // element is name or nested list
// NAME     : ('a'..'z'|'A'..'Z')+ ;   // NAME is sequence of >=1 lette

// We need two state variables to keep track of the parse state: an input token
// stream and a lookahead buffer. In this case we can use a single lookahead
// variable instead of a buffer. To report parse errors we could panic, but here
// we'll just use a variable to track it, though this isn't the optimal solution
// (it only reports the last error and does not stop the parser).
type Parser struct {
	input     *Lexer
	lookahead Token
	err       error
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{input: l}
	// initialize the parser with the first token, otherwise it'll be the
	// zero-value for Token which is EOF
	p.lookahead, p.err = p.input.Next()
	return p
}

func (p *Parser) list() {
	p.match(LBrack)
	p.elements()
	p.match(RBrack)
}

func (p *Parser) elements() {
	p.element()
	for p.lookahead.Type == Comma {
		p.match(Comma)
		p.element()
	}
}

var SyntaxError = errors.New("syntax error")

func (p *Parser) element() {
	switch p.lookahead.Type {
	case Name:
		p.match(Name)
	case LBrack: // we've found a sublist
		p.list()
	default:
		p.err = fmt.Errorf("%w: expecting name or list, found %+v", SyntaxError, p.lookahead)
	}
}

// match checks if the current lookahead token if of the type we're looking for.
// Goes to the next token if it is or reports an error if it isn't.
func (p *Parser) match(typ TokenType) {
	if p.lookahead.Type == typ {
		// go to next token
		p.consume()
	} else {
		p.err = fmt.Errorf("%w: expecting %v, got %v", SyntaxError, typ, p.lookahead.Type)
	}
}

func (p *Parser) consume() {
	tok, err := p.input.Next()
	// if at the end of token input stream, stop consuming. Cannot assign to
	// err otherwise we overwrite the last error
	if tok.Type == EOF {
		return
	}

	p.lookahead, p.err = tok, err
}
