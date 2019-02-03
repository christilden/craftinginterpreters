package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)


func main() {
	if (len(os.Args) == 1) {
		repl()
	} else if (len(os.Args) == 2) {
		runFile(os.Args[1]);
	} else {
		fmt.Printf("Usage: glox [path]\n")
		os.Exit(64)
	}
}

// SimpleReadLine is simple replacement for GNU readline.
// prompt is the command prompt to print before reading input.
func SimpleReadLine(prompt string) (string, error) {
	fmt.Printf(prompt)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')

	if err == nil {
		line = strings.TrimRight(line, "\r\n")
	}
	return line, err
}

func repl() {
	for {
		line, err := SimpleReadLine("> ")
		if err != nil {
			log.Fatal(err)
		}

		interpret(&line)
	}
}

func runFile(path string) {
	var source string = readFile(path)
	result := interpret(&source)

	if (result == INTERPRET_COMPILE_ERROR) { os.Exit(65) }
	if (result == INTERPRET_RUNTIME_ERROR) { os.Exit(70) }
}

func readFile(path string) string {
	dat, err := ioutil.ReadFile("/tmp/dat")
	check(err)
	return string(dat)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
