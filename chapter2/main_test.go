package main

import (
	"fmt"
	"io"
	"log"
)

func ExampleOne() {
	ex := `[a,b,c]`
	l := NewLexer(ex)
	log.SetOutput(io.Discard)
	for tok, err := l.Next(); tok.Type != EOF; tok, err = l.Next() {
		if err != nil {
			panic(err)
		}
		fmt.Println(tok.Text)
	}
	// output:
	//[
	//a
	//,
	//b
	//,
	//c
	//]
}

func ExampleTwo() {
	ex := `[a,[b,c],d]`
	l := NewLexer(ex)
	log.SetOutput(io.Discard)
	for tok, err := l.Next(); tok.Type != EOF; tok, err = l.Next() {
		if err != nil {
			panic(err)
		}
		fmt.Println(tok.Text)
	}
	// output:
	//[
	//a
	//,
	//[
	//b
	//,
	//c
	//]
	//,
	//d
	//]
}
