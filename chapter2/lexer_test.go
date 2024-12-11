package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Token
	}{
		{name: "empty input", input: "", want: nil},
		{
			name:  "simple list",
			input: "[a,b,c]",
			want: []Token{
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "a"},
				{Type: Comma, Text: ","},
				{Type: Name, Text: "b"},
				{Type: Comma, Text: ","},
				{Type: Name, Text: "c"},
				{Type: RBrack, Text: "]"},
			},
		},
		{
			name:  "complicated list",
			input: "[a,[b,c],d]",
			want: []Token{
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "a"},
				{Type: Comma, Text: ","},
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "b"},
				{Type: Comma, Text: ","},
				{Type: Name, Text: "c"},
				{Type: RBrack, Text: "]"},
				{Type: Comma, Text: ","},
				{Type: Name, Text: "d"},
				{Type: RBrack, Text: "]"},
			},
		},
		{
			name:  "ignore whitespaces",
			input: "    [a]",
			want: []Token{
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "a"},
				{Type: RBrack, Text: "]"},
			},
		},
		{
			name:  "list within a list",
			input: "[[a,b]]",
			want: []Token{
				{Type: LBrack, Text: "["},
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "a"},
				{Type: Comma, Text: ","},
				{Type: Name, Text: "b"},
				{Type: RBrack, Text: "]"},
				{Type: RBrack, Text: "]"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			var tokens []Token
			for tok, err := l.Next(); tok.Type != EOF; tok, err = l.Next() {
				if err != nil {
					t.Fatal(err)
				}
				tokens = append(tokens, tok)
			}
			if !cmp.Equal(tokens, tc.want) {
				t.Error(cmp.Diff(tokens, tc.want))
			}
		})
	}
}
