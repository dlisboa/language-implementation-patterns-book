package main

import (
	"errors"
	"testing"
)

func TestParseList(t *testing.T) {
	cases := []struct {
		name  string
		input string
		err   error
	}{
		{name: "simple list", input: "[a]", err: nil},
		{name: "long list", input: "[a,b,c,d]", err: nil},
		{name: "list within a list", input: "[a,[b],c]", err: nil},
		{name: "incomplete list", input: "[a, ]", err: SyntaxError},
		{name: "incomplete list", input: "[[a, ]", err: SyntaxError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			p := NewParser(l)
			p.list()
			if !errors.Is(p.err, tc.err) {
				t.Errorf("expected %v, got: %v", tc.err, p.err)
			}
		})
	}
}
