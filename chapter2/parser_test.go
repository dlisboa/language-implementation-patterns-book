package main

import "testing"

func TestParseList(t *testing.T) {
	cases := []struct {
		name  string
		input string
		err   error
	}{
		{name: "simple list", input: "[a]", err: nil},
		{name: "long list", input: "[a,b,c,d]", err: nil},
		{name: "list within a list", input: "[a,[b],c]", err: nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			p := NewParser(l)
			p.list()
			if p.err != tc.err {
				t.Errorf("expected no error, got: %v", p.err)
			}
		})
	}
}
