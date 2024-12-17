package main

import (
	"errors"
	"fmt"
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
	pos       int     // position into lookahead buffer
	markers   []int   // stack of positions into lookahead buffer
}

// Returns a new Backtracking Parser with k lookahead symbols (length of the buffer)
func NewBacktrackingParser(l *Lexer) *BacktrackingParser {
	p := &BacktrackingParser{input: l}

	p.sync(1)
	return p
}

func (p *BacktrackingParser) stat() {
	if p.speculateList() {
		p.list()
		p.match(EOF)
	} else if p.speculateAssign() {
		p.assign()
		p.match(EOF)
	} else {
		tok := p.peek(1)
		err := fmt.Errorf("%w: expecting list or assign, found %v", SyntaxError, tok.Type)
		panic(err)
	}
}

func (p *BacktrackingParser) speculateList() bool {
	success := true
	p.mark()

	defer func() {
		if r := recover(); r != nil {
			success = false
			p.release()
		}
	}()

	p.list()
	p.match(EOF)

	p.release()
	return success

}

func (p *BacktrackingParser) speculateAssign() bool {
	success := true
	p.mark()

	defer func() {
		if r := recover(); r != nil {
			success = false
			p.release()
		}
	}()

	p.assign()
	p.match(EOF)

	p.release()
	return success

}

func (p *BacktrackingParser) assign() {
	p.list()
	p.match(Equals)
	p.list()
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

// mark pushes the currenct position into the stack so we can backtrack to it
// later
func (p *BacktrackingParser) mark() {
	p.markers = append(p.markers, p.pos)
}

// release pops the last marker from the stack and backtracks to it
func (p *BacktrackingParser) release() {
	position := p.markers[len(p.markers)-1]
	p.markers = p.markers[:len(p.markers)-1] // pop
	p.seek(position)
}

func (p *BacktrackingParser) seek(position int) {
	p.pos = position
}

func (p *BacktrackingParser) isSpeculating() bool {
	return len(p.markers) > 0
}

func (p *BacktrackingParser) sync(i int) {
	if p.pos+i-1 > (len(p.lookahead) - 1) {
		n := p.pos + i - 1 - (len(p.lookahead) - 1)
		p.fill(n)
	}
}

func (p *BacktrackingParser) fill(n int) {
	for range n {
		tok, err := p.input.Next()
		if err != nil {
			panic(fmt.Errorf("fill: error reading next token: %w", err))
		}
		p.lookahead = append(p.lookahead, tok)
	}
}

// peek returns the nth next Token in the lookahead buffer.
func (p *BacktrackingParser) peek(n int) Token {
	p.sync(n)
	index := p.pos + n - 1
	if index == len(p.lookahead) {
		index = 0
	}
	return p.lookahead[index]
}

// match checks if the current lookahead token if of the type we're looking for.
// Goes to the next token if it is or reports an error if it isn't.
func (p *BacktrackingParser) match(typ TokenType) {
	// log.Printf("lookahead buf: %v, position: %d, want to match: %s", p.lookahead, p.pos, typ)
	tok := p.peek(1)
	// log.Printf("peeked: %v", tok)
	if tok.Type == typ {
		// go to next token
		p.consume()
	} else {
		err := fmt.Errorf("match: %w: expecting %v, got %v", SyntaxError, typ, tok.Type)
		panic(err)
	}
}

func (p *BacktrackingParser) consume() {
	p.pos++
	// we're not speculating and we've hit the end of the lookahead buffer?
	if !p.isSpeculating() && p.pos == len(p.lookahead) {
		p.pos = 0
		p.lookahead = p.lookahead[:0] // reset lookahead buffer
	}
	p.sync(1)
}
