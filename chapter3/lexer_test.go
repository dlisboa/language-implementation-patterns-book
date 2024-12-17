package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLexerGoodInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "empty list",
			input: "",
			want: []Token{
				{Type: EOF, Text: ""},
			},
		},
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
				{Type: EOF, Text: ""},
			},
		},
		{
			name:  "complicated list",
			input: "[a,[b,c]]",
			want: []Token{
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "a"},
				{Type: Comma, Text: ","},
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "b"},
				{Type: Comma, Text: ","},
				{Type: Name, Text: "c"},
				{Type: RBrack, Text: "]"},
				{Type: RBrack, Text: "]"},
				{Type: EOF, Text: ""},
			},
		},
		{
			name:  "ignore whitespaces",
			input: "    [a]",
			want: []Token{
				{Type: LBrack, Text: "["},
				{Type: Name, Text: "a"},
				{Type: RBrack, Text: "]"},
				{Type: EOF, Text: ""},
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
				{Type: EOF, Text: ""},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			var tokens []Token
			for l.Scan() {
				token, err := l.Next()
				if err != nil {
					t.Error(err)
				}
				tokens = append(tokens, token)
			}
			if !cmp.Equal(tokens, tc.want) {
				t.Error(cmp.Diff(tokens, tc.want))
			}
		})
	}
}

func TestLexerBadInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "non ASCII input",
			input: "🙈",
		},
		{
			name:  "non letters",
			input: "[1,2,3]",
		},
		{
			name:  "unrecognizable text",
			input: "\x05R\xDF\xD8",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			var err error
			for l.Scan() {
				_, err = l.Next()
				t.Log(err)
			}
			if err == nil {
				t.Error("want: error, got: nil")
			}
		})
	}
}
