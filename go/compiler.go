package main

import (
	"fmt"
)


func compile(source *string) {
	initScanner(source)
	var line int = -1
	for {
		var token Token = scanToken()
		if (token.line != line) {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%d %d'\n", token.tokenType, token.length, token.start)

		if (token.tokenType == TOKEN_EOF) { break }
	}
}
