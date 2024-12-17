package main

import (
	"bytes"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestParserGoodInput(t *testing.T) {
	ar, err := txtar.ParseFile("testdata/good.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range ar.Files {
		lines := bytes.Split(file.Data, []byte("\n"))

		for _, line := range lines {
			if len(line) == 0 {
				continue
			}

			t.Run(file.Name, func(t *testing.T) {
				testcase := string(line)
				t.Logf("parse string: %q\n", testcase)

				lexer := NewLexer(testcase)
				parser := NewBacktrackingParser(lexer)
				defer func() {
					err := recover()
					if err != nil {
						t.Errorf("got error on parse string: %q, error: %q", testcase, err)
					}
				}()
				parser.stat()
			})
		}
	}
}

func TestParserBadInput(t *testing.T) {
	ar, err := txtar.ParseFile("testdata/bad.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range ar.Files {
		lines := bytes.Split(file.Data, []byte("\n"))

		for _, line := range lines {
			if len(line) == 0 {
				continue
			}

			t.Run(file.Name, func(t *testing.T) {
				testcase := string(line)
				t.Logf("parse string: %q\n", testcase)

				lexer := NewLexer(testcase)
				parser := NewBacktrackingParser(lexer)
				defer func() {
					err := recover()
					if err == nil {
						t.Errorf("want error on parse string: %q, got none", testcase)
					}
				}()
				parser.stat()
			})
		}
	}
}
