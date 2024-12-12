package main

import (
	"errors"
	"testing"
)

func TestParseAssignment(t *testing.T) {
	cases := []struct {
		name  string
		input string
		err   error
	}{
		{name: "simple assignment", input: "[a,b=c,[d,e]]", err: nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			// change this to 1 as an exercise and see that the input is no
			// longer parsed correctly
			p := NewLLkParser(l, 2)
			p.list()
			if !errors.Is(p.err, tc.err) {
				t.Errorf("expected %v, got: %v", tc.err, p.err)
			}
		})
	}
}
