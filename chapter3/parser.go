package main

import (
	"errors"
	"fmt"
	"log"
)

// page 53, Pattern 5:
// Backtracking Parser

// Grammar to be parsed (ANTLR syntax):
//
// grammar NestedNameListWithParallelAssign;
// stat 	: list EOF | assign EOF ;
// assign	: list '=' list ;
// list     : '[' elements ']' ;       		// match bracketed list
// elements : element (',' element)* ;		// match comma-separated list
// element  : NAME '=' NAME | NAME | list ;	// match assignment such as a=b
// NAME     : ('a'..'z'|'A'..'Z')+ ;   		// NAME is sequence of >=1 letter

var SyntaxError = errors.New("syntax error")

type BacktrackingParser struct {
	input     *Lexer
	lookahead []Token // circular lookahead buffer
	pos       int     // circular index of next token position to fill
}

// Returns a new Backtracking Parser with k lookahead symbols (length of the buffer)
func NewBacktrackingParser(l *Lexer, k int) *BacktrackingParser {
	buf := make([]Token, k)
	p := &BacktrackingParser{input: l, lookahead: buf}

	// fill the buffer with first k tokens
	for range k {
		p.consume()
	}
	log.Printf("initial: %+v\n", p)
	return p
}

func (p *BacktrackingParser) stat() {
	p.list()
	p.match(EOF)
}

func (p *BacktrackingParser) list() {
	p.match(LBrack)
	p.elements()
	p.match(RBrack)
}

func (p *BacktrackingParser) elements() {
	p.element()
	for p.peek(1).Type == Comma {
		p.match(Comma)
		p.element()
	}
}

// element needs 2 lookahead tokens to make a decision on whether it's an
// assignment or not.
func (p *BacktrackingParser) element() {
	first, second := p.peek(1), p.peek(2)

	if first.Type == Name && second.Type == Equals {
		p.match(Name)
		p.match(Equals)
		p.match(Name)
	} else if first.Type == Name {
		p.match(Name)
	} else if first.Type == LBrack && second.Type != EOF {
		p.list()
	} else {
		err := fmt.Errorf("%w: expecting name or list, found %+v", SyntaxError, first.Type)
		panic(err)
	}
}

// peek returns the nth next Token in the lookahead buffer.
func (p *BacktrackingParser) peek(n int) Token {
	index := p.pos + n - 1
	if index == len(p.lookahead) {
		index = 0
	}
	return p.lookahead[index]
}

// match checks if the current lookahead token if of the type we're looking for.
// Goes to the next token if it is or reports an error if it isn't.
func (p *BacktrackingParser) match(typ TokenType) {
	log.Printf("lookahead buf: %v, position: %d, want to match: %s", p.lookahead, p.pos, typ)
	tok := p.peek(1)
	log.Printf("peeked: %v", tok)
	if tok.Type == typ {
		// go to next token
		p.consume()
	} else {
		err := fmt.Errorf("match: %w: expecting %v, got %v", SyntaxError, typ, tok.Type)
		panic(err)
	}
}

func (p *BacktrackingParser) consume() {
	tok, err := p.input.Next()
	if err != nil {
		panic(err)
	}

	p.lookahead[p.pos] = tok

	// add 1 until we reach end of the lookahead buffer, then wrap around to 0
	p.pos++
	if p.pos == len(p.lookahead) {
		p.pos = 0
	}
}
