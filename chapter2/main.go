package main

import (
	"fmt"
)

func main() {
	ex := `  [  a, 		b,c]`
	l := NewLexer(ex)
	for tok, err := l.Next(); tok.Type != EOF; tok, err = l.Next() {
		if err != nil {
			panic(err)
		}
		fmt.Println(tok.Text)
	}
}
