package main

import (
	"fmt"
)

func main() {
	ex := `  [  a, 		b,c]`
	l := NewLexer(ex)
	for l.Scan() {
		tok, err := l.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(tok.Type, tok.Text)
	}
}
